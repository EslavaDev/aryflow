#!/bin/bash
# AryFlow context monitor — tracks tool usage, context health, and workflow compliance
# Runs on PostToolUse — receives tool info on stdin

INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty' 2>/dev/null)
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // empty' 2>/dev/null)

# Track tool call count for context estimation
MONITOR_FILE="/tmp/aryflow-context-${SESSION_ID:-default}"
if [ -f "$MONITOR_FILE" ]; then
  COUNT=$(cat "$MONITOR_FILE" | tr -d '[:space:]')
  [ -z "$COUNT" ] && COUNT=0
else
  COUNT=0
fi
COUNT=$((COUNT + 1))
echo "$COUNT" > "$MONITOR_FILE"

# Context health warnings based on tool call count
# ~200 calls ≈ 50% context, ~350 calls ≈ 80%, ~400+ ≈ danger zone
CONTEXT_MSG=""
if [ "$COUNT" -eq 200 ]; then
  CONTEXT_MSG="ARYFLOW CONTEXT WARNING: ~50% context used ($COUNT tool calls). Consider saving important state to engram before compaction."
elif [ "$COUNT" -eq 350 ]; then
  CONTEXT_MSG="ARYFLOW CONTEXT WARNING: ~80% context used ($COUNT tool calls). Save progress to engram NOW. Compaction is imminent."
elif [ "$COUNT" -ge 400 ] && [ $((COUNT % 25)) -eq 0 ]; then
  CONTEXT_MSG="ARYFLOW CONTEXT CRITICAL: $COUNT tool calls. Context may compact at any moment. Ensure all progress is saved to engram."
fi

# Workflow compliance reminders
WORKFLOW_MSG=""
case "$TOOL_NAME" in
  "Skill")
    SKILL_NAME=$(echo "$INPUT" | jq -r '.tool_input.skill // empty' 2>/dev/null)
    case "$SKILL_NAME" in
      "spec-it")
        WORKFLOW_MSG="ARYFLOW: spec-it invoked. Remember: auto-review with superpowers BEFORE showing to user, save to engram AFTER approval."
        ;;
      "execute-spec")
        WORKFLOW_MSG="ARYFLOW: execute-spec invoked. Remember: mem_session_start, save wave progress to engram, launch post-spec-docs agent at end."
        ;;
      "commit")
        WORKFLOW_MSG="ARYFLOW: Committing. Did you run /simplify and superpowers:verification-before-completion first?"
        ;;
    esac
    ;;
esac

# Combine messages
if [ -n "$CONTEXT_MSG" ] && [ -n "$WORKFLOW_MSG" ]; then
  echo "{\"systemMessage\":\"$CONTEXT_MSG\\n\\n$WORKFLOW_MSG\"}"
elif [ -n "$CONTEXT_MSG" ]; then
  echo "{\"systemMessage\":\"$CONTEXT_MSG\"}"
elif [ -n "$WORKFLOW_MSG" ]; then
  echo "{\"systemMessage\":\"$WORKFLOW_MSG\"}"
fi
