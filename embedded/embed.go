// Package embedded provides access to all embedded AryFlow files
// (skills, agents, rules, hooks) via Go's embed package.
package embedded

import (
	"embed"
	"io/fs"
)

//go:embed skills agents rules hooks
var content embed.FS

// Skills returns the embedded skills filesystem rooted at "skills/".
func Skills() fs.FS {
	sub, _ := fs.Sub(content, "skills")
	return sub
}

// Agents returns the embedded agents filesystem rooted at "agents/".
func Agents() fs.FS {
	sub, _ := fs.Sub(content, "agents")
	return sub
}

// Rules returns the embedded rules filesystem rooted at "rules/".
func Rules() fs.FS {
	sub, _ := fs.Sub(content, "rules")
	return sub
}

// Hooks returns the embedded hooks filesystem rooted at "hooks/".
func Hooks() fs.FS {
	sub, _ := fs.Sub(content, "hooks")
	return sub
}

// SkillNames lists all skill directories (e.g., "spec-it", "execute-spec").
func SkillNames() []string {
	return listDirs(Skills())
}

// AgentFiles lists all agent markdown files (e.g., "merge-wave.md").
func AgentFiles() []string {
	return listFiles(Agents())
}

// RuleFiles lists all rule markdown files.
func RuleFiles() []string {
	return listFiles(Rules())
}

// HookFiles lists all hook script files.
func HookFiles() []string {
	return listFiles(Hooks())
}

// ReadSkill reads a skill file by path relative to skills/ (e.g., "spec-it/SKILL.md").
func ReadSkill(path string) ([]byte, error) {
	return fs.ReadFile(Skills(), path)
}

// ReadAgent reads an agent file by name (e.g., "merge-wave.md").
func ReadAgent(name string) ([]byte, error) {
	return fs.ReadFile(Agents(), name)
}

// ReadRule reads a rule file by name (e.g., "aryflow.md").
func ReadRule(name string) ([]byte, error) {
	return fs.ReadFile(Rules(), name)
}

// ReadHook reads a hook file by name (e.g., "aryflow-session-start.sh").
func ReadHook(name string) ([]byte, error) {
	return fs.ReadFile(Hooks(), name)
}

// Content returns the full embedded FS (for walking all files).
func Content() embed.FS {
	return content
}

// ManagedFile describes an embedded file and its project destination.
type ManagedFile struct {
	EmbedPath   string // path inside embedded FS (e.g., "skills/spec-it/SKILL.md")
	ProjectPath string // path relative to project root (e.g., ".claude/skills/spec-it/SKILL.md")
}

// ManagedFiles returns the list of all files managed by AryFlow with their
// embedded source path and project destination path.
func ManagedFiles() []ManagedFile {
	var files []ManagedFile

	// Skills
	for _, name := range SkillNames() {
		files = append(files, ManagedFile{
			EmbedPath:   "skills/" + name + "/SKILL.md",
			ProjectPath: ".claude/skills/" + name + "/SKILL.md",
		})
	}

	// Agents
	for _, name := range AgentFiles() {
		files = append(files, ManagedFile{
			EmbedPath:   "agents/" + name,
			ProjectPath: ".claude/agents/" + name,
		})
	}

	// Rules
	for _, name := range RuleFiles() {
		files = append(files, ManagedFile{
			EmbedPath:   "rules/" + name,
			ProjectPath: ".claude/rules/" + name,
		})
	}

	// Hooks
	for _, name := range HookFiles() {
		files = append(files, ManagedFile{
			EmbedPath:   "hooks/" + name,
			ProjectPath: ".claude/hooks/" + name,
		})
	}

	return files
}

// ReadEmbedded reads a file from the embedded FS by its embed path.
func ReadEmbedded(embedPath string) ([]byte, error) {
	return fs.ReadFile(content, embedPath)
}

func listDirs(fsys fs.FS) []string {
	var dirs []string
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}

func listFiles(fsys fs.FS) []string {
	var files []string
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	return files
}
