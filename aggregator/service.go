package main

import (
	"fmt"

	"tollCalculator.com/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}
type Storer interface {
	Insert(types.Distance) error
}
type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage")
	return i.store.Insert(distance)
}
