package main

import (
	"tollCalculator.com/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}
type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
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
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	totalDistance, err := i.store.Get(obuID)
	if err != nil {
		return &types.Invoice{}, err
	}
	invoice := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: totalDistance,
		TotalAmount:   totalDistance * basePrice,
	}
	return invoice, nil
}
