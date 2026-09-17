package trading

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tuanta7/nofomo/pkg/lightspeed"
	"resty.dev/v3"
)

const (
	baseURL        = "https://openapi.dnse.com.vn"
	sandboxBaseURL = "https://sb-openapi.dnse.com.vn"
)

type Client struct {
	apiKey    string
	apiSecret string
	client    *resty.Client
}

func NewClient(apiKey, apiSecret string) *Client {
	client := resty.New()
	client.SetTransport(&http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	})

	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    client,
	}
}

func (c *Client) Get(ctx context.Context, path string, result any) error {
	date, signature, err := c.Sign("GET", path)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-API-Key", c.apiKey).
		SetHeader("X-Aux-Date", date).
		SetHeader("X-Signature", signature).
		SetHeader("version", lightspeed.APIVersion).
		SetResult(result).
		Get(baseURL + path)
	if err != nil {
		return fmt.Errorf("dnse: GET %s: %w", path, err)
	} else if resp.IsStatusFailure() {
		return fmt.Errorf("dnse: GET %s: %s: %s", path, resp.Status(), resp.String())
	}
	return nil
}
