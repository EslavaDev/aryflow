#!/bin/bash
# AryFlow SubagentStop hook — reminds orchestrator to verify subagent saved to engram
cat <<'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "SubagentStop",
    "additionalContext": "ARYFLOW CHECK: A subagent just completed. Verify:\n1. Did the subagent save a work summary to engram? (mem_save with wave-N/agent-task topic)\n2. Did the subagent save any new knowledge discoveries? (mem_save to knowledge/*)\n3. If this was the last agent in a wave, save wave progress to engram.\n4. If this was merge-wave agent, check its summary for unresolved conflicts."
  }
}
EOF
