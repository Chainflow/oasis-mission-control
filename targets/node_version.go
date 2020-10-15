package targets

import (
	"log"
	"os/exec"
	"regexp"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// NodeVersion returns the version of the oasis node
func NodeVersion(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cmd := exec.Command("oasis-node", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		_ = writeToInfluxDb(c, bp, "oasis_node_version", map[string]string{}, map[string]interface{}{"v": "NA"})
		return
	}

	resp := string(out)

	r := regexp.MustCompile(`Software version: ([0-9]{2}.[0-9]{1})`)
	matches := r.FindAllStringSubmatch(resp, -1)
	if len(matches) == 0 {
		_ = writeToInfluxDb(c, bp, "oasis_node_version", map[string]string{}, map[string]interface{}{"v": "NA"})
		return
	}
	if len(matches[0]) != 2 {
		_ = writeToInfluxDb(c, bp, "oasis_node_version", map[string]string{}, map[string]interface{}{"v": "NA"})
		return
	}
	_ = writeToInfluxDb(c, bp, "oasis_node_version", map[string]string{}, map[string]interface{}{"v": matches[0][1]})
	log.Printf("Version: %s", matches[0][1])
}
