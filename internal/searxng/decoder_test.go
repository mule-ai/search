// Package searxng provides tests for optimized JSON decoding.
package searxng

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestOptimizedDecoder verifies the optimized decoder works correctly.
func TestOptimizedDecoder(t *testing.T) {
	data := []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	decoder := NewOptimizedDecoder(bytes.NewReader(data))
	defer decoder.Close()

	var resp SearchResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if resp.Query != "test" {
		t.Errorf("Query = %v, want 'test'", resp.Query)
	}
}

// TestDecodeResponse verifies the DecodeResponse convenience function.
func TestDecodeResponse(t *testing.T) {
	data := []byte(`{"query":"golang","results":[{"title":"Test","url":"https://example.com","content":"Content"}],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":1}`)

	resp, err := DecodeResponse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("DecodeResponse() error = %v", err)
	}

	if resp.Query != "golang" {
		t.Errorf("Query = %v, want 'golang'", resp.Query)
	}

	if len(resp.Results) != 1 {
		t.Errorf("Results length = %v, want 1", len(resp.Results))
	}

	if resp.Results[0].Title != "Test" {
		t.Errorf("Result title = %v, want 'Test'", resp.Results[0].Title)
	}
}

// TestUnmarshalResponse verifies the UnmarshalResponse function.
func TestUnmarshalResponse(t *testing.T) {
	data := []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	var resp SearchResponse
	if err := UnmarshalResponse(data, &resp); err != nil {
		t.Fatalf("UnmarshalResponse() error = %v", err)
	}

	if resp.Query != "test" {
		t.Errorf("Query = %v, want 'test'", resp.Query)
	}
}

// TestUnmarshalResults verifies the UnmarshalResults function.
func TestUnmarshalResults(t *testing.T) {
	data := []byte(`{"query":"test","results":[{"title":"Result 1"},{"title":"Result 2"}],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":2}`)

	results, err := UnmarshalResults(data)
	if err != nil {
		t.Fatalf("UnmarshalResults() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Results length = %v, want 2", len(results))
	}

	if results[0].Title != "Result 1" {
		t.Errorf("Result 0 title = %v, want 'Result 1'", results[0].Title)
	}

	if results[1].Title != "Result 2" {
		t.Errorf("Result 1 title = %v, want 'Result 2'", results[1].Title)
	}
}

// TestStreamingDecoder verifies the streaming decoder.
func TestStreamingDecoder(t *testing.T) {
	data := []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	decoder := NewStreamingDecoder(bytes.NewReader(data))

	var resp SearchResponse
	if err := decoder.DecodeNext(&resp); err != nil {
		t.Fatalf("DecodeNext() error = %v", err)
	}

	if resp.Query != "test" {
		t.Errorf("Query = %v, want 'test'", resp.Query)
	}

	if decoder.More() {
		t.Error("More() = true, want false")
	}
}

// TestDecoderWithBufferPooling verifies buffer pooling works correctly.
func TestDecoderWithBufferPooling(t *testing.T) {
	// Get multiple buffers and return them to the pool
	for i := 0; i < 10; i++ {
		buf := GetBuffer()
		buf.WriteString("test data")
		PutBuffer(buf)
	}

	// Verify we can still get and use a buffer
	buf := GetBuffer()
	defer PutBuffer(buf)

	if buf.String() != "" {
		t.Errorf("Buffer not reset, got = %v", buf.String())
	}

	buf.WriteString("hello")
	if buf.String() != "hello" {
		t.Errorf("Buffer write failed, got = %v", buf.String())
	}
}

// TestByteSlicePooling verifies byte slice pooling works correctly.
func TestByteSlicePooling(t *testing.T) {
	// Get multiple slices and return them to the pool
	for i := 0; i < 10; i++ {
		data := GetByteSlice()
		*data = append(*data, []byte("test")...)
		PutByteSlice(data)
	}

	// Verify we can still get and use a slice
	data := GetByteSlice()
	defer PutByteSlice(data)

	if len(*data) != 0 {
		t.Errorf("Slice not reset, got length = %v", len(*data))
	}

	*data = append(*data, []byte("hello")...)
	if string(*data) != "hello" {
		t.Errorf("Slice append failed, got = %v", string(*data))
	}
}

// TestNewDecoderWithOptions verifies decoder with custom options.
func TestNewDecoderWithOptions(t *testing.T) {
	data := []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	opts := &DecoderOptions{
		UseStreaming:   false,
		BufferSize:     8192,
		DisablePooling: false,
	}

	decoder := NewDecoderWithOptions(bytes.NewReader(data), opts)
	defer decoder.Close()

	var resp SearchResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if resp.Query != "test" {
		t.Errorf("Query = %v, want 'test'", resp.Query)
	}
}

// TestSearchResponseUnmarshalJSON verifies the custom unmarshaler handles edge cases.
func TestSearchResponseUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		check   func(*SearchResponse) bool
	}{
		{
			name: "standard response",
			data: []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":100}`),
			wantErr: false,
			check: func(r *SearchResponse) bool {
				return r.Query == "test" && r.NumberOfResults == 100
			},
		},
		{
			name: "number_of_results as string",
			data: []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":"100"}`),
			wantErr: false,
			check: func(r *SearchResponse) bool {
				return r.NumberOfResults == 100
			},
		},
		{
			name: "number_of_results missing",
			data: []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[]}`),
			wantErr: false,
			check: func(r *SearchResponse) bool {
				return r.NumberOfResults == 0
			},
		},
		{
			name: "number_of_results null",
			data: []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":null}`),
			wantErr: false,
			check: func(r *SearchResponse) bool {
				return r.NumberOfResults == 0
			},
		},
		{
			name: "slices are initialized",
			data: []byte(`{"query":"test","results":[{"title":"Test"}],"answers":[],"infoboxes":[],"suggestions":[]}`),
			wantErr: false,
			check: func(r *SearchResponse) bool {
				return r.Results != nil && r.Answers != nil && r.Infoboxes != nil && r.Suggestions != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp SearchResponse
			err := resp.UnmarshalJSON(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.check(&resp) {
				t.Error("Check function failed")
			}
		})
	}
}

