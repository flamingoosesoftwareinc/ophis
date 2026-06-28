package claudecode

import (
	"fmt"
	"os"

	"github.com/njayp/ophis/internal/cfgmgr/manager"
	"github.com/spf13/cobra"
)

type disableFlags struct {
	configPath string
	serverName string
}

// disableCommand creates a Cobra command for removing an MCP server from Claude Code.
func disableCommand() *cobra.Command {
	f := &disableFlags{}
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Remove server from Claude Code config",
		Long:  "Remove this application from Claude Code MCP servers",
		RunE: func(_ *cobra.Command, _ []string) error {
			return f.run()
		},
	}

	// Add flags
	flags := cmd.Flags()
	flags.StringVar(&f.configPath, "config-path", "", "Path to Claude Code config file")
	flags.StringVar(&f.serverName, "server-name", "", "Name of the MCP server to remove (default: derived from executable name)")

	return cmd
}

func (f *disableFlags) run() error {
	if f.serverName == "" {
		executablePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to determine executable path: %w", err)
		}

		f.serverName = manager.DeriveServerName(executablePath)
	}

	m, err := manager.NewClaudeCodeManager(f.configPath)
	if err != nil {
		return err
	}

	return m.DisableServer(f.serverName)
}
