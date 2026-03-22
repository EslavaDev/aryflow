#!/bin/bash
# AryFlow status line — shows project state + context usage

PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -z "$PROJECT_ROOT" ]; then
  echo "AryFlow | No git repo"
  exit 0
fi

PROJECT_NAME=$(basename "$PROJECT_ROOT" | tr '_' '-' | tr '[:upper:]' '[:lower:]')
BRANCH=$(git branch --show-current 2>/dev/null || echo "detached")

# Check AryFlow version
VERSION=""
if [ -f "$PROJECT_ROOT/.aryflow/version" ]; then
  VERSION="v$(cat "$PROJECT_ROOT/.aryflow/version")"
fi

# Find active spec — highest numbered spec with unchecked TODO items
ACTIVE_SPEC=""
WAVE_PROGRESS=""
for spec_dir in $(ls -1d "$PROJECT_ROOT/specifications/"[0-9]*/ 2>/dev/null | sort -rV); do
  if [ -f "$spec_dir/TODO.md" ]; then
    UNCHECKED=$(grep -c '^\- \[ \]' "$spec_dir/TODO.md" 2>/dev/null | tr -d '[:space:]')
    CHECKED=$(grep -c '^\- \[x\]' "$spec_dir/TODO.md" 2>/dev/null | tr -d '[:space:]')
    [ -z "$UNCHECKED" ] && UNCHECKED=0
    [ -z "$CHECKED" ] && CHECKED=0
    TOTAL=$((UNCHECKED + CHECKED))
    if [ "$UNCHECKED" -gt 0 ] && [ "$TOTAL" -gt 0 ]; then
      ACTIVE_SPEC=$(basename "$spec_dir")
      WAVE_PROGRESS="$CHECKED/$TOTAL"
      break
    fi
  fi
done

# Context usage — read tool call count from most recent monitor file
TOOL_COUNT=0
CONTEXT_PCT=""
LATEST_MONITOR=$(ls -t /tmp/aryflow-context-* 2>/dev/null | head -1)
if [ -n "$LATEST_MONITOR" ] && [ -f "$LATEST_MONITOR" ]; then
  TOOL_COUNT=$(cat "$LATEST_MONITOR" 2>/dev/null | tr -d '[:space:]')
  [ -z "$TOOL_COUNT" ] && TOOL_COUNT=0
  # Show tool call count as context indicator
  # Note: exact token % is not available to hooks — count is a rough proxy
  if [ "$TOOL_COUNT" -gt 0 ]; then
    CONTEXT_PCT="${TOOL_COUNT} calls"
  fi
fi

# Build status line
STATUS="AryFlow"
[ -n "$VERSION" ] && STATUS="$STATUS $VERSION"
STATUS="$STATUS | $PROJECT_NAME | $BRANCH"
[ -n "$ACTIVE_SPEC" ] && STATUS="$STATUS | $ACTIVE_SPEC [$WAVE_PROGRESS]"
[ -n "$CONTEXT_PCT" ] && STATUS="$STATUS | ctx:$CONTEXT_PCT"

echo "$STATUS"
