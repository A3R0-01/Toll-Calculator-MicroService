package aggendpoint

import "github.com/go-kit/kit/endpoint"

type Set struct {
	AggregatorEndPoint endpoint.Endpoint
	CalculateEndPoint  endpoint.Endpoint
}
