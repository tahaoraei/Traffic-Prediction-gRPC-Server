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
	tehranML := ml.New("timemachine-lightgbm-l1-tehran-evening-20231122.txt", .4, .6)
	tehranSvc := timeservice.New(repo, tehranML, 1)
	mashhadML := ml.New("timemachine-lightgbm-l1-mashhad-20231122.txt", .3, .7)
	mashhadSvc := timeservice.New(repo, mashhadML, 2)

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cronTehran := scheduler.New(&tehranSvc, &mashhadSvc)
		wg.Add(1)
		cronTehran.Start(done, &wg)
	}()

	go func() {
		server := httpserver.New(cfg_httpserver)
		server.Serve()
	}()

	grpc := grpcserver.New(&tehranSvc, &mashhadSvc)
	grpc.Start()

	done <- true
	wg.Wait()
}
