#!/bin/bash
# AryFlow Stop hook -- only triggers summary/knowledge agents when execute-spec just finished
# Uses .aryflow/.executing marker file to detect active spec execution
#
# States:
#   1. .aryflow/.executing EXISTS → mid-spec execution → block agents
#   2. .aryflow/.spec-completed EXISTS → spec just finished → allow agents, then cleanup marker
#   3. Neither exists → normal conversation → block agents (not needed)

EXECUTING_MARKER=".aryflow/.executing"
COMPLETED_MARKER=".aryflow/.spec-completed"

if [ -f "$EXECUTING_MARKER" ]; then
  # Active spec execution -- block everything
  cat <<'EOF'
{
  "continue": false,
  "stopReason": "ARYFLOW: Mid-spec execution. Skipping summary and knowledge -- orchestrator handles progress."
}
EOF
elif [ -f "$COMPLETED_MARKER" ]; then
  # Spec just completed -- allow summary + knowledge extraction, cleanup marker
  rm -f "$COMPLETED_MARKER"
  cat <<'EOF'
{
  "systemMessage": "ARYFLOW: Spec completed. Saving summary to claude-mem and extracting knowledge to engram."
}
EOF
else
  # Normal conversation -- don't waste agents
  cat <<'EOF'
{
  "continue": false,
  "stopReason": "ARYFLOW: Normal conversation. Skipping summary and knowledge agents."
}
EOF
fi
