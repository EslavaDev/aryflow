// Package initialize implements the "aryflow init" command.
package initialize

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/EslavaDev/aryflow/embedded"
	"github.com/EslavaDev/aryflow/internal/checks"
	"github.com/EslavaDev/aryflow/internal/ui"
)

// Run executes the init command, setting up AryFlow in the current project.
func Run(force bool, skipClaudeMD bool, verbose bool, version string) error {
	// 1. Check preconditions: must be in git repo
	gitRoot, err := gitRepoRoot()
	if err != nil {
		return fmt.Errorf("Not inside a git repository. Run `git init` first.")
	}

	// 2. Check that setup has been done (key prerequisites)
	if err := checkPrerequisites(); err != nil {
		return err
	}

	// 3. Detect project name
	projectName := toKebabCase(filepath.Base(gitRoot))

	ui.Header(fmt.Sprintf("AryFlow Init — Setting up project %q", projectName))

	// 4. Check if already initialized
	versionFile := filepath.Join(gitRoot, ".aryflow", "version")
	if data, err := os.ReadFile(versionFile); err == nil {
		existingVer := strings.TrimSpace(string(data))
		ui.Warning(fmt.Sprintf("Project already initialized (v%s). Re-initialize?", existingVer))
		if !ui.PromptDefaultNo("Re-initialize?") {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// 5. Create .aryflow/version
	if err := os.MkdirAll(filepath.Join(gitRoot, ".aryflow"), 0o755); err != nil {
		return fmt.Errorf("failed to create .aryflow/: %w", err)
	}
	if err := os.WriteFile(versionFile, []byte(version+"\n"), 0o644); err != nil {
		return fmt.Errorf("failed to write .aryflow/version: %w", err)
	}
	ui.Success(fmt.Sprintf("Created .aryflow/version (v%s)", version))

	// 6. Check for existing non-AryFlow files in .claude/skills/
	warnExistingSkills(gitRoot)

	// 7. Copy all managed files (skills, agents, rules, hooks)
	managedFiles := embedded.ManagedFiles()
	for _, mf := range managedFiles {
		destPath := filepath.Join(gitRoot, mf.ProjectPath)
		friendlyName := friendlyLabel(mf)

		// Check if file already exists
		if _, err := os.Stat(destPath); err == nil && !force {
			if !ui.PromptDefaultNo(fmt.Sprintf("Overwrite %s?", friendlyName)) {
				ui.Info(fmt.Sprintf("Skipped %s", friendlyName))
				continue
			}
		}

		// Read embedded content
		data, err := embedded.ReadEmbedded(mf.EmbedPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to read embedded %s: %v", mf.EmbedPath, err))
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			ui.Error(fmt.Sprintf("Failed to create directory for %s: %v", friendlyName, err))
			continue
		}

		// Write file
		perm := os.FileMode(0o644)
		// Make hook files executable
		if strings.Contains(mf.ProjectPath, "hooks/") && strings.HasSuffix(mf.ProjectPath, ".sh") {
			perm = 0o755
		}
		if err := os.WriteFile(destPath, data, perm); err != nil {
			ui.Error(fmt.Sprintf("Failed to write %s: %v", friendlyName, err))
			continue
		}

		ui.Success(fmt.Sprintf("Copied %s", friendlyName))
	}

	// 8. Merge settings.json (hooks + statusLine)
	if err := mergeSettings(gitRoot); err != nil {
		ui.Error(fmt.Sprintf("Failed to merge settings.json: %v", err))
	}

	// 9. Handle CLAUDE.md
	claudeMDPath := filepath.Join(gitRoot, "CLAUDE.md")
	if skipClaudeMD {
		ui.Info("CLAUDE.md — skipped (--skip-claude-md)")
	} else if _, err := os.Stat(claudeMDPath); err == nil {
		ui.Success("CLAUDE.md found — skipping (manual configuration)")
	} else {
		// Create minimal template
		template := fmt.Sprintf(`# CLAUDE.md — %s

## Commands

<!-- Add your project's dev commands here -->

## Architecture

<!-- Describe your project structure -->

## Conventions

<!-- List coding conventions and rules -->

> Configure this file with your project's conventions.
> See: https://github.com/EslavaDev/aryflow
`, projectName)
		if err := os.WriteFile(claudeMDPath, []byte(template), 0o644); err != nil {
			ui.Error(fmt.Sprintf("Failed to create CLAUDE.md: %v", err))
		} else {
			ui.Success("Created CLAUDE.md template")
		}
	}

	// 10. Create specifications/ directory
	specsDir := filepath.Join(gitRoot, "specifications")
	if _, err := os.Stat(specsDir); err == nil {
		ui.Success("specifications/ already exists")
	} else {
		if err := os.MkdirAll(specsDir, 0o755); err != nil {
			ui.Error(fmt.Sprintf("Failed to create specifications/: %v", err))
		} else {
			ui.Success("Created specifications/")
		}
	}

	// 11. Final message
	fmt.Println()
	ui.Success("AryFlow initialized. Start with: /spec-it {feature-name}")

	return nil
}

// gitRepoRoot returns the root of the current git repository.
func gitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// checkPrerequisites verifies that critical dependencies are available.
func checkPrerequisites() error {
	// Check the most critical ones
	criticalChecks := []struct {
		name  string
		check func() (string, error)
	}{
		{"Git", checks.CheckGit},
		{"Claude Code", checks.CheckClaude},
		{"Engram", checks.CheckEngram},
	}

	var missing []string
	for _, c := range criticalChecks {
		if _, err := c.check(); err != nil {
			missing = append(missing, c.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing prerequisites: %s. Run `aryflow setup` first", strings.Join(missing, ", "))
	}
	return nil
}

// toKebabCase converts a string to kebab-case.
func toKebabCase(s string) string {
	// Replace underscores, spaces, and camelCase boundaries with hyphens
	re := regexp.MustCompile(`[_\s]+`)
	s = re.ReplaceAllString(s, "-")

	// Insert hyphens before uppercase letters (camelCase → kebab-case)
	re2 := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re2.ReplaceAllString(s, "${1}-${2}")

	return strings.ToLower(s)
}

// warnExistingSkills checks if .claude/skills/ has non-AryFlow files.
func warnExistingSkills(gitRoot string) {
	skillsDir := filepath.Join(gitRoot, ".claude", "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return // directory doesn't exist yet, that's fine
	}

	aryflowSkills := make(map[string]bool)
	for _, name := range embedded.SkillNames() {
		aryflowSkills[name] = true
	}

	for _, entry := range entries {
		if entry.IsDir() && !aryflowSkills[entry.Name()] {
			ui.Warning("Existing skills found. AryFlow will add its skills alongside them.")
			return
		}
	}
}

// mergeSettings merges the AryFlow settings template into the project's .claude/settings.json.
// If settings.json doesn't exist, it writes the template directly.
// If it exists, it merges hooks (adding missing ones) and statusLine (if not present).
func mergeSettings(gitRoot string) error {
	templateData, err := embedded.ReadSettingsTemplate()
	if err != nil {
		return fmt.Errorf("reading settings template: %w", err)
	}

	var template map[string]interface{}
	if err := json.Unmarshal(templateData, &template); err != nil {
		return fmt.Errorf("parsing settings template: %w", err)
	}

	settingsPath := filepath.Join(gitRoot, ".claude", "settings.json")

	// Ensure .claude/ directory exists
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return fmt.Errorf("creating .claude directory: %w", err)
	}

	existingData, err := os.ReadFile(settingsPath)
	if err != nil {
		// File doesn't exist — write template directly
		formatted, err := json.MarshalIndent(template, "", "  ")
		if err != nil {
			return fmt.Errorf("formatting settings: %w", err)
		}
		if err := os.WriteFile(settingsPath, append(formatted, '\n'), 0o644); err != nil {
			return err
		}
		ui.Success("Created .claude/settings.json with AryFlow hooks")
		return nil
	}

	// File exists — merge
	var existing map[string]interface{}
	if err := json.Unmarshal(existingData, &existing); err != nil {
		return fmt.Errorf("parsing existing settings.json: %w", err)
	}

	merged := false

	// Merge hooks
	templateHooks, _ := template["hooks"].(map[string]interface{})
	existingHooks, ok := existing["hooks"].(map[string]interface{})
	if !ok {
		// No hooks in existing — add all template hooks
		existing["hooks"] = template["hooks"]
		merged = true
	} else {
		// Add missing hook events from template
		for key, val := range templateHooks {
			if _, exists := existingHooks[key]; !exists {
				existingHooks[key] = val
				merged = true
			}
		}
	}

	// Merge statusLine if not present
	if _, ok := existing["statusLine"]; !ok {
		if tmplStatus, ok := template["statusLine"]; ok {
			existing["statusLine"] = tmplStatus
			merged = true
		}
	}

	if !merged {
		ui.Info("settings.json already has all AryFlow hooks")
		return nil
	}

	formatted, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("formatting merged settings: %w", err)
	}
	if err := os.WriteFile(settingsPath, append(formatted, '\n'), 0o644); err != nil {
		return err
	}
	ui.Success("Merged AryFlow hooks into .claude/settings.json")
	return nil
}

// friendlyLabel returns a human-readable label for a managed file.
func friendlyLabel(mf embedded.ManagedFile) string {
	if strings.HasPrefix(mf.EmbedPath, "skills/") {
		// "skills/spec-it/SKILL.md" → "spec-it skill"
		parts := strings.Split(mf.EmbedPath, "/")
		if len(parts) >= 2 {
			return parts[1] + " skill"
		}
	}
	if strings.HasPrefix(mf.EmbedPath, "agents/") {
		// "agents/merge-wave.md" → "merge-wave agent"
		name := strings.TrimPrefix(mf.EmbedPath, "agents/")
		name = strings.TrimSuffix(name, ".md")
		return name + " agent"
	}
	if strings.HasPrefix(mf.EmbedPath, "rules/") {
		name := strings.TrimPrefix(mf.EmbedPath, "rules/")
		name = strings.TrimSuffix(name, ".md")
		return name + " rule"
	}
	if strings.HasPrefix(mf.EmbedPath, "hooks/") {
		name := strings.TrimPrefix(mf.EmbedPath, "hooks/")
		return name + " hook"
	}
	return mf.EmbedPath
}
