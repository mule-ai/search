//go:build benchmark
// +build benchmark

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"gopkg.in/yaml.v3"
)

// BenchmarkConfigLoading benchmarks loading config from file
func BenchmarkConfigLoading(b *testing.B) {
	// Create a temporary config file
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
instance: "https://search.butler.ooo"
results: 10
format: "text"
language: "en"
safe_search: 1
timeout: 30
categories:
  - general
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to create config file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadConfigFromFile(configPath)
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

// BenchmarkConfigLoadingWithDefaults benchmarks loading config and applying defaults
func BenchmarkConfigLoadingWithDefaults(b *testing.B) {
	// Create a minimal config file
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
instance: "https://search.butler.ooo"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to create config file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config, err := LoadConfigFromFile(configPath)
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
		// Apply defaults
		if config.Results == 0 {
			config.Results = 10
		}
		if config.Format == "" {
			config.Format = "text"
		}
		if config.Language == "" {
			config.Language = "en"
		}
		if config.SafeSearch == 0 {
			config.SafeSearch = 1
		}
		if config.Timeout == 0 {
			config.Timeout = 30
		}
	}
}

// BenchmarkConfigValidation benchmarks config validation
func BenchmarkConfigValidation(b *testing.B) {
	config := &Config{
		Instance:   "https://search.butler.ooo",
		Results:    10,
		Format:     "text",
		Language:   "en",
		SafeSearch: 1,
		Timeout:    30,
		Categories: []string{"general"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := config.Validate(); err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}

// BenchmarkConfigMerge benchmarks merging config with overrides
func BenchmarkConfigMerge(b *testing.B) {
	base := &Config{
		Instance:   "https://search.butler.ooo",
		Results:    10,
		Format:     "text",
		Language:   "en",
		SafeSearch: 1,
		Timeout:    30,
		Categories: []string{"general"},
	}

	overrides := map[string]interface{}{
		"results": 20,
		"format":  "json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := *base
		if v, ok := overrides["results"].(int); ok {
			result.Results = v
		}
		if v, ok := overrides["format"].(string); ok {
			result.Format = v
		}
	}
}

// BenchmarkDefaultConfig benchmarks creating default config
func BenchmarkDefaultConfig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := DefaultConfig()
		if config.Instance == "" {
			b.Fatal("Default config should have instance")
		}
	}
}

// BenchmarkConfigYAMLParsing benchmarks YAML parsing for config
func BenchmarkConfigYAMLParsing(b *testing.B) {
	configContent := []byte(`
instance: "https://search.butler.ooo"
results: 10
format: "text"
language: "en"
safe_search: 1
timeout: 30
categories:
  - general
  - images
  - videos
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config Config
		if err := yaml.Unmarshal(configContent, &config); err != nil {
			b.Fatalf("Failed to parse YAML: %v", err)
		}
	}
}

// BenchmarkConfigComplexYAML benchmarks parsing complex config with many options
func BenchmarkConfigComplexYAML(b *testing.B) {
	configContent := []byte(`
instance: "https://search.butler.ooo"
results: 50
format: "json"
language: "en"
safe_search: 2
timeout: 60
categories:
  - general
  - images
  - videos
  - news
  - music
  - it
  - science
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config Config
		if err := yaml.Unmarshal(configContent, &config); err != nil {
			b.Fatalf("Failed to parse YAML: %v", err)
		}
	}
}

// BenchmarkConfigFileExists benchmarks checking if config file exists
func BenchmarkConfigFileExists(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `instance: "https://search.butler.ooo"`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to create config file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := os.Stat(configPath)
		if err != nil {
			b.Fatalf("Failed to stat config: %v", err)
		}
	}
}

// BenchmarkGetConfigPath benchmarks getting config path with home directory lookup
func BenchmarkGetConfigPath(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		homeDir, _ := os.UserHomeDir()
		_ = filepath.Join(homeDir, ".search", "config.yaml")
	}
}

// BenchmarkConfigDirCreation benchmarks creating config directory
func BenchmarkConfigDirCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempDir := b.TempDir()
		configDir := filepath.Join(tempDir, ".search")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			b.Fatalf("Failed to create config dir: %v", err)
		}
	}
}

// BenchmarkLoadFromEnv benchmarks loading config from environment variables
func BenchmarkLoadFromEnv(b *testing.B) {
	// Set environment variables
	os.Setenv("SEARCH_INSTANCE", "https://search.butler.ooo")
	os.Setenv("SEARCH_RESULTS", "20")
	os.Setenv("SEARCH_FORMAT", "json")
	defer func() {
		os.Unsetenv("SEARCH_INSTANCE")
		os.Unsetenv("SEARCH_RESULTS")
		os.Unsetenv("SEARCH_FORMAT")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := DefaultConfig()
		if v := os.Getenv("SEARCH_INSTANCE"); v != "" {
			config.Instance = v
		}
		if v := os.Getenv("SEARCH_RESULTS"); v != "" {
			if results, err := strconv.Atoi(v); err == nil {
				config.Results = results
			}
		}
		if v := os.Getenv("SEARCH_FORMAT"); v != "" {
			config.Format = v
		}
	}
}
