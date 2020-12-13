package cmd

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"

	"github.com/chrisdoc/homewizard-p1-prometheus/internal/exporter"
	"github.com/chrisdoc/homewizard-p1-prometheus/internal/homewizard"
	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type config struct {
	Host string `env:"HOST,required"`
	Tick uint64 `env:"TICK" envDefault:"10"`
	Port uint64 `env:"PORT" envDefault:"9898"`
}

// Start the homewizard exporter
func Start() {

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

func homeWizardsTask(cfg config, exporter exporter.Prometheus) {
	client := homewizard.NewP1Client(cfg.Host)
	home, err := client.Retrieve()
	if err != nil {
		return
	}
	exporter.SetData(home)
}

func executeCronJob(cfg config) {
	s := gocron.NewScheduler()
	exporter := exporter.Prometheus{}
	s.Every(cfg.Tick).Second().Do(homeWizardsTask, cfg, exporter)
	<-s.Start()
}
