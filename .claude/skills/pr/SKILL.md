---
name: pr
description: Create a pull request following the project template. Handles branch creation, commit generation, and PR targeting. Use when asked to create a PR, open a pull request, or submit changes for review.
allowed-tools: Bash
model: claude-haiku-4-5-20251001
---

Create a pull request for the current changes following the project workflow.

**IMPORTANT:** All git and GitHub operations MUST use `gh` CLI. Never use raw `git push` — SAML SSO blocks SSH/HTTPS pushes. Use `gh` for pushing and PR creation.

## Step 1 — Gather git state

Run these commands in a single bash call to understand the current state:

```bash
git branch --show-current
git diff --staged --name-only
git diff --name-only
git log --oneline origin/main..HEAD 2>/dev/null || git log --oneline -5
git status --short
```

## Step 2 — Handle main branch

If the current branch is `main`:
- Ask the user: **"You're on main. What should the new branch be named?"**
- Create and switch: `git checkout -b <branch-name>`
- Continue with the new branch

## Step 3 — Handle staged changes with no commits

If there are staged changes AND no commits ahead of main (i.e. nothing to PR yet):
- Generate a conventional commit message from `git diff --staged` following these types:
  `feat` | `fix` | `update` | `build` | `docs` | `breaking` | `upgrade` | `chore`
- Run: `git commit -m "<type>: <message>"`

If there are no staged changes and no commits ahead of main → tell the user there is nothing to PR and stop.

## Step 4 — Push branch via gh

Ensure the branch is pushed to the remote before creating the PR:

```bash
gh repo sync --source . 2>/dev/null; git push -u origin "$(git branch --show-current)" 2>&1 || true
```

If push fails, try setting remote to HTTPS and retry:

```bash
gh repo view --json url -q '.url' | xargs -I {} git remote set-url origin {}.git
git push -u origin "$(git branch --show-current)"
```

## Step 5 — Build PR body from template

Analyze all commits ahead of main (`git log --oneline origin/main..HEAD`) and the diff (`git diff origin/main..HEAD --stat`) to fill in the template:

```markdown
## Description

> <summary of what this PR does>

## Type of change

- [ ] Bug fix
- [ ] Feature
- [ ] Hotfix
- [ ] Enhancement
- [ ] Chore (config, deps, infra)

## Scope

- [ ] `apps/core` — FastAPI backend
- [ ] `apps/cms` — React frontend
- [ ] `apps/board` — COO Board SPA
- [ ] `packages/` — types / api
- [ ] Infrastructure (Docker, Makefile, CI/CD)

## Checklists

### Development

- [ ] Tested locally with `make dev` or `make init`
- [ ] `VITE_*` env changes reflected in `apps/cms/.env.example`
- [ ] New core env vars added to `apps/core/.env.example`
- [ ] Lint passed (`make lint` / `make typecheck`)

### Backend (`apps/core`)

- [ ] New endpoints documented in `DEPLOYMENT.md` or `CLAUDE.md`
- [ ] DB schema changes handled with Alembic migration
- [ ] Valid model IDs used (`claude-opus-4-6` / `claude-haiku-4-5-20251001`)

### Deployment

- [ ] Deploy target is correct (`apps/core/.deploy.yml` / `apps/cms/.deploy.yml`)
- [ ] New secrets added to GitHub Actions if needed

### Code review

- [ ] PR has descriptive title and context useful to a reviewer
- [ ] Screenshots or screencasts attached if UI changes
```

Check the relevant boxes (`[x]`) based on what files changed:
- If `apps/core/` files changed → check `apps/core` scope
- If `apps/cms/` files changed → check `apps/cms` scope
- If `apps/board/` files changed → check `apps/board` scope
- If `packages/` files changed → check `packages/` scope
- If `Makefile`, `docker-compose*`, `.github/` changed → check Infrastructure scope
- Infer type of change from commit messages and diff

## Step 6 — Ask target branch

Ask the user: **"Target branch: `main` or `qa`?"**

Wait for the answer before creating the PR.

## Step 7 — Create the PR

Generate a short descriptive title (under 70 chars) from the commits, then run:

```bash
gh pr create \
  --title "<title>" \
  --body "$(cat <<'EOF'
<filled template from step 5>
EOF
)" \
  --base <main or qa>
```

Print the PR URL when done.
