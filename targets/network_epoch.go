package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// GetNetworkEpoch returns the work epoch number of the network
func GetNetworkEpoch(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var networkEpoch NetworkEpochNumber
	err = json.Unmarshal(resp.Body, &networkEpoch)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if &networkEpoch == nil {
		log.Printf("NetworkEpoch res is empty : %v", networkEpoch.Result)
		return
	}

	epoch := networkEpoch.Result

	err = writeToInfluxDb(c, bp, "oasis_network_epoch_number", map[string]string{}, map[string]interface{}{"epoch_number": epoch})
	if err != nil {
		log.Println("Error while storing network epoch number ", err)
		return
	}
	log.Printf("Network epoch number : %d", epoch)

	// Calling function to get validator epoch number from db
	valEpoch := GetValidatorEpoch(cfg, c)
	log.Printf("Validator epoch number : %s", valEpoch)

	valEpochNumber, _ := strconv.Atoi(valEpoch)
	epochDiff := valEpochNumber - epoch
	log.Printf("Epoch difference : %d", epochDiff)

	err = writeToInfluxDb(c, bp, "oasis_epoch_diff", map[string]string{}, map[string]interface{}{"difference": epochDiff})
	if err != nil {
		log.Printf("Error while storing epoch number difference : %v ", err)
	}

	if int64(epochDiff) >= cfg.AlertsThreshold.EpochDiffThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Epoch difference between validator and network has exceeded %d", cfg.AlertsThreshold.EpochDiffThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Epoch difference between validator and network has exceeded %d", cfg.AlertsThreshold.EpochDiffThreshold), cfg)

		log.Println("Sent alert of worker epoch height difference")
	}

	return
}

// GetValidatorEpoch returns the epoch number of the validator from db
func GetValidatorEpoch(cfg *config.Config, c client.Client) string {
	var epochNumber string
	q := client.NewQuery("SELECT last(epoch_number) FROM oasis_worker_epoch_number", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						epoch := r.Series[0].Values[0][idx]
						epochNumber = fmt.Sprintf("%v", epoch)
						break
					}
				}
			}
		}
	}

	return epochNumber
}
