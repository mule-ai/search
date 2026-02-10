package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()
	if cmd == nil {
		t.Fatal("NewRootCommand() returned nil")
	}
	if cmd.Command == nil {
		t.Error("NewRootCommand() Command is nil")
	}
	if cmd.Use != "search" {
		t.Errorf("Expected Use 'search', got '%s'", cmd.Use)
	}
	if cmd.Version == "" {
		t.Error("Version is empty")
	}
}

func TestRootCommandFlags(t *testing.T) {
	cmd := NewRootCommand()
	flags := cmd.Flags()

	// Test that all required flags are defined
	requiredFlags := []string{
		"instance", "i",
		"results", "n",
		"format", "f",
		"category", "c",
		"timeout", "t",
		"language", "l",
		"safe", "s",
		"config",
		"verbose", "v",
		"page",
		"time",
		"open",
		"open-all",
		"no-color",
	}

	flagFound := make(map[string]bool)
	flags.VisitAll(func(flag *pflag.Flag) {
		flagFound[flag.Name] = true
		if flag.Shorthand != "" {
			flagFound[flag.Shorthand] = true
		}
	})

	for _, flag := range requiredFlags {
		if !flagFound[flag] {
			t.Errorf("Required flag '%s' not found", flag)
		}
	}
}

func TestRootCommandDefaultValues(t *testing.T) {
	cmd := NewRootCommand()
	flags := cmd.Flags()

	tests := []struct {
		name     string
		flag     string
		expected string
	}{
		{"instance", "instance", "https://search.butler.ooo"},
		{"format", "format", "text"},
		{"category", "category", "general"},
		{"language", "language", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flags.Lookup(tt.flag).DefValue
			if got != tt.expected {
				t.Errorf("Flag %s default value = %s, want %s", tt.flag, got, tt.expected)
			}
		})
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	cmd := NewRootCommand()

	// Check for expected subcommands
	expectedCommands := []string{"version", "categories", "completion"}
	for _, expected := range expectedCommands {
		found := false
		for _, subcmd := range cmd.Commands() {
			if subcmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	cmd := NewRootCommand()
	versionCmd, _, err := cmd.Find([]string{"version"})
	if err != nil {
		t.Fatalf("Failed to find version command: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := versionCmd.RunE(versionCmd, []string{}); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	output := buf.String()
	if !strings.Contains(output, "search version") {
		t.Errorf("Version output missing version string: %s", output)
	}
}

func TestCategoriesCommand(t *testing.T) {
	cmd := NewRootCommand()
	categoriesCmd, _, err := cmd.Find([]string{"categories"})
	if err != nil {
		t.Fatalf("Failed to find categories command: %v", err)
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := categoriesCmd.RunE(categoriesCmd, []string{}); err != nil {
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("categories command failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	output := buf.String()
	if !strings.Contains(output, "Available Search Categories") {
		t.Errorf("Categories output missing header: %s", output)
	}

	// Check for known categories
	knownCategories := []string{"general", "images", "videos", "news"}
	for _, cat := range knownCategories {
		if !strings.Contains(output, cat) {
			t.Errorf("Categories output missing category '%s'", cat)
		}
	}
}

func TestRootCommandRequiresArgs(t *testing.T) {
	cmd := NewRootCommand()

	// Test with no arguments
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}
}

func TestRootCommandHelp(t *testing.T) {
	cmd := NewRootCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Search is a powerful") {
		t.Error("Help output missing description")
	}
	if !strings.Contains(output, "Usage:") {
		t.Error("Help output missing usage section")
	}
	if !strings.Contains(output, "Flags:") {
		t.Error("Help output missing flags section")
	}
}

func TestRootCommandLongDescription(t *testing.T) {
	cmd := NewRootCommand()
	long := cmd.Long

	expectedContent := []string{
		"SearXNG",
		"search -n 20",
		"search -f json",
	}

	for _, content := range expectedContent {
		if !strings.Contains(long, content) {
			t.Errorf("Long description missing '%s'", content)
		}
	}
}

func TestConfigFlagsDefaults(t *testing.T) {
	cfg := ConfigFlags{}

	if cfg.Instance != "" {
		t.Errorf("Expected empty Instance, got '%s'", cfg.Instance)
	}
	if cfg.Results != 0 {
		t.Errorf("Expected Results to be 0, got %d", cfg.Results)
	}
	if cfg.Format != "" {
		t.Errorf("Expected empty Format, got '%s'", cfg.Format)
	}
}

func TestAddGlobalFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	var cfg ConfigFlags

	addGlobalFlags(cmd.Flags(), &cfg)

	// Verify flags were added
	if cmd.Flags().Lookup("instance") == nil {
		t.Error("instance flag not added")
	}
	if cmd.Flags().Lookup("results") == nil {
		t.Error("results flag not added")
	}
	if cmd.Flags().Lookup("format") == nil {
		t.Error("format flag not added")
	}
	if cmd.Flags().Lookup("verbose") == nil {
		t.Error("verbose flag not added")
	}

	// Verify shorthand flags
	if cmd.Flags().ShorthandLookup("i") == nil {
		t.Error("instance shorthand 'i' not found")
	}
	if cmd.Flags().ShorthandLookup("n") == nil {
		t.Error("results shorthand 'n' not found")
	}
	if cmd.Flags().ShorthandLookup("v") == nil {
		t.Error("verbose shorthand 'v' not found")
	}
}

func TestRootCommandVersionTemplate(t *testing.T) {
	cmd := NewRootCommand()
	template := cmd.VersionTemplate()

	if !strings.Contains(template, "search version") {
		t.Error("Version template missing version placeholder")
	}
}

func TestRootCommandArgsValidation(t *testing.T) {
	cmd := NewRootCommand()

	// Test with too many arguments
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"arg1", "arg2"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when too many arguments provided")
	}
}

func TestPersistentPreRun(t *testing.T) {
	cfg := &ConfigFlags{Verbose: true}
	preRun := persistentPreRun(cfg)

	cmd := &cobra.Command{}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := preRun(cmd, []string{"test"})
	w.Close()
	os.Stderr = oldStderr

	if err != nil {
		t.Errorf("persistentPreRun failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	// Verbose mode should output something
	if buf.Len() == 0 {
		t.Error("Expected output in verbose mode")
	}
}

func TestPersistentPreRunNonVerbose(t *testing.T) {
	cfg := &ConfigFlags{Verbose: false}
	preRun := persistentPreRun(cfg)

	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetErr(&out)

	err := preRun(cmd, []string{"test"})
	if err != nil {
		t.Errorf("persistentPreRun failed: %v", err)
	}

	// Non-verbose mode should not output to stderr
	if out.Len() > 0 {
		t.Error("Expected no output in non-verbose mode")
	}
}

// Mock the Execute function for basic testing
func TestExecuteFunction(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with --help
	os.Args = []string{"search", "--help"}
	err := Execute()
	if err != nil {
		t.Errorf("Execute with --help failed: %v", err)
	}
}

func TestCategoriesCommandDescription(t *testing.T) {
	cmd := NewRootCommand()
	categoriesCmd, _, err := cmd.Find([]string{"categories"})
	if err != nil {
		t.Fatalf("Failed to find categories command: %v", err)
	}

	if categoriesCmd.Short == "" {
		t.Error("Categories command has no Short description")
	}
	if categoriesCmd.Long == "" {
		t.Error("Categories command has no Long description")
	}

	long := categoriesCmd.Long
	expectedContent := []string{"general", "images", "videos", "news", "map", "music"}
	for _, content := range expectedContent {
		if !strings.Contains(long, content) {
			t.Errorf("Categories Long description missing '%s'", content)
		}
	}
}
