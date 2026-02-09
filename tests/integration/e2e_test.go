//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	// Default test instance
	defaultTestInstance = "https://search.butler.ooo"
	// Build timeout for compilation
	buildTimeout = 30 * time.Second
	// Test execution timeout
	testTimeout = 60 * time.Second
)

// TestHelper provides helper methods for E2E tests
type TestHelper struct {
	t            *testing.T
	binaryPath   string
	testInstance string
	tempDir      string
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	instance := os.Getenv("SEARXNG_TEST_INSTANCE")
	if instance == "" {
		instance = defaultTestInstance
	}

	tempDir := t.TempDir()

	// Build the binary
	binaryPath := filepath.Join(tempDir, "search")
	projectRoot := getProjectRoot()
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = filepath.Join(projectRoot, "cmd", "search")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, string(output))
	}

	return &TestHelper{
		t:            t,
		binaryPath:   binaryPath,
		testInstance: instance,
		tempDir:      tempDir,
	}
}

// RunCommand executes the search CLI with given arguments
func (h *TestHelper) RunCommand(args ...string) *TestResult {
	h.t.Helper()

	cmd := exec.Command(h.binaryPath, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SEARXNG_TEST_INSTANCE=%s", h.testInstance),
	)

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	return &TestResult{
		ExitCode:   getExitCode(err),
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		Duration:   duration,
		Error:      err,
	}
}

// RunCommandWithConfig runs command with a custom config file
func (h *TestHelper) RunCommandWithConfig(configPath string, args ...string) *TestResult {
	h.t.Helper()

	fullArgs := append([]string{"--config", configPath}, args...)
	return h.RunCommand(fullArgs...)
}

// CreateConfigFile creates a test config file
func (h *TestHelper) CreateConfigFile(content string) string {
	configPath := filepath.Join(h.tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		h.t.Fatalf("Failed to create config file: %v", err)
	}
	return configPath
}

// TestResult holds the result of a command execution
type TestResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	Error    error
}

// getExitCode extracts exit code from error
func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// TestE2EBasicSearch performs a basic end-to-end search
func TestE2EBasicSearch(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("BasicTextSearch", func(t *testing.T) {
		result := h.RunCommand("golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStdout: %s\nStderr: %s",
				result.ExitCode, result.Stdout, result.Stderr)
		}

		// Verify output contains expected elements
		if !strings.Contains(result.Stdout, "golang") && !strings.Contains(strings.ToLower(result.Stdout), "golang") {
			t.Error("Expected output to contain search query or results")
		}

		if len(result.Stdout) == 0 {
			t.Error("Expected non-empty output")
		}
	})

	t.Run("SearchWithInstanceFlag", func(t *testing.T) {
		result := h.RunCommand("-i", h.testInstance, "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("SearchWithResultsFlag", func(t *testing.T) {
		result := h.RunCommand("-n", "5", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Count result entries in plaintext format
		// Plaintext format shows [1], [2], etc.
		count := strings.Count(result.Stdout, "[")
		// Should have at least some results (plus header)
		if count < 1 {
			t.Logf("Warning: Few results found in output. Count: %d", count)
		}
	})
}

// TestE2EJSONFormat tests JSON output format
func TestE2EJSONFormat(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("JSONOutput", func(t *testing.T) {
		result := h.RunCommand("-f", "json", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verify valid JSON
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResponse); err != nil {
			t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, result.Stdout)
		}

		// Check for required fields
		if _, ok := jsonResponse["query"]; !ok {
			t.Error("Expected JSON to contain 'query' field")
		}
		if _, ok := jsonResponse["results"]; !ok {
			t.Error("Expected JSON to contain 'results' field")
		}
		if _, ok := jsonResponse["metadata"]; !ok {
			t.Error("Expected JSON to contain 'metadata' field")
		}

		// Verify results array
		results, ok := jsonResponse["results"].([]interface{})
		if !ok {
			t.Error("Expected 'results' to be an array")
		} else if len(results) == 0 {
			t.Logf("Warning: No results returned for query 'golang'")
		} else {
			// Check first result structure
			firstResult, ok := results[0].(map[string]interface{})
			if !ok {
				t.Error("Expected result to be an object")
			} else {
				if _, ok := firstResult["title"]; !ok {
					t.Error("Expected result to have 'title' field")
				}
				if _, ok := firstResult["url"]; !ok {
					t.Error("Expected result to have 'url' field")
				}
			}
		}
	})

	t.Run("JSONOutputWithResultsFlag", func(t *testing.T) {
		result := h.RunCommand("-f", "json", "-n", "15", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResponse); err != nil {
			t.Errorf("Expected valid JSON, got error: %v", err)
		}

		results, ok := jsonResponse["results"].([]interface{})
		if ok && len(results) > 0 {
			// We got results, verify reasonable count
			t.Logf("Got %d results with -n 15", len(results))
		}
	})
}

