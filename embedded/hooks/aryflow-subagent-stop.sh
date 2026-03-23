#!/bin/bash
# AryFlow SubagentStop hook — reminds orchestrator to check wave completion
cat <<'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "SubagentStop",
    "additionalContext": "ARYFLOW CHECK: A subagent just completed. If this was the last agent in a wave: UPDATE TODO.md marking tasks [x], save minimal wave progress to engram with '[ACTIVE] YYYY-MM-DD — ' prefix (mark previous wave progress as [DEPRECATED]), THEN commit. Subagents do NOT save work summaries — only technical discoveries (prefixed with '[ACTIVE] YYYY-MM-DD — ')."
  }
}
EOF
