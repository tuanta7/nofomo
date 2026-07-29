package lightspeed

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	"resty.dev/v3"
)

const (
	// DNSE date-versions its OpenAPI; required on every request via the "version" header.
	apiVersion = "2026-05-07"

	baseURL                    = "https://openapi.dnse.com.vn"
	pathAccounts               = "/accounts"
	pathAccountBalances        = "/accounts/%s/balances"
	pathAccountLoanPackages    = "/accounts/%s/loan-packages"
	pathAccountPurchasingPower = "/accounts/%s/ppse"
)

type Client struct {
	client    *resty.Client
	apiKey    string
	apiSecret string
	logger    *zap.Logger
}

// get sends an HMAC-signed GET request to a path and decodes the JSON response into a result.
func (c *Client) get(ctx context.Context, path string, result any) error {
	date, signature, err := c.sign("GET", path)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
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

// sign builds the HMAC-SHA256 auth headers DNSE requires on every request.
// https://developers.dnse.com.vn/docs/guide/intro/authentication
func (c *Client) sign(method, path string) (date, signature string, err error) {
	date = time.Now().UTC().Format(time.RFC1123Z)

	nonceBytes := make([]byte, 16)
	if _, err = rand.Read(nonceBytes); err != nil {
		return "", "", fmt.Errorf("dnse: generate nonce: %w", err)
	}
	nonce := hex.EncodeToString(nonceBytes)

	sigString := fmt.Sprintf(
		"(request-target): %s %s\ndate: %s\nnonce: %s",
		strings.ToLower(method), path, date, nonce)

	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(sigString))
	sig := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	signature = fmt.Sprintf(
		`Signature keyId="%s",algorithm="hmac-sha256",headers="(request-target) date",signature="%s",nonce="%s"`,
		c.apiKey, sig, nonce,
	)
	return date, signature, nil
}