// TestSearchResultUnmarshalJSON verifies the custom unmarshaler for results.
func TestSearchResultUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		check   func(*SearchResult) bool
	}{
		{
			name: "standard result",
			data: []byte(`{"title":"Test","url":"https://example.com","content":"Content","engine":"google","category":"general","score":0.9}`),
			wantErr: false,
			check: func(r *SearchResult) bool {
				return r.Title == "Test" && r.Score == 0.9
			},
		},
		{
			name: "score missing",
			data: []byte(`{"title":"Test","url":"https://example.com","content":"Content"}`),
			wantErr: false,
			check: func(r *SearchResult) bool {
				return r.Score == 0.0
			},
		},
		{
			name: "parsed_url initialized",
			data: []byte(`{"title":"Test","url":"https://example.com"}`),
			wantErr: false,
			check: func(r *SearchResult) bool {
				return r.ParsedURL != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result SearchResult
			err := result.UnmarshalJSON(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.check(&result) {
				t.Error("Check function failed")
			}
		})
	}
}

// TestDecoderClose verifies Close() properly returns buffers to pool.
func TestDecoderClose(t *testing.T) {
	data := []byte(`{"query":"test","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	decoder := NewOptimizedDecoder(bytes.NewReader(data))

	var resp SearchResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	// Close should return buffer to pool
	decoder.Close()

	// Verify we can get another buffer (means pool is working)
	buf := GetBuffer()
	PutBuffer(buf)
}

// TestMultipleDecodes verifies multiple decode operations work correctly.
func TestMultipleDecodes(t *testing.T) {
	data1 := []byte(`{"query":"test1","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)
	data2 := []byte(`{"query":"test2","results":[],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":0}`)

	// First decode
	decoder1 := NewOptimizedDecoder(bytes.NewReader(data1))
	var resp1 SearchResponse
	if err := decoder1.Decode(&resp1); err != nil {
		t.Fatalf("First Decode() error = %v", err)
	}
	decoder1.Close()

	if resp1.Query != "test1" {
		t.Errorf("First query = %v, want 'test1'", resp1.Query)
	}

	// Second decode
	decoder2 := NewOptimizedDecoder(bytes.NewReader(data2))
	var resp2 SearchResponse
	if err := decoder2.Decode(&resp2); err != nil {
		t.Fatalf("Second Decode() error = %v", err)
	}
	decoder2.Close()

	if resp2.Query != "test2" {
		t.Errorf("Second query = %v, want 'test2'", resp2.Query)
	}
}

// TestOptimizedDecoderVsStandard compares results between optimized and standard decoders.
func TestOptimizedDecoderVsStandard(t *testing.T) {
	data := []byte(`{"query":"test","results":[{"title":"Test Result","url":"https://example.com","content":"Test content","engine":"google","category":"general","score":0.95}],"answers":[],"infoboxes":[],"suggestions":["suggestion1","suggestion2"],"number_of_results":100}`)

	// Decode with standard decoder
	var standardResp SearchResponse
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&standardResp); err != nil {
		t.Fatalf("Standard decoder error = %v", err)
	}

	// Decode with optimized decoder
	decoder := NewOptimizedDecoder(bytes.NewReader(data))
	defer decoder.Close()

	var optimizedResp SearchResponse
	if err := decoder.Decode(&optimizedResp); err != nil {
		t.Fatalf("Optimized decoder error = %v", err)
	}

	// Compare results
	if standardResp.Query != optimizedResp.Query {
		t.Errorf("Query mismatch: standard=%v, optimized=%v", standardResp.Query, optimizedResp.Query)
	}

	if len(standardResp.Results) != len(optimizedResp.Results) {
		t.Errorf("Results length mismatch: standard=%v, optimized=%v", len(standardResp.Results), len(optimizedResp.Results))
	}

	if len(standardResp.Suggestions) != len(optimizedResp.Suggestions) {
		t.Errorf("Suggestions length mismatch: standard=%v, optimized=%v", len(standardResp.Suggestions), len(optimizedResp.Suggestions))
	}

	if standardResp.NumberOfResults != optimizedResp.NumberOfResults {
		t.Errorf("NumberOfResults mismatch: standard=%v, optimized=%v", standardResp.NumberOfResults, optimizedResp.NumberOfResults)
	}
}

// TestDecoderWithLargeResponse verifies decoder handles large responses.
func TestDecoderWithLargeResponse(t *testing.T) {
	// Create a large response with many results
	var buf bytes.Buffer
	buf.WriteString(`{"query":"test","results":[`)
	for i := 0; i < 100; i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"title":"Result ` + string(rune('0'+i%10)) + `","url":"https://example.com/` + string(rune('0'+i%10)) + `","content":"Content"}`)
	}
	buf.WriteString(`],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":100}`)

	decoder := NewOptimizedDecoder(&buf)
	defer decoder.Close()

	var resp SearchResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if len(resp.Results) != 100 {
		t.Errorf("Results length = %v, want 100", len(resp.Results))
	}
}