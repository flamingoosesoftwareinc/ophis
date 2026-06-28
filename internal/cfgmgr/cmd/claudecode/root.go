package claudecode

import (
	"github.com/spf13/cobra"
)

// Command creates a new Cobra command for managing Claude Code MCP servers.
// commandName is the Use name of the ophis root command (e.g. "mcp" or "agent"),
// threaded through to enableCommand so that GetCmdPath can locate it.
// defaultEnv is merged into the server env on enable; user --env values take precedence.
func Command(commandName string, defaultEnv map[string]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claude-code",
		Short: "Manage Claude Code MCP servers",
		Long:  "Manage MCP server configuration for Claude Code",
	}

	// Add subcommands
	cmd.AddCommand(enableCommand(commandName, defaultEnv), disableCommand(), listCommand())
	return cmd
}
