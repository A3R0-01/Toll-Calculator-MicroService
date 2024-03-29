package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"tollCalculator.com/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, requestAgg *types.AggregateRequest) error {

	b, err := json.Marshal(requestAgg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))

	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service responded with non 200 status code: %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := &types.GetInvoiceRequest{
		ObuID: int32(id),
	}

	b, err := json.Marshal(invReq)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/invoice?obu=%d", c.Endpoint, id)
	logrus.Infof("requesting get invoice -> %s", endpoint)
	req, err := http.NewRequest("POST", c.Endpoint+"/invoice", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with %s", resp.Status)
	}

	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return &inv, nil
}
