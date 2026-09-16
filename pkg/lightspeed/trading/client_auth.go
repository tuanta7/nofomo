package trading

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Sign builds the HMAC-SHA256 auth headers DNSE requires on every request.
// https://developers.dnse.com.vn/docs/guide/intro/authentication
func (c *Client) Sign(method, path string) (date, signature string, err error) {
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
