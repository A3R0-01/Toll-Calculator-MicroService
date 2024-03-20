package main

import "tollCalculator.com/types"

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}
func (s *GRPCAggregatorServer) AggregateDistance(request *types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(request.ObuID),
		Value: float64(request.Value),
		Unix:  int64(request.Unix),
	}

	return s.svc.AggregateDistance(distance)
}
