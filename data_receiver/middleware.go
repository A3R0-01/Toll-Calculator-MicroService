package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"tollCalculator.com/types"
)

type LogMiddleware struct {
	next DataProducer
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID:": data.OBUID,
			"lat:":   data.Lat,
			"long:":  data.Long,
			"took:":  time.Since(start),
		}).Info("producing to kafka")
	}(start)
	return l.next.ProduceData(data)
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}
