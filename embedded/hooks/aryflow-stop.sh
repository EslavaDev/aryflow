#!/bin/bash
# AryFlow Stop hook -- blocks summary/knowledge agents during active spec execution
# If execute-spec is active (TODO.md has unchecked items), block subsequent hooks.
# If no spec active, allow agents to save summary + extract knowledge.

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
  "continue": false,
  "stopReason": "ARYFLOW: execute-spec session detected (unchecked TODO items). Skipping summary and knowledge extraction -- wave progress is tracked by the orchestrator."
}
EOF
else
  cat <<'EOF'
{
  "systemMessage": "ARYFLOW SESSION END: Saving summary to claude-mem and extracting knowledge to engram."
}
EOF
fi
