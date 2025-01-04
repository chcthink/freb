package pool

import (
	"bytes"
	"strings"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func BufferGet() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func BufferPut(buf *bytes.Buffer) {
	if buf != nil {
		buf.Reset()
		bufferPool.Put(buf)
	}
}

var strBuilderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

func strBuilderGet() *strings.Builder {
	return strBuilderPool.Get().(*strings.Builder)
}

func StrBuilderPut(buf *strings.Builder) {
	if buf != nil {
		buf.Reset()
		strBuilderPool.Put(buf)
	}
}
