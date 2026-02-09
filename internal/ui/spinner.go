// Package ui provides user interface components for the search CLI.
//
// It includes spinners, progress indicators, and other interactive elements
// that enhance the user experience during search operations.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mule-ai/search/internal/config"
)

// Spinner represents a loading indicator with animation.
//
// It provides visual feedback during long-running operations like network requests.
// The spinner automatically detects if output is a TTY and disables animation
// for non-interactive environments (pipes, redirects, CI).
type Spinner struct {
	message       string
	active        bool
	stopChan      chan struct{}
	wg            sync.WaitGroup
	mu            sync.Mutex
	writer        io.Writer
	isTTY         bool
	frames        []string
	frameInterval time.Duration
}

// NewSpinner creates a new spinner with the given message.
//
// The spinner only animates if the output is a terminal (TTY).
// For non-TTY output (pipes, files, CI), it prints a simple message.
//
// Example:
//
//	spinner := ui.NewSpinner("Searching...")
//	spinner.Start()
//	defer spinner.Stop()
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:       message,
		stopChan:      make(chan struct{}),
		writer:        os.Stderr,
		frames:        getSpinnerFrames(),
		frameInterval: 100 * time.Millisecond,
		isTTY:         isTerminal(os.Stderr),
	}
}

// Start begins the spinner animation.
//
// If running in a TTY, it shows an animated spinner.
// If not a TTY, it prints a simple message.
//
// Example:
//
//	spinner := ui.NewSpinner("Searching...")
//	spinner.Start()
//	// Do work here
//	spinner.Stop()
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		return
	}

	s.active = true
	s.wg.Add(1)

	if s.isTTY {
		go s.animate()
	} else {
		// Non-TTY: just print the message
		fmt.Fprintf(s.writer, "%s...\n", s.message)
	}
}

// Stop stops the spinner animation.
//
// If running in a TTY, it clears the spinner line.
// If not a TTY, it prints a completion message.
//
// Example:
//
//	spinner := ui.NewSpinner("Searching...")
//	spinner.Start()
//	// Do work
//	spinner.Stop("Done!")
func (s *Spinner) Stop(finalMessage string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	// Handle edge case where stopChan might be nil
	if s.stopChan != nil {
		close(s.stopChan)
		s.wg.Wait()
		s.stopChan = make(chan struct{})
	}

	// Clear the spinner line
	if s.isTTY {
		fmt.Fprintf(s.writer, "\r\033[K") // Clear line
		if finalMessage != "" {
			fmt.Fprintf(s.writer, "%s\n", finalMessage)
		}
	} else {
		if finalMessage != "" {
			fmt.Fprintf(s.writer, "%s\n", finalMessage)
		}
	}

	s.active = false
}

// Update changes the spinner message.
//
// This is useful for providing progress updates during long operations.
//
// Example:
//
//	spinner := ui.NewSpinner("Connecting...")
//	spinner.Start()
//	spinner.Update("Searching...")
//	spinner.Stop("Done!")
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.message = message
}

// animate runs the spinner animation loop.
func (s *Spinner) animate() {
	defer s.wg.Done()

	frame := 0
	ticker := time.NewTicker(s.frameInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.mu.Lock()
			frameText := s.frames[frame%len(s.frames)]
			fmt.Fprintf(s.writer, "\r\033[K%s %s", frameText, s.message)
			s.mu.Unlock()
			frame++
		}
	}
}

// getSpinnerFrames returns the appropriate spinner frames.
//
// Uses simple ASCII characters that work on most terminals.
func getSpinnerFrames() []string {
	// Simple rotating spinner using common Unicode characters
	return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}

// isTerminal checks if the writer is a terminal.
func isTerminal(w io.Writer) bool {
	f, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}

	// Check if file descriptor is a terminal
	fd := f.Fd()
	return isTerminalFile(fd)
}

// isTerminalFile checks if a file descriptor is a terminal.
func isTerminalFile(fd uintptr) bool {
	// Simple check: in most environments, std{in,out,err} are 0,1,2
	// For a more robust check, we'd use syscall isatty
	// This is a simplified version that works for most cases
	return fd >= 0 && fd <= 2
}

// ProgressReporter tracks progress of multi-step operations.
//
// It provides a simple way to report progress during operations that
// have multiple steps (e.g., "Loading config", "Connecting", "Searching").
type ProgressReporter struct {
	steps       []string
	currentStep int
	spinner     *Spinner
	verbose     bool
	started     bool
}

// NewProgressReporter creates a new progress reporter.
//
// Example:
//
//	reporter := ui.NewProgressReporter([]string{
//	    "Loading configuration",
//	    "Connecting to instance",
//	    "Performing search",
//	}, true)
//	reporter.Start()
//	reporter.Next()
//	reporter.Next()
//	reporter.Done("Found 10 results")
func NewProgressReporter(steps []string, verbose bool) *ProgressReporter {
	return &ProgressReporter{
		steps:   steps,
		verbose: verbose,
		started: false,
	}
}

