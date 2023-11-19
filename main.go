package main

import (
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/service/timeservice"
)

func main() {
	cfg := httpserver.Config{Port: 8080}
	timeSvc := timeservice.New("timemachine-lightgbm-20231118.txt")

	grpc := grpcserver.New(&timeSvc)
	go func() {
		grpc.Start()
	}()

	server := httpserver.New(cfg, timeSvc)
	server.Serve()
}
