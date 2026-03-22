#!/usr/bin/env node
// AryFlow Context Monitor — PostToolUse hook
// Reads real context metrics from the statusline bridge file and injects
// warnings when context usage is high + workflow compliance reminders.

const fs = require('fs');
const path = require('path');
const os = require('os');

const WARNING_THRESHOLD = 35;  // remaining_percentage <= 35%
const CRITICAL_THRESHOLD = 25; // remaining_percentage <= 25%
const STALE_SECONDS = 60;
const DEBOUNCE_CALLS = 5;

let input = '';
const stdinTimeout = setTimeout(() => process.exit(0), 3000);
process.stdin.setEncoding('utf8');
process.stdin.on('data', chunk => input += chunk);
process.stdin.on('end', () => {
  clearTimeout(stdinTimeout);
  try {
    const data = JSON.parse(input);
    const sessionId = data.session_id;
    const toolName = data.tool_name;
    const messages = [];

    // --- Workflow compliance reminders ---
    if (toolName === 'Skill') {
      const skillName = data.tool_input?.skill;
      if (skillName === 'spec-it') {
        messages.push('ARYFLOW: spec-it invoked. Auto-review with superpowers BEFORE showing to user. Save to engram AFTER approval. Update TODO.md after each wave.');
      } else if (skillName === 'execute-spec') {
        messages.push('ARYFLOW: execute-spec invoked. mem_session_start, save wave progress to engram, update TODO.md [x] after each wave, launch post-spec-docs at end.');
      } else if (skillName === 'commit') {
        messages.push('ARYFLOW: Committing. Did you run /simplify and superpowers:verification-before-completion first? Is TODO.md updated?');
      }
    }

    // --- Context window monitoring ---
    if (sessionId) {
      const bridgePath = path.join(os.tmpdir(), `aryflow-ctx-${sessionId}.json`);

      if (fs.existsSync(bridgePath)) {
        const metrics = JSON.parse(fs.readFileSync(bridgePath, 'utf8'));
        const now = Math.floor(Date.now() / 1000);

        // Skip stale metrics
        if (!metrics.timestamp || (now - metrics.timestamp) <= STALE_SECONDS) {
          const remaining = metrics.remaining_percentage;
          const usedPct = metrics.used_pct;

          if (remaining <= WARNING_THRESHOLD) {
            // Debounce logic
            const warnPath = path.join(os.tmpdir(), `aryflow-ctx-${sessionId}-warned.json`);
            let warnData = { callsSinceWarn: 0, lastLevel: null };

            if (fs.existsSync(warnPath)) {
              try { warnData = JSON.parse(fs.readFileSync(warnPath, 'utf8')); } catch (e) {}
            }

            warnData.callsSinceWarn = (warnData.callsSinceWarn || 0) + 1;

            const isCritical = remaining <= CRITICAL_THRESHOLD;
            const currentLevel = isCritical ? 'critical' : 'warning';
            const severityEscalated = currentLevel === 'critical' && warnData.lastLevel === 'warning';

            if (warnData.callsSinceWarn >= DEBOUNCE_CALLS || severityEscalated || warnData.lastLevel === null) {
              warnData.callsSinceWarn = 0;
              warnData.lastLevel = currentLevel;

              if (isCritical) {
                messages.push(
                  `ARYFLOW CONTEXT CRITICAL: Usage at ${usedPct}%. Remaining: ${remaining}%. ` +
                  'Context nearly exhausted. Save ALL progress to engram NOW (mem_save for wave progress, knowledge). ' +
                  'Inform the user context is critical.'
                );
              } else {
                messages.push(
                  `ARYFLOW CONTEXT WARNING: Usage at ${usedPct}%. Remaining: ${remaining}%. ` +
                  'Context getting limited. Save important state to engram. Avoid starting new complex work.'
                );
              }
            }

            fs.writeFileSync(warnPath, JSON.stringify(warnData));
          }
        }
      }
    }

    // Output combined messages
    if (messages.length > 0) {
      const output = {
        hookSpecificOutput: {
          hookEventName: 'PostToolUse',
          additionalContext: messages.join('\n\n')
        }
      };
      process.stdout.write(JSON.stringify(output));
    }
  } catch (e) {
    process.exit(0);
  }
});
