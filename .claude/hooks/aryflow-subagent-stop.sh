#!/bin/bash
# AryFlow SubagentStop hook — reminds orchestrator to check wave completion
cat <<'EOF'
{
  "systemMessage": "ARYFLOW CHECK: Subagent completed. If this was the LAST agent in a wave: UPDATE TODO.md marking tasks [x], save minimal wave progress to engram, THEN commit. Subagents do NOT save work summaries — only technical discoveries."
}
EOF
