package pool

import (
	"bytes"
	"sync"
)

var _bufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// GetBufferFromPool get buffer
func GetBufferFromPool() *bytes.Buffer {
	return _bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer2Pool return buffer after reset
func PutBuffer2Pool(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	_bufferPool.Put(buf)
}
