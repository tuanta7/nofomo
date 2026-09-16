package websocket

import "time"

const (
	DefaultTimeout = 60 * time.Second
)

type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
}
