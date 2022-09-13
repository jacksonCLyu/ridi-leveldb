package jsonc

import (
	"encoding/json"

	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec"
)

type JSONCodec[K any, V any] struct{}

// NewJSONCodec new json codec
func NewJSONCodec[K any, V any]() *JSONCodec[K, V] {
	return &JSONCodec[K, V]{}
}

var _ codec.Lcodec[any, any] = (*JSONCodec[any, any])(nil)

// EncodeVal json encode value
func (c *JSONCodec[K, V]) EncodeVal(data V) (v []byte, err error) {
	v, err = json.Marshal(data)
	return
}

// DecodeVal json decode value
func (c *JSONCodec[K, V]) DecodeVal(data []byte) (v V, err error) {
	err = json.Unmarshal(data, &v)
	return
}

// EncodeKey json encode key
func (c *JSONCodec[K, V]) EncodeKey(data K) (k []byte, err error) {
	k, err = json.Marshal(data)
	return
}

// DecodeKey json decode key
func (c *JSONCodec[K, V]) DecodeKey(data []byte) (k K, err error) {
	err = json.Unmarshal(data, &k)
	return
}
