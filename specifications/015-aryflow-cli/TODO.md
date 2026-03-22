# TODO — AryFlow CLI

> **Repo:** github.com/EslavaDev/aryflow (separate from kiwi_nexus)
> **Language:** Go
> **Distribution:** Homebrew via aryflow/homebrew-tap
>
> **RULE: After completing a wave, mark tasks `[x]` IMMEDIATELY before doing anything else. The statusline and engram progress depend on this.**

## Wave 1 — Project scaffold + shared utilities <!-- independent, branch -->
- [x] Create repo structure: cmd/aryflow/main.go, internal/*, embedded/*, go.mod
- [x] Implement internal/ui/ui.go — colored output helpers (success, error, warning, spinner, prompt)
- [x] Implement internal/checks/checks.go — dependency check functions (CheckGit, CheckNode, CheckBun, CheckClaude, CheckEngram, CheckClaudeMem, CheckSuperpowers)
- [x] Copy current skill/agent files to embedded/ directory
- [x] Implement global flag handling (--verbose, --yes, --version) in main.go
- [x] Implement ARYFLOW_NO_COLOR and ARYFLOW_YES env var support in ui package
- [x] Write unit tests for checks and ui packages

## Wave 2 — Setup + Init commands <!-- depends:wave1, branch -->
- [x] Implement internal/setup/setup.go — orchestrate all checks, prompt for installs, execute install commands
- [x] Implement internal/init/init.go — detect project, copy embedded files, create .aryflow/version, create specifications/, handle CLAUDE.md
- [x] Implement cmd/aryflow/main.go — CLI entry point with cobra or flag-based routing (setup, init, doctor, update, --version)
- [x] Implement --force and --skip-claude-md flags for init command
- [x] Handle edge cases: not in git repo, pre-existing .claude/, double init
- [x] Write unit tests for setup and init packages

## Wave 3 — Doctor + Update commands <!-- depends:wave1, branch -->
(runs PARALLEL with Wave 2 — both depend only on Wave 1)
- [x] Implement internal/doctor/doctor.go — run all project checks, report pass/fail/warning, suggest fixes
- [x] Implement internal/update/update.go — version comparison, file diff, self-update via brew, project file update with embedded versions
- [x] Implement --dry-run flag for update command
- [x] Implement --self flag with install method detection (brew vs manual)
- [x] Write unit tests for doctor and update packages

## Wave 4 — Integration + Build <!-- depends:wave2,wave3, branch -->
- [x] Wire all commands in main.go, test full CLI flow end-to-end
- [x] Create Makefile with targets: build, test, build-all (cross-compile darwin-arm64, darwin-amd64, linux-amd64)
- [x] Create .goreleaser.yml or equivalent for automated release builds
- [x] Write integration tests: setup in clean env, init in git repo, doctor after init, update with version mismatch
- [x] Create `make sync-embedded` target to copy skills/agents from canonical source
- [x] Write README.md with install instructions, usage, examples

## Wave 5 — Hooks, Rules, Settings <!-- depends:wave1, branch -->
(runs PARALLEL with Wave 2 and 3 — depends only on Wave 1)
- [x] Copy hooks to embedded/: aryflow-session-start.sh, aryflow-stop.sh, aryflow-subagent-stop.sh, aryflow-statusline.js, aryflow-context-monitor.js
- [x] Copy rules to embedded/: aryflow.md
- [x] Create embedded/settings-template.json with AryFlow hooks config
- [x] Update init command to copy hooks, rules, and merge settings.json
- [x] Update doctor to check for hooks and rules files

## Wave 6 — Homebrew + Release <!-- depends:wave4,wave5, branch -->
- [x] Create aryflow/homebrew-tap repo with Formula/aryflow.rb
- [x] Create GitHub Actions workflow: tag → build → release → update formula
- [ ] Tag v0.1.0, verify `brew tap EslavaDev/aryflow && brew install aryflow` works
- [ ] Test full flow: setup → init → doctor → update on a clean project
