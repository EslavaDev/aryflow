// AryFlow CLI — AI-powered spec-driven development workflow.
package main

import (
	"fmt"
	"os"

	"github.com/EslavaDev/aryflow/internal/doctor"
	aryinit "github.com/EslavaDev/aryflow/internal/init"
	"github.com/EslavaDev/aryflow/internal/setup"
	"github.com/EslavaDev/aryflow/internal/update"
)

// Version is set at build time via -ldflags.
var version = "0.1.0"

// Global flags
var (
	verbose bool
	yes     bool
)

func main() {
	args := os.Args[1:]

	// Parse global flags
	var command string
	var cmdArgs []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--version", "-v":
			fmt.Printf("aryflow v%s\n", version)
			os.Exit(0)
		case "--verbose":
			verbose = true
		case "--yes", "-y":
			yes = true
			os.Setenv("ARYFLOW_YES", "1")
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		default:
			if len(args[i]) > 0 && args[i][0] == '-' {
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", args[i])
				printUsage()
				os.Exit(1)
			}
			// First non-flag arg is the command
			if command == "" {
				command = args[i]
				cmdArgs = args[i+1:]
			}
			// Stop parsing globals once we hit the command
			goto dispatch
		}
	}

dispatch:
	if command == "" {
		printUsage()
		os.Exit(0)
	}

	switch command {
	case "setup":
		if err := setup.Run(verbose); err != nil {
			os.Exit(1)
		}
	case "init":
		force, skipClaudeMD := parseInitFlags(cmdArgs)
		if err := aryinit.Run(force, skipClaudeMD, verbose, version); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "doctor":
		exitCode := doctor.Run(verbose, version)
		os.Exit(exitCode)
	case "update":
		selfUpdate, force, dryRun := parseUpdateFlags(cmdArgs)
		if selfUpdate {
			os.Exit(update.RunSelf(verbose, version))
		}
		os.Exit(update.Run(force, dryRun, verbose, version))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// parseInitFlags parses init-specific flags from the command arguments.
func parseInitFlags(args []string) (force bool, skipClaudeMD bool) {
	for _, arg := range args {
		switch arg {
		case "--force":
			force = true
		case "--skip-claude-md":
			skipClaudeMD = true
		case "--verbose":
			verbose = true
		case "--yes", "-y":
			yes = true
			os.Setenv("ARYFLOW_YES", "1")
		}
	}
	return
}

// parseUpdateFlags parses update-specific flags from the command arguments.
func parseUpdateFlags(args []string) (selfUpdate bool, force bool, dryRun bool) {
	for _, arg := range args {
		switch arg {
		case "--self":
			selfUpdate = true
		case "--force":
			force = true
		case "--dry-run":
			dryRun = true
		case "--verbose":
			verbose = true
		case "--yes", "-y":
			yes = true
			os.Setenv("ARYFLOW_YES", "1")
		}
	}
	return
}

func printUsage() {
	fmt.Printf(`AryFlow v%s — AI-powered spec-driven development CLI

Usage:
  aryflow <command> [flags]

Commands:
  setup     Validate and install system prerequisites
  init      Initialize AryFlow in the current project
  doctor    Check project health and configuration
  update    Update CLI or project files

Global Flags:
  --verbose   Show detailed command output
  --yes, -y   Auto-accept all prompts (CI mode)
  --version   Show CLI version
  --help      Show this help

Init Flags:
  --force          Overwrite existing files without asking
  --skip-claude-md Skip CLAUDE.md creation even if missing

Update Flags:
  --self       Update the CLI binary itself
  --force      Update project files without asking
  --dry-run    Show what would change without applying

Environment Variables:
  ARYFLOW_NO_COLOR=1   Disable colored output
  ARYFLOW_YES=1        Auto-accept all prompts

`, version)
}
