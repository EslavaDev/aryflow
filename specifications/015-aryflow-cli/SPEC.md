# SPEC — AryFlow CLI

## 1. Overview

### WHY
AryFlow is a methodology combining multiple tools (engram, claude-mem, superpowers, custom skills/agents). Setting it up in a new project requires manually installing prerequisites, copying skill/agent files, and configuring CLAUDE.md. This is error-prone and not scalable. A CLI tool automates the entire setup and maintenance lifecycle.

### WHAT
A single Go binary (`aryflow`) distributed via Homebrew that provides 4 commands:
- `aryflow setup` — validate and install all system-level prerequisites
- `aryflow init` — configure a project with AryFlow skills, agents, and CLAUDE.md
- `aryflow doctor` — validate that a project has everything it needs
- `aryflow update` — self-update the CLI and update project files to latest versions

### HOW
Go binary compiled for macOS (arm64 + amd64) and Linux. Distributed via a Homebrew tap (`aryflow/tap/aryflow`). The binary embeds the latest versions of all skill and agent files using Go's `embed` package. Commands shell out to check/install dependencies (brew, npm, claude CLI).

### Scope
**IN scope:**
- 4 CLI commands (setup, init, doctor, update)
- Homebrew tap + formula
- Embedded skill/agent files
- Version tracking per project (`.aryflow/version`)
- Colored terminal output with status indicators

**OUT of scope:**
- GUI / web interface
- MCP server functionality
- Runtime orchestration (that's Claude Code's job)
- Windows support (macOS + Linux only for v1)

---

## 2. Architecture

```
aryflow/                          ← new repo: github.com/EslavaDev/aryflow
├── cmd/
│   └── aryflow/
│       └── main.go               ← entry point
├── internal/
│   ├── setup/
│   │   └── setup.go              ← setup command logic
│   ├── init/
│   │   └── init.go               ← init command logic
│   ├── doctor/
│   │   └── doctor.go             ← doctor command logic
│   ├── update/
│   │   └── update.go             ← update command logic
│   ├── checks/
│   │   └── checks.go             ← shared dependency check functions
│   └── ui/
│       └── ui.go                 ← terminal output helpers (colors, spinners, checks)
├── embedded/
│   ├── skills/
│   │   ├── spec-it/SKILL.md
│   │   ├── execute-spec/SKILL.md
│   │   ├── commit/SKILL.md
│   │   └── pr/SKILL.md
│   └── agents/
│       ├── merge-wave.md
│       ├── post-spec-docs.md
│       └── knowledge-gc.md
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

> **Note:** The Homebrew formula lives in `aryflow/homebrew-tap`, not in this repo.

This is a **standalone repo** — not inside kiwi_nexus. The CLI is project-agnostic.

---

## 3. CLI Commands

### Global Flags
- `--verbose` — show stdout/stderr of shelled-out commands
- `--yes` — auto-accept all prompts (alias for ARYFLOW_YES=1)
- `--version` — show CLI version

### 3.1 `aryflow setup`

Validates and installs all system-level prerequisites. Run once per machine.

**Checks (in order):**

| # | Dependency | Check command | Install command | Required |
|---|-----------|---------------|-----------------|----------|
| 0 | Homebrew | `brew --version` | Show install instructions: https://brew.sh | Yes (macOS) / Optional (Linux) |
| 1 | Git | `git --version` | — (must be pre-installed) | Yes |
| 2 | Node.js 18+ | `node --version` → parse semver ≥ 18 | `brew install node` | Yes |
| 3 | Bun | `bun --version` | `brew install oven-sh/bun/bun` (primary) or `curl -fsSL https://bun.sh/install \| bash` (fallback with user consent) | Yes |
| 4 | Claude Code | `claude --version` | Prompt user to install from claude.ai/code | Yes |
| 5 | Engram | `engram --version` | `brew install gentleman-programming/tap/engram` | Yes |
| 6 | Claude-Mem | Check if plugin installed via `claude plugin list \| grep claude-mem` | `claude plugin marketplace add thedotmack/claude-mem && claude plugin install claude-mem` | Yes |
| 7 | Superpowers | Check if plugin installed via `claude plugin list \| grep superpowers` | `claude plugin install superpowers` | Yes |

**Behavior:**
- For each dependency: show check status (pass/fail/warning)
- If a dependency is missing and auto-installable: ask user "Install {dep}? [Y/n]"
- If a dependency is missing and NOT auto-installable: show instructions
- At the end: summary of what's installed and what's missing
- Exit code 0 if all required deps are present, 1 if any missing

**Output example:**
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

### 3.2 `aryflow init`

Initializes AryFlow in the current project directory. Run once per project.

