// Package update implements the "aryflow update" command.
package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/EslavaDev/aryflow/embedded"
	"github.com/EslavaDev/aryflow/internal/ui"
)

// FileChange describes a difference between embedded and on-disk files.
type FileChange struct {
	ProjectPath    string // relative to project root
	EmbedPath      string // path in embedded FS
	Status         string // "changed", "added", "removed"
	LocallyModified bool  // file differs from both old and new embedded versions
}

// CompareVersions returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
//
// Handles semver strings like "0.1.0", "1.2.3".
func CompareVersions(a, b string) int {
	aParts := parseSemver(a)
	bParts := parseSemver(b)

	for i := 0; i < 3; i++ {
		if aParts[i] < bParts[i] {
			return -1
		}
		if aParts[i] > bParts[i] {
			return 1
		}
	}
	return 0
}

func parseSemver(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	var parts [3]int
	segs := strings.SplitN(v, ".", 3)
	for i, s := range segs {
		if i >= 3 {
			break
		}
		// Strip any pre-release suffix (e.g., "1.0.0-rc1")
		s = strings.SplitN(s, "-", 2)[0]
		n, _ := strconv.Atoi(s)
		parts[i] = n
	}
	return parts
}

// DiffFiles compares embedded files against on-disk project files and returns changes.
func DiffFiles(projectRoot string) []FileChange {
	managed := embedded.ManagedFiles()
	var changes []FileChange

	for _, mf := range managed {
		embeddedData, err := embedded.ReadEmbedded(mf.EmbedPath)
		if err != nil {
			continue
		}

		diskPath := filepath.Join(projectRoot, mf.ProjectPath)
		diskData, err := os.ReadFile(diskPath)
		if err != nil {
			// File doesn't exist on disk — it's new/added
			changes = append(changes, FileChange{
				ProjectPath: mf.ProjectPath,
				EmbedPath:   mf.EmbedPath,
				Status:      "added",
			})
			continue
		}

		if !bytes.Equal(embeddedData, diskData) {
			changes = append(changes, FileChange{
				ProjectPath:    mf.ProjectPath,
				EmbedPath:      mf.EmbedPath,
				Status:         "changed",
				LocallyModified: true, // conservative: disk differs from new embedded
			})
		}
		// If equal, no change needed
	}

	return changes
}

// Run executes the project update command.
func Run(force bool, dryRun bool, verbose bool, cliVersion string) int {
	gitRoot := findGitRoot()
	if gitRoot == "" {
		ui.Error("Not inside a git repository.")
		return 1
	}

	// Check if project is initialized
	versionFile := filepath.Join(gitRoot, ".aryflow/version")
	data, err := os.ReadFile(versionFile)
	if err != nil {
		ui.Error(".aryflow/version not found — project not initialized.")
		ui.Suggestion("Run: aryflow init")
		return 1
	}

	projectVersion := strings.TrimSpace(string(data))
	ui.Header("AryFlow Update — Checking for updates...")
	fmt.Printf("  Current: v%s\n", projectVersion)
	fmt.Printf("  Latest:  v%s\n", cliVersion)
	fmt.Println()

	if CompareVersions(projectVersion, cliVersion) == 0 {
		ui.Success("Project is up to date.")
		return 0
	}

	// Show diff summary
	changes := DiffFiles(gitRoot)
	if len(changes) == 0 {
		ui.Success("All files are up to date.")
		// Still update version marker
		if !dryRun {
			writeVersion(versionFile, cliVersion)
			ui.Success(fmt.Sprintf("Version updated to v%s", cliVersion))
		}
		return 0
	}

	fmt.Println("  Changes:")
	hasLocalMods := false
	for _, c := range changes {
		prefix := "~"
		if c.Status == "added" {
			prefix = "+"
		}
		line := fmt.Sprintf("    %s %s", prefix, c.ProjectPath)
		if c.LocallyModified {
			line += " (locally modified)"
			hasLocalMods = true
		}
		fmt.Println(line)
	}
	fmt.Println()

	if dryRun {
		ui.Info("Dry run — no changes applied.")
		return 0
	}

	// Prompt for confirmation
	if !force {
		if hasLocalMods {
			ui.Warning("Some files have been locally modified — update will overwrite your changes.")
			if !ui.PromptDefaultNo("Proceed anyway?") {
				fmt.Println("Update cancelled.")
				return 0
			}
		} else {
			if !ui.Prompt(fmt.Sprintf("Update project files to v%s?", cliVersion)) {
				fmt.Println("Update cancelled.")
				return 0
			}
		}
	}

	// Apply changes
	for _, c := range changes {
		embeddedData, err := embedded.ReadEmbedded(c.EmbedPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to read embedded %s: %v", c.EmbedPath, err))
			continue
		}

		destPath := filepath.Join(gitRoot, c.ProjectPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			ui.Error(fmt.Sprintf("Failed to create directory for %s: %v", c.ProjectPath, err))
			continue
		}

		if err := os.WriteFile(destPath, embeddedData, 0o644); err != nil {
			ui.Error(fmt.Sprintf("Failed to write %s: %v", c.ProjectPath, err))
			continue
		}

		action := "Updated"
		if c.Status == "added" {
			action = "Added"
		}
		ui.Success(fmt.Sprintf("%s %s", action, c.ProjectPath))
	}

	// Update version file
	writeVersion(versionFile, cliVersion)
	ui.Success(fmt.Sprintf("Version updated to v%s", cliVersion))

	fmt.Println()
	fmt.Printf("Project updated to v%s.\n", cliVersion)
	return 0
}