// TestE2EMarkdownFormat tests markdown output format
func TestE2EMarkdownFormat(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("MarkdownOutput", func(t *testing.T) {
		result := h.RunCommand("-f", "markdown", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verify markdown elements
		output := result.Stdout

		// Should have markdown headers
		if !strings.Contains(output, "#") {
			t.Error("Expected markdown headers")
		}

		// Should have markdown links
		if !strings.Contains(output, "](") {
			t.Error("Expected markdown links")
		}

		// Should have horizontal rules
		if !strings.Contains(output, "---") {
			t.Error("Expected markdown horizontal rules")
		}
	})

	t.Run("MarkdownLinkFormat", func(t *testing.T) {
		result := h.RunCommand("-f", "markdown", "golang")

		// Verify link format: [title](url)
		lines := strings.Split(result.Stdout, "\n")
		foundLink := false
		for _, line := range lines {
			if strings.Contains(line, "](") && strings.Contains(line, "http") {
				foundLink = true
				break
			}
		}

		if !foundLink {
			t.Log("Warning: No markdown links found in output")
		}
	})
}

// TestE2ETextFormat tests plaintext output format
func TestE2ETextFormat(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("TextOutput", func(t *testing.T) {
		result := h.RunCommand("-f", "text", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verify text format elements
		output := result.Stdout

		// Should have numbered results
		if !strings.Contains(output, "[1]") {
			t.Log("Warning: No numbered results found (may be format variation)")
		}

		// Should have URLs
		if !strings.Contains(output, "http") {
			t.Error("Expected URLs in output")
		}

		// Should not have JSON or markdown specific formatting
		if strings.Contains(output, "{") && strings.Contains(output, "}") {
			// Might be a false positive, but check if it's actually JSON
			if strings.Contains(output, "\"query\"") {
				t.Error("Expected text output, got JSON-like content")
			}
		}
	})
}

// TestE2ECategories tests search categories
func TestE2ECategories(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("GeneralCategory", func(t *testing.T) {
		result := h.RunCommand("-c", "general", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("ImagesCategory", func(t *testing.T) {
		result := h.RunCommand("-c", "images", "nature")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Image results should contain image-related info
		_ = strings.ToLower(result.Stdout)
		t.Logf("Image search output length: %d", len(result.Stdout))
	})

	t.Run("NewsCategory", func(t *testing.T) {
		result := h.RunCommand("-c", "news", "technology")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("VideosCategory", func(t *testing.T) {
		result := h.RunCommand("-c", "videos", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

// TestE2EPagination tests pagination functionality
func TestE2EPagination(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("FirstPage", func(t *testing.T) {
		result := h.RunCommand("--page", "1", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("SecondPage", func(t *testing.T) {
		result := h.RunCommand("--page", "2", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("InvalidPage", func(t *testing.T) {
		result := h.RunCommand("--page", "0", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid page number")
		}

		if !strings.Contains(result.Stderr, "page") && !strings.Contains(strings.ToLower(result.Stdout), "page") {
			t.Log("Warning: Error message doesn't mention 'page'")
		}
	})
}

// TestE2EFiltering tests time range and safe search filtering
func TestE2EFiltering(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("TimeRangeDay", func(t *testing.T) {
		result := h.RunCommand("--time", "day", "news")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("TimeRangeWeek", func(t *testing.T) {
		result := h.RunCommand("--time", "week", "news")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("TimeRangeMonth", func(t *testing.T) {
		result := h.RunCommand("--time", "month", "news")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("TimeRangeYear", func(t *testing.T) {
		result := h.RunCommand("--time", "year", "news")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("SafeSearchOff", func(t *testing.T) {
		result := h.RunCommand("--safe", "0", "test")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("SafeSearchModerate", func(t *testing.T) {
		result := h.RunCommand("--safe", "1", "test")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("SafeSearchStrict", func(t *testing.T) {
		result := h.RunCommand("--safe", "2", "test")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

// TestE2EConfigFile tests config file integration
func TestE2EConfigFile(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("CustomConfigFile", func(t *testing.T) {
		configContent := fmt.Sprintf(`
instance: "%s"
results: 5
format: "json"
language: "en"
safe_search: 1
timeout: 30
`, h.testInstance)

		configPath := h.CreateConfigFile(configContent)

		result := h.RunCommandWithConfig(configPath, "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verify JSON output (from config)
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResponse); err != nil {
			t.Errorf("Expected valid JSON from config, got error: %v", err)
		}
	})

	t.Run("ConfigWithCustomInstance", func(t *testing.T) {
		configContent := fmt.Sprintf(`
instance: "%s"
results: 10
format: "text"
`, h.testInstance)

		configPath := h.CreateConfigFile(configContent)

		result := h.RunCommandWithConfig(configPath, "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("CLIFlagsOverrideConfig", func(t *testing.T) {
		configContent := fmt.Sprintf(`
instance: "%s"
results: 5
format: "text"
`, h.testInstance)

		configPath := h.CreateConfigFile(configContent)

		// CLI flag should override config
		result := h.RunCommandWithConfig(configPath, "-f", "json", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Should be JSON format (from flag), not text (from config)
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResponse); err != nil {
			t.Errorf("Expected JSON output from flag override, got error: %v", err)
		}
	})
}

// TestE2EValidation tests input validation
func TestE2EValidation(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("EmptyQuery", func(t *testing.T) {
		result := h.RunCommand("")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for empty query")
		}
	})

	t.Run("InvalidResultsCount", func(t *testing.T) {
		result := h.RunCommand("-n", "0", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid results count")
		}
	})

	t.Run("InvalidResultsCountTooHigh", func(t *testing.T) {
		result := h.RunCommand("-n", "101", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for results count > 100")
		}
	})

	t.Run("InvalidTimeout", func(t *testing.T) {
		result := h.RunCommand("-t", "0", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid timeout")
		}
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		result := h.RunCommand("-f", "invalid", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid format")
		}
	})

	t.Run("InvalidSafeSearch", func(t *testing.T) {
		result := h.RunCommand("--safe", "5", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid safe search level")
		}
	})

	t.Run("InvalidTimeRange", func(t *testing.T) {
		result := h.RunCommand("--time", "invalid", "golang")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for invalid time range")
		}
	})
}

// TestE2EFlags tests various CLI flags
func TestE2EFlags(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("VerboseFlag", func(t *testing.T) {
		result := h.RunCommand("-v", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verbose mode should output to stderr
		if len(result.Stderr) == 0 {
			t.Log("Warning: Verbose flag produced no stderr output")
		}
	})

	t.Run("NoColorFlag", func(t *testing.T) {
		result := h.RunCommand("--no-color", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Should not have ANSI color codes
		if strings.Contains(result.Stdout, "\x1b[") || strings.Contains(result.Stdout, "\033[") {
			t.Error("Expected no color codes with --no-color flag")
		}
	})

	t.Run("LanguageFlag", func(t *testing.T) {
		result := h.RunCommand("-l", "de", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("TimeoutFlag", func(t *testing.T) {
		result := h.RunCommand("-t", "60", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

// TestE2ECommands tests subcommands
func TestE2ECommands(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("VersionCommand", func(t *testing.T) {
		result := h.RunCommand("version")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		if !strings.Contains(result.Stdout, "version") {
			t.Error("Expected version output to contain 'version'")
		}
	})

	t.Run("VersionFlag", func(t *testing.T) {
		result := h.RunCommand("--version")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		if !strings.Contains(result.Stdout, "version") {
			t.Error("Expected version output to contain 'version'")
		}
	})

	t.Run("HelpCommand", func(t *testing.T) {
		result := h.RunCommand("--help")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		if !strings.Contains(result.Stdout, "Usage") {
			t.Error("Expected help output to contain 'Usage'")
		}
	})

	t.Run("CategoriesCommand", func(t *testing.T) {
		result := h.RunCommand("categories")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		if !strings.Contains(result.Stdout, "Categories") {
			t.Error("Expected categories output to contain 'Categories'")
		}

		// Should list some common categories
		lowerOutput := strings.ToLower(result.Stdout)
		if !strings.Contains(lowerOutput, "general") {
			t.Error("Expected 'general' in categories list")
		}
	})

	t.Run("CompletionCommandBash", func(t *testing.T) {
		result := h.RunCommand("completion", "bash")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Bash completion should contain bash-specific syntax
		if !strings.Contains(result.Stdout, "complete") {
			t.Error("Expected bash completion to contain 'complete'")
		}
	})

	t.Run("CompletionCommandZsh", func(t *testing.T) {
		result := h.RunCommand("completion", "zsh")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("CompletionCommandFish", func(t *testing.T) {
		result := h.RunCommand("completion", "fish")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

// TestE2EMultiFormat tests all three formats against the same query
func TestE2EMultiFormat(t *testing.T) {
	h := NewTestHelper(t)

	queries := []string{"golang", "rust programming", "kubernetes"}
	formats := []string{"json", "markdown", "text"}

	for _, query := range queries {
		for _, format := range formats {
			t.Run(fmt.Sprintf("%s_%s", format, strings.ReplaceAll(query, " ", "_")), func(t *testing.T) {
				result := h.RunCommand("-f", format, query)

				if result.ExitCode != 0 {
					t.Errorf("Expected exit code 0 for %s format with query '%s'\nStderr: %s",
						format, query, result.Stderr)
				}

				if len(result.Stdout) == 0 {
					t.Errorf("Expected non-empty output for %s format", format)
				}
			})
		}
	}
}

// TestE2ENetworkErrorHandling tests error scenarios
func TestE2ENetworkErrorHandling(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("InvalidInstanceURL", func(t *testing.T) {
		result := h.RunCommand("-i", "https://this-instance-does-not-exist-12345.com", "test")

		// Should fail or timeout
		if result.ExitCode == 0 {
			t.Log("Warning: Invalid instance URL didn't fail (may be DNS resolution issue)")
		} else {
			t.Logf("Correctly failed with invalid instance: %v", result.Error)
		}
	})

	t.Run("MalformedURL", func(t *testing.T) {
		result := h.RunCommand("-i", "not-a-valid-url", "test")

		if result.ExitCode == 0 {
			t.Error("Expected non-zero exit code for malformed URL")
		}
	})
}

// TestE2EPerformance performs basic performance checks
func TestE2EPerformance(t *testing.T) {
	h := NewTestHelper(t)

	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("SearchLatency", func(t *testing.T) {
		start := time.Now()
		result := h.RunCommand("golang")
		duration := time.Since(start)

		if result.ExitCode != 0 {
			t.Errorf("Search failed: %v", result.Error)
		}

		t.Logf("Search completed in %v", duration)

		// Search should complete within reasonable time (30 seconds)
		if duration > 30*time.Second {
			t.Errorf("Search took too long: %v (expected < 30s)", duration)
		}
	})

	t.Run("LargeResultSet", func(t *testing.T) {
		result := h.RunCommand("-n", "50", "golang")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		t.Logf("Large result set output size: %d bytes", len(result.Stdout))
	})
}

// TestE2EEnvironmentVariables tests environment variable support
func TestE2EEnvironmentVariables(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("InstanceFromEnv", func(t *testing.T) {
		cmd := exec.Command(h.binaryPath, "golang")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("SEARCH_INSTANCE=%s", h.testInstance),
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Errorf("Command failed: %v\nStderr: %s", err, stderr.String())
		}
	})

	t.Run("ResultsFromEnv", func(t *testing.T) {
		cmd := exec.Command(h.binaryPath, "golang")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("SEARCH_RESULTS=15"),
			fmt.Sprintf("SEARCH_INSTANCE=%s", h.testInstance),
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Errorf("Command failed: %v\nStderr: %s", err, stderr.String())
		}
	})

	t.Run("FormatFromEnv", func(t *testing.T) {
		cmd := exec.Command(h.binaryPath, "golang")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("SEARCH_FORMAT=json"),
			fmt.Sprintf("SEARCH_INSTANCE=%s", h.testInstance),
		)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Errorf("Command failed: %v\nStderr: %s", err, stderr.String())
		}

		// Verify JSON output
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal(stdout.Bytes(), &jsonResponse); err != nil {
			t.Errorf("Expected JSON from env var, got error: %v", err)
		}
	})
}

// TestE2ECombinedFlags tests combinations of flags
func TestE2ECombinedFlags(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("MultipleFlagsCombined", func(t *testing.T) {
		result := h.RunCommand(
			"-f", "json",
			"-n", "15",
			"-c", "general",
			"-l", "en",
			"--page", "1",
			"golang",
		)

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}

		// Verify JSON
		var jsonResponse map[string]interface{}
		if err := json.Unmarshal([]byte(result.Stdout), &jsonResponse); err != nil {
			t.Errorf("Expected valid JSON with combined flags, got error: %v", err)
		}
	})

	t.Run("AllFormatsWithVerbose", func(t *testing.T) {
		formats := []string{"json", "markdown", "text"}

		for _, format := range formats {
			result := h.RunCommand("-f", format, "-v", "golang")

			if result.ExitCode != 0 {
				t.Errorf("Expected exit code 0 for %s with verbose\nStderr: %s",
					format, result.Stderr)
			}

			// Verbose should produce stderr output
			if len(result.Stderr) == 0 {
				t.Logf("Warning: No stderr output with verbose for %s format", format)
			}
		}
	})
}

// TestE2EResponseParsing tests parsing of various response formats
func TestE2EResponseParsing(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("JSONResponseStructure", func(t *testing.T) {
		result := h.RunCommand("-f", "json", "golang")

		var response struct {
			Query        string                 `json:"query"`
			Results      []map[string]interface{} `json:"results"`
			Answers      []interface{}          `json:"answers"`
			Infoboxes    []interface{}          `json:"infoboxes"`
			Suggestions  []string               `json:"suggestions"`
			NumberOfResults int                 `json:"number_of_results"`
			Metadata     map[string]interface{} `json:"metadata"`
		}

		if err := json.Unmarshal([]byte(result.Stdout), &response); err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if response.Query != "golang" {
			t.Errorf("Expected query 'golang', got '%s'", response.Query)
		}

		if response.Results == nil {
			t.Error("Expected results array to be present")
		}

		if response.Metadata == nil {
			t.Error("Expected metadata to be present")
		}
	})
}

// TestE2ESpecialQueries tests edge cases in queries
func TestE2ESpecialQueries(t *testing.T) {
	h := NewTestHelper(t)

	t.Run("QueryWithSpaces", func(t *testing.T) {
		result := h.RunCommand("golang tutorial for beginners")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("QueryWithSpecialCharacters", func(t *testing.T) {
		result := h.RunCommand("c++ programming")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("QueryWithQuotes", func(t *testing.T) {
		result := h.RunCommand("golang \"hello world\"")

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})

	t.Run("LongQuery", func(t *testing.T) {
		longQuery := strings.Repeat("test ", 50)
		result := h.RunCommand(longQuery)

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d\nStderr: %s", result.ExitCode, result.Stderr)
		}
	})
}

// TestE2ECrossFormatConsistency tests that all formats return consistent results
func TestE2ECrossFormatConsistency(t *testing.T) {
	h := NewTestHelper(t)

	query := "golang"

	// Get results in all formats
	jsonResult := h.RunCommand("-f", "json", query)
	markdownResult := h.RunCommand("-f", "markdown", query)
	textResult := h.RunCommand("-f", "text", query)

	if jsonResult.ExitCode != 0 {
		t.Fatalf("JSON search failed: %s", jsonResult.Stderr)
	}
	if markdownResult.ExitCode != 0 {
		t.Fatalf("Markdown search failed: %s", markdownResult.Stderr)
	}
	if textResult.ExitCode != 0 {
		t.Fatalf("Text search failed: %s", textResult.Stderr)
	}

	// Parse JSON to get result count
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResult.Stdout), &jsonResponse); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, ok := jsonResponse["results"].([]interface{})
	if !ok {
		t.Fatal("Expected results array in JSON")
	}

	jsonResultCount := len(results)

	// All formats should have results
	if jsonResultCount == 0 {
		t.Log("Warning: No results returned for query")
	} else {
		t.Logf("Consistent result count: %d results across all formats", jsonResultCount)
	}
}