**Preconditions:**
- Must be inside a git repository. If not in a git repo, print: "Not inside a git repository. Run `git init` first." and exit with code 1.
- `aryflow setup` must have been run (all prerequisites present)
- If `.aryflow/version` already exists, warn: "Project already initialized (v{X}). Re-initialize? [y/N]". If user accepts, overwrite version. If not, exit.

**Actions:**

1. **Detect project** — read git repo root, derive project name (kebab-case)
2. **Create `.aryflow/` directory** with:
   - `version` file containing the CLI version that initialized
3. **Copy skills** to `.claude/skills/`:
   - `spec-it/SKILL.md`
   - `execute-spec/SKILL.md`
   - `commit/SKILL.md`
   - `pr/SKILL.md`
   - If any already exist: ask "Overwrite {skill}? [y/N]" (default: no)
   - If `.claude/skills/` already contains files not managed by AryFlow, print warning: "Existing skills found. AryFlow will add its skills alongside them."
4. **Copy agents** to `.claude/agents/`:
   - `merge-wave.md`
   - `post-spec-docs.md`
   - `knowledge-gc.md`
   - If any already exist: ask "Overwrite {agent}? [y/N]" (default: no)
5. **Check CLAUDE.md**:
   - If exists: do NOT modify. Print "CLAUDE.md found — skipping (manual configuration)"
   - If missing: create a minimal CLAUDE.md template with:
     - Project name
     - Placeholder sections (Commands, Architecture, Conventions)
     - Note: "Configure this file with your project's conventions"
6. **Copy rules** to `.claude/rules/`:
   - `aryflow.md` — mandatory workflow rules
   - If already exists: ask "Overwrite? [y/N]"
7. **Copy hooks** to `.claude/hooks/`:
   - `aryflow-session-start.sh` — workflow reminder on session start
   - `aryflow-stop.sh` — session end check
   - `aryflow-subagent-stop.sh` — subagent completion check
   - `aryflow-statusline.js` — status bar with project state + context %
   - `aryflow-context-monitor.js` — context warnings + workflow compliance
   - If any already exist: ask "Overwrite? [y/N]"
8. **Merge `.claude/settings.json`**:
   - If exists: merge AryFlow hooks config into existing settings (preserve existing permissions, add hooks)
   - If missing: create with AryFlow hooks config (SessionStart, PostToolUse, SubagentStop, Stop, statusLine)
9. **Create `specifications/` directory** if it doesn't exist
7. **Initialize engram session**:
   - Print: "Run your first spec with: /spec-it {feature-name}"

**Output example:**
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

**Flags:**
- `--force` — overwrite existing skills/agents without asking
- `--skip-claude-md` — skip CLAUDE.md creation even if missing

### 3.3 `aryflow doctor`

Validates that a project is correctly set up for AryFlow. Run anytime to diagnose issues.

**Checks:**

| # | Check | Pass condition | Severity |
|---|-------|---------------|----------|
| 1 | Git repo | Inside a git repository | Error |
| 2 | .aryflow/version | File exists | Error |
| 3 | CLAUDE.md | File exists | Error |
| 4 | spec-it skill | `.claude/skills/spec-it/SKILL.md` exists | Error |
| 5 | execute-spec skill | `.claude/skills/execute-spec/SKILL.md` exists | Error |
| 6 | commit skill | `.claude/skills/commit/SKILL.md` exists | Warning |
| 7 | pr skill | `.claude/skills/pr/SKILL.md` exists | Warning |
| 8 | merge-wave agent | `.claude/agents/merge-wave.md` exists | Error |
| 9 | post-spec-docs agent | `.claude/agents/post-spec-docs.md` exists | Warning |
| 10 | knowledge-gc agent | `.claude/agents/knowledge-gc.md` exists | Warning |
| 11 | specifications/ dir | Directory exists | Warning |
| 12 | Engram installed | `engram --version` succeeds | Error |
| 13 | Claude-Mem installed | Check plugin list | Warning |
| 14 | Superpowers installed | Check plugin list | Warning |
| 16 | AryFlow rules | `.claude/rules/aryflow.md` exists | Error |
| 17 | SessionStart hook | `.claude/hooks/aryflow-session-start.sh` exists | Error |
| 18 | StatusLine hook | `.claude/hooks/aryflow-statusline.js` exists | Warning |
| 19 | Context monitor | `.claude/hooks/aryflow-context-monitor.js` exists | Warning |
| 20 | Settings.json hooks | `.claude/settings.json` contains AryFlow hooks | Warning |
| 21 | Active TODO.md | If spec has wave comments with completed code but unchecked `[ ]` items → warn stale TODO | Warning |
| 15 | Version match | `.aryflow/version` matches CLI version | Warning |

**Behavior:**
- Run all checks, don't stop on first failure
- Show pass/fail/warning for each
- At the end: summary count (errors, warnings, passed)
- Exit code 0 if no errors, 1 if any errors
- Suggest fix for each failure: "Run `aryflow init` to fix" or "Run `aryflow setup` to install"

