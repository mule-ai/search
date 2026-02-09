// Package searxng provides benchmark tests for JSON decoding performance.
package searxng

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

// sampleSearchResponse is a realistic SearXNG response for benchmarking.
var sampleSearchResponse = []byte(`{
	"query": "golang tutorial",
	"results": [
		{
			"title": "A Tour of Go",
			"url": "https://go.dev/tour/",
			"content": "Welcome to a tour of the Go programming language. The Tour is divided into a list of modules that you can access by clicking on A Tour of Go on the top left of the page.",
			"engine": "google",
			"category": "general",
			"score": 0.95,
			"parsed_url": ["go.dev", "tour"],
			"template": "default.html"
		},
		{
			"title": "The Go Programming Language",
			"url": "https://go.dev/",
			"content": "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
			"engine": "duckduckgo",
			"category": "general",
			"score": 0.92,
			"parsed_url": ["go.dev"],
			"template": "default.html"
		},
		{
			"title": "Go by Example",
			"url": "https://gobyexample.com/",
			"content": "Go by Example is a hands-on introduction to Go using annotated example programs. Check out the first example or browse the full list below.",
			"engine": "bing",
			"category": "general",
			"score": 0.88,
			"parsed_url": ["gobyexample.com"],
			"template": "default.html"
		}
	],
	"answers": [],
	"infoboxes": [],
	"suggestions": ["golang tutorial pdf", "golang tutorial w3schools", "golang tutorial point"],
	"number_of_results": 1250000
}`)

