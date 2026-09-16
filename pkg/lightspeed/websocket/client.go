package websocket

import (
	"errors"

	gorilla "github.com/gorilla/websocket"
)

const (
	maxConnections = 10
	maxStreams     = 200

	baseURL        = "wss://ws-openapi.dnse.com.vn"
	sandboxBaseURL = "wss://ws-sb-openapi.dnse.com.vn"
)

var (
	ErrAlreadyConnected = errors.New("dnse websocket: client is already connected")
	ErrNotConnected     = errors.New("dnse websocket: client is not connected")
	ErrClosed           = errors.New("dnse websocket: client is closed")
)

type Client struct {
	conn     *gorilla.Conn
	encoding Encoding
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(url string) error {
	return nil
}

func (c *Client) Close() error {
	return nil
}
