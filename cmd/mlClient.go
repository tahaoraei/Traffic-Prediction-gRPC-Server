package main

import (
	"fmt"
	"os"
	"sync"
	"timeMachine/delivery/httpserver"
	"timeMachine/ml"
	"timeMachine/param"
	"timeMachine/pkg/logger"
	"timeMachine/repository/postgres"
	"timeMachine/scheduler"
	"timeMachine/service/timeservice"
)

var log = logger.Get()

func main() {
	cfg_httpserver := httpserver.Config{
		Port: 7182,
	}

	cfg_db := postgres.Config{
		Host:   "localhost",
		Port:   5435,
		User:   os.Getenv("SECRETS_DBUSER"),
		Pass:   os.Getenv("SECRETS_DBPASS"),
		DBName: "traffic",
	}

	repo := postgres.New(cfg_db)

	tehranML, err := ml.New(repo, "tehran", 1, .4, .6)
	if err != nil {
		log.Fatal().Msgf("faild to load tehran ml model: %v", err)
	}

	mashhadML, err := ml.New(repo, "mashhad", 2, .3, .7)
	if err != nil {
		log.Fatal().Msgf("faild to load mashhad ml model: %v", err)
	}

	var wg sync.WaitGroup
	done := make(chan bool)
	go func() {
		cronTehran := scheduler.New(tehranML, mashhadML)
		wg.Add(1)
		cronTehran.Start(done, &wg)
		fmt.Println("end of scheduler")
	}()

	ml := timeservice.New(tehranML, mashhadML)
	resp := ml.GetETA(&param.ETARequest{
		CurrentETA:    1000,
		Distance:      10000,
		Sx:            0,
		Sy:            0,
		Dx:            0,
		Dy:            0,
		Time:          1000,
		TrafficLength: 1000,
	})
	log.Info().Msgf("%+v", resp)

	server := httpserver.New(cfg_httpserver)
	server.Serve()
	fmt.Println("end of httpserver")

	done <- true
	wg.Wait()

}