// largeSearchResponse simulates a response with many results.
func generateLargeResponse(numResults int) []byte {
	var buf bytes.Buffer
	buf.WriteString(`{"query":"test","results":[`)
	for i := 0; i < numResults; i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{
			"title":"Result ` + string(rune('0'+i%10)) + `",
			"url":"https://example.com/` + string(rune('0'+i%10)) + `",
			"content":"This is a test result with some content",
			"engine":"google",
			"category":"general",
			"score":0.` + string(rune('0'+(9-i%10))) + `
		}`)
	}
	buf.WriteString(`],"answers":[],"infoboxes":[],"suggestions":[],"number_of_results":` + string(rune('0'+numResults%10)) + `00000}`)
	return buf.Bytes()
}

// Benchmark_StandardJSONUnmarshal benchmarks the standard json.Unmarshal.
func Benchmark_StandardJSONUnmarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var resp SearchResponse
		if err := json.Unmarshal(sampleSearchResponse, &resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_OptimizedDecoder benchmarks the optimized decoder.
func Benchmark_OptimizedDecoder(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := NewOptimizedDecoder(bytes.NewReader(sampleSearchResponse))
		var resp SearchResponse
		if err := decoder.Decode(&resp); err != nil {
			b.Fatal(err)
		}
		decoder.Close()
	}
}

// Benchmark_UnmarshalResponse benchmarks the UnmarshalResponse function.
func Benchmark_UnmarshalResponse(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var resp SearchResponse
		if err := UnmarshalResponse(sampleSearchResponse, &resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_StandardJSONDecoder benchmarks the standard json.Decoder.
func Benchmark_StandardJSONDecoder(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := json.NewDecoder(bytes.NewReader(sampleSearchResponse))
		var resp SearchResponse
		if err := decoder.Decode(&resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_LargeResponse_StandardUnmarshal benchmarks standard unmarshal on large response.
func Benchmark_LargeResponse_StandardUnmarshal(b *testing.B) {
	largeData := generateLargeResponse(100)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var resp SearchResponse
		if err := json.Unmarshal(largeData, &resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_LargeResponse_Optimized benchmarks optimized decoder on large response.
func Benchmark_LargeResponse_Optimized(b *testing.B) {
	largeData := generateLargeResponse(100)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := NewOptimizedDecoder(bytes.NewReader(largeData))
		var resp SearchResponse
		if err := decoder.Decode(&resp); err != nil {
			b.Fatal(err)
		}
		decoder.Close()
	}
}

// Benchmark_DecodeResponse benchmarks the DecodeResponse convenience function.
func Benchmark_DecodeResponse(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := DecodeResponse(bytes.NewReader(sampleSearchResponse))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_StreamingDecoder benchmarks the streaming decoder.
func Benchmark_StreamingDecoder(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := NewStreamingDecoder(bytes.NewReader(sampleSearchResponse))
		var resp SearchResponse
		if err := decoder.DecodeNext(&resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_WithPoolUsage benchmarks the effect of buffer pooling.
func Benchmark_WithPoolUsage(b *testing.B) {
	b.Run("WithPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := GetBuffer()
			buf.Write(sampleSearchResponse)
			decoder := NewOptimizedDecoder(buf)
			var resp SearchResponse
			if err := decoder.Decode(&resp); err != nil {
				b.Fatal(err)
			}
			decoder.Close()
			PutBuffer(buf)
		}
	})

	b.Run("WithoutPool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf := new(bytes.Buffer)
			buf.Write(sampleSearchResponse)
			decoder := json.NewDecoder(buf)
			var resp SearchResponse
			if err := decoder.Decode(&resp); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark_DecoderReuse benchmarks reusing decoder vs creating new ones.
func Benchmark_DecoderReuse(b *testing.B) {
	b.Run("ReuseDecoder", func(b *testing.B) {
		b.ReportAllocs()
		data := sampleSearchResponse
		decoder := NewOptimizedDecoder(bytes.NewReader(data))
		defer decoder.Close()
		var resp SearchResponse

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// In practice, you'd need to reset the reader
			reader := bytes.NewReader(data)
			decoder.decoder = json.NewDecoder(reader)
			if err := decoder.Decode(&resp); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("NewDecoderEachTime", func(b *testing.B) {
		b.ReportAllocs()
		data := sampleSearchResponse

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			decoder := NewOptimizedDecoder(bytes.NewReader(data))
			var resp SearchResponse
			if err := decoder.Decode(&resp); err != nil {
				b.Fatal(err)
			}
			decoder.Close()
		}
	})
}

// Benchmark_RealWorld benchmarks simulating real-world usage with multiple decodes.
func Benchmark_RealWorld(b *testing.B) {
	responses := make([][]byte, 10)
	for i := range responses {
		responses[i] = generateLargeResponse(50)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Simulate processing multiple responses
		for _, data := range responses {
			decoder := NewOptimizedDecoder(bytes.NewReader(data))
			var resp SearchResponse
			if err := decoder.Decode(&resp); err != nil {
				b.Fatal(err)
			}
			decoder.Close()
		}
	}
}

// Benchmark_CustomUnmarshal benchmarks the custom unmarshal methods.
func Benchmark_CustomUnmarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var resp SearchResponse
		if err := resp.UnmarshalJSON(sampleSearchResponse); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_ReadAllAndUnmarshal benchmarks reading all data then unmarshaling.
func Benchmark_ReadAllAndUnmarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(sampleSearchResponse)
		data, err := io.ReadAll(reader)
		if err != nil {
			b.Fatal(err)
		}
		var resp SearchResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_DirectDecoder benchmarks direct decoder usage (no intermediate read).
func Benchmark_DirectDecoder(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := json.NewDecoder(bytes.NewReader(sampleSearchResponse))
		var resp SearchResponse
		if err := decoder.Decode(&resp); err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_OptimizedDecoderWithOptions benchmarks decoder with custom options.
func Benchmark_OptimizedDecoderWithOptions(b *testing.B) {
	opts := &DecoderOptions{
		UseStreaming:   false,
		BufferSize:     8192,
		DisablePooling: false,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		decoder := NewDecoderWithOptions(bytes.NewReader(sampleSearchResponse), opts)
		var resp SearchResponse
		if err := decoder.Decode(&resp); err != nil {
			b.Fatal(err)
		}
		decoder.Close()
	}
}