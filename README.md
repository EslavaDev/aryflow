# AryFlow

**AryFlow by Alejandro Eslava**

AI-powered spec-driven development CLI. AryFlow automates the setup and maintenance of a structured AI development workflow combining persistent memory, custom skills, agents, and Claude Code.

A single Go binary that provides four commands to go from zero to a fully configured AI development environment.

## Installation

### Homebrew (recommended)

```bash
brew tap aryflow/tap
brew install aryflow
```

### Go install

```bash
go install github.com/EslavaDev/aryflow/cmd/aryflow@latest
```

### Manual download

Download the latest binary from [GitHub Releases](https://github.com/EslavaDev/aryflow/releases) for your platform:

- `aryflow-darwin-arm64.tar.gz` (macOS Apple Silicon)
- `aryflow-darwin-amd64.tar.gz` (macOS Intel)
- `aryflow-linux-amd64.tar.gz` (Linux x86_64)

## Quick Start

```bash
# 1. Install and validate all prerequisites
aryflow setup

# 2. Navigate to your project and initialize AryFlow
cd my-project
aryflow init

# 3. Verify everything is configured correctly
aryflow doctor

# 4. Start building with specs
# (inside Claude Code)
/spec-it my-feature
```

## Commands

### Global Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Show detailed command output |
| `--yes`, `-y` | Auto-accept all prompts (CI mode) |
| `--version` | Show CLI version |

### `aryflow setup`

Validates and installs all system-level prerequisites. Run once per machine.

Checks for: Git, Node.js 18+, Bun, Claude Code, Engram, Claude-Mem plugin, and Superpowers plugin. Missing dependencies that are auto-installable will prompt for confirmation.

```
AryFlow Setup — Checking prerequisites...

  ✓ Git 2.43.0
  ✓ Node.js 22.1.0
  ✓ Bun 1.2.5
  ✓ Claude Code 1.0.33
  ✗ Engram — not found
    → Install with: brew install gentleman-programming/tap/engram
    Install now? [Y/n] y
    Installing engram... ✓ engram 0.4.2
  ✓ Claude-Mem plugin
  ✓ Superpowers plugin

Setup complete. All prerequisites installed.
```

### `aryflow init`

Initializes AryFlow in the current project directory. Run once per project. Must be inside a git repository.

Copies skills (`spec-it`, `execute-spec`, `commit`, `pr`), agents (`merge-wave`, `post-spec-docs`, `knowledge-gc`), rules, hooks, and configures `.claude/settings.json`. Creates a `CLAUDE.md` template if one doesn't exist, or leaves your existing one untouched.

| Flag | Description |
|------|-------------|
| `--force` | Overwrite existing skills/agents without asking |
| `--skip-claude-md` | Skip CLAUDE.md creation even if missing |

```
AryFlow Init — Setting up project "my-saas-app"

  ✓ Created .aryflow/version (v0.1.0)
  ✓ Copied spec-it skill
  ✓ Copied execute-spec skill
  ✓ Copied commit skill
  ✓ Copied pr skill
  ✓ Copied merge-wave agent
  ✓ Copied post-spec-docs agent
  ✓ Copied knowledge-gc agent
  ✓ CLAUDE.md found — skipping
  ✓ Created specifications/

AryFlow initialized. Start with: /spec-it {feature-name}
```

### `aryflow doctor`

Validates that a project is correctly set up for AryFlow. Run anytime to diagnose issues.

Checks for: git repo, `.aryflow/version`, `CLAUDE.md`, all skills and agents, `specifications/` directory, Engram, Claude-Mem, Superpowers, rules, hooks, settings, and version match. Suggests fix commands for each failure.

```
AryFlow Doctor — Checking project "my-saas-app"

  ✓ Git repository
  ✓ .aryflow/version (v0.1.0)
  ✓ CLAUDE.md
  ✓ spec-it skill
  ✓ execute-spec skill
  ✓ commit skill
  ✗ pr skill — missing
    → Run: aryflow init
  ✓ merge-wave agent
  ✓ post-spec-docs agent
  ✓ knowledge-gc agent
  ✓ specifications/ directory
  ✓ Engram available
  ✓ Claude-Mem plugin
  ✓ Superpowers plugin
  ⚠ Version mismatch — project: v0.1.0, CLI: v0.2.0
    → Run: aryflow update

Results: 13 passed, 1 warning, 1 error
```

### `aryflow update`

Updates project files to the latest embedded versions, or self-updates the CLI binary.

| Flag | Description |
|------|-------------|
| `--self` | Update the CLI binary itself (via Homebrew or GitHub releases) |
| `--force` | Update project files without asking |
| `--dry-run` | Show what would change without applying |

```bash
# Update project files
aryflow update

# Update the CLI binary
aryflow update --self
```

```
AryFlow Update — Checking for updates...

  Current: v0.1.0
  Latest:  v0.2.0

  Changes:
    ~ execute-spec/SKILL.md — updated wave execution logic
    + knowledge-gc.md — new agent (added in v0.2.0)

  Update project files? [Y/n] y
  ✓ Updated execute-spec skill
  ✓ Added knowledge-gc agent
  ✓ Version updated to v0.2.0

Project updated to v0.2.0.
```

## What AryFlow Sets Up

AryFlow configures your project with:

- **Skills**: `spec-it`, `execute-spec`, `commit`, `pr` -- structured development workflows for Claude Code
- **Agents**: `merge-wave`, `post-spec-docs`, `knowledge-gc` -- automated maintenance tasks
- **Rules**: Mandatory workflow rules (`aryflow.md`)
- **Hooks**: Session management, status line, and context monitoring
- **Project structure**: `specifications/` directory, `.aryflow/` config, `CLAUDE.md` template

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ARYFLOW_NO_COLOR=1` | Disable colored terminal output |
| `ARYFLOW_YES=1` | Auto-accept all prompts (same as `--yes`) |

## Development

### Prerequisites

- Go 1.22+

### Building

```bash
make build          # Build for current platform
make build-all      # Cross-compile for macOS (arm64/amd64) and Linux
make clean          # Remove build artifacts
```

### Testing

```bash
make test           # Run all tests
go test ./... -v    # Verbose test output
```

### Project Structure

```
aryflow/
├── cmd/aryflow/main.go         # CLI entry point, flag parsing, subcommand routing
├── internal/
│   ├── checks/                 # Dependency check functions
│   ├── setup/                  # aryflow setup command
│   ├── init/                   # aryflow init command
│   ├── doctor/                 # aryflow doctor command
│   ├── update/                 # aryflow update command
│   └── ui/                     # Terminal output helpers (colors, prompts)
├── embedded/                   # Files embedded in binary via go:embed
│   ├── skills/                 # spec-it, execute-spec, commit, pr
│   ├── agents/                 # merge-wave, post-spec-docs, knowledge-gc
│   ├── rules/                  # aryflow.md
│   └── hooks/                  # statusline, context-monitor, session hooks
└── specifications/             # Specs for this CLI itself
```

### Conventions

- Standard library only -- no external Go dependencies
- `internal/` packages are not importable by external code
- Every package needs tests in `*_test.go`

## Acknowledgments

AryFlow stands on the shoulders of excellent open-source tools:

- **[Engram](https://github.com/gentleman-programming/engram)** by Gentleman Programming -- structured persistent memory for AI development workflows
- **[Claude-Mem](https://github.com/thedotmack/claude-mem)** by thedotmack -- automatic passive memory with semantic search
- **[Superpowers](https://github.com/obra/superpowers)** by obra -- development skills framework
- **[Claude Code](https://claude.ai/code)** by Anthropic -- the AI coding assistant that powers everything

## License

MIT
