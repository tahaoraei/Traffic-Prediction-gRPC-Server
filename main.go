package main

import (
	"os"
	"sync"
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/ml"
	"timeMachine/repository/postgres"
	"timeMachine/scheduler"
	"timeMachine/service/timeservice"
)

func main() {
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
	mlModel := ml.New("timemachine-lightgbm-tehran-evening-20231122.txt")
	timeSvc := timeservice.New(repo, mlModel)

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cron := scheduler.New(&timeSvc)
		wg.Add(1)
		cron.Start(done, &wg)
	}()

	go func() {
		server := httpserver.New(cfg_httpserver)
		server.Serve()
	}()

	grpc := grpcserver.New(&timeSvc)
	grpc.Start()

	done <- true
	wg.Wait()
}
