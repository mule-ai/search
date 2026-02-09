package profiling

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfilerCreation(t *testing.T) {
	// Test disabled profiler
	p := New(Options{Enabled: false})
	if p.enabled {
		t.Error("Expected profiler to be disabled")
	}

	// Test enabled profiler with no output files
	p = New(Options{Enabled: true})
	if !p.enabled {
		t.Error("Expected profiler to be enabled")
	}
	p.Stop()
}

func TestCPUProfile(t *testing.T) {
	tmpDir := t.TempDir()
	cpuPath := filepath.Join(tmpDir, "cpu.prof")

	p := New(Options{
		Enabled:    true,
		CPUProfile: cpuPath,
	})

	if !p.enabled {
		t.Fatal("Expected profiler to be enabled")
	}

	// Simulate some work
	for i := 0; i < 1000; i++ {
		_ = i * i
	}

	p.Stop()

	// Check that file was created
	if _, err := os.Stat(cpuPath); os.IsNotExist(err) {
		t.Errorf("CPU profile file was not created: %s", cpuPath)
	}
}

func TestMemProfile(t *testing.T) {
	tmpDir := t.TempDir()
	memPath := filepath.Join(tmpDir, "mem.prof")

	p := New(Options{Enabled: true})

	err := p.WriteMemProfile(memPath)
	if err != nil {
		t.Fatalf("WriteMemProfile failed: %v", err)
	}

	p.Stop()

	// Check that file was created
	if _, err := os.Stat(memPath); os.IsNotExist(err) {
		t.Errorf("Memory profile file was not created: %s", memPath)
	}
}

func TestTrace(t *testing.T) {
	tmpDir := t.TempDir()
	tracePath := filepath.Join(tmpDir, "trace.out")

	p := New(Options{
		Enabled:   true,
		TraceFile: tracePath,
	})

	if !p.enabled {
		t.Fatal("Expected profiler to be enabled")
	}

	// Simulate some work
	for i := 0; i < 100; i++ {
		_ = i
	}

	p.Stop()

	// Check that file was created
	if _, err := os.Stat(tracePath); os.IsNotExist(err) {
		t.Errorf("Trace file was not created: %s", tracePath)
	}
}

func TestTimed(t *testing.T) {
	tim := StartTimed("test_operation")
	
	// Simulate work
	sum := 0
	for i := 0; i < 100; i++ {
		sum += i
	}
	
	duration := tim.StopSilent()
	
	if duration <= 0 {
		t.Error("Expected positive duration")
	}
	
	if sum != 4950 { // Sum of 0-99
		t.Errorf("Unexpected sum: %d", sum)
	}
}
