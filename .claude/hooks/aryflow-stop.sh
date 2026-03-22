#!/bin/bash
# AryFlow Stop hook — reminds to save session summary before ending
cat <<'EOF'
{
  "systemMessage": "ARYFLOW SESSION END: If engram is available, call mem_session_summary() and mem_session_end(). If mid-execute-spec, save progress to engram for resume."
}
EOF
