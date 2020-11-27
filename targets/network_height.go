package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// GetNetworkLatestHeight to get the latest height of the network
func GetNetworkLatestHeight(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var networkBlock NetworkLatestBlock
	err = json.Unmarshal(resp.Body, &networkBlock)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if &networkBlock != nil {
		networkBlockHeight := networkBlock.Result.Height

		_ = writeToInfluxDb(c, bp, "oasis_network_latest_block", map[string]string{}, map[string]interface{}{"block_height": networkBlockHeight})
		log.Printf("Network height: %d", networkBlockHeight)

		// Calling function to get validator latest
		// block height
		validatorHeight := GetValidatorBlock(cfg, c)
		if validatorHeight == "" {
			log.Println("Error while fetching validator block height from db ", validatorHeight)
			return
		}

		vaidatorBlockHeight, _ := strconv.Atoi(validatorHeight)

		heightDiff := networkBlockHeight - vaidatorBlockHeight

		_ = writeToInfluxDb(c, bp, "oasis_height_difference", map[string]string{}, map[string]interface{}{"difference": heightDiff})
		log.Printf("Network height: %d and Validator Height: %d", networkBlockHeight, vaidatorBlockHeight)

		// Send alert
		if int64(heightDiff) >= cfg.AlertsThreshold.BlockDiffThreshold {
			_ = SendTelegramAlert(fmt.Sprintf("Block difference between validator and network has exceeded %d", cfg.AlertsThreshold.BlockDiffThreshold), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Block difference between validator and network has exceeded %d", cfg.AlertsThreshold.BlockDiffThreshold), cfg)

			log.Println("Sent alert of block height difference")
		}
	}

	return
}
