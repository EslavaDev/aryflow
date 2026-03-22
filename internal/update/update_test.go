package update

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EslavaDev/aryflow/embedded"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"0.1.0", "0.1.0", 0},
		{"0.1.0", "0.2.0", -1},
		{"0.2.0", "0.1.0", 1},
		{"1.0.0", "0.9.9", 1},
		{"0.9.9", "1.0.0", -1},
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", -1},
		{"1.2.4", "1.2.3", 1},
		{"2.0.0", "1.99.99", 1},
		{"0.0.1", "0.0.2", -1},
		// With v prefix
		{"v0.1.0", "0.1.0", 0},
		{"0.1.0", "v0.1.0", 0},
		{"v1.0.0", "v0.9.0", 1},
	}

	for _, tt := range tests {
		got := CompareVersions(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestParseSemver(t *testing.T) {
	tests := []struct {
		input string
		want  [3]int
	}{
		{"0.1.0", [3]int{0, 1, 0}},
		{"1.2.3", [3]int{1, 2, 3}},
		{"v10.20.30", [3]int{10, 20, 30}},
		{"0.0.0", [3]int{0, 0, 0}},
		{"1.0.0-rc1", [3]int{1, 0, 0}},
		{"2.1", [3]int{2, 1, 0}},
	}

	for _, tt := range tests {
		got := parseSemver(tt.input)
		if got != tt.want {
			t.Errorf("parseSemver(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestDiffFiles_AllMissing(t *testing.T) {
	// In an empty dir, all embedded files should show as "added"
	dir := t.TempDir()

	changes := DiffFiles(dir)
	if len(changes) == 0 {
		t.Fatal("expected changes for empty project dir")
	}

	for _, c := range changes {
		if c.Status != "added" {
			t.Errorf("expected status 'added' for %s, got %q", c.ProjectPath, c.Status)
		}
	}
}

func TestDiffFiles_AllPresent_MatchingContent(t *testing.T) {
	dir := t.TempDir()

	// Write all embedded files to the temp dir with matching content
	managed := embedded.ManagedFiles()
	for _, mf := range managed {
		data, err := embedded.ReadEmbedded(mf.EmbedPath)
		if err != nil {
			t.Fatalf("cannot read embedded %s: %v", mf.EmbedPath, err)
		}
		destPath := filepath.Join(dir, mf.ProjectPath)
		os.MkdirAll(filepath.Dir(destPath), 0o755)
		os.WriteFile(destPath, data, 0o644)
	}

	changes := DiffFiles(dir)
	if len(changes) != 0 {
		t.Errorf("expected no changes when files match, got %d changes", len(changes))
		for _, c := range changes {
			t.Logf("  %s: %s", c.Status, c.ProjectPath)
		}
	}
}

func TestDiffFiles_ModifiedFile(t *testing.T) {
	dir := t.TempDir()

	// Write all files matching, then modify one
	managed := embedded.ManagedFiles()
	if len(managed) == 0 {
		t.Skip("no managed files")
	}

	for _, mf := range managed {
		data, _ := embedded.ReadEmbedded(mf.EmbedPath)
		destPath := filepath.Join(dir, mf.ProjectPath)
		os.MkdirAll(filepath.Dir(destPath), 0o755)
		os.WriteFile(destPath, data, 0o644)
	}

	// Modify the first file
	first := managed[0]
	modPath := filepath.Join(dir, first.ProjectPath)
	os.WriteFile(modPath, []byte("modified content"), 0o644)

	changes := DiffFiles(dir)
	if len(changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(changes))
	}
	if len(changes) > 0 {
		if changes[0].Status != "changed" {
			t.Errorf("expected status 'changed', got %q", changes[0].Status)
		}
		if !changes[0].LocallyModified {
			t.Error("expected LocallyModified to be true")
		}
	}
}

func TestWriteVersion(t *testing.T) {
	dir := t.TempDir()
	versionFile := filepath.Join(dir, ".aryflow", "version")

	writeVersion(versionFile, "0.2.0")

	data, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}
	if string(data) != "0.2.0\n" {
		t.Errorf("expected '0.2.0\\n', got %q", string(data))
	}
}
