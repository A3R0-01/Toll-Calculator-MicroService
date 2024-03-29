package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"tollCalculator.com/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":         time.Since(start),
			"err":          err,
			"obuID":        obuID,
			"totalDist":    invoice.TotalDistance,
			"totalAmouunt": invoice.TotalAmount,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obuID)
	return
}
