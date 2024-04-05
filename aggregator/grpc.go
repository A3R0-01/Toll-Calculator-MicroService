package main

import (
	"context"

	"tollCalculator.com/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}
func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, request *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(request.ObuID),
		Value: float64(request.Value),
		Unix:  int64(request.Unix),
	}

	err := s.svc.AggregateDistance(distance)
	if err != nil {
		return nil, err
	}
	return &types.None{}, nil
}

func (s *GRPCAggregatorServer) GetInvoice(ctx context.Context, request *types.GetInvoiceRequest) (*types.GetInvoiceResponse, error) {
	obuID := request.ObuID

	invoice, err := s.svc.CalculateInvoice(int(obuID))
	if err != nil {
		return nil, err
	}

	return &types.GetInvoiceResponse{
		ObuID:         request.ObuID,
		TotalDistance: float32(invoice.TotalDistance),
		TotalAmount:   float32(invoice.TotalAmount),
	}, nil
}
