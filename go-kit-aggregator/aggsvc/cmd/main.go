package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
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
	var (
		httpListener, err = net.Listen("tcp", httpAddr)
		basicService      = aggservice.New()
		endpoints         = aggendpoint.New(basicService, logger)
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
