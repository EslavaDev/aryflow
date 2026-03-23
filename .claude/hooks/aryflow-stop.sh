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
  "systemMessage": "ARYFLOW: execute-spec session detected (unchecked TODO items exist). Do NOT save a session summary — wave progress is already tracked by the orchestrator."
}
EOF
else
  cat <<'EOF'
{
  "systemMessage": "ARYFLOW SESSION END: The Stop hook agent handles dual memory cleanup. Session summaries go to claude-mem (HTTP API, chronological, no lifecycle tags). Discoveries that pass strict criteria get extracted to engram as [ACTIVE] knowledge entries with topic_key '{project}/knowledge/{category}'. All engram mem_save content MUST start with '[ACTIVE] YYYY-MM-DD — '."
}
EOF
fi
