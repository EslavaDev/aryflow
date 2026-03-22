package ui

import (
	"os"
	"testing"
)

func TestNoColor(t *testing.T) {
	// Save original value and restore after test
	orig := os.Getenv("ARYFLOW_NO_COLOR")
	defer os.Setenv("ARYFLOW_NO_COLOR", orig)

	// With ARYFLOW_NO_COLOR=1, colors should be disabled
	os.Setenv("ARYFLOW_NO_COLOR", "1")
	if !noColor() {
		t.Error("expected noColor() to return true when ARYFLOW_NO_COLOR=1")
	}

	// Without the env var, colors should be enabled
	os.Unsetenv("ARYFLOW_NO_COLOR")
	origNoColor := os.Getenv("NO_COLOR")
	os.Unsetenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", origNoColor)

	if noColor() {
		t.Error("expected noColor() to return false when env vars are unset")
	}
}

func TestNoColorDisablesANSI(t *testing.T) {
	orig := os.Getenv("ARYFLOW_NO_COLOR")
	defer os.Setenv("ARYFLOW_NO_COLOR", orig)

	// With colors disabled, colorize should return plain text
	os.Setenv("ARYFLOW_NO_COLOR", "1")
	result := colorize(colorGreen, "hello")
	if result != "hello" {
		t.Errorf("expected plain 'hello', got %q", result)
	}

	// With colors enabled, colorize should wrap with ANSI codes
	os.Unsetenv("ARYFLOW_NO_COLOR")
	origNoColor := os.Getenv("NO_COLOR")
	os.Unsetenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", origNoColor)

	result = colorize(colorGreen, "hello")
	expected := colorGreen + "hello" + colorReset
	if result != expected {
		t.Errorf("expected ANSI-wrapped string, got %q", result)
	}
}

func TestAutoYes(t *testing.T) {
	orig := os.Getenv("ARYFLOW_YES")
	defer os.Setenv("ARYFLOW_YES", orig)

	// With ARYFLOW_YES=1, autoYes should return true
	os.Setenv("ARYFLOW_YES", "1")
	if !autoYes() {
		t.Error("expected autoYes() to return true when ARYFLOW_YES=1")
	}

	// Without env var, autoYes should return false
	os.Unsetenv("ARYFLOW_YES")
	if autoYes() {
		t.Error("expected autoYes() to return false when ARYFLOW_YES is unset")
	}
}

func TestPromptAutoAccepts(t *testing.T) {
	orig := os.Getenv("ARYFLOW_YES")
	defer os.Setenv("ARYFLOW_YES", orig)

	os.Setenv("ARYFLOW_YES", "1")
	if !Prompt("Test?") {
		t.Error("expected Prompt to return true when ARYFLOW_YES=1")
	}
}

func TestPromptDefaultNoAutoAccepts(t *testing.T) {
	orig := os.Getenv("ARYFLOW_YES")
	defer os.Setenv("ARYFLOW_YES", orig)

	os.Setenv("ARYFLOW_YES", "1")
	if !PromptDefaultNo("Test?") {
		t.Error("expected PromptDefaultNo to return true when ARYFLOW_YES=1")
	}
}
