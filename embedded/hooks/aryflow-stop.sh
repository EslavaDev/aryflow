#!/bin/bash
# AryFlow Stop hook — only reminds about session summary for non-spec sessions
# If execute-spec is active (TODO.md has unchecked items), skip — progress is already saved by the orchestrator.

TODO_FILE="specifications/*/TODO.md"
MID_SPEC=false

for f in $TODO_FILE; do
  if [ -f "$f" ] && grep -q '^\s*- \[ \]' "$f" 2>/dev/null; then
    MID_SPEC=true
    break
  fi
done

if [ "$MID_SPEC" = true ]; then
  cat <<'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": "ARYFLOW: execute-spec session detected (unchecked TODO items exist). Do NOT save a session summary — wave progress is already tracked by the orchestrator."
  }
}
EOF
else
  cat <<'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": "ARYFLOW SESSION END: The Stop hook agent will handle session summary. REMINDER: All mem_save content MUST start with '[ACTIVE] YYYY-MM-DD — '. Search for and mark previous session summaries as [DEPRECATED] before saving a new one."
  }
}
EOF
fi
