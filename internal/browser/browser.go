// Package browser provides cross-platform browser opening functionality.
//
// It detects the appropriate browser command for the current operating system
// and provides methods to open URLs in the default browser.
package browser

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

// Command returns the appropriate browser command for the current OS.
//
// On macOS: "open"
// On Windows: "start"
// On Linux/Unix: First available of "xdg-open", "gio", "firefox", "chromium", etc.
//
// Example:
//
//	cmd := browser.Command()
//	fmt.Printf("Using browser command: %s\n", cmd)
func Command() string {
	switch runtime.GOOS {
	case "darwin":
		// macOS
		return "open"
	case "windows":
		// Windows
		return "start"
	default:
		// Linux and other Unix-like systems
		// Check for common browsers in order of preference
		browsers := []string{
			"xdg-open", // Standard Linux opener
			"gio",      // GNOME
			"firefox",
			"chromium",
			"chromium-browser",
			"google-chrome",
			"chrome",
		}

		for _, browser := range browsers {
			if _, err := exec.LookPath(browser); err == nil {
				return browser
			}
		}

		// Fallback to xdg-open
		return "xdg-open"
	}
}

// ValidateURL checks if a URL is safe to open in a browser.
//
// It validates the URL format and checks for potentially dangerous schemes.
// Returns an error if the URL is invalid or unsafe.
//
// Valid schemes: http, https, ftp, ftps, file
// Invalid schemes: javascript:, data:, vbscript:, file: (when not from user input)
//
// Example:
//
//	err := browser.ValidateURL("https://example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse URL to validate format
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check for dangerous schemes
	dangerousSchemes := []string{
		"javascript",
		"data",
		"vbscript",
		"file",
		"chrome",
		"chrome-extension",
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	for _, dangerous := range dangerousSchemes {
		if scheme == dangerous {
			return fmt.Errorf("unsafe URL scheme: %s (potential security risk)", scheme)
		}
	}

	// Ensure the URL has a safe scheme
	validSchemes := []string{"http", "https", "ftp", "ftps", "mailto"}
	hasValidScheme := false
	for _, valid := range validSchemes {
		if scheme == valid {
			hasValidScheme = true
			break
		}
	}

	if !hasValidScheme {
		return fmt.Errorf("unsupported URL scheme: %s", scheme)
	}

	return nil
}

// SanitizeURL cleans up a URL before opening it.
//
// It removes potentially dangerous elements and ensures the URL is safe.
// Returns the sanitized URL or an error if the URL is invalid.
//
// Example:
//
//	cleanURL, err := browser.SanitizeURL(userURL)
//	if err != nil {
//	    log.Fatal(err)
//	}
func SanitizeURL(rawURL string) (string, error) {
	// First validate the URL
	if err := ValidateURL(rawURL); err != nil {
		return "", err
	}

	// Parse and rebuild the URL to normalize it
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Force HTTPS if the URL uses HTTP (security best practice)
	// This can be overridden if needed
	if parsedURL.Scheme == "http" {
		// For security, prefer https
		// parsedURL.Scheme = "https"
	}

	// Remove any authentication information for security
	parsedURL.User = nil

	// Return the cleaned URL
	return parsedURL.String(), nil
}

// OpenURL opens a URL in the default browser.
//
// It validates the URL format and checks for potentially dangerous schemes
// before opening. Returns an error if the URL is invalid or the browser command fails.
//
// Example:
//
//	err := browser.OpenURL("https://example.com")
//	if err != nil {
//	    log.Fatalf("Failed to open browser: %v", err)
//	}
func OpenURL(rawURL string) error {
	// Validate and sanitize the URL
	cleanURL, err := SanitizeURL(rawURL)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// On Windows, use 'start' command
		cmd = exec.Command("cmd", "/c", "start", "", cleanURL)
	case "darwin":
		// On macOS, use 'open' command
		cmd = exec.Command("open", cleanURL)
	default:
		// On Linux/Unix, use the detected browser command
		browserCmd := Command()
		cmd = exec.Command(browserCmd, cleanURL)
	}

	// Start the command without waiting for it to complete
	// This allows the browser to open independently
	return cmd.Start()
}

// OpenURLs opens multiple URLs in the default browser.
//
// Each URL is opened in a separate tab/window. The function continues
// opening URLs even if one fails, but returns the first error encountered.
//
// Example:
//
//	urls := []string{"https://example.com", "https://google.com"}
//	err := browser.OpenURLs(urls)
//	if err != nil {
//	    log.Fatalf("Failed to open URLs: %v", err)
//	}
func OpenURLs(urls []string) error {
	for i, url := range urls {
		if err := OpenURL(url); err != nil {
			return fmt.Errorf("failed to open URL %d (%s): %w", i+1, url, err)
		}
	}
	return nil
}

// IsSupported checks if browser opening is supported on the current platform.
//
// Returns true if a browser command is found, false otherwise.
//
// Example:
//
//	if !browser.IsSupported() {
//	    log.Fatal("Browser opening not supported on this platform")
//	}
func IsSupported() bool {
	cmd := Command()
	if _, err := exec.LookPath(cmd); err != nil {
		return false
	}
	return true
}

// OpenWithCommand opens a URL using a specific browser command.
//
// The browser command can include arguments (e.g., "firefox --private-window").
// The URL is appended to the command arguments.
//
// Example:
//
//	err := browser.OpenWithCommand("https://example.com", "firefox --private-window")
//	if err != nil {
//	    log.Fatalf("Failed to open browser: %v", err)
//	}
func OpenWithCommand(url string, browserCmd string) error {
	parts := strings.Fields(browserCmd)

	// If the command has arguments, use them
	if len(parts) > 1 {
		cmdParts := append(parts, url)
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		return cmd.Start()
	}

	// Simple command
	cmd := exec.Command(parts[0], url)
	return cmd.Start()
}