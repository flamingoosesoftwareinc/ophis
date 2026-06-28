package claudecode

import (
	"os"
	"path/filepath"
)

// ConfigPath returns the default path for Claude Code's project MCP configuration.
// Claude Code reads project-level MCP servers from .mcp.json in the current working
// directory. Falls back to a relative ".mcp.json" if cwd cannot be determined.
func ConfigPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		// Fallback to relative path
		return ".mcp.json"
	}

	return filepath.Join(workingDir, ".mcp.json")
}
