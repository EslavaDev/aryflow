// Package checks provides dependency check functions for AryFlow CLI.
package checks

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// runCommand executes a command and returns its trimmed stdout.
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s not found or failed: %w", name, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CheckHomebrew checks if Homebrew is installed and returns its version.
func CheckHomebrew() (string, error) {
	out, err := runCommand("brew", "--version")
	if err != nil {
		return "", fmt.Errorf("Homebrew not found. Install from https://brew.sh")
	}
	// "Homebrew 4.2.0" — extract first line
	lines := strings.Split(out, "\n")
	if len(lines) > 0 {
		return strings.TrimPrefix(lines[0], "Homebrew "), nil
	}
	return out, nil
}

// CheckGit checks if Git is installed and returns its version.
func CheckGit() (string, error) {
	out, err := runCommand("git", "--version")
	if err != nil {
		return "", fmt.Errorf("Git not found. Please install Git")
	}
	// "git version 2.43.0" → "2.43.0"
	version := strings.TrimPrefix(out, "git version ")
	return strings.TrimSpace(version), nil
}

// CheckNode checks if Node.js is installed and validates >= 18.
func CheckNode() (string, error) {
	out, err := runCommand("node", "--version")
	if err != nil {
		return "", fmt.Errorf("Node.js not found")
	}
	// "v22.1.0" → "22.1.0"
	version := strings.TrimPrefix(out, "v")
	version = strings.TrimSpace(version)

	// Parse major version
	parts := strings.Split(version, ".")
	if len(parts) < 1 {
		return version, fmt.Errorf("cannot parse Node.js version: %s", out)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return version, fmt.Errorf("cannot parse Node.js major version: %s", parts[0])
	}
	if major < 18 {
		return version, fmt.Errorf("Node.js %s found but >= 18 required", version)
	}

	return version, nil
}

// CheckBun checks if Bun is installed and returns its version.
func CheckBun() (string, error) {
	out, err := runCommand("bun", "--version")
	if err != nil {
		return "", fmt.Errorf("Bun not found")
	}
	return strings.TrimSpace(out), nil
}

// CheckClaude checks if Claude Code CLI is installed and returns its version.
func CheckClaude() (string, error) {
	out, err := runCommand("claude", "--version")
	if err != nil {
		return "", fmt.Errorf("Claude Code not found. Install from https://claude.ai/code")
	}
	return strings.TrimSpace(out), nil
}

// CheckEngram checks if Engram is installed and returns its version.
func CheckEngram() (string, error) {
	out, err := runCommand("engram", "--version")
	if err != nil {
		return "", fmt.Errorf("Engram not found. Install with: brew install gentleman-programming/tap/engram")
	}
	return strings.TrimSpace(out), nil
}

// CheckClaudeMem checks if the claude-mem plugin is installed.
func CheckClaudeMem() (string, error) {
	out, err := runCommand("claude", "plugin", "list")
	if err != nil {
		return "", fmt.Errorf("cannot list Claude plugins: %w", err)
	}
	if !strings.Contains(out, "claude-mem") {
		return "", fmt.Errorf("claude-mem plugin not found. Install with: claude plugin marketplace add thedotmack/claude-mem && claude plugin install claude-mem")
	}
	return "installed", nil
}

// CheckSuperpowers checks if the superpowers plugin is installed.
func CheckSuperpowers() (string, error) {
	out, err := runCommand("claude", "plugin", "list")
	if err != nil {
		return "", fmt.Errorf("cannot list Claude plugins: %w", err)
	}
	if !strings.Contains(out, "superpowers") {
		return "", fmt.Errorf("superpowers plugin not found. Install with: claude plugin install superpowers")
	}
	return "installed", nil
}
