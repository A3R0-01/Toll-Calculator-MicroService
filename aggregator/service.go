package main

import "tollCalculator.com/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
}
type Storer interface {
	Insert(types.Distance) error
}
type InvoiceAggregator struct {
	store Storer
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	i.store.Insert(distance)
	return nil
}
