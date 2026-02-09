// Package profiling provides profiling utilities for the search CLI.
package profiling

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// Profiler manages CPU and memory profiling.
type Profiler struct {
	cpuFile   *os.File
	memFile   *os.File
	traceFile *os.File
	enabled   bool
}

// Options configures the profiler.
type Options struct {
	CPUProfile  string // Path to write CPU profile
	MemProfile  string // Path to write memory profile
	TraceFile   string // Path to write execution trace
	Enabled     bool   // Whether profiling is enabled
}

// New creates a new Profiler with the given options.
func New(opts Options) *Profiler {
	if !opts.Enabled {
		return &Profiler{enabled: false}
	}

	p := &Profiler{enabled: true}

	var err error

	if opts.CPUProfile != "" {
		p.cpuFile, err = os.Create(opts.CPUProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create CPU profile: %v\n", err)
		} else {
			if err := pprof.StartCPUProfile(p.cpuFile); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not start CPU profile: %v\n", err)
				p.cpuFile.Close()
				p.cpuFile = nil
			}
		}
	}

	if opts.TraceFile != "" {
		p.traceFile, err = os.Create(opts.TraceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create trace file: %v\n", err)
		} else {
			if err := trace.Start(p.traceFile); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not start trace: %v\n", err)
				p.traceFile.Close()
				p.traceFile = nil
			}
		}
	}

	return p
}

// Stop stops profiling and writes profiles to disk.
func (p *Profiler) Stop() {
	if !p.enabled {
		return
	}

	if p.cpuFile != nil {
		pprof.StopCPUProfile()
		if err := p.cpuFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not close CPU profile: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "CPU profile written to %s\n", p.cpuFile.Name())
	}

	if p.traceFile != nil {
		trace.Stop()
		if err := p.traceFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not close trace file: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "Trace written to %s\n", p.traceFile.Name())
	}

	if p.memFile != nil {
		if err := pprof.WriteHeapProfile(p.memFile); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not write memory profile: %v\n", err)
		}
		if err := p.memFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not close memory profile: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "Memory profile written to %s\n", p.memFile.Name())
	}
}

// WriteMemProfile writes a heap profile to the specified file.
func (p *Profiler) WriteMemProfile(filename string) error {
	if !p.enabled {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Memory profile written to %s\n", filename)
	return nil
}

// Timed tracks the duration of an operation.
type Timed struct {
	name      string
	startTime time.Time
}

// StartTimed begins timing an operation with the given name.
func StartTimed(name string) *Timed {
	return &Timed{
		name:      name,
		startTime: time.Now(),
	}
}

// Stop stops the timer and prints the duration.
func (t *Timed) Stop() time.Duration {
	duration := time.Since(t.startTime)
	fmt.Fprintf(os.Stderr, "[profiling] %s took %v\n", t.name, duration)
	return duration
}

// StopSilent stops the timer without printing.
func (t *Timed) StopSilent() time.Duration {
	return time.Since(t.startTime)
}
