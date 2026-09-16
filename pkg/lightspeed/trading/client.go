package trading

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	// DNSE date-versions its OpenAPI; required on every request via the "version" header.
	apiVersion = "2026-05-07"

	baseURL        = "https://openapi.dnse.com.vn"
	sandboxBaseURL = "https://sb-openapi.dnse.com.vn"
)

type Client struct {
	apiKey    string
	apiSecret string
	client    *resty.Client
	logger    *zap.Logger
}

// Request sends an HMAC-signed GET request to a path and decodes the JSON response into a result.
func (c *Client) Request(ctx context.Context, method, path string, result any) error {
	date, signature, err := c.Sign("GET", path)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetMethod(method).
		SetContext(ctx).
		SetHeader("x-api-key", c.apiKey).
		SetHeader("Date", date).
		SetHeader("X-Signature", signature).
		SetHeader("version", apiVersion).
		SetResult(result).
		Get(baseURL + path)
	if err != nil {
		return fmt.Errorf("dnse: GET %s: %w", path, err)
	}
	if resp.IsStatusFailure() {
		return fmt.Errorf("dnse: GET %s: %s: %s", path, resp.Status(), resp.String())
	}
	return nil
}
