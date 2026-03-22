// Package doctor implements the "aryflow doctor" command.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EslavaDev/aryflow/internal/checks"
	"github.com/EslavaDev/aryflow/internal/ui"
)

// Severity represents how critical a check failure is.
type Severity int

const (
	SeverityError   Severity = iota
	SeverityWarning
)

// CheckResult holds the outcome of a single doctor check.
type CheckResult struct {
	Name     string
	Passed   bool
	Severity Severity
	Message  string // detail on failure
	Fix      string // suggested fix command
}

// Summary holds aggregated check counts.
type Summary struct {
	Passed   int
	Warnings int
	Errors   int
}

// RunChecks executes all 15 doctor checks and returns the results.
// It never stops on first failure.
func RunChecks(cliVersion string) []CheckResult {
	gitRoot := findGitRoot()
	results := make([]CheckResult, 0, 15)

	// 1. Git repo
	results = append(results, checkGitRepo())

	// For file checks we need the git root; if not in a repo, use cwd.
	root := gitRoot
	if root == "" {
		root, _ = os.Getwd()
	}

	// 2. .aryflow/version exists
	results = append(results, checkFileExists(root, ".aryflow/version", ".aryflow/version", SeverityError, "Run: aryflow init"))

	// 3. CLAUDE.md exists
	results = append(results, checkFileExists(root, "CLAUDE.md", "CLAUDE.md", SeverityError, "Run: aryflow init"))

	// 4-5. Required skills (Error)
	results = append(results, checkFileExists(root, ".claude/skills/spec-it/SKILL.md", "spec-it skill", SeverityError, "Run: aryflow init"))
	results = append(results, checkFileExists(root, ".claude/skills/execute-spec/SKILL.md", "execute-spec skill", SeverityError, "Run: aryflow init"))

	// 6-7. Optional skills (Warning)
	results = append(results, checkFileExists(root, ".claude/skills/commit/SKILL.md", "commit skill", SeverityWarning, "Run: aryflow init"))
	results = append(results, checkFileExists(root, ".claude/skills/pr/SKILL.md", "pr skill", SeverityWarning, "Run: aryflow init"))

	// 8. Required agent (Error)
	results = append(results, checkFileExists(root, ".claude/agents/merge-wave.md", "merge-wave agent", SeverityError, "Run: aryflow init"))

	// 9-10. Optional agents (Warning)
	results = append(results, checkFileExists(root, ".claude/agents/post-spec-docs.md", "post-spec-docs agent", SeverityWarning, "Run: aryflow init"))
	results = append(results, checkFileExists(root, ".claude/agents/knowledge-gc.md", "knowledge-gc agent", SeverityWarning, "Run: aryflow init"))

	// 11. specifications/ directory
	results = append(results, checkDirExists(root, "specifications", "specifications/ directory", SeverityWarning, "Run: aryflow init"))

	// 12. Engram installed
	results = append(results, checkEngram())

	// 13. Claude-Mem installed
	results = append(results, checkClaudeMem())

	// 14. Superpowers installed
	results = append(results, checkSuperpowers())

	// 15. Version match
	results = append(results, checkVersionMatch(root, cliVersion))

	return results
}

// Summarize computes pass/warning/error counts from results.
func Summarize(results []CheckResult) Summary {
	var s Summary
	for _, r := range results {
		if r.Passed {
			s.Passed++
		} else if r.Severity == SeverityWarning {
			s.Warnings++
		} else {
			s.Errors++
		}
	}
	return s
}

// Run executes the doctor command, printing results and returning exit code.
func Run(verbose bool, cliVersion string) int {
	gitRoot := findGitRoot()
	projectName := "unknown"
	if gitRoot != "" {
		projectName = filepath.Base(gitRoot)
	}

	ui.Header(fmt.Sprintf("AryFlow Doctor — Checking project %q", projectName))

	results := RunChecks(cliVersion)

	for _, r := range results {
		if r.Passed {
			ui.Success(r.Name)
		} else if r.Severity == SeverityWarning {
			ui.Warning(fmt.Sprintf("%s — %s", r.Name, r.Message))
			if r.Fix != "" {
				ui.Suggestion(r.Fix)
			}
		} else {
			ui.Error(fmt.Sprintf("%s — %s", r.Name, r.Message))
			if r.Fix != "" {
				ui.Suggestion(r.Fix)
			}
		}
	}

	s := Summarize(results)
	fmt.Println()
	fmt.Printf("Results: %d passed, %d warnings, %d errors\n", s.Passed, s.Warnings, s.Errors)

	if s.Errors > 0 {
		return 1
	}
	return 0
}

