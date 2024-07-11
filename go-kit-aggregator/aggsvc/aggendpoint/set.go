package aggendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"tollCalculator.com/go-kit-aggregator/aggsvc/aggservice"
	"tollCalculator.com/types"
)

type Set struct {
	AggregateEndPoint endpoint.Endpoint
	CalculateEndPoint endpoint.Endpoint
}
type AggregateRequest struct {
	OBUID int     `json:"obuID"`
	Value float64 `json:"value"`
	Unix  int64   `json:"unix"`
}
type AggregateResponse struct {
	Err error `json:"err"`
}
type CalculateRequest struct {
	OBUID int `json:"obuID"`
}
type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s *Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndPoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})

	return err
}
func (s *Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndPoint(ctx, CalculateRequest{OBUID: obuID})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, err
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix:  req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

// MakeConcatEndpoint constructs a Concat endpoint wrapping the service.
func MakeCalculateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount:   inv.TotalAmount,
			Err:           err,
		}, nil
	}
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc aggservice.Service, logger log.Logger) Set {
	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = MakeAggregateEndpoint(svc)
		// Sum is limited to 1 request per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		aggregateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(aggregateEndpoint)
		// aggregateEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(aggregateEndpoint)
		// aggregateEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(aggregateEndpoint)
	}
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = MakeCalculateEndpoint(svc)
		// Concat is limited to 1 request per second with burst of 100 requests.
		// Note, rate is defined as a number of requests per second.
		calculateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(calculateEndpoint)
		// calculateEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(calculateEndpoint)
		// calculateEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(calculateEndpoint)
	}
	return Set{
		AggregateEndPoint: aggregateEndpoint,
		CalculateEndPoint: calculateEndpoint,
	}
}

// func New(s aggservice.Service) Set {
// 	return Set{
// 		AggregateEndPoint: MakeAggregateEndpoint(s),
// 		CalculateEndPoint: MakeCalculateEndpoint(s),
// 	}
// }
