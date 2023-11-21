package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"sync"
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	prometh "timeMachine/pkg/prometheus"
	"timeMachine/repository/postgres"
	"timeMachine/scheduler"
	"timeMachine/service/timeservice"
)

func main() {
	prometheus.MustRegister(prometh.ResponseHistogram)
	cfg_httpserver := httpserver.Config{
		Port: 7182,
	}

	cfg_db := postgres.Config{
		Host:   "172.20.11.137",
		Port:   5432,
		User:   os.Getenv("SECRETS_DBUSER"),
		Pass:   os.Getenv("SECRETS_DBPASS"),
		DBName: "traffic",
	}

	repo := postgres.New(cfg_db)
	timeSvc := timeservice.New("timemachine-lightgbm-20231120.txt", repo)

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cron := scheduler.New(&timeSvc)
		wg.Add(1)
		cron.Start(done, &wg)
	}()

	grpc := grpcserver.New(&timeSvc)
	go func() {
		grpc.Start()
	}()

	server := httpserver.New(cfg_httpserver)
	server.Serve()

	done <- true
	wg.Wait()
}
