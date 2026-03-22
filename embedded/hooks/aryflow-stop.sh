#!/bin/bash
# AryFlow Stop hook — reminds to save session summary before ending
cat <<'EOF'
{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": "ARYFLOW SESSION END CHECK:\n1. If engram is available, did you call mem_session_summary() and mem_session_end()?\n2. If you were in the middle of execute-spec, did you save progress to engram so the next session can resume?\n3. Were there any knowledge discoveries that should be saved before the session ends?"
  }
}
EOF
