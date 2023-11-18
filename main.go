package main

import (
	"timeMachine/delivery/httpserver"
	"timeMachine/service/timeservice"
)

func main() {
	cfg := httpserver.Config{Port: 8080}
	timeSvc := timeservice.New("timemachine-lightgbm-20231118.txt")
	server := httpserver.New(cfg, timeSvc)
	server.Serve()
	//req := param.ETARequest{
	//	CurrentETA:  1667.0,
	//	StraightETA: 1601.0,
	//	Distance:    16088.0,
	//	Sx:          5712845.0,
	//	Sy:          4262054.0,
	//	Dx:          5724172.0,
	//	Dy:          4263342.0,
	//	Time:        879.0,
	//}
	//
	//timeSvc := timeservice.New("timemachine-lightgbm-20231118.txt")
	//p := timeSvc.GetETA(req)
	//fmt.Printf("Prediction for %v: %d\n", req, p)
}