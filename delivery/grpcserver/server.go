package grpcserver

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"timeMachine/contract/goproto/time"
	"timeMachine/param"
	"timeMachine/pkg/logger"
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
	log.Info().Msgf("taha req %+v", req)
	eta := s.svc.GetETA(param.ETARequest{
		CurrentETA: req.CurrentETA,
		Distance:   req.Distance,
		Sx:         req.Sx,
		Sy:         req.Sy,
		Dx:         req.Dx,
		Dy:         req.Dy,
		Time:       req.Time,
	})
	resp := time.TravelResponse{ETA: eta.ETA}
	return &resp, nil
}

func (s Server) Start() {
	log := logger.Get()
	address := fmt.Sprintf(":%d", 9090)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Msgf("grpc listener problem: ", err)
		panic(err)
	}

	timeServer := Server{svc: s.svc}
	grpcServer := grpc.NewServer()

	time.RegisterGetETAServer(grpcServer, &timeServer)

	log.Info().Msgf("ETA grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("couldn't server presence grpc server")
	}
}
