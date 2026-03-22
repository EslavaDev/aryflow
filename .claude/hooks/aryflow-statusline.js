#!/usr/bin/env node
// AryFlow Statusline — shows project state + real context usage
// Receives context_window metrics from Claude Code via stdin

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const AUTO_COMPACT_BUFFER_PCT = 16.5;

let input = '';
const stdinTimeout = setTimeout(() => process.exit(0), 3000);
process.stdin.setEncoding('utf8');
process.stdin.on('data', chunk => input += chunk);
process.stdin.on('end', () => {
  clearTimeout(stdinTimeout);
  try {
    const data = JSON.parse(input);
    const model = data.model?.display_name || 'Claude';
    const dir = data.workspace?.current_dir || process.cwd();
    const session = data.session_id || '';
    const remaining = data.context_window?.remaining_percentage;

    // Project info
    const dirname = path.basename(dir);
    let branch = '';
    try {
      branch = execSync('git branch --show-current 2>/dev/null', { cwd: dir, encoding: 'utf8' }).trim();
    } catch (e) {
      branch = '';
    }

    // Context window display
    let ctx = '';
    if (remaining != null) {
      const usableRemaining = Math.max(0, ((remaining - AUTO_COMPACT_BUFFER_PCT) / (100 - AUTO_COMPACT_BUFFER_PCT)) * 100);
      const used = Math.max(0, Math.min(100, Math.round(100 - usableRemaining)));

      // Write bridge file for context-monitor hook
      if (session) {
        try {
          const bridgePath = path.join(os.tmpdir(), `aryflow-ctx-${session}.json`);
          fs.writeFileSync(bridgePath, JSON.stringify({
            session_id: session,
            remaining_percentage: remaining,
            used_pct: used,
            timestamp: Math.floor(Date.now() / 1000)
          }));
        } catch (e) {}
      }

      // Progress bar (10 segments)
      const filled = Math.floor(used / 10);
      const bar = '█'.repeat(filled) + '░'.repeat(10 - filled);

      if (used < 50) {
        ctx = ` \x1b[32m${bar} ${used}%\x1b[0m`;
      } else if (used < 65) {
        ctx = ` \x1b[33m${bar} ${used}%\x1b[0m`;
      } else if (used < 80) {
        ctx = ` \x1b[38;5;208m${bar} ${used}%\x1b[0m`;
      } else {
        ctx = ` \x1b[5;31m${bar} ${used}%\x1b[0m`;
      }
    }

    // Active spec — find highest numbered spec with unchecked items
    let spec = '';
    const specsDir = path.join(dir, 'specifications');
    if (fs.existsSync(specsDir)) {
      try {
        const dirs = fs.readdirSync(specsDir)
          .filter(d => /^\d{3}-/.test(d))
          .sort()
          .reverse();

        for (const d of dirs) {
          const todoPath = path.join(specsDir, d, 'TODO.md');
          if (fs.existsSync(todoPath)) {
            const content = fs.readFileSync(todoPath, 'utf8');
            const unchecked = (content.match(/^- \[ \]/gm) || []).length;
            const checked = (content.match(/^- \[x\]/gm) || []).length;
            const total = unchecked + checked;
            if (unchecked > 0 && total > 0) {
              spec = ` │ \x1b[1m${d}\x1b[0m [${checked}/${total}]`;
              break;
            }
          }
        }
      } catch (e) {}
    }

    // Output
    const parts = [`\x1b[36mAryFlow\x1b[0m`];
    parts.push(`\x1b[2m${model}\x1b[0m`);
    if (branch) parts.push(`\x1b[2m${dirname}\x1b[0m:\x1b[35m${branch}\x1b[0m`);
    else parts.push(`\x1b[2m${dirname}\x1b[0m`);

    process.stdout.write(parts.join(' │ ') + spec + ctx);
  } catch (e) {
    process.exit(0);
  }
});
