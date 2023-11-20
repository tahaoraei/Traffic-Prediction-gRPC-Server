package main

import (
	"os"
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/repository/postgres"
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
	timeSvc := timeservice.New("timemachine-lightgbm-20231120.txt", repo)

	grpc := grpcserver.New(&timeSvc)
	go func() {
		grpc.Start()
	}()

	server := httpserver.New(cfg_httpserver, timeSvc)
	server.Serve()
}
