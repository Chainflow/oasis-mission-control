package rpc

import (
	"fmt"
	"log"

	beacon "github.com/oasisprotocol/oasis-core/go/beacon/api"
	cmnGrpc "github.com/oasisprotocol/oasis-core/go/common/grpc"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	scheduler "github.com/oasisprotocol/oasis-core/go/scheduler/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
	"google.golang.org/grpc"
)

// ConsensusClient - initiate new consensus client
func ConsensusClient(address string) (*grpc.ClientConn,
	consensus.ClientBackend, error) {
	conn, err := Connect(address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to establish connection "+
			"with node %s", address)
	}

	client := consensus.NewConsensusClient(conn)
	return conn, client, nil
}

// ConsensusLightClientBackend - initiate new consensus light client backend
func ConsensusLightClientBackend(address string) (*grpc.ClientConn,
	consensus.LightClientBackend, error) {
	conn, err := Connect(address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to establish connection "+
			"with node %s", address)
	}

	client := consensus.NewConsensusLightClient(conn)
	return conn, client, nil
}

// SchedulerClient - initiate new scheduler client
func SchedulerClient(address string) (*grpc.ClientConn, scheduler.Backend,
	error) {

	conn, err := Connect(address)
	if err != nil {
		log.Println("Failed to establish Scheduler "+
			"Client Connection with node %s", address)
	}

	client := scheduler.NewSchedulerClient(conn)
	return conn, client, nil
}

// StakingClient - initiate new staking client
func StakingClient(address string) (*grpc.ClientConn, staking.Backend, error) {
	conn, err := Connect(address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to establish connection "+
			"with node %s", address)
	}

	client := staking.NewStakingClient(conn)
	return conn, client, nil
}

// StakingClient - initiate new staking client
func BeaconClient(address string) (*grpc.ClientConn, beacon.Backend, error) {
	conn, err := Connect(address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to establish connection "+
			"with node %s", address)
	}

	client := beacon.NewBeaconClient(conn)
	return conn, client, nil
}

// Connect to grpc
func Connect(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	opts = append(opts, grpc.WithDefaultCallOptions(
		grpc.WaitForReady(false)))

	conn, err := cmnGrpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
