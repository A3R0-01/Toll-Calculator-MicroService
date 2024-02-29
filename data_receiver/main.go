package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"tollCalculator.com/types"
)

const kafkaTopic = "obudata"

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed %v\n ", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := kafkaTopic

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte("test producing"),
	}, nil)
	p.Flush(15 * 1000)

	return
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.HandleWS)
	http.ListenAndServe(":30000", nil)
	fmt.Println("data rec")
}

type DataReceiver struct {
	msgChan chan types.OBUData
	conn    *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgChan: make(chan types.OBUData, 128),
	}
}
func (dr *DataReceiver) HandleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()

}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU Connected: client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		fmt.Printf("received data from [%d] :: lat %.2f, long %.2f \n", data.OBUID, data.Lat, data.Long)
		// dr.msgChan <- data
	}
}
