package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type encodingPayload struct {
	Symbol string `json:"symbol" msgpack:"symbol"`
}

func TestEncoderEncode(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		data     any
		want     []byte
	}{
		{"json struct", EncodingJSON, encodingPayload{Symbol: "ABC"}, []byte(`{"symbol":"ABC"}`)},
		{"msgpack struct", EncodingMsgPack, encodingPayload{Symbol: "ABC"}, []byte{0x81, 0xa6, 's', 'y', 'm', 'b', 'o', 'l', 0xa3, 'A', 'B', 'C'}},
		{"json nil", EncodingJSON, nil, []byte("null")},
		{"msgpack nil", EncodingMsgPack, nil, []byte{0xc0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEncoder(tt.encoding).Encode(tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncoderEncodeUnsupportedType(t *testing.T) {
	for _, encoding := range []Encoding{EncodingJSON, EncodingMsgPack} {
		t.Run(string(encoding), func(t *testing.T) {
			_, err := NewEncoder(encoding).Encode(make(chan int))
			assert.Error(t, err)
		})
	}
}

func TestDecoderDecode(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		data     []byte
		want     encodingPayload
	}{
		{"json struct", EncodingJSON, []byte(`{"symbol":"ABC"}`), encodingPayload{Symbol: "ABC"}},
		{"msgpack struct", EncodingMsgPack, []byte{0x81, 0xa6, 's', 'y', 'm', 'b', 'o', 'l', 0xa3, 'A', 'B', 'C'}, encodingPayload{Symbol: "ABC"}},
		{"json empty object", EncodingJSON, []byte("{}"), encodingPayload{}},
		{"msgpack empty map", EncodingMsgPack, []byte{0x80}, encodingPayload{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got encodingPayload
			err := NewDecoder(tt.encoding).Decode(tt.data, &got)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecoderDecodeInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		encoding Encoding
		data     []byte
	}{
		{"json empty", EncodingJSON, nil},
		{"msgpack empty", EncodingMsgPack, nil},
		{"json malformed", EncodingJSON, []byte(`{"symbol":`)},
		{"msgpack malformed", EncodingMsgPack, []byte{0xc1}},
		{"json wrong field type", EncodingJSON, []byte(`{"symbol":123}`)},
		{"msgpack wrong field type", EncodingMsgPack, []byte{0x81, 0xa6, 's', 'y', 'm', 'b', 'o', 'l', 0x7b}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got encodingPayload
			assert.Error(t, NewDecoder(tt.encoding).Decode(tt.data, &got))
		})
	}
}

func TestUnsupportedEncoding(t *testing.T) {
	for _, encoding := range []Encoding{"", "xml"} {
		t.Run(string(encoding), func(t *testing.T) {
			wantError := "unsupported encoding: " + string(encoding)
			data, err := NewEncoder(encoding).Encode(encodingPayload{})
			assert.EqualError(t, err, wantError)
			assert.Nil(t, data)

			var got encodingPayload
			assert.EqualError(t, NewDecoder(encoding).Decode([]byte("{}"), &got), wantError)
		})
	}
}
