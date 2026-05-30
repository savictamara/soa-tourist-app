package rpc

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const CodecName = "json"

type JSONCodec struct{}

func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Name() string {
	return CodecName
}

func RegisterCodec() {
	encoding.RegisterCodec(JSONCodec{})
}
