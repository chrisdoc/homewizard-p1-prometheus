package homewizard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client is responsible to interact with Homewizard P1 energy monitor
type Client struct {
	Host   string
	Client *http.Client
}

// Data data
type Data struct {
	SmrVersion 						int64 `json:"smr_version"`
	MeterModel            string  `json:"meter_model"`
	WifiSSID              string  `json:"wifi_ssid"`
	WifiStrength          float64 `json:"wifi_strength"`
	TotalPowerImportT1Kwh float64 `json:"total_power_import_t1_kwh"`
	TotalPowerImportT2Kwh float64 `json:"total_power_import_t2_kwh"`
	TotalPowerExportT1Kwh float64 `json:"total_power_export_t1_kwh"`
	TotalPowerExportT2Kwh float64 `json:"total_power_export_t2_kwh"`
	ActivePowerW          float64 `json:"active_power_w"`
	ActivePowerL1W        float64 `json:"active_power_l1_w"`
	ActivePowerL2W        float64 `json:"active_power_l2_w"`
	ActivePowerL3W        float64 `json:"active_power_l3_w"`
	TotalGasM3            float64 `json:"total_gas_m3"`
}

// NewP1Client createdsa new client with a defautl http client
func NewP1Client(host string) Client {
	defaultClient := http.Client{
		Timeout: time.Second * 5,
	}
	return NewP1ClientWithHTTPClient(host, &defaultClient)
}

// NewP1ClientWithHTTPClient creates a new client with a specified http client
func NewP1ClientWithHTTPClient(host string, httpClient *http.Client) Client {
	return Client{
		Host:   host,
		Client: httpClient,
	}
}

// Retrieve data from P1
func (c *Client) Retrieve() (home *Data, err error) {
	url := fmt.Sprintf("http://%s/api/v1/data", c.Host)

	spaceClient := http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{"url": url}).Error("Coudln't create new http request", url, err)
		return nil, err
	}

	res, err := spaceClient.Do(req)
	if err != nil {
		log.Error("Couldn't execute request", err)
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{"body": body}).Error("Couldn't read body", err)
		return nil, err
	}

	home = &Data{}
	err = json.Unmarshal(body, &home)
	if err != nil {
		log.WithFields(log.Fields{"body": body}).Error("Couldn't parse body", err)
		return nil, err
	}

	return home, nil
}
