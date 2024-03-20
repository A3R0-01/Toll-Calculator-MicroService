package client

import (
	"context"

	"tollCalculator.com/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
