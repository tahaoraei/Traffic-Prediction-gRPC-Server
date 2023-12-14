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

	tehranML, err := ml.New(repo, "tehran", 1, .4, .6)
	if err != nil {
		log.Fatalf("faild to load tehran ml model: %v", err)
	}

	mashhadML, err := ml.New(repo, "mashhad", 2, .3, .7)
	if err != nil {
		log.Fatalf("faild to load mashhad ml model: %v", err)
	}

	svc := timeservice.New(tehranML, mashhadML)

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cronTehran := scheduler.New(tehranML, mashhadML)
		wg.Add(1)
		cronTehran.Start(done, &wg)
	}()

	go func() {
		server := httpserver.New(cfg_httpserver)
		server.Serve()
	}()

	grpc := grpcserver.New(&svc)
	grpc.Start()

	done <- true
	wg.Wait()
}
