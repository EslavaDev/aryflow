#!/bin/bash
# AryFlow SubagentStop hook — reminds orchestrator to verify subagent saved to engram
cat <<'EOF'
{
  "systemMessage": "ARYFLOW CHECK: Subagent completed. Verify: 1) Saved work summary to engram (wave-N/agent-task topic). 2) Saved knowledge discoveries. 3) If this was the LAST agent in a wave: UPDATE TODO.md marking tasks [x], save wave progress to engram, THEN commit."
}
EOF
