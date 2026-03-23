#!/bin/bash
# AryFlow SessionStart hook -- reminds the agent of mandatory workflow steps
cat <<'EOF'
{
  "systemMessage": "ARYFLOW WORKFLOW REMINDER -- MANDATORY STEPS (only brainstorming is optional):\n1. /brainstorm (OPTIONAL)\n2. /spec-it -> writes SPEC.md + TODO.md\n3. MANDATORY: Auto-review SPEC.md with superpowers:requesting-code-review BEFORE showing to user\n4. MANDATORY: Review TODO.md for wave dependencies BEFORE showing to user\n5. MANDATORY: Save spec + tasks to engram (mem_save) -- warn if engram unavailable\n6. /execute-spec -> waves with engram session, progress saves, agent work summaries\n7. MANDATORY: Launch post-spec-docs agent after all waves\n8. MANDATORY: Save wave progress + agent summaries to engram\n9. superpowers:verification-before-completion\n10. /simplify\n11. superpowers:finishing-a-development-branch\n12. /commit -> /pr\n\nNEVER skip steps 3-8. If engram is available, ALWAYS use it. Warn user if degraded mode.\n\nCRITICAL: After EVERY completed wave, update TODO.md marking tasks as [x] BEFORE doing anything else. The statusline reads TODO.md for progress -- stale TODO = wrong status."
}
EOF
