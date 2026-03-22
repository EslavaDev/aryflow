---
name: commit
description: Generate a conventional commit message from staged changes and commit. Use when asked to commit, create a commit, or generate a commit message.
allowed-tools: Bash
model: claude-haiku-4-5-20251001
context: fork
metadata:
  author: EslavaDev
  version: "1.0.0"
---

Generate a commit message in English following conventional commits format from the staged changes.

## Process

1. Run `git diff --staged` to understand what changed
2. Generate a commit message based only on staged changes
3. Run `git commit -m "<type>: <message>" --trailer "Co-Authored-By: Claude <noreply@anthropic.com>"`

## Valid types

- `feat` — new features
- `fix` — bug fixes
- `update` — updates to existing features
- `build` — build system changes
- `docs` — documentation only
- `breaking` — breaking changes
- `upgrade` — dependency upgrades
- `chore` — maintenance tasks

## Rules

- Message in English, lowercase, concise
- No assumptions about unstaged changes
- No scope needed unless it adds clarity

## Examples

- `feat: add social media audit endpoint`
- `fix: correct haiku model id in app.py`
- `chore: update docker compose profiles`
- `build: add makefile deploy targets`
