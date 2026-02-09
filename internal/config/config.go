// Package config provides configuration management for the search CLI.
//
// It handles loading configuration from files, environment variables, and CLI flags,
// with proper precedence handling and validation. Configuration is stored in YAML format
// at ~/.search/config.yaml by default.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	defaultConfigDir = ".search"
	configFileName   = "config.yaml"
)

// Config holds the application configuration.
//
// The configuration includes settings for the SearXNG instance, search parameters,
// output formatting, and other user preferences. Fields can be configured via
// config file, environment variables, or CLI flags (in that order of precedence).
type Config struct {
	Instance     string   `yaml:"instance" mapstructure:"instance"`
	Results      int      `yaml:"results" mapstructure:"results"`
	Format       string   `yaml:"format" mapstructure:"format"`
	APIKey       string   `yaml:"api_key,omitempty" mapstructure:"api_key"`
	Timeout      int      `yaml:"timeout" mapstructure:"timeout"`
	Categories   []string `yaml:"categories,omitempty" mapstructure:"categories"`
	Language     string   `yaml:"language" mapstructure:"language"`
	SafeSearch   int      `yaml:"safe_search" mapstructure:"safe_search"`
	Verbose      bool     `yaml:"verbose,omitempty" mapstructure:"verbose"`
	// Cache settings
	CacheEnabled bool `yaml:"cache_enabled,omitempty" mapstructure:"cache_enabled"`
	CacheSize    int  `yaml:"cache_size,omitempty" mapstructure:"cache_size"`
	CacheTTL     int  `yaml:"cache_ttl,omitempty" mapstructure:"cache_ttl"` // in seconds
}

// NewConfig creates a new Config with default values.
//
// This is the recommended way to create a new configuration instance.
// Returns a Config populated with sensible defaults for all fields.
//
// Example:
//
//	cfg := config.NewConfig()
//	fmt.Println(cfg.Instance) // "https://search.butler.ooo"
//	fmt.Println(cfg.Results)  // 10
func NewConfig() *Config {
	return DefaultConfig()
}

// DefaultConfig returns a Config with default values.
//
// The defaults include:
//   - Instance: https://search.butler.ooo
//   - Results: 10
//   - Format: text
//   - Timeout: 30 seconds
//   - Categories: general
//   - Language: en
//   - SafeSearch: 1 (moderate)
//   - CacheEnabled: true
//   - CacheSize: 100
//   - CacheTTL: 300 (5 minutes)
//
// Example:
//
//	cfg := config.DefaultConfig()
//	// cfg.Instance == "https://search.butler.ooo"
//	// cfg.Results == 10
func DefaultConfig() *Config {
	return &Config{
		Instance:     "https://search.butler.ooo",
		Results:      10,
		Format:       "text",
		Timeout:      30,
		Categories:   []string{"general"},
		Language:     "en",
		SafeSearch:   1,
		CacheEnabled: true,
		CacheSize:    100,
		CacheTTL:     300,
	}
}

