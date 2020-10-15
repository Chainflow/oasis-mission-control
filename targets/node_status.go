package targets

import (
	"fmt"
	"log"
	"os"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// NodeStatus returns weather the node is up or not
func NodeStatus(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	path := strings.Split(cfg.SocketPath, ":")

	if len(path) >= 2 {
		if _, err := os.Stat(path[1]); err == nil {
			_ = writeToInfluxDb(c, bp, "oasis_node_status", map[string]string{}, map[string]interface{}{"status": 1})
			log.Println("File exists!")
		} else {
			_ = SendTelegramAlert(fmt.Sprintf("Oasis node on your validator instance is not running: \n%v", err), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Oasis node on your validator instance is not running: \n%v", err), cfg)
			_ = writeToInfluxDb(c, bp, "oasis_node_status", map[string]string{}, map[string]interface{}{"status": 0})
			log.Println("File does not exist!")
		}
	} else {
		log.Println("Invalid socket path!")
		return
	}
}
