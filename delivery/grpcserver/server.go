package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"timeMachine/contract/goproto/time"
	"timeMachine/param"
	"timeMachine/service/timeservice"
)

type Server struct {
	time.UnimplementedGetETAServer
	svc *timeservice.Service
}

func New(svc *timeservice.Service) Server {
	return Server{
		UnimplementedGetETAServer: time.UnimplementedGetETAServer{},
		svc:                       svc,
	}
}

func (s Server) GetETA(c context.Context, req *time.TravelRequest) (*time.TravelResponse, error) {
	eta := s.svc.GetETA(param.ETARequest{
		CurrentETA:  req.CurrentETA,
		StraightETA: req.StraightETA,
		Distance:    req.Distance,
		Sx:          req.Sx,
		Sy:          req.Sy,
		Dx:          req.Dx,
		Dy:          req.Dy,
		Time:        req.Time,
	})
	resp := time.TravelResponse{ETA: eta.ETA}
	return &resp, nil
}

func (s Server) Start() {
	address := fmt.Sprintf(":%d", 8086)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	timeServer := Server{svc: s.svc}
	grpcServer := grpc.NewServer()

	time.RegisterGetETAServer(grpcServer, &timeServer)

	log.Println("ETA grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("couldn't server presence grpc server")
	}
}
