package targets

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// GetStatus returns the current status overview and alerts the user based on user configuration
func GetStatus(ops HTTPOptions, cfg *config.Config, c client.Client) {
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
		log.Printf("Failed to establish connection using socket: " +
			socket)
		return
	}

	// Retrieve the current status overview
	status, err := co.GetStatus(context.Background())
	if err != nil {
		log.Printf("Failed to retrieve Status : %v", err)
		return
	}

	numPeers := len(status.NodePeers)

	// Sent alert if no.of peers dropped below given threshold
	if int64(numPeers) < cfg.AlertsThreshold.NumPeersThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.AlertsThreshold.NumPeersThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.AlertsThreshold.NumPeersThreshold), cfg)
	}

	// Write data into influxdb
	_ = writeToInfluxDb(c, bp, "oasis_num_peers", map[string]string{}, map[string]interface{}{"num_peers": numPeers})
	_ = writeToInfluxDb(c, bp, "oasis_peer_addresses", map[string]string{"addresses_count": strconv.Itoa(numPeers)}, map[string]interface{}{"addresses": status.NodePeers})

	bh := status.LatestHeight
	height := int(bh)

	log.Printf("Peers count : %d and validator latest height : %d", numPeers, height)

	// Write latest block height
	_ = writeToInfluxDb(c, bp, "oasis_current_block_height", map[string]string{}, map[string]interface{}{"height": height})

	// Calling function to alert about validator status
	ValidatorStatusAlert(cfg, c, status.IsValidator)
}

// ValidatorStatusAlert is to alert about validator status weather it's active or jailed
func ValidatorStatusAlert(cfg *config.Config, c client.Client, isVoting bool) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	alertTime1 := cfg.DailyAlerts.AlertTime1
	alertTime2 := cfg.DailyAlerts.AlertTime2

	t1, _ := time.Parse(time.Kitchen, alertTime1)
	t2, _ := time.Parse(time.Kitchen, alertTime2)

	now := time.Now().UTC()
	t := now.Format(time.Kitchen)

	a1 := t1.Format(time.Kitchen)
	a2 := t2.Format(time.Kitchen)

	log.Println("validator status...", isVoting)
	log.Println("a1, a2 and present time : ", a1, a2, t)

	validatorStatus := isVoting

	// Calling function to get alert count to avoid duplicate alerts
	ac := GetValAlertCount(cfg, c)

	alertCount, _ := strconv.Atoi(ac)

	log.Println("val alert count..", alertCount)

	if validatorStatus {
		if t == a1 || t == a2 {
			if alertCount == 0 {
				_ = SendTelegramAlert(fmt.Sprintf("Your oasis validator is currently voting"), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Your oasis validator is currently voting"), cfg)
				log.Println("Sent validator status alert")

				_ = writeToInfluxDb(c, bp, "oasis_val_alert_count", map[string]string{}, map[string]interface{}{"count": 1})
			}
		} else {
			_ = writeToInfluxDb(c, bp, "oasis_val_alert_count", map[string]string{}, map[string]interface{}{"count": 0})
		}

		_ = writeToInfluxDb(c, bp, "oasis_validator_status", map[string]string{}, map[string]interface{}{"status": 1})
	} else {
		if t == a1 || t == a2 {
			if alertCount == 0 {
				_ = SendTelegramAlert(fmt.Sprintf("Your oasis validator is in jailed status"), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Your oasis validator is in jailed status"), cfg)
				log.Println("Sent validator status alert")

				_ = writeToInfluxDb(c, bp, "oasis_val_alert_count", map[string]string{}, map[string]interface{}{"count": 1})
			}
		} else {
			_ = writeToInfluxDb(c, bp, "oasis_val_alert_count", map[string]string{}, map[string]interface{}{"count": 0})
		}
		_ = writeToInfluxDb(c, bp, "oasis_validator_status", map[string]string{}, map[string]interface{}{"status": 0})
	}

	return
}

// GetValAlertCount returns count of the val status alert
func GetValAlertCount(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(count) FROM oasis_val_alert_count", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						countValue := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", countValue)
						break
					}
				}
			}
		}
	}
	return count
}
