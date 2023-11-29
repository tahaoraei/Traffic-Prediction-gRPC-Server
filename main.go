package main

import (
	"log"
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
	// TODO: main should be write in 3 line :)
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

	tehranML, err := ml.New("tehran-20231125.txt", .4, .6)
	if err != nil {
		log.Fatalf("faild to load tehran ml model: %v", err)
	}
	tehranSvc := timeservice.New(repo, tehranML, 1)

	mashhadML, err := ml.New("mashhad-20231125.txt", .3, .7)
	if err != nil {
		log.Fatalf("faild to load mashhad ml model: %v", err)
	}
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
