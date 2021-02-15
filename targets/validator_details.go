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
		log.Printf("Failed to establish connection to light client "+
			" backend : %v ", err)
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
	socket := cfg.ValidatorDetails.SocketPath

	// Attempt to load connection with consensus client
	connection, co := loadLightClientBackend(socket)

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop here faild to establish connection and reply
		log.Printf("Failed to establish connection using socket: %s" +
			socket)
		return
	}

	var height int64 = consensus.HeightLatest

	blk := GetBlockDetails(cfg, height)
	if blk == nil {
		return
	}

	// Write validator hex address and publick ky into database
	_ = writeToInfluxDb(c, bp, "oasis_validator_desc", map[string]string{}, map[string]interface{}{"hex_address": cfg.ValidatorDetails.ValidatorHexAddress, "val_address": cfg.ValidatorDetails.ValidatorAddress})

	validator, err := co.GetLightBlock(context.Background(), blk.Height)
	if err != nil {
		log.Printf("Failed to query GetLightBlock : %v", err)
		return
	}

	var protoLb tmproto.LightBlock
	if err = protoLb.Unmarshal(validator.Meta); err != nil {
		log.Printf("Error while unmarshelling the lb proto data : %v", err)
		return
	}

	vals, err := tmtypes.ValidatorSetFromProto(protoLb.ValidatorSet)
	if err != nil {
		log.Printf("Error while unmarshelling the validator set data :%v ", err)
		return
	}

	for _, val := range vals.Validators {
		if val.Address.String() == cfg.ValidatorDetails.ValidatorHexAddress {

			log.Printf("val hex address from val set : %s", cfg.ValidatorDetails.ValidatorHexAddress)

			var vp string
			log.Printf("VOTING POWER: %d", val.VotingPower)

			valVotingPow := strconv.FormatInt(val.VotingPower, 10)

			if valVotingPow != "" {
				vp = valVotingPow
			} else {
				vp = "0"
			}
			_ = writeToInfluxDb(c, bp, "oasis_voting_power", map[string]string{}, map[string]interface{}{"power": vp})

			votingPower, err := strconv.Atoi(vp)
			if err != nil {
				log.Printf("Error wile converting string to int of voting power : %v", err)
			}

			if int64(votingPower) < cfg.AlertsThreshold.VotingPowerThreshold {
				_ = SendTelegramAlert(fmt.Sprintf("Your oasis validator's voting power has dropped below %d", cfg.AlertsThreshold.VotingPowerThreshold), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Your oasis validator's voting power has dropped below %d", cfg.AlertsThreshold.VotingPowerThreshold), cfg)
			}
		}
	}
}
