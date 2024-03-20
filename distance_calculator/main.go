package main

import (
	"log"

	"tollCalculator.com/aggregator/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var (
	kafkatopic      = "obuData"
	httpAggEndpoint = "http://127.0.0.1:3000/aggregate"
	grpcAggEndpoint = "127.0.0.1:3001"
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
	// httpClient := client.NewHTTPClient(httpAggEndpoint)
	grpcClient := client.NewGRPCClient(grpcAggEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkatopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
