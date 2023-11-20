package main

import (
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/service/timeservice"
)

func main() {
	cfg := httpserver.Config{Port: 7182}
	timeSvc := timeservice.New("timemachine-lightgbm-20231120.txt")

	grpc := grpcserver.New(&timeSvc)
	go func() {
		grpc.Start()
	}()

	server := httpserver.New(cfg, timeSvc)
	server.Serve()
}
