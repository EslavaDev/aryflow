# AryFlow

**AryFlow by Alejandro Eslava**

AI-powered spec-driven development CLI. AryFlow automates the setup and maintenance of a structured AI development workflow combining persistent memory, custom skills, agents, and Claude Code.

## Install

```bash
brew tap aryflow/tap && brew install aryflow
```

## Quick Start

```bash
aryflow setup   # Install and validate all prerequisites
aryflow init    # Initialize AryFlow in your project
```

## Commands

| Command | Description |
|---------|-------------|
| `aryflow setup` | Validate and install system-level prerequisites (Git, Node.js, Bun, Claude Code, Engram, Claude-Mem, Superpowers) |
| `aryflow init` | Initialize AryFlow in the current project -- copies skills, agents, and creates config files |
| `aryflow doctor` | Check project health and diagnose configuration issues |
| `aryflow update` | Update CLI binary (`--self`) or project files to the latest embedded versions |

### Global Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Show detailed command output |
| `--yes`, `-y` | Auto-accept all prompts (CI mode) |
| `--version` | Show CLI version |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `ARYFLOW_NO_COLOR=1` | Disable colored terminal output |
| `ARYFLOW_YES=1` | Auto-accept all prompts |

## What AryFlow Sets Up

AryFlow configures your project with:

- **Skills**: `spec-it`, `execute-spec`, `commit`, `pr` -- structured development workflows for Claude Code
- **Agents**: `merge-wave`, `post-spec-docs`, `knowledge-gc` -- automated maintenance tasks
- **Rules and hooks**: Session management and workflow automation
- **Project structure**: `specifications/` directory, `.aryflow/` config, `CLAUDE.md` template

## Development

```bash
make build       # Build for current platform
make test        # Run all tests
make build-all   # Cross-compile for macOS (arm64/amd64) and Linux
make clean       # Remove build artifacts
```

## Acknowledgments

AryFlow stands on the shoulders of excellent open-source tools:

- **[Engram](https://github.com/gentleman-programming/engram)** by Gentleman Programming -- structured persistent memory for AI development workflows
- **[Claude-Mem](https://github.com/thedotmack/claude-mem)** by thedotmack -- automatic passive memory with semantic search
- **[Superpowers](https://github.com/obra/superpowers)** by obra -- development skills framework
- **[Claude Code](https://claude.ai/code)** by Anthropic -- the AI coding assistant that powers everything

## License

MIT
