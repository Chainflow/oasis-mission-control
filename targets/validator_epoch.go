package targets

import (
	"context"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	beacon "github.com/oasisprotocol/oasis-core/go/beacon/api"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"google.golang.org/grpc"

	"github.com/Chainflow/oasis-mission-control/config"
	"github.com/Chainflow/oasis-mission-control/rpc"
)

// loadBeaconClient loads beacon client and returns it
func loadBeaconClient(socket string) (*grpc.ClientConn, beacon.Backend) {

	// Attempt to load connection with staking client
	connection, beaconClient, err := rpc.BeaconClient(socket)
	if err != nil {
		log.Println("Failed to establish connection to beacon client : ",
			err)
		return nil, nil
	}
	return connection, beaconClient
}

// GetValEpoch returns the work epoch number
func GetValEpoch(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

		// Stop code here faild to establish connection and reply
		log.Printf("Failed to establish connection using socket: %s" +
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

	if &blk == nil {
		log.Printf("Got empty block res : %v", blk)
		return
	}

	bh := blk.Height

	var backend beacon.Backend
	timeSource := (backend).(beacon.SetableBackend)
	ep, err := timeSource.GetEpoch(context.Background(), height)
	log.Fatalf("ep, err", ep, err)

	connection, bo := loadBeaconClient(socket)

	// Return epcoh of specific height
	epoch, err := bo.GetEpoch(context.Background(), bh)
	if err != nil {
		log.Printf("Failed to retrieve Epoch of Block : %v", err)
		return
	}

	err = writeToInfluxDb(c, bp, "oasis_worker_epoch_number", map[string]string{}, map[string]interface{}{"epoch_number": epoch})
	if err != nil {
		log.Printf("Error while storing worker epoch number : %v", err)
		return
	}
	log.Printf("validator worker epoch number : %v", epoch)
}
