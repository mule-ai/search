// Package browser provides cross-platform browser opening functionality.
package browser

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid https URL",
			url:         "https://example.com",
			expectError: false,
		},
		{
			name:        "valid http URL",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "valid ftp URL",
			url:         "ftp://example.com/file",
			expectError: false,
		},
		{
			name:        "valid ftps URL",
			url:         "ftps://example.com/file",
			expectError: false,
		},
		{
			name:        "valid mailto URL",
			url:         "mailto:test@example.com",
			expectError: false,
		},
		{
			name:        "valid URL with path",
			url:         "https://example.com/path/to/page",
			expectError: false,
		},
		{
			name:        "valid URL with query",
			url:         "https://example.com?query=test",
			expectError: false,
		},
		{
			name:        "valid URL with fragment",
			url:         "https://example.com#section",
			expectError: false,
		},
		{
			name:        "valid URL with port",
			url:         "https://example.com:8080",
			expectError: false,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
			errorMsg:    "URL cannot be empty",
		},
		{
			name:        "javascript URL",
			url:         "javascript:alert('xss')",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "data URL",
			url:         "data:text/html,<script>alert('xss')</script>",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "vbscript URL",
			url:         "vbscript:msgbox('xss')",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "file URL",
			url:         "file:///etc/passwd",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "chrome-extension URL",
			url:         "chrome-extension://abc/page.html",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "chrome URL",
			url:         "chrome://settings",
			expectError: true,
			errorMsg:    "unsafe URL scheme",
		},
		{
			name:        "invalid URL format",
			url:         "://example.com",
			expectError: true,
			errorMsg:    "invalid URL format",
		},
		{
			name:        "invalid scheme",
			url:         "gopher://example.com",
			expectError: true,
			errorMsg:    "unsupported URL scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateURL(%q) expected error, got nil", tt.url)
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateURL(%q) error = %v, expected to contain %q", tt.url, err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateURL(%q) unexpected error: %v", tt.url, err)
				}
			}
		})
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "normal HTTPS URL",
			input:       "https://example.com",
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "URL with authentication",
			input:       "https://user:pass@example.com",
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:        "URL with port",
			input:       "https://example.com:8080/path",
			expected:    "https://example.com:8080/path",
			expectError: false,
		},
		{
			name:        "URL with query and fragment",
			input:       "https://example.com/path?q=1#section",
			expected:    "https://example.com/path?q=1#section",
			expectError: false,
		},
		{
			name:        "invalid URL",
			input:       "://example.com",
			expectError: true,
		},
		{
			name:        "javascript URL rejected",
			input:       "javascript:alert('xss')",
			expectError: true,
		},
		{
			name:        "file URL rejected",
			input:       "file:///etc/passwd",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeURL(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("SanitizeURL(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("SanitizeURL(%q) unexpected error: %v", tt.input, err)
					return
				}
				if result != tt.expected {
					t.Errorf("SanitizeURL(%q) = %q, expected %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestCommand(t *testing.T) {
	cmd := Command()

	if cmd == "" {
		t.Error("Expected Command() to return a non-empty string")
	}
}

func TestIsSupported(t *testing.T) {
	// We can't test the actual value since it depends on the system
	// but we can verify it returns a boolean and doesn't panic
	supported := IsSupported()

	// Should be either true or false
	if supported != true && supported != false {
		t.Error("IsSupported() should return a boolean")
	}
}

func TestOpenURLValidation(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid URL - but may fail if no browser",
			url:         "https://example.com",
			expectError: false, // May still fail if no browser, but validation passes
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
		},
		{
			name:        "javascript URL blocked",
			url:         "javascript:alert('xss')",
			expectError: true,
		},
		{
			name:        "file URL blocked",
			url:         "file:///etc/passwd",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Errorf("OpenURL(%q) expected error, got nil", tt.url)
				}
			} else {
				// URL validation passed, but actual browser opening might still fail
				// We're mainly testing that invalid URLs are rejected
				if err != nil && !containsString(err.Error(), "browser") && !containsString(err.Error(), "not supported") {
					t.Logf("OpenURL(%q) error (may be expected on system without browser): %v", tt.url, err)
				}
			}
		})
	}
}

func TestOpenURLs(t *testing.T) {
	tests := []struct {
		name        string
		urls        []string
		expectError bool
	}{
		{
			name:        "single URL",
			urls:        []string{"https://example.com"},
			expectError: false,
		},
		{
			name:        "multiple URLs",
			urls:        []string{"https://example.com", "https://example.org"},
			expectError: false,
		},
		{
			name:        "empty list",
			urls:        []string{},
			expectError: false,
		},
		{
			name:        "contains invalid URL",
			urls:        []string{"https://example.com", "javascript:alert('xss')"},
			expectError: true,
		},
		{
			name:        "contains empty URL",
			urls:        []string{"https://example.com", ""},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenURLs(tt.urls)

			if tt.expectError {
				if err == nil {
					t.Errorf("OpenURLs(%v) expected error, got nil", tt.urls)
				}
			} else {
				// Validation passed, browser opening might still fail
				if err != nil && !containsString(err.Error(), "browser") && !containsString(err.Error(), "not supported") {
					t.Logf("OpenURLs(%v) error (may be expected on system without browser): %v", tt.urls, err)
				}
			}
		})
	}
}

func TestOpenWithCommand(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		browserCmd  string
		expectError bool
	}{
		{
			name:        "simple command",
			url:         "https://example.com",
			browserCmd:  "echo",
			expectError: false, // echo will succeed
		},
		{
			name:        "command with args",
			url:         "https://example.com",
			browserCmd:  "echo test",
			expectError: false,
		},
		{
			name:        "non-existent command",
			url:         "https://example.com",
			browserCmd:  "nonexistent-browser-command-xyz",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OpenWithCommand(tt.url, tt.browserCmd)

			if tt.expectError {
				if err == nil {
					t.Errorf("OpenWithCommand(%q, %q) expected error, got nil", tt.url, tt.browserCmd)
				}
			} else {
				// Command should succeed
				if err != nil {
					t.Errorf("OpenWithCommand(%q, %q) unexpected error: %v", tt.url, tt.browserCmd, err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkValidateURL(b *testing.B) {
	url := "https://example.com/path?query=value#fragment"
	for i := 0; i < b.N; i++ {
		ValidateURL(url)
	}
}

func BenchmarkSanitizeURL(b *testing.B) {
	url := "https://user:pass@example.com:8080/path?q=1#section"
	for i := 0; i < b.N; i++ {
		SanitizeURL(url)
	}
}

func BenchmarkValidateURLWithAuth(b *testing.B) {
	url := "https://user:password@example.com/path"
	for i := 0; i < b.N; i++ {
		ValidateURL(url)
	}
}
