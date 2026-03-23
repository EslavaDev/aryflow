---
name: spec-it
description: Write a formal specification for a feature so context can be reset and implementation can continue without losing details. Use when asked to spec, specify, or write a specification for a feature.
argument-hint: <feature description>
allowed-tools: Read, Glob, Grep, Write, Bash
model: claude-opus-4-6
---

Write a formal specification for: **$ARGUMENTS**

## MANDATORY STEPS — DO NOT SKIP ANY

```
Step 0: Branch setup
Step 1: Determine spec number
Step 2: Read context (CLAUDE.md)
Step 3: Load engram knowledge ← MANDATORY if engram available
→ Write SPEC.md
→ Review SPEC.md with superpowers ← MANDATORY, auto-run before showing to user
→ Fix issues from review
→ Present to user → user approves
→ Write TODO.md
→ Review TODO.md ← MANDATORY
→ Fix issues
→ Present to user → user approves
→ Save to engram ← MANDATORY if engram available
```

**If you skip any mandatory step, the spec is INVALID.** Every step produces output that the next step depends on.

---

## Step 0 — Branch setup

Before doing anything else, ask the user how they want to work on this feature:

**where `{NNN}` is a zero-padded 3-digit sequential number**

1. **New branch** — create a git branch from the current branch (e.g. `feat/{NNN}-{feature-name}` `)
2. **Worktree** — create a git worktree for isolated work
3. **Current branch** — stay on the current branch, no branching

Wait for the user's answer before proceeding. Then:
- If **new branch**: run `git checkout -b feat/{feature-name}` (use the kebab-case feature name)
- If **worktree**: run `git worktree add .claude/worktrees/feat-{NNN}-{feature-name} -b feat/{NNN}-{feature-name}`
- If **current branch**: do nothing, continue

## Step 1 — Determine spec number

Before anything else, list the existing `specifications/` directory to find the highest spec number:

```bash
ls specifications/
```

If the directory doesn't exist, create it and start with `001`.

Specs follow the naming convention `{NNN}-{feature-name}` (e.g. `001-voice-audit`, `015-credit-disputes`).

Parse the highest existing number and increment by 1 to get the next number.

## Step 2 — Read context

Read these files first — they define what's already built, the target architecture, and non-negotiable conventions:

| File | Primary location | Fallback | If neither exists |
|------|-----------------|----------|-------------------|
| Project instructions | `./CLAUDE.md` | — | **Stop** — cannot proceed without CLAUDE.md |

## Step 3 — Load project knowledge from Engram

Before writing the spec, load accumulated project knowledge:

1. Run `mem_context("{project}")` to load all project knowledge (patterns, decisions, conventions, integrations, bugs)
2. Run `mem_search("{project}/{change}/explore")` to check if a brainstorming exploration exists for this change
3. If exploration exists, use `mem_get_observation(id)` to read the full content
4. Incorporate relevant knowledge into the spec — don't re-discover what previous specs already learned

> **Project detection:** `{project}` is auto-detected from the git repo root directory name, converted to lowercase kebab-case (e.g., `kiwi_nexus` → `kiwi-nexus`, `my-saas-app` → `my-saas-app`).

> `{change}` is derived from the spec folder name (e.g., `agent-registry-phase1`).

**Dual memory read (always both, NOT fallback):**
After loading from engram, ALSO search claude-mem for historical context:
- `claude-mem search(query: "{change}", project: "{project}")` — finds prior work, discussions, decisions from past sessions
- This is NOT a fallback — always query both systems. Engram has structured knowledge, claude-mem has historical narrative.

If engram tools are not available, **warn the user**: "Engram not available. Run `aryflow setup` to install. Continuing without persistent memory — knowledge will not be saved." Proceed with local context only, but flag this as degraded mode.

## Output

Create `specifications/{NNN}-{feature-name}/SPEC.md` where:
- `{NNN}` is the next sequential 3-digit number (determined in Step 1)
- `{feature-name}` is a kebab-case name resembling the feature

Example: `specifications/002-payment-reminders/SPEC.md`

The spec MUST be a single `SPEC.md` file — no sub-files.

---

## SPEC.md structure

### 1. Overview
- **WHY** — business/product reason for this feature
- **WHAT** — what it does from the user's perspective
- **HOW** — high-level technical approach
- Scope: what's explicitly IN and OUT of scope

### 2. Architecture

Which layers are touched and how they connect. Read the project structure from CLAUDE.md to determine what layers exist (backend, frontend, shared packages, mobile, etc.).

### 3. Backend spec

Read the backend structure from CLAUDE.md and the codebase. Document:
- **File structure** — where new files go, following existing patterns
- **Functions/methods** — signatures with input/output types and business rules
- **Endpoints** — method, path, request/response shapes, status codes
- **Validation** — input validation and output shaping
- **Error cases** — domain exceptions (not HTTP errors)

Follow the layer separation rules defined in CLAUDE.md. If no rules exist, enforce:
- Route/controller layer: HTTP only
- Service/logic layer: business logic only
- Data/repository layer: database only

### 4. Frontend spec

Read the frontend structure from CLAUDE.md and the codebase. Document:
- **File structure** — where new files go, following existing patterns
- **State management** — how state is managed (read from codebase, don't assume)
- **Data fetching** — how API calls are made (read from codebase, don't assume)
- **Component tree** — page components and their props
- **Types** — TypeScript interfaces or equivalent
- **Routing** — where to add new routes

Follow the patterns already established in the codebase. Do not introduce new patterns without explicit justification.

### 5. Shared packages (if applicable)

If the project uses shared packages (monorepo), document new types and API wrappers needed. Skip this section for single-app projects.

### 6. Testing spec

Read the testing setup from CLAUDE.md and the codebase. Document:
- **Backend tests** — test file locations and framework (from CLAUDE.md)
- **Frontend tests** — test file locations and framework (from CLAUDE.md)
- **E2E tests** — if applicable, locations and framework (from CLAUDE.md)

Provide example test cases for the key behaviors specified in this spec.

### 7. Database changes (if applicable)

List any new tables, columns, or migrations:
- Table name, columns, types, constraints
- Indexes needed
- Migration approach (read from CLAUDE.md — Alembic, Prisma, Knex, raw SQL, etc.)

### 8. Environment variables (if applicable)

List any new env vars needed. Reference the project's .env.example for naming conventions.

### 9. Open questions

List anything ambiguous that must be resolved before implementation starts. **Never guess — surface unknowns explicitly.**

---

## Guidelines

1. **Be thorough** — the spec must be complete enough to implement without asking questions
2. **Fail loud** — every error case must have a specific exception, never silent failures
3. **No `any`** — if TypeScript, use `unknown` + type guard at boundaries
4. **Respect existing patterns** — follow the layer rules, conventions, and patterns already in the codebase and CLAUDE.md
5. **DRY + KISS + SOLID** — always
6. **NEVER assume** — if something is ambiguous, add it to Open Questions and ask before finishing
7. **Follow existing patterns** — don't introduce new frameworks or patterns unless explicitly justified
8. **Keep models in sync** — when adding fields, specify changes across all affected layers
9. **Leverage project knowledge** — check engram for existing patterns and decisions before inventing new ones

---

## After writing SPEC.md — Review (MANDATORY, DO NOT SKIP)

**IMMEDIATELY** after writing SPEC.md, invoke `superpowers:requesting-code-review` to review it. Do NOT present the spec to the user first. Do NOT ask for approval first. Review THEN present.

Review checklist:
1. **Completeness**: are all layers covered?
2. **Consistency**: do types match between layers?
3. **Gaps**: missing error cases, edge cases, undefined behaviors?
4. **Conventions**: does the spec follow CLAUDE.md and project knowledge?

If issues are found → fix them in SPEC.md before showing to the user.

Then present the reviewed spec to the user for approval.

---

## After the spec is approved — Write TODO.md

Create `specifications/{NNN}-{feature-name}/TODO.md` (same directory as the SPEC.md) with implementation tasks organized for wave-based parallel execution.

### TODO.md structure

The TODO must be **project-agnostic** — derive the structure from the SPEC.md and CLAUDE.md, not from a hardcoded template. Read the project structure to determine what layers exist (backend, frontend, shared packages, etc.).

Each wave header MUST include:
1. **Description** — what this wave accomplishes
2. **Dependency** — which wave it depends on (or `independent`)
3. **Isolation** — `branch` (default) or `worktree` (exception)

> **Isolation is per-wave, not per-task.** Either the entire wave runs on the same branch, or the entire wave uses worktrees. You cannot mix isolation strategies within a wave.

### Format

```markdown
# TODO — {Feature Name}

> **RULE: After completing a wave, mark tasks `[x]` IMMEDIATELY before doing anything else.**

## Wave 1 — {description} <!-- independent, branch -->
- [ ] {task description}
- [ ] {task description}
- [ ] {task description}

## Wave 2 — {description} <!-- depends:wave1, branch -->
- [ ] {task description}
- [ ] {task description}

## Wave N — {description} <!-- depends:waveN-1, branch -->
- [ ] {task description}
```

There is NO fixed number of waves. Create as many as needed based on actual dependencies. A small feature might have 2 waves; a large feature might have 10+.

### How to assign waves

**Wave assignment rules:**
- Tasks that touch DIFFERENT files → same wave (parallel)
- Tasks that depend on another task's output → later wave
- Max 4-6 tasks per wave (system resource limit)
- Waves are numbered sequentially: Wave 1, Wave 2, ... Wave N

**Waves CAN run in parallel with each other.** Dependencies are between specific waves, not strictly sequential. If Wave 3 depends on Wave 1 (but NOT Wave 2), then Wave 2 and Wave 3 can run simultaneously after Wave 1 completes.

Example — monorepo with independent apps:
```
Wave 1 — Backend models        <!-- independent, branch -->
Wave 2 — Backend services      <!-- depends:wave1, branch -->
Wave 3 — Frontend types        <!-- depends:wave1, branch -->  ← runs PARALLEL with Wave 2
Wave 4 — Frontend components   <!-- depends:wave3, branch -->  ← runs PARALLEL with Wave 2
Wave 5 — Integration tests     <!-- depends:wave2,wave4, branch -->  ← waits for both
```

execute-spec reads the `depends:` hints and builds a dependency graph. Waves without mutual dependencies launch simultaneously.

**Wave isolation (branch vs worktree):**
- **Default: `branch`** — all agents in the wave work on the same feature branch. This is the normal workflow.
- **Exception: `worktree`** — ONLY when tasks in the wave might modify the SAME file. The entire wave switches to worktree mode. This is rare — if waves are grouped by file independence, you almost never need worktrees.

**How to decide isolation:** Look at the files each task will touch. If any two tasks in the wave could edit the same file → mark the wave as `worktree`. Otherwise → `branch`.

### Example (from a specific project — adapt structure to YOUR project)

```markdown
# TODO — Agent Registry

## Wave 1 — Core models (3 agents, different files) <!-- independent, branch -->
- [ ] Create agents table + migration
- [ ] Create tools table + migration
- [ ] Create knowledge_bases table + migration

## Wave 2 — Relations + config <!-- depends:wave1, branch -->
- [ ] Create junction tables agent_tools + agent_knowledge_bases
- [ ] Create agent_sync_state + service_registry tables

## Wave 3 — Service layer <!-- depends:wave2, worktree -->
(worktree because both tasks may touch shared service utilities)
- [ ] Implement agent service + repository
- [ ] Implement ElevenLabs runtime (sync + execute)

## Wave 4 — API + UI <!-- depends:wave3, branch -->
- [ ] Create API routes /api/agents/*
- [ ] Create CMS agent list + detail pages
- [ ] Create CMS KB management + sync UI

## Wave 5 — Migration + Tests <!-- depends:wave4, branch -->
- [ ] Migrate 31 DPD agents from env vars to registry
- [ ] Write integration tests for full flow
```

---

## Save to Engram

After the spec and TODO are approved, persist to engram for agent access:

1. `mem_save(topic: "{project}/{change}/spec", content: <SPEC.md content>, project: "{project}")`
2. `mem_save(topic: "{project}/{change}/tasks", content: <TODO.md content>, project: "{project}")`
3. If this spec makes architectural decisions not already in project knowledge:
   - `mem_save(topic: "{project}/knowledge/decisions", content: "[ACTIVE] {date} — {decision description}", project: "{project}")`
4. If this spec establishes new conventions:
   - `mem_save(topic: "{project}/knowledge/conventions", content: "[ACTIVE] {date} — {convention description}", project: "{project}")`

Local files (`specifications/{NNN}-{name}/SPEC.md` and `TODO.md`) remain the human-readable source. Engram entries are the agent-readable source.

If engram is not available, **warn the user**: "Engram not available — spec and tasks were NOT saved to persistent memory. Future agents and sessions will not have access to this spec via engram. Run `aryflow setup` to install."

---

## Review TODO.md (MANDATORY, DO NOT SKIP)

After writing TODO.md, review it before presenting to user:
1. Are wave dependencies correct? Can parallel waves be optimized?
2. Are ALL tasks from the spec represented?
3. Is the wave ordering logical?
4. Is isolation (branch/worktree) correctly assigned per wave?

Fix any issues found, then present to user for approval.
