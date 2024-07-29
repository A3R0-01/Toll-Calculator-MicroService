package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggendpoint"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggservice"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggtransport"
)

func main() {
	httpAddr := ":8002"
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "TollCalculator",
			Subsystem: "AggregatorService",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	var (
		httpListener, err = net.Listen("tcp", httpAddr)
		basicService      = aggservice.New(logger)
		endpoints         = aggendpoint.New(basicService, logger, duration)
		httpHandler       = aggtransport.NewHTTPHandler(endpoints, logger)
	)

	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", httpAddr)
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}

}
