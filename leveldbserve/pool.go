package leveldbserve

import (
	"bytes"
	"sync"
)

var _encodeBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

var _decodeBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getDecodeBuf() *bytes.Buffer {
	buf := _decodeBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putDecodeBuf(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	_decodeBufPool.Put(buf)
}

func getEncodeBuf() *bytes.Buffer {
	buf := _encodeBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putEncodeBuf(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	_encodeBufPool.Put(buf)
}
