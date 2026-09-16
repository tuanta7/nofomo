package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type Encoding string

const (
	EncodingJSON    Encoding = "json"
	EncodingMsgPack Encoding = "msgpack"
)

type Encoder struct {
	encoding Encoding
}

func NewEncoder(encoding Encoding) *Encoder {
	return &Encoder{encoding: encoding}
}

func (e *Encoder) Encode(data interface{}) ([]byte, error) {
	switch e.encoding {
	case EncodingJSON:
		return json.Marshal(data)
	case EncodingMsgPack:
		return msgpack.Marshal(data)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", e.encoding)
	}
}

type Decoder struct {
	encoding Encoding
}

func NewDecoder(encoding Encoding) *Decoder {
	return &Decoder{encoding: encoding}
}

func (d *Decoder) Decode(data []byte, v any) error {
	switch d.encoding {
	case EncodingJSON:
		return json.Unmarshal(data, &v)
	case EncodingMsgPack:
		return msgpack.Unmarshal(data, &v)
	default:
		return fmt.Errorf("unsupported encoding: %s", d.encoding)
	}
}
