package targets

import (
	"context"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	tmamino "github.com/tendermint/go-amino"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"

	"github.com/Chainflow/oasis-mission-control/config"
	"github.com/Chainflow/oasis-mission-control/rpc"
)

var (
	aminoCodec = tmamino.NewCodec()
)

// GetValidatorsList  loads light client and returns it
func loadLightClientBackend(socket string) (*grpc.ClientConn, consensus.LightClientBackend) {

	connection, lightClient, err := rpc.ConsensusLightClientBackend(socket)
	if err != nil {
		log.Println("Failed to establish connection to light client "+
			" backend : ", err)
		return nil, nil
	}
	return connection, lightClient
}

// GetValidatorsList gives the details of your validator
func GetValidatorsList(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	// func GetValidatorsList(cfg *config.Config, height int64) {
	socket := cfg.SocketPath

	// Attempt to load connection with consensus client
	connection, co := loadLightClientBackend(socket)

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop code here faild to establish connection and reply
		log.Fatalf("Failed to establish connection using socket: " +
			socket)
		return
	}

	var height int64 = consensus.HeightLatest

	blk := GetBlockDetails(cfg, height)

	validator, err := co.GetValidatorBlock(context.Background(), blk.Height)
	if err != nil {
		log.Println("Failed to get Validators!", err)
		return
	}

	var protoVals tmproto.ValidatorSet
	if err = protoVals.Unmarshal(validator.Meta); err != nil {
		log.Println("Error while unmarshelling the validator set data ", err)
		return
	}
	vals, err := tmtypes.ValidatorSetFromProto(&protoVals)
	if err != nil {
		log.Println("Error while unmarshelling the validator set data ", err)
		return
	}

	// Write validator hex address and publick ky into database
	_ = writeToInfluxDb(c, bp, "oasis_validator_desc", map[string]string{}, map[string]interface{}{"hex_address": cfg.ValidatorHexAddress, "val_address": cfg.ValidatorAddress})

	for _, val := range vals.Validators {
		if val.Address.String() == cfg.ValidatorHexAddress {

			log.Println("val desc..", cfg.ValidatorHexAddress)

			var vp string
			fmt.Println("VOTING POWER: \n", val.VotingPower)

			valVotingPow := strconv.FormatInt(val.VotingPower, 10)

			if valVotingPow != "" {
				vp = valVotingPow
			} else {
				vp = "0"
			}
			_ = writeToInfluxDb(c, bp, "oasis_voting_power", map[string]string{}, map[string]interface{}{"power": vp})
			log.Println("Voting Power \n", vp)

			votingPower, err := strconv.Atoi(vp)
			if err != nil {
				log.Println("Error wile converting string to int of voting power \t", err)
			}

			if int64(votingPower) < cfg.VotingPowerThreshold {
				_ = SendTelegramAlert(fmt.Sprintf("Your oasis validator's voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Your oasis validator's voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
			}
		}
	}
}
