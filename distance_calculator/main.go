package main

import (
	"log"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

var (
	kafkatopic = "obuData"
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
	kafkaConsumer, err := NewKafkaConsumer(kafkatopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
