package targets

import (
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

// GetBlockTimeDifference to calculate block time difference of prev block and current block
func GetBlockTimeDiff(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Calling function to get validator latest
	// block height
	currentBlockHeight := GetValidatorBlock(cfg, c)
	if currentBlockHeight == "" {
		log.Println("Error while fetching current block height from db ", currentBlockHeight)
		return
	}

	currentHeight, _ := strconv.ParseInt(currentBlockHeight, 10, 64)

	// callig GetBlockDetails to get validator block details
	currentBlock := GetBlockDetails(cfg, currentHeight)
	if currentBlock == nil {
		return
	}
	currentBlockTime := currentBlock.Time

	prevHeight := currentHeight - 1
	prevBlockResp := GetBlockDetails(cfg, prevHeight)

	prevBlockTime := prevBlockResp.Time
	PrevBlockTime := prevBlockTime
	timeDiff := currentBlockTime.Sub(PrevBlockTime)
	diffSeconds := fmt.Sprintf("%.2f", timeDiff.Seconds())

	_ = writeToInfluxDb(c, bp, "oasis_block_time_diff", map[string]string{}, map[string]interface{}{"time_diff": diffSeconds})
	log.Printf("time diff: %d", diffSeconds)

	return
}

// GetValidatorBlock returns current block height of a validator from db
func GetValidatorBlock(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height) FROM oasis_current_block_height", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						validatorHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}

	return validatorHeight
}
