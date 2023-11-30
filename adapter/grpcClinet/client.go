package grpcClinet

import (
	"context"
	"google.golang.org/grpc"
	"traffic-prediction-grpc-server/contract/goproto/time"
	"traffic-prediction-grpc-server/param"
)

type Client struct {
	address string
}

func New(address string) Client {
	return Client{
		address: address,
	}
}

func (c Client) GetNewETA(ctx context.Context, request param.ETARequest) (param.ETAResponse, error) {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return param.ETAResponse{}, err
	}
	defer conn.Close()

	client := time.NewGetETAClient(conn)

	travelRequest := time.TravelRequest{
		CurrentETA:   request.CurrentETA,
		Distance:     request.Distance,
		SourceX:      request.Sx,
		SourceY:      request.Sy,
		DestinationX: request.Dx,
		DestinationY: request.Dy,
		Time:         request.Time,
	}

	resp, err := client.GetETA(ctx, &travelRequest)
	if err != nil {
		return param.ETAResponse{}, err
	}

	return param.ETAResponse{ETA: resp.ETA}, nil
}
