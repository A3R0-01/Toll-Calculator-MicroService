package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"tollCalculator.com/types"
)

func main() {
	listenAddress := flag.String("listenAddress", ":3000", "the listen address of the http server")
	flag.Parse()
	store := NewMemoryStore()
	var (
		svc = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHttpTransport(*listenAddress, svc)
	fmt.Println("this is working just fine")
}
func makeHttpTransport(listenAddress string, svc Aggregator) {
	fmt.Println("HTTP transport running on port ", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddress, nil)
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
