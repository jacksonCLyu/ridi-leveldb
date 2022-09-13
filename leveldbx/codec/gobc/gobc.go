package gobc

import (
	"encoding/gob"
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/codec"
	"github.com/jacksonCLyu/ridi-leveldb/leveldbx/pool"
)

type GobCodec[K any, V any] struct{}

func NewGobCodec[K any, V any]() *GobCodec[K, V] {
	return &GobCodec[K, V]{}
}

var _ codec.Lcodec[any, any] = (*GobCodec[any, any])(nil)

// EncodeVal gob encode value
func (c *GobCodec[K, V]) EncodeVal(data V) (v []byte, err error) {
	buf := pool.GetBufferFromPool()
	defer pool.PutBuffer2Pool(buf)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	// 解决 leveldb: not found 问题，原因：buf.Bytes() 返回内存强关联 pool.Buffer 的底层切片
	v = make([]byte, 0)
	v = append(v, buf.Bytes()...)
	return
}

// DecodeVal gob decode value
func (c *GobCodec[K, V]) DecodeVal(data []byte) (v V, err error) {
	buf := pool.GetBufferFromPool()
	defer pool.PutBuffer2Pool(buf)
	if _, err = buf.Write(data); err != nil {
		return
	}
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&v)
	return
}

func (c *GobCodec[K, V]) EncodeKey(data K) (k []byte, err error) {
	buf := pool.GetBufferFromPool()
	defer pool.PutBuffer2Pool(buf)
	enc := gob.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	// 解决 leveldb: not found 问题，原因：buf.Bytes() 返回内存强关联 pool.Buffer 的底层切片
	k = make([]byte, 0)
	k = append(k, buf.Bytes()...)
	return
}

// DecodeKey gob decode key
func (c *GobCodec[K, V]) DecodeKey(data []byte) (k K, err error) {
	buf := pool.GetBufferFromPool()
	defer pool.PutBuffer2Pool(buf)
	if _, err = buf.Write(data); err != nil {
		return
	}
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&k)
	return
}