// RunSelf updates the CLI binary itself.
func RunSelf(verbose bool, cliVersion string) int {
	ui.Header("AryFlow Update — Self-update")

	// Detect installation method
	binaryPath, err := os.Executable()
	if err != nil {
		ui.Error(fmt.Sprintf("Cannot determine binary path: %v", err))
		return 1
	}

	// Resolve symlinks for accurate path detection
	resolved, err := filepath.EvalSymlinks(binaryPath)
	if err == nil {
		binaryPath = resolved
	}

	isHomebrew := strings.Contains(binaryPath, "Cellar") || strings.Contains(binaryPath, "homebrew")

	// Check latest version from GitHub
	latestVersion, err := fetchLatestVersion()
	if err != nil {
		if verbose {
			ui.Warning(fmt.Sprintf("Cannot check latest version: %v", err))
		}
		ui.Warning("Cannot reach GitHub to check for updates. Proceeding with install method detection.")
	} else {
		fmt.Printf("  Current: v%s\n", cliVersion)
		fmt.Printf("  Latest:  v%s\n", latestVersion)
		fmt.Println()

		if CompareVersions(cliVersion, latestVersion) >= 0 {
			ui.Success("Already up to date.")
			return 0
		}
	}

	if isHomebrew {
		ui.Info("Detected Homebrew installation. Running: brew upgrade aryflow")
		cmd := exec.Command("brew", "upgrade", "aryflow")
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			ui.Error(fmt.Sprintf("brew upgrade failed: %v", err))
			ui.Suggestion("Try manually: brew upgrade aryflow")
			return 1
		}
		ui.Success("CLI updated via Homebrew.")
	} else {
		ui.Info("Manual installation detected.")
		fmt.Println()
		fmt.Println("  To update, download the latest release from:")
		fmt.Println("  https://github.com/EslavaDev/aryflow/releases/latest")
		fmt.Println()
		fmt.Println("  Or if installed via go install:")
		fmt.Println("  go install github.com/EslavaDev/aryflow/cmd/aryflow@latest")
	}

	return 0
}

// fetchLatestVersion queries GitHub API for the latest release tag.
func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/EslavaDev/aryflow/releases/latest")
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("cannot parse GitHub API response: %w", err)
	}

	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

func writeVersion(path, version string) {
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, []byte(version+"\n"), 0o644)
}

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
