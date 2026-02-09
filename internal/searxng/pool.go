// Package searxng provides buffer pooling for performance optimization.
//
// Buffer pooling reduces memory allocation overhead by reusing buffers
// across multiple operations. This is particularly useful for formatting
// and JSON encoding operations.
package searxng

import (
	"bytes"
	"sync"
)

// bufferPool is a sync.Pool for reusing bytes.Buffer instances.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer retrieves a buffer from the pool.
//
// The returned buffer is reset and ready for use. When done with the buffer,
// return it to the pool using PutBuffer to reduce memory allocations.
//
// Example:
//
//	buf := searxng.GetBuffer()
//	defer searxng.PutBuffer(buf)
//	buf.WriteString("Hello, World!")
//	// Use buffer...
func GetBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool.
//
// Buffers larger than 1MB are discarded to prevent excessive memory usage.
// Always return buffers obtained from GetBuffer to enable reuse.
//
// Example:
//
//	buf := searxng.GetBuffer()
//	// ... use buffer ...
//	searxng.PutBuffer(buf)
func PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() < 1<<20 { // Only pool buffers smaller than 1MB
		bufferPool.Put(buf)
	}
}

// byteSlicePool is a sync.Pool for reusing byte slices.
var byteSlicePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 4096)
		return &b
	},
}

// GetByteSlice retrieves a byte slice from the pool.
//
// The returned slice has zero length but non-zero capacity. When done with
// the slice, return it to the pool using PutByteSlice.
//
// Example:
//
//	data := searxng.GetByteSlice()
//	defer searxng.PutByteSlice(data)
//	*data = append(*data, 'H', 'e', 'l', 'l', 'o')
//	// Use data...
func GetByteSlice() *[]byte {
	p := byteSlicePool.Get().(*[]byte)
	*p = (*p)[:0]
	return p
}

// PutByteSlice returns a byte slice to the pool.
//
// Slices with capacity larger than 1MB are discarded to prevent excessive
// memory usage. Always return slices obtained from GetByteSlice.
//
// Example:
//
//	data := searxng.GetByteSlice()
//	// ... use data ...
//	searxng.PutByteSlice(data)
func PutByteSlice(p *[]byte) {
	if cap(*p) < 1<<20 { // Only pool slices smaller than 1MB
		byteSlicePool.Put(p)
	}
}
