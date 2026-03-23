# CLAUDE.md — AryFlow CLI

## AryFlow Workflow (MANDATORY)

**Read `.claude/rules/aryflow.md` before any implementation task.**

## Project

- **Repo:** github.com/EslavaDev/aryflow
- **Language:** Go 1.22+
- **Binary:** `aryflow`
- **Distribution:** Homebrew via aryflow/homebrew-tap

## Commands

```bash
go build ./cmd/aryflow          # build
go test ./... -v                # test all
make build                     # build with version
make build-all                 # cross-compile (darwin-arm64, darwin-amd64, linux-amd64)
make test                      # run tests
make sync-embedded             # copy skills/agents from kiwi_nexus
```

## Architecture

```
aryflow/
├── cmd/aryflow/main.go         # CLI entry point, flag parsing, subcommand routing
├── internal/
│   ├── checks/checks.go        # Dependency check functions (git, node, bun, claude, engram, etc.)
│   ├── setup/setup.go          # `aryflow setup` — validate/install prerequisites
│   ├── init/init.go            # `aryflow init` — configure project with skills/agents
│   ├── doctor/doctor.go        # `aryflow doctor` — validate project health
│   ├── update/update.go        # `aryflow update` — self-update + project file update
│   └── ui/ui.go                # Terminal output helpers (colors, prompts, spinners)
├── embedded/                    # Files embedded in binary via go:embed
│   ├── embed.go                # Central embed package
│   ├── skills/                 # spec-it, execute-spec, commit, pr
│   ├── agents/                 # merge-wave, post-spec-docs, knowledge-gc
│   ├── rules/                  # aryflow.md
│   └── hooks/                  # statusline, context-monitor, session hooks
└── specifications/             # Spec for this CLI itself
```

## Conventions

- Standard library only — no external Go dependencies (no cobra, no viper)
- `internal/` packages — not importable by external code
- Every package needs tests in `*_test.go`
- Embed files via `//go:embed` in `embedded/embed.go`
- Colors support `ARYFLOW_NO_COLOR=1` env var
- Prompts support `ARYFLOW_YES=1` env var (CI mode)

## Embedded Files Rule (non-negotiable)

`.claude/` is the **source of truth**. `embedded/` is a copy for the binary.

1. **ONLY edit `.claude/`** — never edit `embedded/` directly
2. **After editing `.claude/`, copy to `embedded/`**: `cp .claude/{path} embedded/{path}`
3. **Hook format**: always use `hookSpecificOutput` with `hookEventName` + `additionalContext` — never `systemMessage` (which only shows in UI, not to the agent)
4. **Run `make sync-check`** before committing to verify sync

## Verification

```bash
export PATH="$HOME/.goenv/versions/1.22.4/bin:$PATH"
go build ./cmd/aryflow
go test ./internal/checks/ ./internal/doctor/ ./internal/init/ ./internal/setup/ ./internal/ui/ -v
make sync-check
```
