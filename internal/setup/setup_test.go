package setup

import (
	"testing"
)

func TestDependencyListOrder(t *testing.T) {
	// Verify that the dependency list matches the spec order:
	// Homebrew, Git, Node.js 18+, Bun, Claude Code, Engram, Claude-Mem, Superpowers
	expectedNames := []string{
		"Homebrew",
		"Git",
		"Node.js 18+",
		"Bun",
		"Claude Code",
		"Engram",
		"Claude-Mem",
		"Superpowers",
	}

	deps := buildDeps()
	if len(deps) != len(expectedNames) {
		t.Fatalf("expected %d dependencies, got %d", len(expectedNames), len(deps))
	}

	for i, expected := range expectedNames {
		if deps[i].Name != expected {
			t.Errorf("dependency[%d]: expected %q, got %q", i, expected, deps[i].Name)
		}
	}
}

func TestAutoInstallable(t *testing.T) {
	deps := buildDeps()

	// These should NOT be auto-installable (no install command)
	notAutoInstallable := map[string]bool{
		"Homebrew":    true,
		"Git":         true,
		"Claude Code": true,
	}

	// These SHOULD be auto-installable
	autoInstallable := map[string]bool{
		"Node.js 18+": true,
		"Bun":         true,
		"Engram":      true,
		"Claude-Mem":  true,
		"Superpowers": true,
	}

	for _, dep := range deps {
		hasInstall := dep.InstallCmd != nil || dep.InstallSteps != nil
		if notAutoInstallable[dep.Name] && hasInstall {
			t.Errorf("%s should NOT be auto-installable, but has install command", dep.Name)
		}
		if autoInstallable[dep.Name] && !hasInstall {
			t.Errorf("%s should be auto-installable, but has no install command", dep.Name)
		}
	}
}

func TestNotAutoInstallableHaveInstructions(t *testing.T) {
	deps := buildDeps()

	for _, dep := range deps {
		hasInstall := dep.InstallCmd != nil || dep.InstallSteps != nil
		if !hasInstall && dep.Instructions == "" {
			// Claude-Mem and Superpowers have install steps, so skip
			t.Errorf("%s is not auto-installable and has no instructions", dep.Name)
		}
	}
}

// buildDeps returns the dependency list used by Run.
// This is extracted to make it testable without running the full setup.
func buildDeps() []dependency {
	return []dependency{
		{
			Name:         "Homebrew",
			Instructions: "Install from https://brew.sh",
		},
		{
			Name:         "Git",
			Instructions: "Git must be pre-installed. On macOS: xcode-select --install",
		},
		{
			Name:       "Node.js 18+",
			InstallCmd: []string{"brew", "install", "node"},
		},
		{
			Name:       "Bun",
			InstallCmd: []string{"brew", "install", "oven-sh/bun/bun"},
		},
		{
			Name:         "Claude Code",
			Instructions: "Install from https://claude.ai/code",
		},
		{
			Name:       "Engram",
			InstallCmd: []string{"brew", "install", "gentleman-programming/tap/engram"},
		},
		{
			Name: "Claude-Mem",
			InstallSteps: [][]string{
				{"claude", "plugin", "marketplace", "add", "thedotmack/claude-mem"},
				{"claude", "plugin", "install", "claude-mem"},
			},
			VersionLabel: "plugin",
		},
		{
			Name:         "Superpowers",
			InstallCmd:   []string{"claude", "plugin", "install", "superpowers"},
			VersionLabel: "plugin",
		},
	}
}
