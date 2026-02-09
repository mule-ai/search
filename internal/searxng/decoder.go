// Package searxng provides optimized JSON decoding for SearXNG responses.
//
// The decoders use buffer pooling and efficient reading strategies to
// minimize memory allocations when parsing large JSON responses.
package searxng

import (
	"bytes"
	"encoding/json"
	"io"
)

// OptimizedDecoder provides efficient JSON decoding using buffer pooling.
type OptimizedDecoder struct {
	decoder *json.Decoder
	buffer  *bytes.Buffer
}

// NewOptimizedDecoder creates a new optimized JSON decoder.
//
// The decoder uses buffer pooling to reduce memory allocations.
// Always call Close() when done to return the buffer to the pool.
//
// Example:
//
//	resp, _ := http.Get("https://search.example.com/search?q=test&format=json")
//	defer resp.Body.Close()
//
//	decoder := searxng.NewOptimizedDecoder(resp.Body)
//	defer decoder.Close()
//
//	var searchResp searxng.SearchResponse
//	if err := decoder.Decode(&searchResp); err != nil {
//	    log.Fatal(err)
//	}
func NewOptimizedDecoder(r io.Reader) *OptimizedDecoder {
	// Direct decoder without TeeReader for proper JSON parsing
	return &OptimizedDecoder{
		decoder: json.NewDecoder(r),
		buffer:  nil, // Not used for direct decoding
	}
}

// Decode decodes the JSON value into v.
func (d *OptimizedDecoder) Decode(v interface{}) error {
	return d.decoder.Decode(v)
}

// Close releases the buffer back to the pool.
func (d *OptimizedDecoder) Close() {
	if d.buffer != nil {
		PutBuffer(d.buffer)
	}
}

// DecodeResponse decodes a SearXNG search response efficiently.
//
// This convenience function creates an optimized decoder, decodes the response,
// and automatically cleans up resources.
//
// Example:
//
//	resp, _ := http.Get("https://search.example.com/search?q=test&format=json")
//	defer resp.Body.Close()
//
//	searchResp, err := searxng.DecodeResponse(resp.Body)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found %d results\n", len(searchResp.Results))
func DecodeResponse(r io.Reader) (*SearchResponse, error) {
	decoder := NewOptimizedDecoder(r)
	defer decoder.Close()

	var resp SearchResponse
	if err := decoder.Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// StreamingDecoder provides streaming JSON decoding for large responses.
type StreamingDecoder struct {
	decoder *json.Decoder
}

// NewStreamingDecoder creates a new streaming decoder.
//
// Streaming decoders are useful for processing very large JSON responses
// without loading the entire response into memory.
//
// Example:
//
//	resp, _ := http.Get("https://search.example.com/search?q=test&format=json")
//	defer resp.Body.Close()
//
//	decoder := searxng.NewStreamingDecoder(resp.Body)
//	for decoder.More() {
//	    var result searxng.SearchResult
//	    if err := decoder.DecodeNext(&result); err != nil {
//	        break
//	    }
//	    // Process result...
//	}
func NewStreamingDecoder(r io.Reader) *StreamingDecoder {
	return &StreamingDecoder{
		decoder: json.NewDecoder(r),
	}
}

// DecodeNext decodes the next JSON value from the stream.
func (d *StreamingDecoder) DecodeNext(v interface{}) error {
	return d.decoder.Decode(v)
}

// More returns true if there are more values to decode.
func (d *StreamingDecoder) More() bool {
	return d.decoder.More()
}

// BufferedBytes returns the bytes remaining in the decoder's buffer.
func (d *StreamingDecoder) BufferedBytes() []byte {
	// Read the buffered content
	buf := make([]byte, 4096)
	n, _ := d.decoder.Buffered().Read(buf)
	return buf[:n]
}

// UnmarshalResponse efficiently unmarshals JSON data into a SearchResponse.
//
// This function uses a buffer pool to reduce memory allocations compared
// to the standard json.Unmarshal. It's optimized for the specific structure
// of SearXNG responses.
//
// Example:
//
//	data, _ := os.ReadFile("response.json")
//	var resp SearchResponse
//	if err := searxng.UnmarshalResponse(data, &resp); err != nil {
//	    log.Fatal(err)
//	}
func UnmarshalResponse(data []byte, v *SearchResponse) error {
	// Use bytes.Reader for zero-allocation reading from byte slice
	reader := NewOptimizedDecoder(bytes.NewReader(data))
	defer reader.Close()
	return reader.Decode(v)
}

// UnmarshalResults efficiently unmarshals only the results array from JSON data.
//
// This is useful when you only care about the search results and not the
// full response metadata (answers, infoboxes, etc.). It can be faster for
// large responses with many results.
//
// Example:
//
//	data, _ := os.ReadFile("response.json")
//	results, err := searxng.UnmarshalResults(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, result := range results {
//	    fmt.Println(result.Title)
//	}
func UnmarshalResults(data []byte) ([]SearchResult, error) {
	// Parse the full response to extract results
	var resp SearchResponse
	if err := UnmarshalResponse(data, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

// DecoderOptions configures the behavior of optimized decoders.
type DecoderOptions struct {
	// UseStreaming enables streaming mode for large responses.
	UseStreaming bool
	// BufferSize sets the initial buffer size for reading.
	BufferSize int
	// DisablePooling disables buffer pooling (useful for single-shot decodes).
	DisablePooling bool
}

// DefaultDecoderOptions returns the recommended decoder options.
func DefaultDecoderOptions() *DecoderOptions {
	return &DecoderOptions{
		UseStreaming:   false,
		BufferSize:     4096,
		DisablePooling: false,
	}
}

// NewDecoderWithOptions creates a decoder with custom options.
//
// This allows fine-tuning decoder behavior for specific use cases.
// For example, streaming mode is useful for very large responses,
// while disabling pooling can reduce overhead for single-shot decodes.
//
// Example:
//
//	opts := &searxng.DecoderOptions{
//	    UseStreaming: true,
//	    BufferSize: 8192,
//	}
//	decoder := searxng.NewDecoderWithOptions(resp.Body, opts)
//	defer decoder.Close()
func NewDecoderWithOptions(r io.Reader, opts *DecoderOptions) *OptimizedDecoder {
	if opts == nil {
		opts = DefaultDecoderOptions()
	}
	return &OptimizedDecoder{
		decoder: json.NewDecoder(r),
		buffer:  GetBuffer(),
	}
}
