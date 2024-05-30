package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"tollCalculator.com/types"
)

type Set struct {
	AggregatorEndPoint endpoint.Endpoint
	CalculateEndPoint  endpoint.Endpoint
}
type AggregatorRequest struct {
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
	_, err := s.AggregatorEndPoint(ctx, AggregatorRequest{
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
		req := request.(AggregatorRequest)
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
		inv, err := s.Calculate(ctx, req)
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount:   inv.totalAmount,
			Err:           inv.Err,
		}, nil
	}
}
