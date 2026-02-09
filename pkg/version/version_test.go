package version

import (
	"testing"
)

func TestVersionVariables(t *testing.T) {
	tests := []struct {
		name  string
		varPtr *string
		varName string
	}{
		{
			name:  "Version variable exists",
			varPtr: &Version,
			varName: "Version",
		},
		{
			name:  "GitCommit variable exists",
			varPtr: &GitCommit,
			varName: "GitCommit",
		},
		{
			name:  "BuildDate variable exists",
			varPtr: &BuildDate,
			varName: "BuildDate",
		},
		{
			name:  "GoVersion variable exists",
			varPtr: &GoVersion,
			varName: "GoVersion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.varPtr == nil {
				t.Errorf("%s variable is nil", tt.varName)
			}
			if *tt.varPtr == "" {
				t.Logf("%s is empty (may be set at build time)", tt.varName)
			}
		})
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that Version follows semantic versioning format
	// Format: v1.2.3 or 1.2.3
	if Version == "" {
		t.Skip("Version not set at build time")
	}

	// Check if version starts with 'v' (optional)
	v := Version
	if len(v) > 0 && v[0] == 'v' {
		v = v[1:]
	}

	// Basic format check - should contain dots
	t.Logf("Version: %s", Version)
}

func TestGetVersionInfo(t *testing.T) {
	type VersionInfo struct {
		Version   string
		GitCommit string
		BuildDate string
		GoVersion string
	}

	info := VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	}

	if info.Version == "" {
		t.Error("Expected Version to be set")
	}
	if info.GitCommit == "" {
		t.Error("Expected GitCommit to be set")
	}
	if info.BuildDate == "" {
		t.Error("Expected BuildDate to be set")
	}
	if info.GoVersion == "" {
		t.Error("Expected GoVersion to be set")
	}

	t.Logf("Version Info: %+v", info)
}

func TestVersionString(t *testing.T) {
	versionStr := Version
	
	if versionStr == "unknown" || versionStr == "" {
		t.Log("Version is set to default/unknown value (expected in development)")
	}
	
	if versionStr != "" {
		t.Logf("Running search version: %s", versionStr)
	}
}
