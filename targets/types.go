package targets

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

type (
	// QueryParams map of strings
	QueryParams map[string]string

	// HTTPOptions is a structure that holds all http options parameters
	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	// Target is a structure which holds all the parameters of a target
	//this could be used to write endpoints for each functionality
	Target struct {
		ExecutionType string
		HTTPOptions   HTTPOptions
		Name          string
		Func          func(m HTTPOptions, cfg *config.Config, c client.Client)
		ScraperRate   string
	}

	// Targets list of all the targets
	Targets struct {
		List []Target
	}

	// PingResp is a structure which holds the options of a response
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// Peer is a structure which holds the info about a peer address
	Peer struct {
		RemoteIP         string      `json:"remote_ip"`
		ConnectionStatus interface{} `json:"connection_status"`
		IsOutbound       bool        `json:"is_outbound"`
		NodeInfo         struct {
			Moniker string `json:"moniker"`
			Network string `json:"network"`
		} `json:"node_info"`
	}

	// NetworkLatestBlock which holds the information about network latest block
	NetworkLatestBlock struct {
		Result struct {
			ConsensusVersion string      `json:"consensus_version"`
			Backend          string      `json:"backend"`
			Features         int         `json:"features"`
			NodePeers        []string    `json:"node_peers"`
			LatestHeight     int         `json:"latest_height"`
			LatestHash       string      `json:"latest_hash"`
			LatestTime       time.Time   `json:"latest_time"`
			LatestStateRoot  interface{} `json:"latest_state_root"`
			GenesisHeight    int         `json:"genesis_height"`
			GenesisHash      interface{} `json:"genesis_hash"`
			IsValidator      bool        `json:"is_validator"`
		} `json:"result"`
	}

	// NetworkEpochNumber which holds the epoch number of the network
	NetworkEpochNumber struct {
		Result int `json:"result"`
	}
)
