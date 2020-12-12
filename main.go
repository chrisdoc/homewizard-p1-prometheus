package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"

	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type config struct {
	Host string `env:"HOST,required"`
	Tick uint64 `env:"TICK" envDefault:"10"`
	Port uint64 `env:"PORT" envDefault:"9898"`
}

type homeWizard struct {
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

var (
	wifiStrength = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wifi_strength",
		Help: "Wifi strength in Db",
	})
	totalPowerImportT1Kwh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_power_import_t1_kwh",
		Help: "The total power import on T1 in kWh",
	})
	totalPowerImportT2Kwh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_power_import_t2_kwh",
		Help: "The total power import on T2 in kWh",
	})
	totalPowerExportT1Kwh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_power_export_t1_kwh",
		Help: "The total power export on T1 in kWh",
	})
	totalPowerExportT2Kwh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_power_export_t2_kwh",
		Help: "The total power export on T2 in kWh",
	})
	activePowerW = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_power_w",
		Help: "The active power in W",
	})
	activePowerL1W = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_power_l1_w",
		Help: "he active power on L1 in W",
	})
	activePowerL2W = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_power_l2_w",
		Help: "he active power on L2 in W",
	})
	activePowerL3W = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_power_l3_w",
		Help: "The active power on L3 in W",
	})
	totaGasM3 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_gas_m3",
		Help: "The total gas consumption in m3",
	})
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	finish := make(chan bool)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		listenAddress := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
		http.ListenAndServe(listenAddress, nil)
	}()

	go func() {
		executeCronJob(cfg)
	}()

	<-finish
}

func retrieveData(host string) (home *homeWizard, err error) {
	url := fmt.Sprintf("http://%s/api/v1/data", host)

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

	home = &homeWizard{}
	err = json.Unmarshal(body, &home)
	if err != nil {
		log.WithFields(log.Fields{"body": body}).Error("Couldn't parse body", err)
		return nil, err
	}

	return home, nil
}

func setPromotheus(home *homeWizard) {
	wifiStrength.Set(home.WifiStrength)

	totalPowerImportT1Kwh.Set(home.TotalPowerImportT1Kwh)
	totalPowerImportT2Kwh.Set(home.TotalPowerImportT2Kwh)
	totalPowerImportT2Kwh.Set(home.TotalPowerExportT2Kwh)
	totalPowerImportT2Kwh.Set(home.TotalPowerExportT2Kwh)

	activePowerW.Set(home.ActivePowerW)

	activePowerL1W.Set(home.ActivePowerL1W)
	activePowerL2W.Set(home.ActivePowerL2W)
	activePowerL3W.Set(home.ActivePowerL3W)

	totaGasM3.Set(home.TotalGasM3)
}

func homeWizardsTask(cfg config) {
	home, err := retrieveData(cfg.Host)
	if err != nil {
		return
	}
	setPromotheus(home)
}

func executeCronJob(cfg config) {
	s := gocron.NewScheduler()
	s.Every(cfg.Tick).Second().Do(homeWizardsTask, cfg)
	<-s.Start()
}
