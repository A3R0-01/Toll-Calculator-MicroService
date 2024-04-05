package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"tollCalculator.com/aggregator/client"
	"tollCalculator.com/types"
)

const grpcEndpoint = "127.0.0.1:3001"

func main() {
	httpListenAddress := flag.String("httpAddress", ":4000", "the listen address of the http server")
	grpcListenAddress := flag.String("grpcAddress", ":3001", "the listen address of the http server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddress, svc))
	}()
	time.Sleep(time.Second * 2)
	c := client.NewGRPCClient(grpcEndpoint)
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 30.84,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	res, err := c.GetInvoice(context.Background(), &types.GetInvoiceRequest{ObuID: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response successful", res)
	makeHttpTransport(*httpListenAddress, svc)
	fmt.Println("this is working just fine")
}

func makeGRPCTransport(listenAddress string, svc Aggregator) error {
	fmt.Println("GRPC running on port ", listenAddress)
	// make a TCP listeners
	ln, err := net.Listen("tcp", listenAddress)

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
	aggMh := newHTTPMetricHandler("aggregate")
	calculateMh := newHTTPMetricHandler("invoice")
	http.HandleFunc("/aggregate", aggMh.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", calculateMh.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port ", listenAddress)

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
