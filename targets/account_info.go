package targets

import (
	"context"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
	"google.golang.org/grpc"

	"github.com/Chainflow/oasis-mission-control/config"
	"github.com/Chainflow/oasis-mission-control/rpc"
)

// loadStakingClient loads staking client and returns it
func loadStakingClient(socket string) (*grpc.ClientConn, staking.Backend) {

	// Attempt to load connection with staking client
	connection, stakingClient, err := rpc.StakingClient(socket)
	if err != nil {
		log.Println("Failed to establish connection to staking client : ",
			err)
		return nil, nil
	}
	return connection, stakingClient
}

// GetAccount returns the account descriptor for the given account.
func GetAccount(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	socket := cfg.ValidatorDetails.SocketPath
	connection, co := loadStakingClient(socket) // Attempt to load connection with consensus client

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop here faild to establish connection and reply
		log.Printf("Failed to establish connection using socket: " +
			socket)
		return
	}

	var address staking.Address
	valAddress := cfg.ValidatorDetails.ValidatorAddress

	// Unmarshall text into public key object
	err = address.UnmarshalText([]byte(valAddress))
	if err != nil {

		log.Println("Failed to UnmarshalText into Address", err)
		return
	}

	var height int64 = consensus.HeightLatest
	blk := GetBlockDetails(cfg, height) // get block height from db
	if blk == nil {
		return
	}

	if len(address) != 0 {
		// Create an owner query to be able to retrieve data with regards to account
		query := staking.OwnerQuery{Height: blk.Height, Owner: address}

		// Retrieve account information using created query
		account, err := co.Account(context.Background(), &query)
		if err != nil {

			log.Printf("Failed to get Account : %v ", err)
			return
		}

		balance := account.General.Balance

		_ = writeToInfluxDb(c, bp, "oasis_address_balance", map[string]string{}, map[string]interface{}{"balance": balance})
		log.Printf("Address Balance: %s", balance)
	}
}

// GetSelfDelegationBal is to returns the escrow
func GetSelfDelegationBal(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	socket := cfg.ValidatorDetails.SocketPath
	connection, co := loadStakingClient(socket) // Attempt to load connection with consensus client

	// Close connection once code underneath executes
	defer connection.Close()

	// If null object was retrieved send response
	if co == nil {

		// Stop code here faild to establish connection and reply
		log.Printf("Failed to establish connection using socket: " +
			socket)
		return
	}

	var address staking.Address
	valAddress := cfg.ValidatorDetails.ValidatorAddress

	// Unmarshall text into public key object
	err = address.UnmarshalText([]byte(valAddress))
	if err != nil {

		log.Printf("Failed to UnmarshalText into Address : %v", err)
		return
	}

	var height int64 = consensus.HeightLatest
	blk := GetBlockDetails(cfg, height)
	if blk == nil || len(address) == 0 {
		return
	}

	// Create an owner query to be able to retrieve data with regards to account
	query := staking.OwnerQuery{Height: blk.Height, Owner: address}

	delegations, err := co.DelegationsFor(context.Background(), &query)
	if err != nil {

		log.Printf("Failed to get Delegations : %v", err)
		return
	}

	if len(delegations) != 0 {

		value := &delegations[address].Shares
		_ = writeToInfluxDb(c, bp, "oasis_self_delegation_balance", map[string]string{}, map[string]interface{}{"balance": value.String()})
		log.Printf("Self delegations : %s", value)
	}
}
