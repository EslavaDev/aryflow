#!/bin/bash
# AryFlow Stop hook -- only triggers summary/knowledge agents when a spec just completed
# Three states:
#   1. Mid-spec (unchecked TODO items) → block agents
#   2. Spec just completed (all TODO checked) → allow agents (summary + knowledge)
#   3. No spec at all (normal conversation) → block agents (not needed)

HAS_SPEC=false
MID_SPEC=false

for f in specifications/*/TODO.md; do
  if [ -f "$f" ]; then
    HAS_SPEC=true
    if grep -q '^\s*- \[ \]' "$f" 2>/dev/null; then
      MID_SPEC=true
      break
    fi
  fi
done

if [ "$MID_SPEC" = true ]; then
  # Active spec with pending tasks -- block everything
  cat <<'EOF'
{
  "continue": false,
  "stopReason": "ARYFLOW: Mid-spec execution. Skipping summary and knowledge -- orchestrator handles progress."
}
EOF
elif [ "$HAS_SPEC" = true ]; then
  # Spec exists and all tasks done -- allow summary + knowledge extraction
  cat <<'EOF'
{
  "systemMessage": "ARYFLOW: Spec completed. Saving summary to claude-mem and extracting knowledge to engram."
}
EOF
else
  # No spec -- normal conversation, don't waste agents
  cat <<'EOF'
{
  "continue": false,
  "stopReason": "ARYFLOW: No active spec. Skipping summary and knowledge agents."
}
EOF
fi
