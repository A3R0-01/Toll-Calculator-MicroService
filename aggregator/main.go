package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"tollCalculator.com/types"
)

func main() {
	httpListenAddress := flag.String("httpAddress", ":3000", "the listen address of the http server")
	grpcListenAddress := flag.String("grpcAddress", ":3001", "the listen address of the http server")
	flag.Parse()
	store := NewMemoryStore()
	var (
		svc = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpcListenAddress, svc)
	makeHttpTransport(*httpListenAddress, svc)
	fmt.Println("this is working just fine")
}

func makeGRPCTransport(listenAddress string, svc Aggregator) error {
	fmt.Println("GRPC running on port ", listenAddress)
	// make a TCP listeners
	ln, err := net.Listen("TCP", listenAddress)

	if err != nil {
		return err
	}
	defer ln.Close()
	// Make a new GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register (our) GRPC server impimentation to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)

}
func makeHttpTransport(listenAddress string, svc Aggregator) {
	fmt.Println("HTTP transport running on port ", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddress, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu ID"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
