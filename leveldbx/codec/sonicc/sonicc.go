package sonicc

import (
	"github.com/bytedance/sonic"
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec"
)

type SonicCodec[K any, V any] struct{}

// NewJSONCodec new json codec
func NewJSONCodec[K any, V any]() *SonicCodec[K, V] {
	return &SonicCodec[K, V]{}
}

var _ codec.LdbCodec[any, any] = (*SonicCodec[any, any])(nil)

// EncodeVal json encode value
func (c *SonicCodec[K, V]) EncodeVal(data V) (v []byte, err error) {
	v, err = sonic.Marshal(data)
	return
}

// DecodeVal json decode value
func (c *SonicCodec[K, V]) DecodeVal(data []byte) (v V, err error) {
	err = sonic.Unmarshal(data, &v)
	return
}

// EncodeKey json encode key
func (c *SonicCodec[K, V]) EncodeKey(data K) (k []byte, err error) {
	k, err = sonic.Marshal(data)
	return
}

// DecodeKey json decode key
func (c *SonicCodec[K, V]) DecodeKey(data []byte) (k K, err error) {
	err = sonic.Unmarshal(data, &k)
	return
}
