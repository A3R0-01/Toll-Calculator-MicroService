package aggservice

import (
	"context"

	"tollCalculator.com/types"
)

type Middleware func(Service) Service
type LoggingMiddleware struct {
	next Service
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &LoggingMiddleware{
			next: next,
		}
	}
}

func (svc *LoggingMiddleware) AggregateDistance(ctx context.Context, distance types.Distance) error {
	return nil
}

func (svc *LoggingMiddleware) CalculateDistance(ctx context.Context, distInt int) (invoice *types.Invoice, err error) {
	return &types.Invoice{}, nil
}

type InstrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &InstrumentationMiddleware{
			next: next,
		}
	}
}

func (svc *InstrumentationMiddleware) AggregateDistance(ctx context.Context, distance types.Distance) error {
	return nil
}

func (svc *InstrumentationMiddleware) CalculateDistance(ctx context.Context, distInt int) (invoice *types.Invoice, err error) {
	return &types.Invoice{}, nil
}
