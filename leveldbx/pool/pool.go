package pool

import (
	"bytes"
	"sync"
)

var _bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBufferFromPool get buffer
func GetBufferFromPool() *bytes.Buffer {
	return _bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer2Pool return buffer after reset **this function return a func() than can be defer call
// Usage:
//
//	 	buf := pool.GetBufferFromPool()
//		defer pool.PutBuffer2Pool(buf)()
func PutBuffer2Pool(buf *bytes.Buffer) func() {
	return func() {
		if buf == nil {
			return
		}
		buf.Reset()
		_bufferPool.Put(buf)
	}
}
