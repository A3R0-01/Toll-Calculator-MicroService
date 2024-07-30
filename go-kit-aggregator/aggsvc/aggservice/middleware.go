package aggservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"tollCalculator.com/types"
)

type Middleware func(Service) Service
type LoggingMiddleware struct {
	log  log.Logger
	next Service
}

// var logger log.Logger
// logger

func newLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &LoggingMiddleware{
			log:  logger,
			next: next,
		}
	}
}

func (svc *LoggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) (err error) {
	defer func(start time.Time) {
		svc.log.Log("method", "Aggregate", "took", time.Since(start), "obuId", distance.OBUID, "dist", distance.Value, "unix", distance.Unix, "err", err)
	}(time.Now())
	err = svc.next.Aggregate(ctx, distance)
	return
}

func (svc *LoggingMiddleware) Calculate(ctx context.Context, distInt int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		svc.log.Log("method", "Calculate", "took", time.Since(start), "obuId", distInt, "TotalDistance", invoice.TotalDistance, "InvoiceAmount", invoice.TotalAmount, "err", err)
	}(time.Now())
	invoice, err = svc.next.Calculate(ctx, distInt)
	return
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
