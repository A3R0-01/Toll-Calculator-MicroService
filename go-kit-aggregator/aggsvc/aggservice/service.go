package aggservice

import (
	"context"

	"github.com/go-kit/log"
	"tollCalculator.com/types"
)

const basePrice = 3.5

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
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

func (svc *BasicService) Aggregate(ctx context.Context, distance types.Distance) error {
	return svc.store.Insert(distance)
}

func (svc *BasicService) Calculate(ctx context.Context, obuID int) (invoice *types.Invoice, err error) {
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
func New(logger log.Logger) Service {
	var svc Service
	{

		svc = NewBasicService(NewMemoryStore())
		svc = newLoggingMiddleware(logger)(svc)
		svc = newInstrumentingMiddleware()(svc)
	}
	return svc
}
