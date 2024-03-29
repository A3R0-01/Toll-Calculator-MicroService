package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"tollCalculator.com/types"
)

func main() {

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.HandleWS)
	http.ListenAndServe(":30000", nil)
	fmt.Println("data rec")
}

type DataReceiver struct {
	msgChan chan types.OBUData
	conn    *websocket.Conn
	prod    DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obuData"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgChan: make(chan types.OBUData, 128),
		prod:    p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
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
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce err: ", err)
		}
		// fmt.Printf("received data from [%d] :: lat %.2f, long %.2f \n", data.OBUID, data.Lat, data.Long)
		// dr.msgChan <- data
	}
}
