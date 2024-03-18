package main

import (
	"log"

	"tollCalculator.com/aggregator/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var (
	kafkatopic         = "obuData"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
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

	kafkaConsumer, err := NewKafkaConsumer(kafkatopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
