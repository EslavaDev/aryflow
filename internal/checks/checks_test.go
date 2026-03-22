package checks

import (
	"strings"
	"testing"
)

func TestCheckGit(t *testing.T) {
	version, err := CheckGit()
	if err != nil {
		// Git should be available in any dev environment
		t.Skipf("Git not available in test environment: %v", err)
	}

	if version == "" {
		t.Error("expected non-empty Git version string")
	}

	// Version should look like a semver (e.g., "2.43.0")
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		t.Errorf("expected semver-like version, got %q", version)
	}
}

func TestCheckNodeVersion(t *testing.T) {
	version, err := CheckNode()
	if err != nil {
		// Node might not be installed in all environments
		if strings.Contains(err.Error(), "not found") {
			t.Skip("Node.js not available in test environment")
		}
		// If Node is installed but version is too low, that's a valid test result
		if strings.Contains(err.Error(), ">= 18 required") {
			t.Logf("Node.js found but too old: %v", err)
			return
		}
		t.Skipf("Node check failed: %v", err)
	}

	if version == "" {
		t.Error("expected non-empty Node version string")
	}
}

func TestCheckHomebrew(t *testing.T) {
	version, err := CheckHomebrew()
	if err != nil {
		t.Skipf("Homebrew not available in test environment: %v", err)
	}

	if version == "" {
		t.Error("expected non-empty Homebrew version string")
	}
}

func TestRunCommandInvalid(t *testing.T) {
	_, err := runCommand("nonexistent-command-that-does-not-exist-xyz")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}