// WriteConfig writes a Config to a specific file path.
//
// It creates any necessary directories with permissions 0755 and writes
// the config file with permissions 0600 (owner read/write only).
// Returns an error if directory creation or file writing fails.
//
// Example:
//
//	cfg := config.DefaultConfig()
//	cfg.Results = 20
//	err := config.WriteConfig("/tmp/test-config.yaml", cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
func WriteConfig(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions (owner read/write only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Load loads the configuration from the default location (~/.search/config.yaml).
//
// If the config file doesn't exist, it will be created with default values.
// If the file exists but can't be read, an error is returned.
// After loading, default values are applied to any empty fields.
func (c *Config) Load() error {
	v := viper.New()

	// Set config name and type
	v.SetConfigName(configFileName)

	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, defaultConfigDir)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v.AddConfigPath(configDir)

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// If config doesn't exist, create default
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Apply defaults to current config before saving
			c.applyDefaults()
			return c.Save()
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal config
	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for any empty fields
	c.applyDefaults()

	return nil
}

// Save saves the configuration to the default location (~/.search/config.yaml).
//
// The config directory will be created if it doesn't exist.
// Returns an error if the directory can't be created or the file can't be written.
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, defaultConfigDir)

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, configFileName)

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SaveTo saves the configuration to a specific file path.
//
// This is useful for creating config files in non-standard locations
// (e.g., for testing or custom setups).
// The parent directory will be created if it doesn't exist.
func (c *Config) SaveTo(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate validates the configuration fields.
//
// It checks that:
//   - Instance URL is not empty
//   - Results is between 1 and 100
//   - Timeout is between 1 and 300
//   - SafeSearch is between 0 and 2
//   - Format is one of: json, markdown, text
//
// Returns an error describing the validation failure, or nil if valid.
func (c *Config) Validate() error {
	if c.Instance == "" {
		return fmt.Errorf("instance URL cannot be empty")
	}
	if c.Results < 1 || c.Results > 100 {
		return fmt.Errorf("results must be between 1 and 100, got %d", c.Results)
	}
	if c.Timeout < 1 || c.Timeout > 300 {
		return fmt.Errorf("timeout must be between 1 and 300, got %d", c.Timeout)
	}
	if c.SafeSearch < 0 || c.SafeSearch > 2 {
		return fmt.Errorf("safe search level must be between 0 and 2, got %d", c.SafeSearch)
	}
	if c.Format != "" && c.Format != "json" && c.Format != "markdown" && c.Format != "text" {
		return fmt.Errorf("invalid format '%s', must be json, markdown, or text", c.Format)
	}
	return nil
}

// ApplyDefaults sets default values for empty fields.
//
// This is a public version of applyDefaults that can be called externally.
// It only sets defaults for fields that are at their zero value.
func (c *Config) ApplyDefaults() {
	c.applyDefaults()
}

// applyDefaults sets default values for empty fields (private version).
//
// Note: SafeSearch is not defaulted here because 0 is a valid value.
// It should be set to 1 only in NewConfig().
func (c *Config) applyDefaults() {
	if c.Instance == "" {
		c.Instance = "https://search.butler.ooo"
	}
	if c.Results == 0 {
		c.Results = 10
	}
	if c.Format == "" {
		c.Format = "text"
	}
	if c.Timeout == 0 {
		c.Timeout = 30
	}
	if len(c.Categories) == 0 {
		c.Categories = []string{"general"}
	}
	if c.Language == "" {
		c.Language = "en"
	}
	if c.CacheSize == 0 {
		c.CacheSize = 100
	}
	if c.CacheTTL == 0 {
		c.CacheTTL = 300
	}
	// Note: CacheEnabled defaults to false here so that it must be explicitly enabled
	// Note: We don't set a default for SafeSearch here because 0 is a valid value
	// It should be set to 1 only in NewConfig()
}

// LoadConfig loads configuration with proper priority handling.
//
// Priority order (highest to lowest):
// 1. CLI flags
// 2. Environment variables
// 3. Config file
// 4. Default values
//
// If cliCfg.ConfigPath is set, that file will be used instead of the default.
// Returns a validated Config or an error if loading/validating fails.
func LoadConfig(cliCfg *CliConfig) (*Config, error) {
	// Start with defaults
	cfg := NewConfig()

	// Load from config file (unless --config is specified with non-existent file)
	if cliCfg.ConfigPath == "" {
		if err := cfg.Load(); err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		v := viper.New()
		v.SetConfigFile(cliCfg.ConfigPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := v.Unmarshal(cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
		cfg.applyDefaults()
	}

	// Apply environment variables (override config file)
	cfg.applyEnvironmentVariables()

	// Apply CLI flags (highest priority)
	cliCfg.ApplyToConfig(cfg)

	return cfg, nil
}

// LoadConfigFromFile loads configuration from a specific file path.
//
// If the file doesn't exist, returns a Config with defaults applied.
// This is useful for loading non-standard config locations.
// Returns an error if the file exists but can't be parsed.
func LoadConfigFromFile(path string) (*Config, error) {
	cfg := NewConfig()

	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		// If file doesn't exist, return defaults
		// Check for os.PathError (file not found)
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			cfg.applyDefaults()
			cfg.applyEnvironmentVariables()
			return cfg, nil
		}
		// Also check for viper's ConfigFileNotFoundError
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg.applyDefaults()
			cfg.applyEnvironmentVariables()
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.applyDefaults()
	cfg.applyEnvironmentVariables()

	return cfg, nil
}

// applyEnvironmentVariables overrides config with environment variables
func (c *Config) applyEnvironmentVariables() {
	if v := os.Getenv("SEARCH_INSTANCE"); v != "" {
		c.Instance = v
	}
	if v := os.Getenv("SEARCH_RESULTS"); v != "" {
		c.Results = parseIntEnv(v)
	}
	if v := os.Getenv("SEARCH_FORMAT"); v != "" {
		c.Format = v
	}
	if v := os.Getenv("SEARCH_TIMEOUT"); v != "" {
		c.Timeout = parseIntEnv(v)
	}
	if v := os.Getenv("SEARCH_LANGUAGE"); v != "" {
		c.Language = v
	}
	if v := os.Getenv("SEARCH_SAFE_SEARCH"); v != "" {
		c.SafeSearch = parseIntEnv(v)
	}
	if v := os.Getenv("SEARCH_API_KEY"); v != "" {
		c.APIKey = v
	}
}

// CliConfig holds CLI-specific configuration overrides.
//
// These values represent command-line flag values that take precedence
// over config file and environment variable settings.
type CliConfig struct {
	Instance     string
	Results      int
	Format       string
	Category     string
	Timeout      int
	Language     string
	SafeSearch   int
	ConfigPath   string
	Page         int
	TimeRange    string
	Open         bool
	OpenAll      bool
	Verbose      bool
	APIKey       string
	// Cache options
	CacheEnabled *bool // Pointer to distinguish between not set, false, and true
	NoCache      bool  // Shortcut for --no-cache to disable caching
	CacheSize    *int  // Pointer to distinguish between not set and 0
	CacheTTL     *int  // Pointer to distinguish between not set and 0
}

// ApplyToConfig applies CLI config values to the main Config.
//
// Only non-zero values from CliConfig are applied, allowing CLI flags
// to selectively override config settings.
func (c *CliConfig) ApplyToConfig(cfg *Config) {
	if c.Instance != "" {
		cfg.Instance = c.Instance
	}
	if c.Results > 0 {
		cfg.Results = c.Results
	}
	if c.Format != "" {
		cfg.Format = c.Format
	}
	if len(c.Category) > 0 {
		cfg.Categories = []string{c.Category}
	}
	if c.Timeout > 0 {
		cfg.Timeout = c.Timeout
	}
	if c.Language != "" {
		cfg.Language = c.Language
	}
	if c.SafeSearch >= 0 {
		cfg.SafeSearch = c.SafeSearch
	}
	if c.APIKey != "" {
		cfg.APIKey = c.APIKey
	}
	cfg.Verbose = c.Verbose
	// Handle cache settings
	if c.CacheEnabled != nil {
		cfg.CacheEnabled = *c.CacheEnabled
	}
	if c.NoCache {
		cfg.CacheEnabled = false
	}
	if c.CacheSize != nil && *c.CacheSize > 0 {
		cfg.CacheSize = *c.CacheSize
	}
	if c.CacheTTL != nil && *c.CacheTTL > 0 {
		cfg.CacheTTL = *c.CacheTTL
	}
	if c.NoCache {
		cfg.CacheEnabled = false
	}
}

func parseIntEnv(v string) int {
	result := 0
	fmt.Sscanf(v, "%d", &result)
	return result
}