**Output example:**
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

### 3.4 `aryflow update`

Self-updates the CLI and updates project files to the latest version.

**Two modes:**

**3.4.1 CLI self-update:**
```
aryflow update --self
```
Detects installation method:
- If installed via Homebrew: runs `brew upgrade aryflow`
- If installed manually (go install, direct binary): downloads latest from GitHub releases

Detection: check if the binary path contains "Cellar" or "homebrew" -> Homebrew install.

- If already latest: print "Already up to date"

**3.4.2 Project update (default):**
```
aryflow update
```
- If `.aryflow/version` is missing, suggest `aryflow init` instead of updating.
- Compare `.aryflow/version` with CLI version
- If different:
  - Show diff summary: which skills/agents changed
  - Ask: "Update project files to v{new}? [Y/n]"
  - If yes: overwrite skills and agents with embedded versions
  - If a file differs from both the old embedded version and the new embedded version, warn: "Locally modified — update will overwrite your changes. Proceed? [y/N]"
  - Update `.aryflow/version` to new version
  - Print changelog if available
- If same: "Project is up to date"

**Flags:**
- `--self` — update the CLI binary itself
- `--force` — update project files without asking
- `--dry-run` — show what would change without applying

**Output example:**
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

---

## 4. Embedded Files

The Go binary embeds all skill and agent files using `//go:embed`:

```go
//go:embed embedded/skills/spec-it/SKILL.md
var specItSkill []byte

//go:embed embedded/skills/execute-spec/SKILL.md
var executeSpecSkill []byte

// ... etc
```

These are the canonical, latest versions of each file. When `aryflow init` or `aryflow update` runs, it copies these embedded files to the project.

The embedded files are the SAME files currently in `.claude/skills/` and `.claude/agents/` in kiwi_nexus. To update them:
1. Edit the skill/agent in any AryFlow project
2. Copy the updated file to `aryflow/embedded/`
3. Rebuild the CLI
4. Publish new version

---

## 5. Homebrew Distribution

### Tap repository: `aryflow/homebrew-tap`

```
homebrew-tap/
└── Formula/
    └── aryflow.rb
```

### Formula (`aryflow.rb`):

```ruby
class Aryflow < Formula
  desc "AryFlow — AI development workflow CLI"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow-darwin-arm64.tar.gz"
      sha256 "TBD"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow-darwin-amd64.tar.gz"
      sha256 "TBD"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow-linux-amd64.tar.gz"
    sha256 "TBD"
  end

  def install
    bin.install "aryflow"
  end

  test do
    system "#{bin}/aryflow", "--version"
  end
end
```

### Install flow:
```bash
brew tap aryflow/tap
brew install aryflow
```

---

## 6. Testing Spec

### Unit tests (Go):
```
internal/checks/checks_test.go    — test dependency detection functions
internal/setup/setup_test.go      — test setup flow (mock exec)
internal/init/init_test.go        — test file copy, directory creation
internal/doctor/doctor_test.go    — test all check functions
internal/update/update_test.go    — test version comparison, file diff
```

**Key test cases:**
- `checks.CheckGit()` — returns version when git available, error when not
- `checks.CheckNode()` — parses semver, rejects < 18
- `init.CopySkills()` — creates directories, copies files, respects --force
- `init.CopySkills()` — does NOT overwrite without --force
- `doctor.RunAll()` — reports correct pass/fail/warning counts
- `update.CompareVersions()` — semver comparison (0.1.0 < 0.2.0 < 1.0.0)
- `update.DiffFiles()` — detects changed, added, removed files

### Integration tests:
- `aryflow setup` in a clean environment — verify all checks run
- `aryflow init` in a git repo — verify all files created
- `aryflow doctor` after init — all checks pass
- `aryflow update` with version mismatch — files updated

### Build/release tests:
- Cross-compile for darwin-arm64, darwin-amd64, linux-amd64
- Verify binary runs on each platform
- Verify Homebrew formula installs correctly

---

## 7. Environment Variables

None required. The CLI reads everything from the filesystem and shelled commands.

Optional:
- `ARYFLOW_NO_COLOR=1` — disable colored output
- `ARYFLOW_YES=1` — auto-accept all prompts (CI mode)

---

## 8. Open Questions

*All resolved:*

1. ~~GitHub org~~ → **EslavaDev/aryflow**
2. ~~Homebrew tap~~ → **aryflow/tap** (`aryflow/homebrew-tap` repo)
3. ~~Auto-install~~ → **Ask for confirmation** (CI mode with `ARYFLOW_YES=1` auto-accepts)
4. ~~Repo location~~ → **Separate repo** (project-agnostic CLI)
