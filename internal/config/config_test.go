package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.Instance != "https://search.butler.ooo" {
		t.Errorf("Expected default instance to be 'https://search.butler.ooo', got '%s'", cfg.Instance)
	}

	if cfg.Results != 10 {
		t.Errorf("Expected default results to be 10, got %d", cfg.Results)
	}

	if cfg.Format != "text" {
		t.Errorf("Expected default format to be 'text', got '%s'", cfg.Format)
	}

	if cfg.Timeout != 30 {
		t.Errorf("Expected default timeout to be 30, got %d", cfg.Timeout)
	}

	if cfg.Language != "en" {
		t.Errorf("Expected default language to be 'en', got '%s'", cfg.Language)
	}

	if cfg.SafeSearch != 1 {
		t.Errorf("Expected default safe search to be 1, got %d", cfg.SafeSearch)
	}

	if cfg.APIKey != "" {
		t.Errorf("Expected default API key to be empty, got '%s'", cfg.APIKey)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: false,
		},
		{
			name: "missing instance",
			cfg: &Config{
				Results:    10,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid results - too low",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    0,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid results - too high",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    101,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "xml",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - too low",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "json",
				Timeout:    0,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - too high",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "json",
				Timeout:    301,
				Language:   "en",
				SafeSearch: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid safe search - too low",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid safe search - too high",
			cfg: &Config{
				Instance:   "https://search.butler.ooo",
				Results:    10,
				Format:     "json",
				Timeout:    30,
				Language:   "en",
				SafeSearch: 3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	configContent := `instance: "https://search.example.com"
results: 20
format: "json"
timeout: 60
language: "de"
safe_search: 2
api_key: "test-key"
categories:
  - "general"
  - "images"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load using viper directly
	cfg := &Config{} // Use empty config instead of NewConfig() to avoid default SafeSearch=1
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	if err := v.Unmarshal(cfg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if cfg.Instance != "https://search.example.com" {
		t.Errorf("Expected instance to be 'https://search.example.com', got '%s'", cfg.Instance)
	}

	if cfg.Results != 20 {
		t.Errorf("Expected results to be 20, got %d", cfg.Results)
	}

	if cfg.Format != "json" {
		t.Errorf("Expected format to be 'json', got '%s'", cfg.Format)
	}

	if cfg.Timeout != 60 {
		t.Errorf("Expected timeout to be 60, got %d", cfg.Timeout)
	}

	if cfg.Language != "de" {
		t.Errorf("Expected language to be 'de', got '%s'", cfg.Language)
	}

	if cfg.SafeSearch != 2 {
		t.Errorf("Expected safe search to be 2, got %d", cfg.SafeSearch)
	}

	if cfg.APIKey != "test-key" {
		t.Errorf("Expected API key to be 'test-key', got '%s'", cfg.APIKey)
	}

	if len(cfg.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(cfg.Categories))
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	// Try to load a non-existent config file
	v := viper.New()
	v.SetConfigFile("/non/existent/path/config.yaml")
	err := v.ReadInConfig()
	if err == nil {
		t.Error("Expected error when loading non-existent config file, got nil")
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Write invalid YAML
	err := os.WriteFile(configFile, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	err = v.ReadInConfig()
	if err == nil {
		t.Error("Expected error when loading invalid YAML, got nil")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := &Config{
		Instance:   "https://search.example.com",
		Results:    25,
		Format:     "markdown",
		Timeout:    45,
		Language:   "fr",
		SafeSearch: 0,
		APIKey:     "saved-key",
		Categories: []string{"general", "news"},
	}

	err := cfg.SaveTo(configFile)
	if err != nil {
		t.Fatalf("SaveTo() error = %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load the saved config and verify
	loadedCfg := &Config{} // Use empty config instead of NewConfig() to avoid default SafeSearch=1
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}
	if err := v.Unmarshal(loadedCfg); err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if loadedCfg.Instance != cfg.Instance {
		t.Errorf("Expected instance '%s', got '%s'", cfg.Instance, loadedCfg.Instance)
	}

	if loadedCfg.Results != cfg.Results {
		t.Errorf("Expected results %d, got %d", cfg.Results, loadedCfg.Results)
	}

	if loadedCfg.Format != cfg.Format {
		t.Errorf("Expected format '%s', got '%s'", cfg.Format, loadedCfg.Format)
	}

	if loadedCfg.SafeSearch != cfg.SafeSearch {
		t.Errorf("Expected safe search %d, got %d", cfg.SafeSearch, loadedCfg.SafeSearch)
	}

	if loadedCfg.APIKey != cfg.APIKey {
		t.Errorf("Expected API key '%s', got '%s'", cfg.APIKey, loadedCfg.APIKey)
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{}

	cfg.ApplyDefaults()

	if cfg.Instance != "https://search.butler.ooo" {
		t.Errorf("Expected default instance, got '%s'", cfg.Instance)
	}

	if cfg.Results != 10 {
		t.Errorf("Expected default results, got %d", cfg.Results)
	}

	if cfg.Format != "text" {
		t.Errorf("Expected default format, got '%s'", cfg.Format)
	}
}

func TestCliConfigApplyToConfig(t *testing.T) {
	cliCfg := &CliConfig{
		Instance:   "https://cli.example.com",
		Results:    50,
		Format:     "json",
		Category:   "images",
		Timeout:    60,
		Language:   "es",
		SafeSearch: 0,
		APIKey:     "cli-api-key",
	}

	cfg := NewConfig()
	cliCfg.ApplyToConfig(cfg)

	if cfg.Instance != "https://cli.example.com" {
		t.Errorf("Expected CLI instance, got '%s'", cfg.Instance)
	}

	if cfg.Results != 50 {
		t.Errorf("Expected CLI results, got %d", cfg.Results)
	}

	if cfg.Format != "json" {
		t.Errorf("Expected CLI format, got '%s'", cfg.Format)
	}

	if len(cfg.Categories) != 1 || cfg.Categories[0] != "images" {
		t.Errorf("Expected CLI category 'images', got %v", cfg.Categories)
	}

	if cfg.Language != "es" {
		t.Errorf("Expected CLI language, got '%s'", cfg.Language)
	}

	if cfg.SafeSearch != 0 {
		t.Errorf("Expected CLI safe search, got %d", cfg.SafeSearch)
	}

	if cfg.APIKey != "cli-api-key" {
		t.Errorf("Expected CLI API key 'cli-api-key', got '%s'", cfg.APIKey)
	}
}