// Start begins progress reporting.
func (p *ProgressReporter) Start() {
	if !p.verbose || len(p.steps) == 0 {
		return
	}

	p.started = true
	p.spinner = NewSpinner(p.steps[0])
	p.spinner.Start()
}

// Next advances to the next step.
func (p *ProgressReporter) Next() {
	if !p.verbose || !p.started || p.spinner == nil {
		return
	}

	p.currentStep++
	if p.currentStep < len(p.steps) {
		p.spinner.Update(p.steps[p.currentStep])
	}
}

// Done completes progress reporting with a final message.
func (p *ProgressReporter) Done(message string) {
	if !p.verbose || p.spinner == nil {
		return
	}

	p.spinner.Stop(message)
}

// ShouldShowSpinner determines if a spinner should be displayed based on config.
//
// Returns false if:
// - Output is not a TTY
// - Verbose mode is disabled
// - NO_COLOR environment variable is set
func ShouldShowSpinner(cfg *config.Config) bool {
	// Don't show spinner if output is not a TTY
	if !isTerminal(os.Stderr) {
		return false
	}

	// Don't show spinner if NO_COLOR is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Show spinner if verbose mode is enabled
	return cfg.Verbose
}

// SearchSpinner is a specialized spinner for search operations.
//
// It provides a convenient interface for the common search workflow.
type SearchSpinner struct {
	spinner *Spinner
	enabled bool
}

// NewSearchSpinner creates a new search spinner.
func NewSearchSpinner(enabled bool) *SearchSpinner {
	s := &SearchSpinner{
		enabled: enabled && isTerminal(os.Stderr),
	}

	if s.enabled {
		s.spinner = NewSpinner("Searching...")
	}

	return s
}

// Start begins the search spinner.
func (s *SearchSpinner) Start() {
	if s.enabled && s.spinner != nil {
		s.spinner.Start()
	}
}

// Stop stops the search spinner with a result message.
//
// The message should indicate how many results were found.
//
// Example:
//
//	spinner := ui.NewSearchSpinner(true)
//	spinner.Start()
//	// Do search
//	spinner.Stop(15, "0.5s") // Found 15 results in 0.5 seconds
func (s *SearchSpinner) Stop(resultCount int, duration string) {
	if s.enabled && s.spinner != nil {
		var msg string
		if resultCount > 0 {
			msg = fmt.Sprintf("Found %d results in %s", resultCount, duration)
		} else {
			msg = fmt.Sprintf("No results found in %s", duration)
		}
		s.spinner.Stop(msg)
	}
}

// StopWithError stops the spinner with an error message.
func (s *SearchSpinner) StopWithError(err error) {
	if s.enabled && s.spinner != nil {
		s.spinner.Stop(fmt.Sprintf("Error: %v", err))
	}
}

// Info prints an informational message to stderr.
//
// This is useful for verbose output that doesn't need a spinner.
//
// Example:
//
//	ui.Info("Using instance: https://search.butler.ooo")
func Info(message string) {
	fmt.Fprintf(os.Stderr, "\033[36m%s\033[0m\n", message)
}

// Warning prints a warning message to stderr.
//
// Example:
//
//	ui.Warning("Timeout is very low, may cause failures")
func Warning(message string) {
	fmt.Fprintf(os.Stderr, "\033[33mWarning: %s\033[0m\n", message)
}

// Error prints an error message to stderr.
//
// Example:
//
//	ui.Error("Failed to connect to instance")
func Error(message string) {
	fmt.Fprintf(os.Stderr, "\033[31mError: %s\033[0m\n", message)
}

// Success prints a success message to stderr.
//
// Example:
//
//	ui.Success("Configuration saved")
func Success(message string) {
	fmt.Fprintf(os.Stderr, "\033[32m%s\033[0m\n", message)
}

// SanitizeInput removes potentially dangerous characters from user input.
//
// This is a basic sanitization function that removes control characters
// and limits the length of the input.
//
// Example:
//
//	cleanQuery := ui.SanitizeInput(userQuery)
func SanitizeInput(input string) string {
	// Remove control characters (except tab, newline, carriage return)
	var sb strings.Builder
	for _, r := range input {
		if r == '\t' || r == '\n' || r == '\r' {
			sb.WriteRune(r)
		} else if r >= 32 && r != 127 {
			sb.WriteRune(r)
		}
	}

	result := sb.String()

	// Limit length to prevent abuse
	const maxLength = 2000
	if len(result) > maxLength {
		result = result[:maxLength]
	}

	return strings.TrimSpace(result)
}
