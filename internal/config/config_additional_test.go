package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	t.Run("write new config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		cfg := Config{
			Instance:   "https://example.com",
			Results:    15,
			Format:     "json",
			Timeout:    60,
			Language:   "de",
			SafeSearch: 1,
		}

		err := WriteConfig(configPath, cfg)
		if err != nil {
			t.Fatalf("WriteConfig() error = %v", err)
		}

		// Verify file exists
		info, err := os.Stat(configPath)
		if err != nil {
			t.Fatalf("Failed to stat config file: %v", err)
		}

		// Check permissions (should be 0600)
		if info.Mode().Perm() != 0600 {
			t.Errorf("Expected file permissions 0600, got %v", info.Mode().Perm())
		}

		// Read and verify content
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "https://example.com") {
			t.Error("Config content missing instance")
		}
		if !strings.Contains(contentStr, "results: 15") {
			t.Error("Config content missing results")
		}
	})

	t.Run("write to nested directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedPath := filepath.Join(tmpDir, "nested", "dir", "config.yaml")

		cfg := Config{
			Instance: "https://example.com",
			Results:  10,
		}

		err := WriteConfig(nestedPath, cfg)
		if err != nil {
			t.Fatalf("WriteConfig() to nested dir error = %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Error("Config file was not created in nested directory")
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("save to default location", func(t *testing.T) {
		// Create a temp home directory
		tmpHome := t.TempDir()

		// Set HOME environment variable for this test
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpHome)
		defer func() { os.Setenv("HOME", oldHome) }()

		cfg := &Config{
			Instance:   "https://saved.example.com",
			Results:    30,
			Format:     "markdown",
			Timeout:    90,
			Language:   "fr",
			SafeSearch: 2,
		}

		err := cfg.Save()
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// Verify file exists at expected location
		configPath := filepath.Join(tmpHome, defaultConfigDir, configFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("Config file not found at %s", configPath)
		}

		// Load and verify
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("Failed to read saved config: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "https://saved.example.com") {
			t.Error("Saved config missing instance")
		}
	})
}

func TestSaveTo(t *testing.T) {
	t.Run("save to specific path", func(t *testing.T) {
		tmpDir := t.TempDir()
		customPath := filepath.Join(tmpDir, "custom", "config.yaml")

		cfg := &Config{
			Instance: "https://custom.example.com",
			Results:  25,
			Format:   "json",
		}

		err := cfg.SaveTo(customPath)
		if err != nil {
			t.Fatalf("SaveTo() error = %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(customPath); os.IsNotExist(err) {
			t.Error("Config file not created at custom path")
		}

		// Load and verify
		content, err := os.ReadFile(customPath)
		if err != nil {
			t.Fatalf("Failed to read saved config: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "https://custom.example.com") {
			t.Error("Saved config missing instance")
		}
		if !strings.Contains(contentStr, "results: 25") {
			t.Error("Saved config missing results")
		}
	})

	t.Run("save to nested directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedPath := filepath.Join(tmpDir, "a", "b", "c", "config.yaml")

		cfg := &Config{
			Instance: "https://nested.example.com",
		}

		err := cfg.SaveTo(nestedPath)
		if err != nil {
			t.Fatalf("SaveTo() nested error = %v", err)
		}

		if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
			t.Error("Config file not created in nested directory")
		}
	})
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Run("load existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `instance: "https://file.example.com"
results: 40
format: "text"
timeout: 50
language: "it"
safe_search: 1
api_key: "file-api-key"
`

		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config: %v", err)
		}

		cfg, err := LoadConfigFromFile(configPath)
		if err != nil {
			t.Fatalf("LoadConfigFromFile() error = %v", err)
		}

		if cfg.Instance != "https://file.example.com" {
			t.Errorf("Expected instance 'https://file.example.com', got '%s'", cfg.Instance)
		}
		if cfg.Results != 40 {
			t.Errorf("Expected results 40, got %d", cfg.Results)
		}
		if cfg.APIKey != "file-api-key" {
			t.Errorf("Expected API key 'file-api-key', got '%s'", cfg.APIKey)
		}
	})

	t.Run("load non-existent file returns defaults", func(t *testing.T) {
		nonExistentPath := "/tmp/nonexistent_config_12345.yaml"

		cfg, err := LoadConfigFromFile(nonExistentPath)
		if err != nil {
			t.Fatalf("LoadConfigFromFile() should not error for non-existent file, got %v", err)
		}

		// Should have defaults applied
		if cfg.Instance != "https://search.butler.ooo" {
			t.Errorf("Expected default instance, got '%s'", cfg.Instance)
		}
		if cfg.Results != 10 {
			t.Errorf("Expected default results, got %d", cfg.Results)
		}
	})
}

func TestLoadConfigIntegration(t *testing.T) {
	t.Run("load with CLI override", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `instance: "https://config.example.com"
results: 20
format: "json"
timeout: 30
language: "en"
safe_search: 1
`

		err := os.WriteFile(configPath, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config: %v", err)
		}

		cliCfg := &CliConfig{
			ConfigPath: configPath,
			Results:    50,  // CLI override
			Format:     "text", // CLI override
		}

		cfg, err := LoadConfig(cliCfg)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Instance from config file
		if cfg.Instance != "https://config.example.com" {
			t.Errorf("Expected instance from config, got '%s'", cfg.Instance)
		}

		// Results from CLI override
		if cfg.Results != 50 {
			t.Errorf("Expected CLI override for results, got %d", cfg.Results)
		}

		// Format from CLI override
		if cfg.Format != "text" {
			t.Errorf("Expected CLI override for format, got '%s'", cfg.Format)
		}

		// Timeout from config file (no CLI override)
		if cfg.Timeout != 30 {
			t.Errorf("Expected timeout from config, got %d", cfg.Timeout)
		}
	})

	t.Run("load with environment variable override", func(t *testing.T) {
		// Set environment variables
		oldInstance := os.Getenv("SEARCH_INSTANCE")
		oldResults := os.Getenv("SEARCH_RESULTS")
		oldFormat := os.Getenv("SEARCH_FORMAT")

		defer func() {
			if oldInstance != "" {
				os.Setenv("SEARCH_INSTANCE", oldInstance)
			} else {
				os.Unsetenv("SEARCH_INSTANCE")
			}
			if oldResults != "" {
				os.Setenv("SEARCH_RESULTS", oldResults)
			} else {
				os.Unsetenv("SEARCH_RESULTS")
			}
			if oldFormat != "" {
				os.Setenv("SEARCH_FORMAT", oldFormat)
			} else {
				os.Unsetenv("SEARCH_FORMAT")
			}
		}()

		os.Setenv("SEARCH_INSTANCE", "https://env.example.com")
		os.Setenv("SEARCH_RESULTS", "75")
		os.Setenv("SEARCH_FORMAT", "markdown")

		cliCfg := &CliConfig{}

		cfg, err := LoadConfig(cliCfg)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Environment variables should override defaults
		if cfg.Instance != "https://env.example.com" {
			t.Errorf("Expected instance from env var, got '%s'", cfg.Instance)
		}
		if cfg.Results != 75 {
			t.Errorf("Expected results from env var, got %d", cfg.Results)
		}
		if cfg.Format != "markdown" {
			t.Errorf("Expected format from env var, got '%s'", cfg.Format)
		}
	})

	t.Run("CLI flags override environment variables", func(t *testing.T) {
		oldInstance := os.Getenv("SEARCH_INSTANCE")
		defer func() {
			if oldInstance != "" {
				os.Setenv("SEARCH_INSTANCE", oldInstance)
			} else {
				os.Unsetenv("SEARCH_INSTANCE")
			}
		}()

		os.Setenv("SEARCH_INSTANCE", "https://env.example.com")

		cliCfg := &CliConfig{
			Instance: "https://cli.example.com", // CLI override
		}

		cfg, err := LoadConfig(cliCfg)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// CLI should override env var
		if cfg.Instance != "https://cli.example.com" {
			t.Errorf("Expected CLI to override env var, got '%s'", cfg.Instance)
		}
	})
}

func TestApplyEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		cfg      *Config
		expected *Config
	}{
		{
			name: "apply all env vars",
			envVars: map[string]string{
				"SEARCH_INSTANCE":    "https://env.example.com",
				"SEARCH_RESULTS":     "50",
				"SEARCH_FORMAT":      "json",
				"SEARCH_TIMEOUT":     "60",
				"SEARCH_LANGUAGE":    "es",
				"SEARCH_SAFE_SEARCH": "0",
			},
			cfg:      &Config{},
			expected: &Config{Instance: "https://env.example.com", Results: 50, Format: "json", Timeout: 60, Language: "es", SafeSearch: 0},
		},
		{
			name: "partial env vars",
			envVars: map[string]string{
				"SEARCH_INSTANCE": "https://partial.example.com",
			},
			cfg:      &Config{},
			expected: &Config{Instance: "https://partial.example.com"},
		},
		{
			name:     "no env vars",
			envVars:  map[string]string{},
			cfg:      &Config{},
			expected: &Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			envKeys := []string{"SEARCH_INSTANCE", "SEARCH_RESULTS", "SEARCH_FORMAT", "SEARCH_TIMEOUT", "SEARCH_LANGUAGE", "SEARCH_SAFE_SEARCH"}
			oldValues := make(map[string]string)
			for _, key := range envKeys {
				if val := os.Getenv(key); val != "" {
					oldValues[key] = val
				}
				os.Unsetenv(key)
			}
			defer func() {
				for key, val := range oldValues {
					os.Setenv(key, val)
				}
				for _, key := range envKeys {
					if _, ok := oldValues[key]; !ok {
						os.Unsetenv(key)
					}
				}
			}()

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Apply
			tt.cfg.applyEnvironmentVariables()

			// Verify
			if tt.cfg.Instance != tt.expected.Instance {
				t.Errorf("Instance = %v, want %v", tt.cfg.Instance, tt.expected.Instance)
			}
			if tt.cfg.Results != tt.expected.Results {
				t.Errorf("Results = %v, want %v", tt.cfg.Results, tt.expected.Results)
			}
			if tt.cfg.Format != tt.expected.Format {
				t.Errorf("Format = %v, want %v", tt.cfg.Format, tt.expected.Format)
			}
			if tt.cfg.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout = %v, want %v", tt.cfg.Timeout, tt.expected.Timeout)
			}
			if tt.cfg.Language != tt.expected.Language {
				t.Errorf("Language = %v, want %v", tt.cfg.Language, tt.expected.Language)
			}
			if tt.cfg.SafeSearch != tt.expected.SafeSearch {
				t.Errorf("SafeSearch = %v, want %v", tt.cfg.SafeSearch, tt.expected.SafeSearch)
			}
		})
	}
}

func TestParseIntEnv(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"valid integer", "42", 42},
		{"zero", "0", 0},
		{"negative", "-5", -5},
		{"with extra text", "42extra", 42},
		{"just text", "abc", 0},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseIntEnv(tt.input)
			if result != tt.expected {
				t.Errorf("parseIntEnv(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadWithErrorHandling(t *testing.T) {
	t.Run("load creates default config when missing", func(t *testing.T) {
		tmpHome := t.TempDir()

		// Set HOME environment variable for this test
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpHome)
		defer func() { os.Setenv("HOME", oldHome) }()

		cfg := &Config{}
		err := cfg.Load()
		if err != nil {
			t.Fatalf("Load() should create default config, got error: %v", err)
		}

		// Should have defaults
		if cfg.Instance == "" {
			t.Error("Expected default instance")
		}

		// Config file should be created
		configPath := filepath.Join(tmpHome, defaultConfigDir, configFileName)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config file should be created")
		}
	})
}