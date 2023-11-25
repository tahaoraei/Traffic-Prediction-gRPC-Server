package grpcserver

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	t "time"
	"timeMachine/contract/goproto/time"
	"timeMachine/param"
	"timeMachine/pkg/logger"
	"timeMachine/pkg/metric"
	"timeMachine/service/timeservice"
)

type Server struct {
	time.UnimplementedGetETAServer
	tehranSvc  *timeservice.Service
	mashhadSvc *timeservice.Service
}

func New(tehranSvc *timeservice.Service, mashhadSvc *timeservice.Service) Server {
	return Server{
		UnimplementedGetETAServer: time.UnimplementedGetETAServer{},
		tehranSvc:                 tehranSvc,
		mashhadSvc:                mashhadSvc,
	}
}

func (s Server) GetETA(c context.Context, req *time.TravelRequest) (*time.TravelResponse, error) {
	startTime := t.Now()

	request := param.ETARequest{
		CurrentETA: req.CurrentETA,
		Distance:   req.Distance,
		Sx:         req.Sx,
		Sy:         req.Sy,
		Dx:         req.Dx,
		Dy:         req.Dy,
		Time:       req.Time,
	}

	var eta *param.ETAResponse
	if req.Sx > 5683551 && req.Sy > 4216113 && req.Dx < 5774664 && req.Dy < 4277874 {
		eta = s.tehranSvc.GetETA(&request)
	} else if req.Sx > 6583521 && req.Sy > 4309213 && req.Dx < 6656388 && req.Dy < 4383510 {
		eta = s.mashhadSvc.GetETA(&request)
	} else {
		return &time.TravelResponse{ETA: req.CurrentETA}, fmt.Errorf("location is not in tehran or mashhad")
	}

	responseDuration := t.Since(startTime).Milliseconds()
	metric.ResponseHistogram.Observe(float64(responseDuration))
	if eta == nil {
		log.Warn().Msgf("cant predict eta and ml didn't response")
		return &time.TravelResponse{ETA: req.CurrentETA}, fmt.Errorf("cant predict eta")
	}

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

	timeServer := Server{tehranSvc: s.tehranSvc, mashhadSvc: s.mashhadSvc}
	grpcServer := grpc.NewServer()

	time.RegisterGetETAServer(grpcServer, &timeServer)

	log.Info().Msgf("ETA grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("couldn't server presence grpc server")
	}
}
