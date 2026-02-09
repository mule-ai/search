// Package ui provides user interface components for the search CLI.
package ui

import (
	"strings"
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "golang tutorials",
			expected: "golang tutorials",
		},
		{
			name:     "text with control characters",
			input:    "test\x00\x01\x02text",
			expected: "testtext",
		},
		{
			name:     "text with tabs and newlines",
			input:    "line1\nline2\ttabbed",
			expected: "line1\nline2\ttabbed",
		},
		{
			name:     "text with null byte",
			input:    "test\x00null",
			expected: "testnull",
		},
		{
			name:     "text with DEL character",
			input:    "test\x7fDEL",
			expected: "testDEL",
		},
		{
			name:     "whitespace trimming",
			input:    "  spaced  ",
			expected: "spaced",
		},
		{
			name:     "very long input",
			input:    strings.Repeat("a", 3000),
			expected: strings.Repeat("a", 2000),
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only control characters",
			input:    "\x00\x01\x02\x03",
			expected: "",
		},
		{
			name:     "mixed valid and control",
			input:    "hello\x00world\x1b!",
			expected: "helloworld!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("Test message")

	if spinner.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", spinner.message)
	}

	if spinner.active {
		t.Error("Expected spinner to be inactive initially")
	}

	if spinner.writer == nil {
		t.Error("Expected writer to be set")
	}
}

func TestSpinnerStartStop(t *testing.T) {
	spinner := NewSpinner("Testing")

	spinner.Start()
	if !spinner.active {
		t.Error("Expected spinner to be active after Start()")
	}

	spinner.Stop("Done!")
	if spinner.active {
		t.Error("Expected spinner to be inactive after Stop()")
	}
}

func TestSpinnerUpdate(t *testing.T) {
	spinner := NewSpinner("Initial message")

	spinner.Update("Updated message")

	if spinner.message != "Updated message" {
		t.Errorf("Expected message 'Updated message', got '%s'", spinner.message)
	}
}

func TestProgressReporter(t *testing.T) {
	steps := []string{
		"Step 1",
		"Step 2",
		"Step 3",
	}

	reporter := NewProgressReporter(steps, true)

	if len(reporter.steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(reporter.steps))
	}

	if reporter.currentStep != 0 {
		t.Errorf("Expected current step to be 0, got %d", reporter.currentStep)
	}

	// Start the reporter first
	reporter.Start()

	reporter.Next()
	if reporter.currentStep != 1 {
		t.Errorf("Expected current step to be 1 after Next(), got %d", reporter.currentStep)
	}

	reporter.Next()
	if reporter.currentStep != 2 {
		t.Errorf("Expected current step to be 2 after second Next(), got %d", reporter.currentStep)
	}

	reporter.Done("Complete")
}

func TestSearchSpinner(t *testing.T) {
	spinner := NewSearchSpinner(true)

	if spinner == nil {
		t.Fatal("Expected spinner to be created")
	}

	spinner.Start()
	spinner.Stop(10, "0.5s")
	spinner.StopWithError(nil)
}

func TestSearchSpinnerDisabled(t *testing.T) {
	spinner := NewSearchSpinner(false)

	if spinner == nil {
		t.Fatal("Expected spinner to be created")
	}

	// Should not panic when disabled
	spinner.Start()
	spinner.Stop(0, "")
	spinner.StopWithError(nil)
}

func TestOutputFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(string)
	}{
		{"Info", Info},
		{"Warning", Warning},
		{"Error", Error},
		{"Success", Success},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should not panic
			tt.function("Test message")
		})
	}
}

func TestSpinnerFrames(t *testing.T) {
	frames := getSpinnerFrames()

	if len(frames) == 0 {
		t.Error("Expected spinner frames to be non-empty")
	}

	// All frames should be non-empty
	for i, frame := range frames {
		if frame == "" {
			t.Errorf("Frame %d is empty", i)
		}
	}
}

func TestSpinnerConcurrent(t *testing.T) {
	// Test that concurrent operations don't cause race conditions
	spinner := NewSpinner("Concurrent test")

	done := make(chan bool)

	// Start spinner in goroutine
	go func() {
		spinner.Start()
		done <- true
	}()

	// Update in another goroutine
	go func() {
		spinner.Update("Updated")
		done <- true
	}()

	// Stop in main goroutine
	<-done
	<-done
	spinner.Stop("Done")
}

func TestProgressReporterSteps(t *testing.T) {
	tests := []struct {
		name       string
		steps      []string
		verbose    bool
		wantSteps  int
	}{
		{
			name:       "three steps",
			steps:      []string{"a", "b", "c"},
			verbose:    true,
			wantSteps:  3,
		},
		{
			name:       "empty steps",
			steps:      []string{},
			verbose:    true,
			wantSteps:  0,
		},
		{
			name:       "not verbose",
			steps:      []string{"a", "b"},
			verbose:    false,
			wantSteps:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewProgressReporter(tt.steps, tt.verbose)
			if len(reporter.steps) != tt.wantSteps {
				t.Errorf("Expected %d steps, got %d", tt.wantSteps, len(reporter.steps))
			}
		})
	}
}

func TestProgressReporterStartStop(t *testing.T) {
	steps := []string{"Step 1", "Step 2"}
	reporter := NewProgressReporter(steps, true)

	// Should not panic
	reporter.Start()
	reporter.Next()
	reporter.Done("Complete")
}

func TestProgressReporterNonVerbose(t *testing.T) {
	steps := []string{"Step 1", "Step 2"}
	reporter := NewProgressReporter(steps, false)

	// Should not panic when not verbose
	reporter.Start()
	reporter.Next()
	reporter.Done("Complete")
}

func TestSpinnerMultipleStops(t *testing.T) {
	spinner := NewSpinner("Test")

	spinner.Start()
	spinner.Stop("First stop")

	// Second stop should not cause issues
	spinner.Stop("Second stop")

	if spinner.active {
		t.Error("Expected spinner to be inactive after multiple stops")
	}
}

func TestSanitizeInputMaxLength(t *testing.T) {
	// Test that input longer than maxLength is truncated
	longInput := strings.Repeat("x", 3000)
	result := SanitizeInput(longInput)

	const maxLength = 2000
	if len(result) > maxLength {
		t.Errorf("Expected sanitized input to be at most %d characters, got %d", maxLength, len(result))
	}
}

func TestSpinnerNilStopChan(t *testing.T) {
	spinner := NewSpinner("Test")

	// Set active without starting (edge case)
	spinner.active = true
	spinner.stopChan = nil

	// Should not panic
	spinner.Stop("Done")
}

func BenchmarkSanitizeInput(b *testing.B) {
	input := "This is a test string with some characters and numbers 12345"
	for i := 0; i < b.N; i++ {
		SanitizeInput(input)
	}
}

func BenchmarkSanitizeInputWithControls(b *testing.B) {
	input := "Test\x00\x01\x02\x03\x04string\nwith\tcontrol\x7fchars"
	for i := 0; i < b.N; i++ {
		SanitizeInput(input)
	}
}

func BenchmarkSanitizeInputLong(b *testing.B) {
	input := strings.Repeat("a", 1000) + string([]byte{0x00, 0x01, 0x02})
	for i := 0; i < b.N; i++ {
		SanitizeInput(input)
	}
}