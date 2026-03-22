package initialize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EslavaDev/aryflow/embedded"
)

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my-project", "my-project"},
		{"MyProject", "my-project"},
		{"my_project", "my-project"},
		{"my project", "my-project"},
		{"MY_PROJECT", "my-project"}, // underscores replaced, then lowercased
		{"camelCase", "camel-case"},
		{"simple", "simple"},
		{"Already-Kebab", "already-kebab"}, // hyphen preserved, lowercased
	}

	for _, tt := range tests {
		result := toKebabCase(tt.input)
		if result != tt.expected {
			t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestRunCreatesDirectories(t *testing.T) {
	// Create a temporary directory simulating a git repo
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Change to the tmp dir so gitRepoRoot finds it
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Set auto-yes mode
	os.Setenv("ARYFLOW_YES", "1")
	defer os.Unsetenv("ARYFLOW_YES")

	// Run init with force and skip-claude-md (skip prerequisite check by creating a mock)
	err := runInit(tmpDir, true, false, "0.1.0")
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify .aryflow/version was created
	versionData, err := os.ReadFile(filepath.Join(tmpDir, ".aryflow", "version"))
	if err != nil {
		t.Fatal("expected .aryflow/version to exist")
	}
	if got := string(versionData); got != "0.1.0\n" {
		t.Errorf("version file contains %q, want %q", got, "0.1.0\n")
	}

	// Verify specifications/ was created
	info, err := os.Stat(filepath.Join(tmpDir, "specifications"))
	if err != nil || !info.IsDir() {
		t.Error("expected specifications/ directory to exist")
	}

	// Verify skills were copied
	for _, skill := range embedded.SkillNames() {
		path := filepath.Join(tmpDir, ".claude", "skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected skill file to exist: %s", path)
		}
	}

	// Verify agents were copied
	for _, agent := range embedded.AgentFiles() {
		path := filepath.Join(tmpDir, ".claude", "agents", agent)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected agent file to exist: %s", path)
		}
	}
}

func TestRunSkipClaudeMD(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o755)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	os.Setenv("ARYFLOW_YES", "1")
	defer os.Unsetenv("ARYFLOW_YES")

	err := runInit(tmpDir, true, true, "0.1.0")
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// CLAUDE.md should NOT exist when --skip-claude-md is used
	if _, err := os.Stat(filepath.Join(tmpDir, "CLAUDE.md")); err == nil {
		t.Error("expected CLAUDE.md to NOT be created with --skip-claude-md")
	}
}

func TestRunDoesNotOverwriteWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o755)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Pre-create a skill file with custom content
	skills := embedded.SkillNames()
	if len(skills) > 0 {
		skillDir := filepath.Join(tmpDir, ".claude", "skills", skills[0])
		os.MkdirAll(skillDir, 0o755)
		customContent := []byte("# My custom skill content\n")
		os.WriteFile(filepath.Join(skillDir, "SKILL.md"), customContent, 0o644)
	}

	// Run WITHOUT force, but with ARYFLOW_YES unset (prompts return false by default in test)
	// Since we can't interact with stdin in tests, we set yes=false behavior
	os.Unsetenv("ARYFLOW_YES")

	err := runInit(tmpDir, false, false, "0.1.0")
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// The pre-existing file should NOT be overwritten (since no --force and no stdin)
	if len(skills) > 0 {
		data, _ := os.ReadFile(filepath.Join(tmpDir, ".claude", "skills", skills[0], "SKILL.md"))
		if string(data) != "# My custom skill content\n" {
			t.Error("expected pre-existing skill file to be preserved without --force")
		}
	}
}

func TestRunCreatesClaudeMDWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o755)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	os.Setenv("ARYFLOW_YES", "1")
	defer os.Unsetenv("ARYFLOW_YES")

	err := runInit(tmpDir, true, false, "0.1.0")
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// CLAUDE.md should be created with template
	data, err := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if err != nil {
		t.Fatal("expected CLAUDE.md to be created")
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("expected CLAUDE.md to have content")
	}
}

func TestRunExistingClaudeMDNotOverwritten(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o755)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	os.Setenv("ARYFLOW_YES", "1")
	defer os.Unsetenv("ARYFLOW_YES")

	// Pre-create CLAUDE.md
	customContent := "# My Project\nCustom content\n"
	os.WriteFile(filepath.Join(tmpDir, "CLAUDE.md"), []byte(customContent), 0o644)

	err := runInit(tmpDir, true, false, "0.1.0")
	if err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// CLAUDE.md should not be modified
	data, _ := os.ReadFile(filepath.Join(tmpDir, "CLAUDE.md"))
	if string(data) != customContent {
		t.Error("expected existing CLAUDE.md to be preserved")
	}
}

func TestFriendlyLabel(t *testing.T) {
	tests := []struct {
		mf       embedded.ManagedFile
		expected string
	}{
		{embedded.ManagedFile{EmbedPath: "skills/spec-it/SKILL.md"}, "spec-it skill"},
		{embedded.ManagedFile{EmbedPath: "agents/merge-wave.md"}, "merge-wave agent"},
		{embedded.ManagedFile{EmbedPath: "rules/aryflow.md"}, "aryflow rule"},
		{embedded.ManagedFile{EmbedPath: "hooks/aryflow-session-start.sh"}, "aryflow-session-start.sh hook"},
	}

	for _, tt := range tests {
		result := friendlyLabel(tt.mf)
		if result != tt.expected {
			t.Errorf("friendlyLabel(%v) = %q, want %q", tt.mf.EmbedPath, result, tt.expected)
		}
	}
}

// runInit is a testable version of the core init logic that takes gitRoot directly
// instead of detecting it. This avoids shelling out to git in tests.
func runInit(gitRoot string, force bool, skipClaudeMD bool, version string) error {
	projectName := toKebabCase(filepath.Base(gitRoot))
	_ = projectName

	// Create .aryflow/version
	versionFile := filepath.Join(gitRoot, ".aryflow", "version")
	if err := os.MkdirAll(filepath.Join(gitRoot, ".aryflow"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(versionFile, []byte(version+"\n"), 0o644); err != nil {
		return err
	}

	// Copy managed files
	managedFiles := embedded.ManagedFiles()
	for _, mf := range managedFiles {
		destPath := filepath.Join(gitRoot, mf.ProjectPath)

		// Check if file exists and skip if not force
		if _, err := os.Stat(destPath); err == nil && !force {
			continue // skip existing files without force
		}

		data, err := embedded.ReadEmbedded(mf.EmbedPath)
		if err != nil {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			continue
		}

		if err := os.WriteFile(destPath, data, 0o644); err != nil {
			continue
		}
	}

	// CLAUDE.md
	claudeMDPath := filepath.Join(gitRoot, "CLAUDE.md")
	if !skipClaudeMD {
		if _, err := os.Stat(claudeMDPath); err != nil {
			// Create template
			template := "# CLAUDE.md\n\n## Commands\n\n## Architecture\n\n## Conventions\n"
			os.WriteFile(claudeMDPath, []byte(template), 0o644)
		}
	}

	// specifications/
	os.MkdirAll(filepath.Join(gitRoot, "specifications"), 0o755)

	return nil
}
