package client

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tollCalculator.com/types"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) *GRPCClient {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}
}

func (g *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := g.client.Aggregate(ctx, req)
	return err
}

func (g *GRPCClient) GetInvoice(ctx context.Context, req *types.GetInvoiceRequest) (*types.Invoice, error) {
	res, err := g.client.GetInvoice(ctx, req)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         int(res.ObuID),
		TotalDistance: float64(res.TotalDistance),
		TotalAmount:   float64(res.TotalAmount),
	}, nil
}
