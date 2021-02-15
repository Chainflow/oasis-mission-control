package targets

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/Chainflow/oasis-mission-control/config"
)

type targetRunner struct{}

// NewRunner returns targetRunner
func NewRunner() *targetRunner {
	return &targetRunner{}
}

// Run to run the request
func (m targetRunner) Run(function func(ops HTTPOptions, cfg *config.Config, c client.Client), ops HTTPOptions, cfg *config.Config, c client.Client) {
	function(ops, cfg, c)
}

// InitTargets which returns the targets
//can write all the endpoints here
func InitTargets(cfg *config.Config) *Targets {
	return &Targets{List: []Target{
		{
			ExecutionType: "Grpc method",
			Name:          "Get Block Details",
			Func:          GetMissedBlock,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get validator meta data",
			Func:          GetValidatorsList,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get node status",
			Func:          GetStatus,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "",
			Name:          "Get block time difference",
			Func:          GetBlockTimeDiff,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get block balance",
			Func:          GetAccount,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get self delegations",
			Func:          GetSelfDelegationBal,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get node version",
			Func:          NodeVersion,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get node status",
			Func:          NodeStatus,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Grpc method",
			Name:          "Get Validator worker epoch number",
			Func:          GetValEpoch,
			ScraperRate:   cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get network latest height",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ValidatorDetails.NetworkURL + "/api/consensus/block?name=" + cfg.ValidatorDetails.NetworkNodeName,
				Method:   http.MethodGet,
			},
			Func:        GetNetworkLatestHeight,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get network epoch number",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ValidatorDetails.NetworkURL + "/api/consensus/epoch?name=" + cfg.ValidatorDetails.NetworkNodeName,
				Method:   http.MethodGet,
			},
			Func:        GetNetworkEpoch,
			ScraperRate: cfg.Scraper.Rate,
		},
	}}
}

func addQueryParameters(req *http.Request, queryParams QueryParams) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()
}

//newHTTPRequest to make a new http request
func newHTTPRequest(ops HTTPOptions) (*http.Request, error) {
	// make new request
	req, err := http.NewRequest(ops.Method, ops.Endpoint, bytes.NewBuffer(ops.Body))
	if err != nil {
		return nil, err
	}

	// Add any query parameters to the URL.
	if len(ops.QueryParams) != 0 {
		addQueryParameters(req, ops.QueryParams)
	}

	return req, nil
}

func makeResponse(res *http.Response) (*PingResp, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &PingResp{}, err
	}

	response := &PingResp{
		StatusCode: res.StatusCode,
		Body:       body,
	}
	_ = res.Body.Close()
	return response, nil
}

// HitHTTPTarget to hit the target and get response
func HitHTTPTarget(ops HTTPOptions) (*PingResp, error) {
	req, err := newHTTPRequest(ops)
	if err != nil {
		return nil, err
	}

	httpcli := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := httpcli.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := makeResponse(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
