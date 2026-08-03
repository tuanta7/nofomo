package binance

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Binance sends a ping frame every 3 minutes; anything longer means a dead link.
const readTimeout = 5 * time.Minute

type WebsocketClient struct {
	mu   sync.Mutex // gorilla allows only one concurrent writer
	id   atomic.Uint64
	conn *websocket.Conn
}

func NewMarketClient(baseEndpoint string) (*WebsocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(baseEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// gorilla replies to pings automatically, but only the handler can extend the deadline.
	// WriteControl is safe to call concurrently with the writes guarded by mu.
	conn.SetReadDeadline(time.Now().Add(readTimeout))
	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	return &WebsocketClient{
		conn: conn,
	}, nil
}

func (c *WebsocketClient) Close() error {
	return c.conn.Close()
}

func (c *WebsocketClient) Subscribe(streams ...StreamName) error {
	return c.send("SUBSCRIBE", streams)
}

func (c *WebsocketClient) Unsubscribe(streams ...StreamName) error {
	return c.send("UNSUBSCRIBE", streams)
}

func (c *WebsocketClient) send(method string, streams []StreamName) error {
	if len(streams) == 0 {
		return nil
	}

	params := make([]string, len(streams))
	for i, s := range streams {
		params[i] = s.String()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(map[string]any{
		"method": method,
		"params": params,
		"id":     c.id.Add(1),
	})
}

type Message struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
	Error  *struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}

func (c *WebsocketClient) ReadMessage() (stream string, data json.RawMessage, err error) {
	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(readTimeout))

		var message Message
		if err := c.conn.ReadJSON(&message); err != nil {
			return "", nil, err
		}

		// A bad stream name is acknowledged with an error, never with data;
		// without this the subscription just goes quiet forever.
		if message.Error != nil {
			return "", nil, fmt.Errorf("binance: %s (code %d)", message.Error.Msg, message.Error.Code)
		}

		if len(message.Data) == 0 {
			continue
		}

		return message.Stream, message.Data, nil
	}
}
