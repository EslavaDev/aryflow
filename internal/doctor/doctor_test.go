package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGitRepo_InGitDir(t *testing.T) {
	// Create a temp dir with .git
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0o755)

	// Change to that dir so findGitRoot works
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	r := checkGitRepo()
	if !r.Passed {
		t.Errorf("expected checkGitRepo to pass in a git dir, got: %s", r.Message)
	}
}

func TestCheckGitRepo_NotInGitDir(t *testing.T) {
	dir := t.TempDir()

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	r := checkGitRepo()
	if r.Passed {
		t.Error("expected checkGitRepo to fail outside a git dir")
	}
	if r.Severity != SeverityError {
		t.Error("expected SeverityError for missing git repo")
	}
}

func TestCheckFileExists_Present(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# test"), 0o644)

	r := checkFileExists(dir, "CLAUDE.md", "CLAUDE.md", SeverityError, "Run: aryflow init")
	if !r.Passed {
		t.Error("expected checkFileExists to pass when file exists")
	}
}

func TestCheckFileExists_Missing(t *testing.T) {
	dir := t.TempDir()

	r := checkFileExists(dir, "CLAUDE.md", "CLAUDE.md", SeverityError, "Run: aryflow init")
	if r.Passed {
		t.Error("expected checkFileExists to fail when file is missing")
	}
	if r.Severity != SeverityError {
		t.Errorf("expected SeverityError, got %d", r.Severity)
	}
	if r.Fix != "Run: aryflow init" {
		t.Errorf("expected fix suggestion, got %q", r.Fix)
	}
}

func TestCheckFileExists_WarningSeverity(t *testing.T) {
	dir := t.TempDir()

	r := checkFileExists(dir, ".claude/skills/pr/SKILL.md", "pr skill", SeverityWarning, "Run: aryflow init")
	if r.Passed {
		t.Error("expected check to fail")
	}
	if r.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %d", r.Severity)
	}
}

func TestCheckDirExists_Present(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "specifications"), 0o755)

	r := checkDirExists(dir, "specifications", "specifications/ directory", SeverityWarning, "Run: aryflow init")
	if !r.Passed {
		t.Error("expected checkDirExists to pass when directory exists")
	}
}

func TestCheckDirExists_Missing(t *testing.T) {
	dir := t.TempDir()

	r := checkDirExists(dir, "specifications", "specifications/ directory", SeverityWarning, "Run: aryflow init")
	if r.Passed {
		t.Error("expected checkDirExists to fail when directory is missing")
	}
}

func TestCheckVersionMatch_Matches(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".aryflow"), 0o755)
	os.WriteFile(filepath.Join(dir, ".aryflow/version"), []byte("0.1.0\n"), 0o644)

	r := checkVersionMatch(dir, "0.1.0")
	if !r.Passed {
		t.Errorf("expected version match to pass, got: %s", r.Message)
	}
}

func TestCheckVersionMatch_Mismatch(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".aryflow"), 0o755)
	os.WriteFile(filepath.Join(dir, ".aryflow/version"), []byte("0.1.0\n"), 0o644)

	r := checkVersionMatch(dir, "0.2.0")
	if r.Passed {
		t.Error("expected version match to fail on mismatch")
	}
	if r.Severity != SeverityWarning {
		t.Error("expected SeverityWarning for version mismatch")
	}
}

func TestCheckVersionMatch_NoFile(t *testing.T) {
	dir := t.TempDir()

	r := checkVersionMatch(dir, "0.1.0")
	if r.Passed {
		t.Error("expected version match to fail when file is missing")
	}
}

func TestSummarize(t *testing.T) {
	results := []CheckResult{
		{Passed: true},
		{Passed: true},
		{Passed: false, Severity: SeverityError},
		{Passed: false, Severity: SeverityWarning},
		{Passed: false, Severity: SeverityWarning},
		{Passed: true},
	}

	s := Summarize(results)
	if s.Passed != 3 {
		t.Errorf("expected 3 passed, got %d", s.Passed)
	}
	if s.Warnings != 2 {
		t.Errorf("expected 2 warnings, got %d", s.Warnings)
	}
	if s.Errors != 1 {
		t.Errorf("expected 1 error, got %d", s.Errors)
	}
}

func TestSummarize_AllPassed(t *testing.T) {
	results := []CheckResult{
		{Passed: true},
		{Passed: true},
		{Passed: true},
	}

	s := Summarize(results)
	if s.Passed != 3 || s.Warnings != 0 || s.Errors != 0 {
		t.Errorf("expected all passed, got passed=%d warnings=%d errors=%d", s.Passed, s.Warnings, s.Errors)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Passed != 0 || s.Warnings != 0 || s.Errors != 0 {
		t.Error("expected all zeros for empty results")
	}
}

func TestRunChecks_InTempDir(t *testing.T) {
	// Run in a temp dir with minimal setup to verify all 21 checks execute
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0o755)

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	results := RunChecks("0.1.0")
	if len(results) != 21 {
		t.Errorf("expected 21 checks, got %d", len(results))
	}

	// First check (git repo) should pass since we created .git
	if !results[0].Passed {
		t.Error("expected git repo check to pass")
	}

	// .aryflow/version should fail (we didn't create it)
	if results[1].Passed {
		t.Error("expected .aryflow/version check to fail")
	}
}

