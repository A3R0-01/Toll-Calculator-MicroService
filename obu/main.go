package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"tollCalculator.com/types"
)

const (
	sendInterval = time.Second
	wsEndPoint   = "ws://127.0.0.1:30000/ws"
)

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}
func main() {
	obuIds := generateOBUID(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIds); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIds[i],
				Lat:   lat,
				Long:  long,
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", data)

		}
		fmt.Println("\n\n")
		time.Sleep(sendInterval)
	}

}

func generateOBUID(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