// --- individual check functions ---

func findGitRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func checkGitRepo() CheckResult {
	root := findGitRoot()
	if root == "" {
		return CheckResult{
			Name:     "Git repository",
			Passed:   false,
			Severity: SeverityError,
			Message:  "not inside a git repository",
			Fix:      "Run: git init",
		}
	}
	return CheckResult{
		Name:   "Git repository",
		Passed: true,
	}
}

func checkFileExists(root, relPath, name string, severity Severity, fix string) CheckResult {
	fullPath := filepath.Join(root, relPath)
	if _, err := os.Stat(fullPath); err != nil {
		return CheckResult{
			Name:     name,
			Passed:   false,
			Severity: severity,
			Message:  "missing",
			Fix:      fix,
		}
	}

	// For .aryflow/version, show the version value
	if relPath == ".aryflow/version" {
		data, err := os.ReadFile(fullPath)
		if err == nil {
			ver := strings.TrimSpace(string(data))
			return CheckResult{
				Name:   fmt.Sprintf("%s (v%s)", name, ver),
				Passed: true,
			}
		}
	}

	return CheckResult{
		Name:   name,
		Passed: true,
	}
}

func checkDirExists(root, relPath, name string, severity Severity, fix string) CheckResult {
	fullPath := filepath.Join(root, relPath)
	info, err := os.Stat(fullPath)
	if err != nil || !info.IsDir() {
		return CheckResult{
			Name:     name,
			Passed:   false,
			Severity: severity,
			Message:  "missing",
			Fix:      fix,
		}
	}
	return CheckResult{
		Name:   name,
		Passed: true,
	}
}

func checkEngram() CheckResult {
	_, err := checks.CheckEngram()
	if err != nil {
		return CheckResult{
			Name:     "Engram installed",
			Passed:   false,
			Severity: SeverityError,
			Message:  "not found",
			Fix:      "Run: aryflow setup",
		}
	}
	return CheckResult{
		Name:   "Engram installed",
		Passed: true,
	}
}

func checkClaudeMem() CheckResult {
	_, err := checks.CheckClaudeMem()
	if err != nil {
		return CheckResult{
			Name:     "Claude-Mem plugin",
			Passed:   false,
			Severity: SeverityWarning,
			Message:  "not found",
			Fix:      "Run: aryflow setup",
		}
	}
	return CheckResult{
		Name:   "Claude-Mem plugin",
		Passed: true,
	}
}

func checkSuperpowers() CheckResult {
	_, err := checks.CheckSuperpowers()
	if err != nil {
		return CheckResult{
			Name:     "Superpowers plugin",
			Passed:   false,
			Severity: SeverityWarning,
			Message:  "not found",
			Fix:      "Run: aryflow setup",
		}
	}
	return CheckResult{
		Name:   "Superpowers plugin",
		Passed: true,
	}
}

func checkVersionMatch(root, cliVersion string) CheckResult {
	data, err := os.ReadFile(filepath.Join(root, ".aryflow/version"))
	if err != nil {
		return CheckResult{
			Name:     "Version match",
			Passed:   false,
			Severity: SeverityWarning,
			Message:  "cannot read .aryflow/version",
			Fix:      "Run: aryflow init",
		}
	}
	projectVersion := strings.TrimSpace(string(data))
	if projectVersion != cliVersion {
		return CheckResult{
			Name:     "Version match",
			Passed:   false,
			Severity: SeverityWarning,
			Message:  fmt.Sprintf("project: v%s, CLI: v%s", projectVersion, cliVersion),
			Fix:      "Run: aryflow update",
		}
	}
	return CheckResult{
		Name:   "Version match",
		Passed: true,
	}
}
