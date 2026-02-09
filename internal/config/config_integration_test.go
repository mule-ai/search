//go:build integration
// +build integration

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIntegrationConfigFile tests config file operations with real filesystem.
func TestIntegrationConfigFile(t *testing.T) {
	// Use a temporary directory for tests
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".search")
	configFile := filepath.Join(configDir, "config.yaml")

	// Set custom config path for testing
	t.Setenv("SEARCH_CONFIG", configFile)

	t.Run("CreateConfigDirectory", func(t *testing.T) {
		// Ensure config directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Verify directory exists
		info, err := os.Stat(configDir)
		if err != nil {
			t.Fatalf("Config directory doesn't exist: %v", err)
		}

		if !info.IsDir() {
			t.Error("Config path is not a directory")
		}
	})

	t.Run("WriteAndReadConfig", func(t *testing.T) {
		// Create a test config
		testConfig := &Config{
			Instance:   "https://search.example.com",
			Results:    20,
			Format:     "json",
			Timeout:    60,
			Categories: []string{"general", "news"},
			Language:   "en",
			SafeSearch: 2,
		}

		// Write config to file
		if err := WriteConfig(configFile, *testConfig); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Fatal("Config file was not created")
		}

		// Read config back
		readConfig, err := LoadConfigFromFile(configFile)
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		// Verify values match
		if readConfig.Instance != testConfig.Instance {
			t.Errorf("Expected instance %s, got %s", testConfig.Instance, readConfig.Instance)
		}
		if readConfig.Results != testConfig.Results {
			t.Errorf("Expected results %d, got %d", testConfig.Results, readConfig.Results)
		}
		if readConfig.Format != testConfig.Format {
			t.Errorf("Expected format %s, got %s", testConfig.Format, readConfig.Format)
		}
		if readConfig.Timeout != testConfig.Timeout {
			t.Errorf("Expected timeout %d, got %d", testConfig.Timeout, readConfig.Timeout)
		}
		if len(readConfig.Categories) != len(testConfig.Categories) {
			t.Errorf("Expected %d categories, got %d", len(testConfig.Categories), len(readConfig.Categories))
		}
	})

	t.Run("LoadWithEnvironmentOverrides", func(t *testing.T) {
		// Set environment variables
		t.Setenv("SEARCH_INSTANCE", "https://env.example.com")
		t.Setenv("SEARCH_RESULTS", "50")
		t.Setenv("SEARCH_FORMAT", "markdown")
		t.Setenv("SEARCH_TIMEOUT", "90")

		// Load config
		cfg, err := LoadConfigFromFile(configFile)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify environment variables override config file
		if cfg.Instance != "https://env.example.com" {
			t.Errorf("Expected instance from env, got %s", cfg.Instance)
		}
		if cfg.Results != 50 {
			t.Errorf("Expected results from env, got %d", cfg.Results)
		}
		if cfg.Format != "markdown" {
			t.Errorf("Expected format from env, got %s", cfg.Format)
		}
		if cfg.Timeout != 90 {
			t.Errorf("Expected timeout from env, got %d", cfg.Timeout)
		}
	})

	t.Run("LoadMissingConfig", func(t *testing.T) {
		missingFile := filepath.Join(tempDir, "missing.yaml")

		// Loading missing config should return default config
		cfg, err := LoadConfigFromFile(missingFile)
		if err != nil {
			t.Logf("LoadConfig returned error for missing file: %v", err)
		}

		// If we got a config, check it has default values
		if cfg != nil {
			if cfg.Instance == "" {
				t.Error("Expected default instance to be set")
			}
			if cfg.Results == 0 {
				t.Error("Expected default results to be set")
			}
			if cfg.Format == "" {
				t.Error("Expected default format to be set")
			}
		}
	})

	t.Run("LoadInvalidYAML", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.yaml")

		// Write invalid YAML
		if err := os.WriteFile(invalidFile, []byte("invalid: yaml: content:"), 0644); err != nil {
			t.Fatalf("Failed to write invalid YAML: %v", err)
		}

		// Loading invalid YAML should return error
		_, err := LoadConfigFromFile(invalidFile)
		if err == nil {
			t.Error("Expected error when loading invalid YAML")
		}
	})

	t.Run("ConfigPermissions", func(t *testing.T) {
		secureConfig := filepath.Join(tempDir, "secure.yaml")

		// Write config
		testConfig := DefaultConfig()
		if err := WriteConfig(secureConfig, *testConfig); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Check file permissions
		info, err := os.Stat(secureConfig)
		if err != nil {
			t.Fatalf("Failed to stat config file: %v", err)
		}

		mode := info.Mode().Perm()
		// Config should not be world-writable
		if mode&0o022 != 0 {
			t.Errorf("Config file has insecure permissions: %o", mode)
		}
	})

	t.Run("PartialConfig", func(t *testing.T) {
		partialConfig := filepath.Join(tempDir, "partial.yaml")

		// Write partial config
		partialContent := `
instance: "https://partial.example.com"
results: 15
`
		if err := os.WriteFile(partialConfig, []byte(partialContent), 0644); err != nil {
			t.Fatalf("Failed to write partial config: %v", err)
		}

		// Load config - should merge with defaults
		cfg, err := LoadConfigFromFile(partialConfig)
		if err != nil {
			t.Fatalf("Failed to load partial config: %v", err)
		}

		// Verify partial values are loaded
		if cfg.Instance != "https://partial.example.com" {
			t.Errorf("Expected instance from config, got %s", cfg.Instance)
		}
		if cfg.Results != 15 {
			t.Errorf("Expected results from config, got %d", cfg.Results)
		}

		// Verify defaults are used for missing values
		if cfg.Format == "" {
			t.Error("Expected default format to be set")
		}
		if cfg.Language == "" {
			t.Error("Expected default language to be set")
		}
	})
}

// TestIntegrationConfigHomeDir tests loading config from user's home directory.
func TestIntegrationConfigHomeDir(t *testing.T) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine home directory")
	}

	configPath := filepath.Join(homeDir, ".search", "config.yaml")

	// Only test if config doesn't exist (to avoid clobbering user's config)
	if _, err := os.Stat(configPath); err == nil {
		t.Skip("User config already exists, skipping home directory test")
	}

	t.Run("CreateDefaultConfig", func(t *testing.T) {
		// Create default config in home directory
		cfg := DefaultConfig()
		
		if err := WriteConfig(configPath, *cfg); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		// Clean up after test
		defer func() {
			os.Remove(configPath)
			os.Remove(filepath.Dir(configPath))
		}()

		// Load config back - use LoadConfigFromFile for simple file loading
		loadedCfg, err := LoadConfigFromFile(configPath)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify
		if loadedCfg.Instance != cfg.Instance {
			t.Errorf("Config mismatch: instance")
		}
		if loadedCfg.Results != cfg.Results {
			t.Errorf("Config mismatch: results")
		}
	})
}
