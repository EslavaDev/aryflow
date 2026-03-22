# AryFlow Workflow Rules

> These rules are MANDATORY for all users of this project. They are enforced by hooks, skills, and this file.

## Update TODO.md After EVERY Wave (non-negotiable)

After completing a wave, the FIRST action is updating TODO.md:
1. Mark completed tasks as `[x]`
2. Commit the TODO.md update
3. Save progress to engram
4. THEN proceed to next wave

The statusline reads TODO.md for progress display. Stale TODO = wrong status for the whole team.

## Mandatory Steps (never skip)

Only brainstorming is optional. Everything else is mandatory:

1. `/brainstorm` — OPTIONAL but recommended
2. `/spec-it {feature}` — writes SPEC.md + TODO.md
3. **Auto-review SPEC.md** with `superpowers:requesting-code-review` — BEFORE showing to user
4. **Review TODO.md** for wave dependencies — BEFORE showing to user
5. **Save to engram** — `mem_save` spec + tasks. If engram unavailable, WARN user loudly
6. `/execute-spec {NNN}-{feature}` — waves with engram session
7. **Each subagent saves** work summary to engram (`{project}/{change}/wave-{N}/agent-{task}`)
8. **Each wave saves** progress to engram (`{project}/{change}/progress`)
9. **Launch post-spec-docs agent** after all waves complete
10. `superpowers:verification-before-completion`
11. `/simplify`
12. `superpowers:finishing-a-development-branch`
13. `/commit` → `/pr`

## Engram is Mandatory

If engram MCP tools are available, you MUST use them:
- `mem_session_start` at session start
- `mem_save` for specs, tasks, progress, knowledge
- `mem_session_summary` + `mem_session_end` at session end
- If engram returns empty → try claude-mem as fallback

If engram is NOT available, warn: "Running in degraded mode. Run `aryflow setup` to install."

## Knowledge Lifecycle

- Mark entries with `[ACTIVE]` or `[DEPRECATED]` (always in English)
- Recency wins: most recent `[ACTIVE]` entry takes precedence
- Never delete `[ACTIVE]` entries
- Run knowledge-gc agent at milestone boundaries

## Wave Execution

- Waves form a dependency graph (not a sequence) — parallel waves when possible
- Default isolation: branch (same feature branch for all agents)
- Worktree exception: only when tasks in the same wave may modify the same file (per-wave, not per-task)
- All wave agents use Opus model
- merge-wave agent uses Opus
- post-spec-docs and knowledge-gc use Sonnet

## Topic Key Convention

```
{project}/{change}/{artifact}     — spec lifecycle
{project}/{change}/wave-{N}/agent-{task}  — agent work summaries
{project}/knowledge/{category}    — permanent project knowledge
```
