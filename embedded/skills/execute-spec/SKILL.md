---
name: execute-spec
description: Execute a feature spec phase by phase with wave-based parallel agents following the TODO checklist. Use when asked to implement, execute, or build a spec from specifications/.
argument-hint: <spec-folder-name>
allowed-tools: Read, Glob, Grep, Write, Edit, Bash, Agent
model: claude-opus-4-6
---

Your task is to execute the spec provided in `specifications/$ARGUMENTS/`.

## Before starting

Read these files first:

| File | Location | If missing |
|------|----------|------------|
| Project instructions | `./CLAUDE.md` | **Stop** — cannot proceed |
| Architecture | `docs/ARCHITECTURE.md` | Skip — proceed without |
| Conventions | `docs/CONVENTIONS.md` | Skip — proceed without |
| Structure | `docs/STRUCTURE.md` | Skip — proceed without |

Then read ALL files in `specifications/$ARGUMENTS/` (SPEC.md, TODO.md, and any other `.md` files).

## Branch and Worktree Strategy

**Default: all agents work on the same feature branch.** This is the normal workflow.

Worktree isolation is the exception — only used when two tasks in the SAME wave modify the SAME file (marked with `<!-- worktree -->` in TODO.md). If waves are grouped correctly by file independence, worktrees are almost never needed.

## Load context from Engram

After reading local files, load context from engram. **This step is MANDATORY if engram is available.** If engram tools are not found, warn: "Engram not available. Running in degraded mode — no persistent memory, no resume capability. Run `aryflow setup` to install."

1. `mem_session_start(project: "{project}", description: "Executing $ARGUMENTS")`
2. `mem_search("{project}/$ARGUMENTS/spec")` → `mem_get_observation(id)` — load spec
3. `mem_search("{project}/$ARGUMENTS/progress")` — check for prior progress (enables resume)
4. `mem_context("{project}")` — load project knowledge
5. If engram returns empty for expected knowledge → try claude-mem semantic search as fallback, then save results to engram with correct topic key

> **Project detection:** `{project}` is auto-detected from git repo root directory name in lowercase kebab-case.

### Resume from interruption

If progress exists: read it, cross-reference with TODO.md checkboxes (`[x]` = done), skip completed waves, announce "Resuming from Wave {N}."

---

## Wave-based parallel execution

### Step 1: Dependency analysis

Read TODO.md and build a **dependency graph** from the wave headers. Waves declare their dependencies via `<!-- depends:waveN -->` hints.

Waves are NOT always sequential. If Wave 3 depends on Wave 1 (not Wave 2), launch Wave 2 and Wave 3 in parallel after Wave 1 completes. Build the graph, identify which waves can run simultaneously, and maximize parallelism.

### Step 2: Group into waves

Organize tasks into waves of **4-6 agents maximum**. Present the wave breakdown to the user before starting.

### Step 3: Execute each wave

**1. Launch subagents in parallel** with `run_in_background: true`. Default: no isolation (same branch). Only add `isolation: "worktree"` for tasks marked `<!-- worktree -->` in TODO.md.

| Task Type | Model |
|-----------|-------|
| All wave execution | opus |

**2. Each subagent receives this context:**

```
CONTEXT:
- Project: {project}
- Change: $ARGUMENTS
- Your task: {specific task description from TODO.md}

INSTRUCTIONS:
1. Read the spec from engram: mem_search("{project}/$ARGUMENTS/spec") → mem_get_observation(id)
2. Read project knowledge: mem_context("{project}")
3. Read CLAUDE.md for project conventions
4. Implement your specific task following the spec exactly
5. Run verification commands from CLAUDE.md after implementation.
6. Save to engram ONLY if you discovered something new (pattern, decision, workaround, bug fix, convention):
   mem_save(topic: "{project}/knowledge/{category}", content: "[ACTIVE] {date} — {discovery}", project: "{project}")
7. If something contradicts existing knowledge:
   mem_update(id: {old_id}) → add [DEPRECATED] + supersedes reference
   mem_save new entry as [ACTIVE]

RULES:
- Follow the spec exactly
- Follow all project conventions from CLAUDE.md
- Respect the layer separation rules defined in CLAUDE.md
- Every new function needs tests
- When conflicting knowledge exists, the most recent [ACTIVE] entry wins. Ignore all [DEPRECATED] entries.
- Do NOT save a summary of what you did — the orchestrator tracks progress.
- Do NOT call mem_session_summary or mem_session_end — the orchestrator handles this.
```

**3. Wait for all agents in the wave to complete.**

**4. Merge and review checkpoint:** Launch the merge-wave agent (`.claude/agents/merge-wave.md`) to handle integration, conflict resolution, dependency installation, and verification. After it returns, present its summary. Stop on unresolved conflicts or verification failures.

**5. Persist wave progress:** Update TODO.md (`[x]`), then save minimal progress checkpoint to engram:
```
mem_save(topic: "{project}/$ARGUMENTS/progress",
  content: "Wave {N} complete. Tasks: {list of task slugs}. Next: Wave {N+1} or Done.",
  project: "{project}")
```
Commit: `feat({domain}): wave {N} — {brief description}`.

**6. Proceed to next wave or finish.**

### Step 4: Fallback to sequential execution

If the Agent tool is unavailable or user prefers sequential: execute one phase at a time, update TODO.md, commit after each phase.

---

## Guidelines

1. SPEC.md has feature details, TODO.md has the checklist. Execute wave by wave.
2. Keep TODO.md updated — check items done (`[x]`) before continuing.
3. **FAIL-HARD-FIRST** — for UI errors always notify the user, do not bypass them.
4. DO NOT build the UI — only run verification commands from CLAUDE.md.
5. STRICTLY FOLLOW CLAUDE.md instructions and conventions.

## Layer rules (non-negotiable)

Follow the layer separation rules defined in CLAUDE.md. If no explicit rules exist, enforce:
- Route/controller layer: HTTP only — parse request, call service, return response
- Service/logic layer: business logic only — no web framework imports
- Data/repository layer: database only — no business logic, no HTTP
- Schema/validation layer: input validation and output shaping

## Testing rules (non-negotiable)

Follow the testing rules from CLAUDE.md. Run tests after each wave to verify nothing is broken. Every new function needs tests — use the framework and file locations specified in CLAUDE.md.

## Session end

After all waves complete:

1. Launch **post-spec-docs agent** (`.claude/agents/post-spec-docs.md`) to update CLAUDE.md, docs, and .env.example if needed.
2. Final verification → announce completion.
3. Do NOT call `mem_session_summary()` or `mem_session_end()` — the Stop hook handles session cleanup.

---

## Complete Workflow (end-to-end reference)

1. `/brainstorm "feature"` — explore problem, save to engram
2. `/spec-it {feature}` — generate SPEC.md + TODO.md, auto-review with superpowers:requesting-code-review
3. `/execute-spec {NNN}-{feature}` — waves of parallel agents, merge, verify, commit per wave
4. `superpowers:verification-before-completion` — verify tests, spec match, no regressions
5. `/simplify` — review changed code for reuse, quality, efficiency
6. `post-spec-docs` agent — update CLAUDE.md, docs, .env.example if needed
7. `superpowers:finishing-a-development-branch` — decide: merge to main, create PR, or cleanup
8. `/commit` — conventional commit with Co-Authored-By trailer
9. `/pr` — create pull request with summary

Each step builds on the previous. Knowledge accumulates across the entire lifecycle.

### Maintenance (at milestone boundaries)
- Launch **knowledge-gc agent** (`.claude/agents/knowledge-gc.md`) to clean up [DEPRECATED] entries from engram
