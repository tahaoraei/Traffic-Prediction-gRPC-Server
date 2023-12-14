package main

import (
	"log"
	"sync"
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/ml"
	"timeMachine/pkg/config"
	"timeMachine/repository/postgres"
	"timeMachine/scheduler"
	"timeMachine/service/timeservice"
)

func main() {
	// TODO: main should be write in 3 line :)
	cfg := config.Load("config.yml")

	repo := postgres.New(cfg.Postgres)

	tehranML, err := ml.New(repo, "tehran", 1, .4)
	if err != nil {
		log.Fatalf("faild to load tehran ml model: %v", err)
	}

	mashhadML, err := ml.New(repo, "mashhad", 2, .3)
	if err != nil {
		log.Fatalf("faild to load mashhad ml model: %v", err)
	}

	svc := timeservice.New(tehranML, mashhadML)

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cron := scheduler.New(cfg.Scheduler, tehranML, mashhadML)
		wg.Add(1)
		cron.Start(done, &wg)
	}()

	go func() {
		server := httpserver.New(cfg.HTTPServer)
		server.Serve()
	}()

	grpc := grpcserver.New(cfg.GRPCServer, &svc)
	grpc.Start()

	done <- true
	wg.Wait()
}
