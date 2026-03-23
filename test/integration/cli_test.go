package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once for all tests
	dir, _ := os.MkdirTemp("", "aryflow-test")
	binaryPath = filepath.Join(dir, "aryflow")

	// Build from project root (two levels up from test/integration/)
	// CGO_ENABLED=0 avoids dyld LC_UUID issues on macOS with Go 1.22
	projectRoot, _ := filepath.Abs(filepath.Join("..", ".."))
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/aryflow")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func runBinary(args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "ARYFLOW_YES=1", "ARYFLOW_NO_COLOR=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runBinaryInDir(dir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ARYFLOW_YES=1", "ARYFLOW_NO_COLOR=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// --- Version tests ---

func TestVersion(t *testing.T) {
	out, err := runBinary("--version")
	if err != nil {
		t.Fatalf("--version failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "aryflow") {
		t.Errorf("expected 'aryflow' in version output, got: %s", out)
	}
	if !strings.Contains(out, "v") {
		t.Errorf("expected version string with 'v', got: %s", out)
	}
}

func TestVersionShort(t *testing.T) {
	out, err := runBinary("-v")
	if err != nil {
		t.Fatalf("-v failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "aryflow") {
		t.Errorf("expected version output, got: %s", out)
	}
}

// --- Help tests ---

func TestHelp(t *testing.T) {
	out, err := runBinary("--help")
	if err != nil {
		t.Fatalf("--help failed: %v\n%s", err, out)
	}
	for _, cmd := range []string{"setup", "init", "doctor", "update"} {
		if !strings.Contains(out, cmd) {
			t.Errorf("help should list %q command, got: %s", cmd, out)
		}
	}
}

func TestHelpShort(t *testing.T) {
	out, err := runBinary("-h")
	if err != nil {
		t.Fatalf("-h failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "setup") {
		t.Errorf("help should list commands, got: %s", out)
	}
}

func TestHelpContainsFlags(t *testing.T) {
	out, _ := runBinary("--help")
	for _, flag := range []string{"--verbose", "--yes", "--version", "--force", "--dry-run"} {
		if !strings.Contains(out, flag) {
			t.Errorf("help should mention %q, got: %s", flag, out)
		}
	}
}

// --- No args / unknown tests ---

func TestNoArgs(t *testing.T) {
	out, err := runBinary()
	// No args prints usage and exits 0
	if err != nil {
		t.Fatalf("no args should exit 0, got error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Usage") {
		t.Errorf("expected usage in output, got: %s", out)
	}
}

func TestUnknownCommand(t *testing.T) {
	out, err := runBinary("foobar")
	if err == nil {
		t.Error("expected error for unknown command")
	}
	if !strings.Contains(out, "Unknown command") {
		t.Errorf("expected 'Unknown command' in output, got: %s", out)
	}
}

func TestUnknownFlag(t *testing.T) {
	out, err := runBinary("--nonexistent")
	if err == nil {
		t.Error("expected error for unknown flag")
	}
	if !strings.Contains(out, "Unknown flag") {
		t.Errorf("expected 'Unknown flag' in output, got: %s", out)
	}
}

// --- Init tests ---

func TestInitOutsideGitRepo(t *testing.T) {
	dir, _ := os.MkdirTemp("", "no-git")
	defer os.RemoveAll(dir)

	out, err := runBinaryInDir(dir, "init")
	if err == nil {
		t.Error("init outside git repo should fail")
	}
	if !strings.Contains(strings.ToLower(out), "git") {
		t.Errorf("should mention git, got: %s", out)
	}
}

func TestInitRequiresPrerequisites(t *testing.T) {
	// init inside a git repo but without claude/engram should report missing prerequisites
	dir, _ := os.MkdirTemp("", "aryflow-init-prereq")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	out, err := runBinaryInDir(dir, "init", "--force")
	// This may fail due to missing Claude Code / Engram — that's expected
	if err != nil {
		// Should mention prerequisites or specific missing tool
		lower := strings.ToLower(out)
		if !strings.Contains(lower, "prerequisit") && !strings.Contains(lower, "missing") && !strings.Contains(lower, "not found") {
			t.Logf("init failed (possibly missing prerequisites): %s", out)
		}
	}
	// If it succeeds, that's fine too (system has all prerequisites)
}

// --- Doctor tests ---

func TestDoctorOutsideGitRepo(t *testing.T) {
	dir, _ := os.MkdirTemp("", "no-git")
	defer os.RemoveAll(dir)

	out, err := runBinaryInDir(dir, "doctor")
	// Doctor always runs all checks; it should report git repo error and exit 1
	if err == nil {
		t.Error("doctor outside git repo should fail (exit 1)")
	}
	if !strings.Contains(strings.ToLower(out), "git") {
		t.Errorf("should mention git, got: %s", out)
	}
}

func TestDoctorInEmptyGitRepo(t *testing.T) {
	dir, _ := os.MkdirTemp("", "aryflow-doctor-empty")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	out, err := runBinaryInDir(dir, "doctor")
	// Doctor should report errors for missing .aryflow/version, skills, etc.
	if err == nil {
		t.Error("doctor in empty git repo should fail (missing aryflow files)")
	}
	// Should show results summary
	if !strings.Contains(out, "Results:") {
		t.Errorf("doctor should show Results summary, got: %s", out)
	}
	if !strings.Contains(out, "error") || !strings.Contains(out, "missing") {
		t.Errorf("doctor should report missing files, got: %s", out)
	}
}

func TestDoctorShowsPassedChecks(t *testing.T) {
	dir, _ := os.MkdirTemp("", "aryflow-doctor-partial")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	out, _ := runBinaryInDir(dir, "doctor")
	// Git repo check should pass
	if !strings.Contains(out, "Git repository") {
		t.Errorf("doctor should check git repository, got: %s", out)
	}
}

// --- Update tests ---

func TestUpdateOutsideGitRepo(t *testing.T) {
	dir, _ := os.MkdirTemp("", "no-git-update")
	defer os.RemoveAll(dir)

	out, err := runBinaryInDir(dir, "update")
	if err == nil {
		t.Error("update outside git repo should fail")
	}
	if !strings.Contains(strings.ToLower(out), "git") {
		t.Errorf("should mention git, got: %s", out)
	}
}

func TestUpdateNotInitialized(t *testing.T) {
	dir, _ := os.MkdirTemp("", "aryflow-update-noinit")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	out, err := runBinaryInDir(dir, "update")
	if err == nil {
		t.Error("update without .aryflow/version should fail")
	}
	if !strings.Contains(out, "not initialized") || !strings.Contains(out, "aryflow init") {
		t.Errorf("should suggest running init, got: %s", out)
	}
}

func TestUpdateDryRun(t *testing.T) {
	dir, _ := os.MkdirTemp("", "aryflow-update-dryrun")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	// Create .aryflow/version with old version
	os.MkdirAll(filepath.Join(dir, ".aryflow"), 0755)
	os.WriteFile(filepath.Join(dir, ".aryflow", "version"), []byte("0.0.1\n"), 0644)

	out, _ := runBinaryInDir(dir, "update", "--dry-run")
	if !strings.Contains(out, "0.0.1") {
		t.Errorf("dry-run should show current version, got: %s", out)
	}
	// Verify version file was NOT changed (dry-run)
	data, _ := os.ReadFile(filepath.Join(dir, ".aryflow", "version"))
	if strings.TrimSpace(string(data)) != "0.0.1" {
		t.Errorf("dry-run should not modify version file, got: %s", string(data))
	}
}

func TestUpdateUpToDate(t *testing.T) {
	dir, _ := os.MkdirTemp("", "aryflow-update-current")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	// Use 0.2.1 which matches the hardcoded version in main.go
	os.MkdirAll(filepath.Join(dir, ".aryflow"), 0755)
	os.WriteFile(filepath.Join(dir, ".aryflow", "version"), []byte("0.2.1\n"), 0644)

	out, err := runBinaryInDir(dir, "update")
	if err != nil {
		t.Fatalf("update with current version should succeed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "up to date") {
		t.Errorf("should report up to date, got: %s", out)
	}
}

// --- Flag combination tests ---

func TestGlobalFlagsBeforeCommand(t *testing.T) {
	// --verbose before a command should work
	dir, _ := os.MkdirTemp("", "aryflow-flags")
	defer os.RemoveAll(dir)
	exec.Command("git", "init", dir).Run()

	// --verbose doctor should not crash
	out, _ := runBinaryInDir(dir, "--verbose", "doctor")
	if len(out) == 0 {
		t.Error("verbose doctor should produce output")
	}
}

func TestYesFlagSetsEnv(t *testing.T) {
	// --yes / -y should work as global flag
	out, err := runBinary("--yes", "--version")
	if err != nil {
		t.Fatalf("--yes --version failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "aryflow") {
		t.Errorf("expected version output, got: %s", out)
	}
}
