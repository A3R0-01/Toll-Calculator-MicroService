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

func (svc *LoggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return svc.next.Aggregate(ctx, distance)
}

func (svc *LoggingMiddleware) Calculate(ctx context.Context, distInt int) (invoice *types.Invoice, err error) {
	return svc.next.Calculate(ctx, distInt)
}

type InstrumentingMiddleware struct {
	next Service
}

func newInstrumentingMiddleware() Middleware {
	return func(next Service) Service {
		return &InstrumentingMiddleware{
			next: next,
		}
	}
}

func (svc *InstrumentingMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {
	return svc.next.Aggregate(ctx, distance)
}

func (svc *InstrumentingMiddleware) Calculate(ctx context.Context, distInt int) (invoice *types.Invoice, err error) {
	return svc.next.Calculate(ctx, distInt)
}
