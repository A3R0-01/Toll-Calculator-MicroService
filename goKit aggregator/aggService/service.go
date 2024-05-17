package main

import (
	"context"

	"tollCalculator.com/types"
)

const basePrice = 3.5

type Service interface {
	AggregateDistance(context.Context, types.Distance) error
	CalculateInvoice(context.Context, int) (*types.Invoice, error)
}
type Storer interface {
	Get(int) (float64, error)
	Insert(types.Distance) error
}

type BasicService struct {
	store Storer
}

func NewBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) AggregateDistance(ctx context.Context, distance types.Distance) error {
	return svc.store.Insert(distance)
}

func (svc *BasicService) CalculateInvoice(ctx context.Context, obuID int) (invoice *types.Invoice, err error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	invoice, err = &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   dist * basePrice,
	}, nil
	return
}

// NewAggregatorService will construct a complete microservice with logging and instrumentation middleware
func NewAggregatorService() Service {
	var svc Service
	{

		svc = NewBasicService(NewMemoryStore())
		svc = newLoggingMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}
	return svc
}
