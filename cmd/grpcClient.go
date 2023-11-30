package main

import (
	"context"
	"fmt"
	"sync"
	"traffic-prediction-grpc-server/adapter/grpcClinet"
	"traffic-prediction-grpc-server/param"
)

func main() {
	client := grpcClinet.New(":9090")

	req := param.ETARequest{
		CurrentETA: 1667.0,
		Distance:   16088.0,
		Sx:         5712845.0,
		Sy:         4262054.0,
		Dx:         5724172.0,
		Dy:         4263342.0,
		Time:       879.0,
	}

	var wg sync.WaitGroup
	counter := 0
	for {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetNewETA(context.Background(), req)
			if err != nil {
				panic(err)
			}

			//fmt.Println(resp.ETA)
		}()
		wg.Wait()
		counter++
		fmt.Println(counter)
	}

}