func TestCheckFileExists_VersionShowsValue(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".aryflow"), 0o755)
	os.WriteFile(filepath.Join(dir, ".aryflow/version"), []byte("0.3.0\n"), 0o644)

	r := checkFileExists(dir, ".aryflow/version", ".aryflow/version", SeverityError, "Run: aryflow init")
	if !r.Passed {
		t.Error("expected check to pass")
	}
	if r.Name != ".aryflow/version (v0.3.0)" {
		t.Errorf("expected name to include version, got %q", r.Name)
	}
}

func TestCheckSettingsJSON_Valid(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	content := `{"hooks":{"SessionStart":["bash .claude/hooks/aryflow-session-start.sh"]},"statusLine":"node .claude/hooks/aryflow-statusline.js"}`
	os.WriteFile(filepath.Join(dir, ".claude/settings.json"), []byte(content), 0o644)

	r := checkSettingsJSON(dir)
	if !r.Passed {
		t.Errorf("expected checkSettingsJSON to pass, got: %s", r.Message)
	}
}

func TestCheckSettingsJSON_Missing(t *testing.T) {
	dir := t.TempDir()

	r := checkSettingsJSON(dir)
	if r.Passed {
		t.Error("expected checkSettingsJSON to fail when file is missing")
	}
	if r.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %d", r.Severity)
	}
}

func TestCheckSettingsJSON_NoHooks(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	os.WriteFile(filepath.Join(dir, ".claude/settings.json"), []byte(`{"statusLine":"test"}`), 0o644)

	r := checkSettingsJSON(dir)
	if r.Passed {
		t.Error("expected checkSettingsJSON to fail when hooks missing")
	}
	if r.Message != "hooks not configured in settings.json" {
		t.Errorf("unexpected message: %s", r.Message)
	}
}

func TestCheckSettingsJSON_NoSessionStart(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	os.WriteFile(filepath.Join(dir, ".claude/settings.json"), []byte(`{"hooks":{},"statusLine":"test"}`), 0o644)

	r := checkSettingsJSON(dir)
	if r.Passed {
		t.Error("expected checkSettingsJSON to fail when SessionStart missing")
	}
	if r.Message != "SessionStart hook not configured" {
		t.Errorf("unexpected message: %s", r.Message)
	}
}

func TestCheckSettingsJSON_NoStatusLine(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	os.WriteFile(filepath.Join(dir, ".claude/settings.json"), []byte(`{"hooks":{"SessionStart":["test"]}}`), 0o644)

	r := checkSettingsJSON(dir)
	if r.Passed {
		t.Error("expected checkSettingsJSON to fail when statusLine missing")
	}
	if r.Message != "statusLine not configured in settings.json" {
		t.Errorf("unexpected message: %s", r.Message)
	}
}

func TestCheckSettingsJSON_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	os.WriteFile(filepath.Join(dir, ".claude/settings.json"), []byte(`{invalid`), 0o644)

	r := checkSettingsJSON(dir)
	if r.Passed {
		t.Error("expected checkSettingsJSON to fail on invalid JSON")
	}
	if r.Message != "invalid JSON in .claude/settings.json" {
		t.Errorf("unexpected message: %s", r.Message)
	}
}

func TestCheckActiveTODO_NoTodos(t *testing.T) {
	dir := t.TempDir()

	r := checkActiveTODO(dir)
	if !r.Passed {
		t.Error("expected checkActiveTODO to pass when no TODOs exist")
	}
}

func TestCheckActiveTODO_StaleProgress(t *testing.T) {
	dir := t.TempDir()
	specDir := filepath.Join(dir, "specifications", "001-feature")
	os.MkdirAll(specDir, 0o755)
	content := "# TODO\n- [x] Done task\n- [ ] Pending task\n"
	os.WriteFile(filepath.Join(specDir, "TODO.md"), []byte(content), 0o644)

	r := checkActiveTODO(dir)
	if r.Passed {
		t.Error("expected checkActiveTODO to fail with stale progress")
	}
	if r.Severity != SeverityWarning {
		t.Errorf("expected SeverityWarning, got %d", r.Severity)
	}
}

func TestCheckActiveTODO_AllDone(t *testing.T) {
	dir := t.TempDir()
	specDir := filepath.Join(dir, "specifications", "001-feature")
	os.MkdirAll(specDir, 0o755)
	content := "# TODO\n- [x] Done task\n- [x] Also done\n"
	os.WriteFile(filepath.Join(specDir, "TODO.md"), []byte(content), 0o644)

	r := checkActiveTODO(dir)
	if !r.Passed {
		t.Errorf("expected checkActiveTODO to pass when all done, got: %s", r.Message)
	}
}

func TestCheckActiveTODO_AllPending(t *testing.T) {
	dir := t.TempDir()
	specDir := filepath.Join(dir, "specifications", "001-feature")
	os.MkdirAll(specDir, 0o755)
	content := "# TODO\n- [ ] Pending task\n- [ ] Also pending\n"
	os.WriteFile(filepath.Join(specDir, "TODO.md"), []byte(content), 0o644)

	r := checkActiveTODO(dir)
	if !r.Passed {
		t.Errorf("expected checkActiveTODO to pass when all pending, got: %s", r.Message)
	}
}
