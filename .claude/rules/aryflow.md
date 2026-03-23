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
7. **Each subagent saves** only technical discoveries to engram (NOT work summaries)
8. **Orchestrator saves** wave progress checkpoint to engram (`{project}/{change}/progress`)
9. **Launch post-spec-docs agent** after all waves complete
10. `superpowers:verification-before-completion`
11. `/simplify`
12. `superpowers:finishing-a-development-branch`
13. `/commit` → `/pr`

## Dual Memory System (non-negotiable)

Two memory systems are ALWAYS active simultaneously. They are NOT fallbacks for each other — both serve distinct roles and both are always consulted.

### Engram — Knowledge

Engram stores structured, curated knowledge with explicit lifecycle management.

- **Writes:** Explicit via `mem_save`. Only knowledge that passes the strict save criteria below.
- **Lifecycle:** Every entry MUST use `[ACTIVE]`/`[DEPRECATED]` markers (see Knowledge Lifecycle section).
- **Content:** Specs, tasks, architectural decisions, bug root causes, conventions, gotchas.
- **Session lifecycle:** `mem_session_start` and `mem_session_end` are handled by hooks — do NOT call these manually. Summaries go to claude-mem, NOT engram.

### Claude-Mem — History

Claude-mem provides automatic chronological capture and semantic search via HTTP API (ChromaDB-backed).

- **Writes:** Automatic capture in background + session summaries via HTTP (`POST http://localhost:3100/api/sessions/summarize`).
- **Lifecycle:** Chronological, NO `[ACTIVE]`/`[DEPRECATED]` tags. Entries are never deprecated — they form a timeline.
- **Content:** Session summaries, work history, everything that happened.

### Reading Flow (always both, NOT fallback)

When searching for context, ALWAYS consult both systems:

1. **Engram first** — search by topic key (`mem_search`, `mem_context`). This gives curated knowledge.
2. **Claude-mem second** — search by semantic query (`search`). This gives historical context and work that may not have been explicitly saved.
3. Both results are complementary. Engram gives you what the team decided is important. Claude-mem gives you what actually happened.

### Writing Rules

| What | Where | How |
|------|-------|-----|
| Knowledge discoveries | Engram | `mem_save` with `[ACTIVE]` prefix, strict criteria |
| Session summaries | Claude-mem | HTTP API, automatic, no lifecycle tags |
| Specs and tasks | Engram | `mem_save` with topic key |
| Wave progress | Engram | `mem_save` with `[ACTIVE]` prefix |

### Knowledge Save Criteria (STRICT) — save to engram ONLY if it matches at least one:

1. Bug fix with root cause — "X crashes because Y, fix: Z"
2. Gotcha/trap — something that wastes time if you don't know it
3. Architectural decision with reasoning — "chose X over Y because Z"
4. Established convention — "in this project we always do X"

**NEVER save to engram:**
- "Implemented X" (that's in git)
- "Used tool Y" (obvious from code)
- "Test passes" (temporal)
- Anything already in CLAUDE.md or derivable from code
- Status updates or progress reports
- Generic best practices everyone knows

**Deduplication Rule:** Before saving, ask: "Is this a NEW discovery or decision that passes the criteria above?" If it's just a status update or summary of work done, do NOT save it. Each datum is saved in ONE place by ONE actor.

### Discovery Extraction

The Stop hook extracts discoveries from session summaries (claude-mem) and promotes qualifying ones to engram as permanent knowledge entries. This bridges the two systems: summaries live in claude-mem, but important discoveries get extracted into engram where they have proper lifecycle management.

If engram is NOT available, warn: "Running in degraded mode. Run `aryflow setup` to install."

## Knowledge Lifecycle (non-negotiable)

Every `mem_save` content MUST start with `[ACTIVE] YYYY-MM-DD — `. No exceptions.

### Format

```
[ACTIVE] 2026-03-22 — {content here}
```

### When saving something that updates/replaces existing knowledge

1. `mem_search` for the existing entry on the same topic
2. `mem_update(id: old_id)` to prepend `[DEPRECATED] YYYY-MM-DD — Superseded by: {new_topic_key}\n` to the existing content
3. `mem_save` the new entry with `[ACTIVE] YYYY-MM-DD — ` prefix

### When reading from engram

- IGNORE any entry whose content starts with `[DEPRECATED]`
- When multiple `[ACTIVE]` entries exist for the same topic, the most recent date wins
- Never delete `[ACTIVE]` entries

### Examples

Saving a new discovery:
```
mem_save(topic: "aryflow/knowledge/go-conventions",
  content: "[ACTIVE] 2026-03-22 — Go 1.22 requires explicit loop variable capture in goroutines",
  project: "aryflow")
```

Updating an existing entry (old_id = 42):
```
mem_update(id: 42)  → prepend "[DEPRECATED] 2026-03-22 — Superseded by: aryflow/knowledge/go-conventions\n"
mem_save(topic: "aryflow/knowledge/go-conventions",
  content: "[ACTIVE] 2026-03-22 — Go 1.23 no longer requires explicit loop variable capture",
  project: "aryflow")
```

### Discovery Extraction from Session Summaries

Session summaries mix temporal info (Goal, Accomplished) with permanent knowledge (Discoveries). When a summary is deprecated, the discoveries would be lost. To prevent this:

1. The Stop hook agent saves the session summary as before (temporal, gets deprecated next session)
2. After saving the summary, it evaluates each Discovery against the **strict criteria checklist**:
   - [ ] Is it a bug fix with root cause? ("X crashes because Y, fix: Z")
   - [ ] Is it a gotcha/trap that wastes time if unknown?
   - [ ] Is it an architectural decision with reasoning? ("chose X over Y because Z")
   - [ ] Is it an established convention? ("in this project we always do X")
   - If NONE checked → do NOT save. If at least one checked → save as `[ACTIVE]` knowledge entry.
3. Also reject if it matches any NEVER-save rule: "Implemented X", "Used tool Y", "Test passes", anything in CLAUDE.md or derivable from code, status updates, generic best practices.
4. These knowledge entries are **independent** — they survive even when the summary is deprecated

Discovery categories use topic keys like:
- `aryflow/knowledge/go-macos` — Go build issues, macOS gotchas
- `aryflow/knowledge/ci-release` — CI/CD pipeline, GoReleaser
- `aryflow/knowledge/homebrew` — Homebrew tap and distribution
- `aryflow/knowledge/engram` — Memory lifecycle rules
- `aryflow/knowledge/{new-category}` — Any new topic as needed

### Maintenance

- Run knowledge-gc agent at milestone boundaries to clean up `[DEPRECATED]` entries

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
