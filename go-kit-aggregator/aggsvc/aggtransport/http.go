package aggtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggendpoint"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggservice"
)

// func NewHTTPHandler(){
// 	options := []httptransport.ServerOption{
// 		httptransport.ServerErrorEncoder(errorEncoder),
// 		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
// 	}

// 	if zipkinTracer != nil {
// 		// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
// 		// provided operation name or a global tracing service can be instantiated
// 		// without an operation name and fed to each Go kit endpoint as ServerOption.
// 		// In the latter case, the operation name will be the endpoint's http method.
// 		// We demonstrate a global tracing service here.
// 		options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
// 	}

// 	m := http.NewServeMux()
// 	m.Handle("/sum", httptransport.NewServer(
// 		endpoints.aggregateEndpoint,
// 		decodeHTTPSumRequest,
// 		encodeHTTPGenericResponse,
// 		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
// 	))
// 	m.Handle("/concat", httptransport.NewServer(
// 		endpoints.calculateEndpoint,
// 		decodeHTTPConcatRequest,
// 		encodeHTTPGenericResponse,
// 		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))...,
// 	))
// 	return m
// }

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	fmt.Println("This is coming from the error Encoder -> ", err)
}

func NewHTTPHandler(endpoints aggendpoint.Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/aggregate", httptransport.NewServer(
		endpoints.AggregateEndPoint,
		decodeHTTPAggregateRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/invoice", httptransport.NewServer(
		endpoints.CalculateEndPoint,
		decodeHTTPCalculateRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

func NewHTTPClient(instance string, logger log.Logger) (aggservice.Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []httptransport.ClientOption

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/aggregate"),
			encodeHTTPGenericRequest,
			decodeHTTPAggregateResponse,
			options...,
		).Endpoint()
		aggregateEndpoint = limiter(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "aggregate",
			Timeout: 30 * time.Second,
		}))(aggregateEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/calculate"),
			encodeHTTPGenericRequest,
			decodeHTTPCalculateResponse,
			options...,
		).Endpoint()
		calculateEndpoint = limiter(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Concat",
			Timeout: 10 * time.Second,
		}))(calculateEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return &aggendpoint.Set{
		AggregateEndPoint: aggregateEndpoint,
		CalculateEndPoint: calculateEndpoint,
	}, nil
}

func decodeHTTPAggregateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.AggregateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
func decodeHTTPCalculateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
func decodeHTTPAggregateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.AggregateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
func decodeHTTPCalculateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.CalculateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
