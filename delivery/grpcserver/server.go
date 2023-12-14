package grpcserver

import (
	"context"
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net"
	t "time"
	"timeMachine/contract/goproto/time"
	"timeMachine/param"
	"timeMachine/pkg/logger"
	"timeMachine/pkg/metric"
	"timeMachine/service/timeservice"
)

var log = logger.Get()

type Server struct {
	time.UnimplementedGetETAServer
	svc *timeservice.Service
}

// TODO: change to option pattern
func New(svc *timeservice.Service) *Server {
	return &Server{
		UnimplementedGetETAServer: time.UnimplementedGetETAServer{},
		svc:                       svc,
	}
}

// TODO: change Server to pointer
func (s *Server) GetETA(c context.Context, req *time.TravelRequest) (*time.TravelResponse, error) {
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

	var eta *param.ETAResponse

	eta = s.svc.GetETA(&request)
	if eta == nil {
		log.Warn().Msgf("cant predict eta and svc didn't response")
		return &time.TravelResponse{ETA: req.CurrentETA}, fmt.Errorf("cant predict eta")
	}

	resp := time.TravelResponse{ETA: eta.ETA}
	return &resp, nil
}

func (s *Server) Start() {
	address := fmt.Sprintf(":%d", 9090)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Msgf("grpc listener problem: ", err)
		// TODO: return error here not panic
		panic(err)
	}

	timeServer := Server{svc: s.svc}
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
