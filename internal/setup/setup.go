// Package setup implements the "aryflow setup" command.
package setup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/EslavaDev/aryflow/internal/checks"
	"github.com/EslavaDev/aryflow/internal/ui"
)

// dependency describes a single prerequisite to check and optionally install.
type dependency struct {
	Name         string
	CheckFunc    func() (string, error)
	InstallCmd   []string   // empty = not auto-installable
	InstallSteps [][]string // multi-step install (used instead of InstallCmd when set)
	Instructions string     // shown when not auto-installable
	VersionLabel string     // e.g. "plugin" for things without semver
}

// Run executes the setup command, checking and installing prerequisites.
func Run(verbose bool) error {
	deps := []dependency{
		{
			Name:         "Homebrew",
			CheckFunc:    checks.CheckHomebrew,
			Instructions: "Install from https://brew.sh",
		},
		{
			Name:      "Git",
			CheckFunc: checks.CheckGit,
			Instructions: "Git must be pre-installed. On macOS: xcode-select --install",
		},
		{
			Name:       "Node.js 18+",
			CheckFunc:  checks.CheckNode,
			InstallCmd: []string{"brew", "install", "node"},
		},
		{
			Name:       "Bun",
			CheckFunc:  checks.CheckBun,
			InstallCmd: []string{"brew", "install", "oven-sh/bun/bun"},
		},
		{
			Name:         "Claude Code",
			CheckFunc:    checks.CheckClaude,
			Instructions: "Install from https://claude.ai/code",
		},
		{
			Name:       "Engram",
			CheckFunc:  checks.CheckEngram,
			InstallCmd: []string{"brew", "install", "gentleman-programming/tap/engram"},
		},
		{
			Name:      "Claude-Mem",
			CheckFunc: checks.CheckClaudeMem,
			InstallSteps: [][]string{
				{"claude", "plugin", "marketplace", "add", "thedotmack/claude-mem"},
				{"claude", "plugin", "install", "claude-mem"},
			},
			VersionLabel: "plugin",
		},
		{
			Name:         "Superpowers",
			CheckFunc:    checks.CheckSuperpowers,
			InstallCmd:   []string{"claude", "plugin", "install", "superpowers"},
			VersionLabel: "plugin",
		},
	}

	ui.Header("AryFlow Setup — Checking prerequisites...")

	passed := 0
	failed := 0
	installed := 0

	for _, dep := range deps {
		ver, err := dep.CheckFunc()
		if err == nil {
			// Already installed
			label := ver
			if dep.VersionLabel != "" && ver == "installed" {
				label = dep.VersionLabel
			}
			ui.Success(fmt.Sprintf("%s %s", dep.Name, label))
			passed++
			continue
		}

		// Not installed — try to install or show instructions
		ui.Error(fmt.Sprintf("%s — not found", dep.Name))

		if dep.InstallCmd != nil || dep.InstallSteps != nil {
			// Auto-installable
			if !ui.Prompt(fmt.Sprintf("Install %s?", dep.Name)) {
				ui.Suggestion("Skipped")
				failed++
				continue
			}

			var installErr error
			if dep.InstallSteps != nil {
				installErr = runMultiStep(dep.InstallSteps, verbose)
			} else {
				installErr = runInstall(dep.InstallCmd, verbose)
			}

			if installErr != nil {
				ui.Error(fmt.Sprintf("Failed to install %s: %v", dep.Name, installErr))
				failed++
				continue
			}

			// Re-check after install
			ver, err = dep.CheckFunc()
			if err != nil {
				ui.Error(fmt.Sprintf("Installed but check still fails: %v", err))
				failed++
				continue
			}

			label := ver
			if dep.VersionLabel != "" && ver == "installed" {
				label = dep.VersionLabel
			}
			ui.Success(fmt.Sprintf("Installed %s %s", dep.Name, label))
			installed++
			passed++
		} else if dep.Instructions != "" {
			ui.Suggestion(dep.Instructions)
			failed++
		} else {
			failed++
		}
	}

	// Summary
	fmt.Println()
	if failed == 0 {
		ui.Success(fmt.Sprintf("Setup complete. All %d prerequisites installed.", passed))
		if installed > 0 {
			ui.Info(fmt.Sprintf("(%d newly installed)", installed))
		}
		return nil
	}

	ui.Warning(fmt.Sprintf("%d passed, %d missing", passed, failed))
	return fmt.Errorf("%d prerequisites missing", failed)
}

// runInstall runs a single install command.
func runInstall(args []string, verbose bool) error {
	fmt.Printf("    Installing...")
	cmd := exec.Command(args[0], args[1:]...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println()
		fmt.Printf("    $ %s\n", strings.Join(args, " "))
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println()
		return err
	}
	if !verbose {
		fmt.Println(" done")
	}
	return nil
}

// runMultiStep runs multiple install commands in sequence.
func runMultiStep(steps [][]string, verbose bool) error {
	for _, step := range steps {
		fmt.Printf("    Running: %s\n", strings.Join(step, " "))
		cmd := exec.Command(step[0], step[1:]...)
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("step '%s' failed: %w", strings.Join(step, " "), err)
		}
	}
	return nil
}
