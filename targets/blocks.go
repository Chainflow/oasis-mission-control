package targets

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fxamacker/cbor"
	client "github.com/influxdata/influxdb1-client/v2"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	mint_api "github.com/oasisprotocol/oasis-core/go/consensus/tendermint/api"
	"google.golang.org/grpc"

	"github.com/Chainflow/oasis-mission-control/config"
	"github.com/Chainflow/oasis-mission-control/rpc"
)

// loadConsensusClient loads consensus client and returns it
func loadConsensusClient(socket string) (*grpc.ClientConn,
	consensus.ClientBackend) {

	// Attempt to load connection with consensus client
	connection, consensusClient, err := rpc.ConsensusClient(socket)
	if err != nil {
		log.Printf("Failed to establish connection to consensus"+
			" client : %v", err)
		return nil, nil
	}
	return connection, consensusClient
}

// GetMissedBlock checks if the validator is missing any blocks, if so it will send
// telegram and email alerts to the user about missed blocks
func GetMissedBlock(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	socket := cfg.ValidatorDetails.SocketPath
	// Attempt to load connection with consensus client
	connection, co := loadConsensusClient(socket)

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop here faild to establish connection and reply
		log.Fatalf("Failed to establish connection using socket: %s" +
			socket)
		return
	}

	var height int64 = consensus.HeightLatest

	// Retrieve block at specific height from consensus client
	blk, err := co.GetBlock(context.Background(), height)
	if err != nil {
		log.Printf("Error while getting block info : %v", err)
		return
	}

	bh := blk.Height
	blockHeight := strconv.FormatInt(bh, 10)

	// Responding with retrieved block
	log.Printf("block height : %s", blockHeight)

	var meta mint_api.BlockMeta
	if err := cbor.Unmarshal(blk.Meta, &meta); err != nil {
		log.Printf("Request failed to Unmarshal Block Metadata : %v", err)
		return
	}

	log.Printf("Blok meta header height : %v", meta.LastCommit.Height)
	// Send missed block alerts
	addrExists := false
	for _, c := range meta.LastCommit.Signatures {
		if c.ValidatorAddress.String() == cfg.ValidatorDetails.ValidatorHexAddress {
			addrExists = true
		}
	}

	if !addrExists {
		// Send emergency missed blocks alert
		err := SendEmeregencyAlerts(cfg, c, blockHeight)
		if err != nil {
			log.Printf("Error while sending emergency missed block alerts : %v", err)
		}

		blocks := GetContinuousMissedBlock(cfg, c)
		currentHeightFromDb := GetlatestCurrentHeightFromDB(cfg, c)
		blocksArray := strings.Split(blocks, ",")
		log.Println("blocks length ", int64(len(blocksArray)), currentHeightFromDb)
		//Calling function to store and send single missed block alerts
		err = SendSingleMissedBlockAlert(cfg, c, blockHeight)
		if err != nil {
			log.Println("Error while sending missed block alert ", err)
			return
		}
		if cfg.AlertsThreshold.MissedBlocksThreshold > 1 {
			if int64(len(blocksArray))-1 >= cfg.AlertsThreshold.MissedBlocksThreshold {
				missedBlocks := strings.Split(blocks, ",")
				_ = SendTelegramAlert(fmt.Sprintf("Validator missed blocks from height %s to %s", missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Validator missed blocks from height %s to %s", missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = writeToInfluxDb(c, bp, "oasis_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
				_ = writeToInfluxDb(c, bp, "oasis_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": blockHeight})
				return
			}
			if len(blocksArray) == 1 {
				blocks = blockHeight + ","
			} else {
				rpcBlockHeight, _ := strconv.Atoi(blockHeight)
				dbBlockHeight, _ := strconv.Atoi(currentHeightFromDb)
				diff := rpcBlockHeight - dbBlockHeight
				if diff == 1 {
					blocks = blocks + blockHeight + ","
				} else if diff > 1 {
					blocks = ""
				}
			}
			_ = writeToInfluxDb(c, bp, "oasis_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": blockHeight})
			return

		}
	}
	return
}

// GetContinuousMissedBlock returns the latest missed block from the db
func GetContinuousMissedBlock(cfg *config.Config, c client.Client) string {
	var blocks string
	q := client.NewQuery("SELECT last(block_height) FROM oasis_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						blocks = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return blocks
}

// GetlatestCurrentHeightFromDB returns latest current height from db
func GetlatestCurrentHeightFromDB(cfg *config.Config, c client.Client) string {
	var currentHeight string
	q := client.NewQuery("SELECT last(current_height) FROM oasis_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						currentHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return currentHeight
}

// SendSingleMissedBlockAlert is to send alert about single missed block alerts
func SendSingleMissedBlockAlert(cfg *config.Config, c client.Client, blockHeight string) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	if cfg.AlertsThreshold.MissedBlocksThreshold == 1 {
		err = SendTelegramAlert(fmt.Sprintf("Validator missed a block at block height %s", blockHeight), cfg)
		err = SendEmailAlert(fmt.Sprintf("Validator missed a block at block height %s", blockHeight), cfg)
		err = writeToInfluxDb(c, bp, "oasis_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blockHeight, "range": blockHeight})
		err = writeToInfluxDb(c, bp, "oasis_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blockHeight, "current_height": blockHeight})
		err = writeToInfluxDb(c, bp, "oasis_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blockHeight, "current_height": blockHeight})

		if err != nil {
			return err
		}
	} else {
		err = writeToInfluxDb(c, bp, "oasis_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blockHeight, "current_height": blockHeight})
		if err != nil {
			log.Println("Error while storing missed block height ", err)
			return err
		}
	}

	return nil
}

// GetBlockDetails returns validator block details
func GetBlockDetails(cfg *config.Config, height int64) *consensus.Block {
	socket := cfg.ValidatorDetails.SocketPath
	// Attempt to load connection with consensus client
	connection, co := loadConsensusClient(socket)

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop code here faild to establish connection and reply
		log.Fatalf("Failed to establish connection using socket: " +
			socket)
		return nil
	}

	// Retrieve block at specific height from consensus client
	blk, err := co.GetBlock(context.Background(), height)
	if err != nil {
		log.Printf("Error while getting block info : %v", err)
		return nil
	}

	return blk
}
