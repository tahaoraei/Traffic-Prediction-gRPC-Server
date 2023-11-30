package grpcserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	t "time"
	"traffic-prediction-grpc-server/contract/goproto/time"
	"traffic-prediction-grpc-server/param"
	"traffic-prediction-grpc-server/pkg/logger"
	"traffic-prediction-grpc-server/pkg/metric"
	"traffic-prediction-grpc-server/service/timeservice"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
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

var log = logger.Get()

type Server struct {
	time.UnimplementedGetETAServer
	tehranSvc  *timeservice.Service
	mashhadSvc *timeservice.Service
}

// TODO: change to option pattern
func New(tehranSvc *timeservice.Service, mashhadSvc *timeservice.Service) Server {
	return Server{
		UnimplementedGetETAServer: time.UnimplementedGetETAServer{},
		tehranSvc:                 tehranSvc,
		mashhadSvc:                mashhadSvc,
	}
}

// TODO: change Server to pointer
func (s Server) GetETA(c context.Context, req *time.TravelRequest) (*time.TravelResponse, error) {
	startTime := t.Now()
	defer metric.ResponseHistogram.WithLabelValues("GetETA").Observe(float64(t.Since(startTime).Milliseconds()))

	request := param.ETARequest{
		CurrentETA: req.CurrentETA,
		Distance:   req.Distance,
		Sx:         req.SourceX,
		Sy:         req.SourceY,
		Dx:         req.DestinationX,
		Dy:         req.DestinationY,
		Time:       req.Time,
	}

	// TODO: move this logic to business layer
	var eta *param.ETAResponse
	if req.SourceX > minXTehran && req.SourceY > minYTehran && req.SourceX < maxXTehran && req.SourceY < maxYTehran &&
		req.DestinationX > minXTehran && req.DestinationY > minYTehran && req.DestinationX < maxXTehran && req.DestinationY < maxYTehran {
		eta = s.tehranSvc.GetETA(&request)
	} else if req.SourceX > minXMashhad && req.SourceY > minYMashhad && req.SourceX < maxXMashhad && req.SourceY < maxYMashhad &&
		req.DestinationX > minXMashhad && req.DestinationY > minYMashhad && req.DestinationX < maxXMashhad && req.DestinationY < maxYMashhad {
		eta = s.mashhadSvc.GetETA(&request)
	} else {
		return &time.TravelResponse{ETA: req.CurrentETA}, nil
	}

	if eta == nil {
		log.Warn().Msgf("cant predict eta and ml didn't response")
		return &time.TravelResponse{ETA: req.CurrentETA}, fmt.Errorf("cant predict eta")
	}

	resp := time.TravelResponse{ETA: eta.ETA}
	return &resp, nil
}

func (s Server) Start() {
	address := fmt.Sprintf(":%d", 9090)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Msgf("grpc listener problem: ", err)
		// TODO: return error here not panic
		panic(err)
	}

	timeServer := Server{tehranSvc: s.tehranSvc, mashhadSvc: s.mashhadSvc}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	grpc_prometheus.Register(grpcServer)

	time.RegisterGetETAServer(grpcServer, &timeServer)

	log.Info().Msgf("ETA grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("couldn't server presence grpc server")
	}
}
