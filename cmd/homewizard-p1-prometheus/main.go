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
		log.Fatalf("Error parsing environment variables %+v\n", err)

	}

	finish := make(chan bool)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		listenAddress := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
		err := http.ListenAndServe(listenAddress, nil)
		if err != nil {
			log.Fatalf("Error starting metrics server %+v\n", err)
		}
	}()

	go func() {
		executeCronJob(cfg)
	}()

	<-finish
}

func homeWizardsTask(cfg config, exporter exporter.Prometheus) error {
	client := homewizard.NewP1Client(cfg.Host)
	home, err := client.Retrieve()
	if err != nil {
		return err
	}
	exporter.SetData(home)
	return nil
}

func executeCronJob(cfg config) {
	s := gocron.NewScheduler()
	prometheus := exporter.Prometheus{}
	err := s.Every(cfg.Tick).Second().Do(homeWizardsTask, cfg, prometheus)
	if err != nil {
		log.Errorf("Error executing cron job %+v\n", err)
	}
	<-s.Start()
}
