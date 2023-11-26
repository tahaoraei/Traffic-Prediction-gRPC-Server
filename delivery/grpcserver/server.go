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

const (
	minXTehran  = 5683551
	minYTehran  = 4216113
	maxXTehran  = 5774664
	maxYTehran  = 4277874
	minXMashhad = 6583521
	minYMashhad = 4309213
	maxXMashhad = 6656388
	maxYMashhad = 4383510
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
	if req.Sx > minXTehran && req.Sy > minYTehran && req.Sx < maxXTehran && req.Sy < maxYTehran &&
		req.Dx > minXTehran && req.Dy > minYTehran && req.Dx < maxXTehran && req.Dy < maxYTehran {
		eta = s.tehranSvc.GetETA(&request)
	} else if req.Sx > minXMashhad && req.Sy > minYMashhad && req.Sx < maxXMashhad && req.Sy < maxYMashhad &&
		req.Dx > minXMashhad && req.Dy > minYMashhad && req.Dx < maxXMashhad && req.Dy < maxYMashhad {
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
