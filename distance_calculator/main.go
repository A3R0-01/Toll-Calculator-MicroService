package main

import (
	"log"

	"tollCalculator.com/aggregator/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var (
	kafkatopic   = "obuData"
	httpEndpoint = "http://127.0.0.1:3000"
	grpcEndpoint = "127.0.0.1:3001"
)

// Transport (HTTP, GRPC, KAFKA) -> attach business Logic

func main() {
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	httpClient := client.NewHTTPClient(httpEndpoint)
	grpcClient := client.NewGRPCClient(grpcEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkatopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	_ = grpcClient
}
