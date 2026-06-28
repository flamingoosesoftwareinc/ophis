package manager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/njayp/ophis/internal/cfgmgr/manager/claude"
	"github.com/njayp/ophis/internal/cfgmgr/manager/claudecode"
	"github.com/njayp/ophis/internal/cfgmgr/manager/vscode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeCodeManager(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mcp.json")

	t.Run("NewClaudeCodeManager creates empty config", func(t *testing.T) {
		m, err := NewClaudeCodeManager(configPath)
		require.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, configPath, m.configPath)
	})

	t.Run("EnableServer writes .mcp.json with server entry", func(t *testing.T) {
		m, err := NewClaudeCodeManager(configPath)
		require.NoError(t, err)

		server := claudecode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", server)
		require.NoError(t, err)

		// Verify the server was added in memory
		assert.True(t, m.config.HasServer("test-server"))

		// Verify the config was persisted
		data, err := os.ReadFile(configPath)
		require.NoError(t, err)

		var savedConfig claudecode.Config
		err = json.Unmarshal(data, &savedConfig)
		require.NoError(t, err)
		assert.True(t, savedConfig.HasServer("test-server"))
	})

	t.Run("EnableServer second server preserves first", func(t *testing.T) {
		cleanPath := filepath.Join(tmpDir, ".mcp_merge.json")

		m, err := NewClaudeCodeManager(cleanPath)
		require.NoError(t, err)

		first := claudecode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/first",
			Args:    []string{"mcp", "start"},
		}
		err = m.EnableServer("first-server", first)
		require.NoError(t, err)

		// Reload and add second server — simulates a fresh invocation
		m2, err := NewClaudeCodeManager(cleanPath)
		require.NoError(t, err)

		second := claudecode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/second",
			Args:    []string{"mcp", "start"},
		}
		err = m2.EnableServer("second-server", second)
		require.NoError(t, err)

		// Reload and verify both servers are present (merge, not clobber)
		m3, err := NewClaudeCodeManager(cleanPath)
		require.NoError(t, err)
		assert.True(t, m3.config.HasServer("first-server"), "first server must be preserved")
		assert.True(t, m3.config.HasServer("second-server"), "second server must be present")
	})

	t.Run("DisableServer removes existing server", func(t *testing.T) {
		m, err := NewClaudeCodeManager(configPath)
		require.NoError(t, err)

		server := claudecode.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("to-remove", server)
		require.NoError(t, err)

		err = m.DisableServer("to-remove")
		require.NoError(t, err)

		assert.False(t, m.config.HasServer("to-remove"))

		// Reload and verify persistence
		m2, err := NewClaudeCodeManager(configPath)
		require.NoError(t, err)
		assert.False(t, m2.config.HasServer("to-remove"))
	})

	t.Run("DisableServer handles non-existent server", func(t *testing.T) {
		m, err := NewClaudeCodeManager(configPath)
		require.NoError(t, err)

		err = m.DisableServer("non-existent")
		require.NoError(t, err)
	})

	t.Run("loadConfig handles missing file", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "non-existent.json")
		m, err := NewClaudeCodeManager(nonExistentPath)
		require.NoError(t, err)
		assert.NotNil(t, m)
	})

	t.Run("loadConfig handles invalid JSON", func(t *testing.T) {
		invalidPath := filepath.Join(tmpDir, "invalid_claudecode.json")
		err := os.WriteFile(invalidPath, []byte("not valid json"), 0o644)
		require.NoError(t, err)

		_, err = NewClaudeCodeManager(invalidPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})
}

func TestClaudeManager(t *testing.T) {
	// Create a temporary directory for test configs
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "claude_config.json")

	t.Run("NewClaudeManager creates empty config", func(t *testing.T) {
		m, err := NewClaudeManager(configPath)
		require.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, configPath, m.configPath)
	})

	t.Run("EnableServer adds new server", func(t *testing.T) {
		m, err := NewClaudeManager(configPath)
		require.NoError(t, err)

		server := claude.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", server)
		require.NoError(t, err)

		// Verify the server was added
		assert.True(t, m.config.HasServer("test-server"))

		// Verify the config was saved
		data, err := os.ReadFile(configPath)
		require.NoError(t, err)

		var savedConfig claude.Config
		err = json.Unmarshal(data, &savedConfig)
		require.NoError(t, err)
		assert.True(t, savedConfig.HasServer("test-server"))
	})

	t.Run("EnableServer updates existing server", func(t *testing.T) {
		m, err := NewClaudeManager(configPath)
		require.NoError(t, err)

		originalServer := claude.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", originalServer)
		require.NoError(t, err)

		updatedServer := claude.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start", "--log-level", "debug"},
		}

		err = m.EnableServer("test-server", updatedServer)
		require.NoError(t, err)

		// Reload and verify the update
		m2, err := NewClaudeManager(configPath)
		require.NoError(t, err)
		assert.True(t, m2.config.HasServer("test-server"))
	})

	t.Run("DisableServer removes existing server", func(t *testing.T) {
		m, err := NewClaudeManager(configPath)
		require.NoError(t, err)

		server := claude.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", server)
		require.NoError(t, err)

		err = m.DisableServer("test-server")
		require.NoError(t, err)

		// Verify the server was removed
		assert.False(t, m.config.HasServer("test-server"))

		// Reload and verify persistence
		m2, err := NewClaudeManager(configPath)
		require.NoError(t, err)
		assert.False(t, m2.config.HasServer("test-server"))
	})

	t.Run("DisableServer handles non-existent server", func(t *testing.T) {
		m, err := NewClaudeManager(configPath)
		require.NoError(t, err)

		err = m.DisableServer("non-existent")
		require.NoError(t, err)
	})

	t.Run("loadConfig handles missing file", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "non-existent.json")
		m, err := NewClaudeManager(nonExistentPath)
		require.NoError(t, err)
		assert.NotNil(t, m)
	})

	t.Run("loadConfig handles invalid JSON", func(t *testing.T) {
		invalidPath := filepath.Join(tmpDir, "invalid.json")
		err := os.WriteFile(invalidPath, []byte("not valid json"), 0o644)
		require.NoError(t, err)

		_, err = NewClaudeManager(invalidPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})

	t.Run("backupConfig creates backup", func(t *testing.T) {
		backupTestPath := filepath.Join(tmpDir, "backup_test.json")
		m, err := NewClaudeManager(backupTestPath)
		require.NoError(t, err)

		server := claude.Server{
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("backup-test", server)
		require.NoError(t, err)

		// Make a change to trigger backup
		err = m.EnableServer("backup-test-2", server)
		require.NoError(t, err)

		// Verify backup exists
		backupPath := filepath.Join(tmpDir, "backup_test.backup.json")
		_, err = os.Stat(backupPath)
		require.NoError(t, err)
	})
}

func TestVSCodeManager(t *testing.T) {
	// Create a temporary directory for test configs
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "vscode_config.json")

	t.Run("NewVSCodeManager creates empty config", func(t *testing.T) {
		m, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, configPath, m.configPath)
	})

	t.Run("EnableServer adds new server", func(t *testing.T) {
		m, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)

		server := vscode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", server)
		require.NoError(t, err)

		// Verify the server was added
		assert.True(t, m.config.HasServer("test-server"))

		// Verify the config was saved
		data, err := os.ReadFile(configPath)
		require.NoError(t, err)

		var savedConfig vscode.Config
		err = json.Unmarshal(data, &savedConfig)
		require.NoError(t, err)
		assert.True(t, savedConfig.HasServer("test-server"))
	})

	t.Run("EnableServer with environment variables", func(t *testing.T) {
		m, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)

		server := vscode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
			Env: map[string]string{
				"DEBUG": "true",
				"PORT":  "8080",
			},
		}

		err = m.EnableServer("env-test", server)
		require.NoError(t, err)

		// Reload and verify
		m2, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)
		assert.True(t, m2.config.HasServer("env-test"))
	})

	t.Run("DisableServer removes existing server", func(t *testing.T) {
		m, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)

		server := vscode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/myapp",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("test-server", server)
		require.NoError(t, err)

		err = m.DisableServer("test-server")
		require.NoError(t, err)

		// Verify the server was removed
		assert.False(t, m.config.HasServer("test-server"))
	})

	t.Run("Multiple servers can coexist", func(t *testing.T) {
		m, err := NewVSCodeManager(configPath, false)
		require.NoError(t, err)

		server1 := vscode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/app1",
			Args:    []string{"mcp", "start"},
		}

		server2 := vscode.Server{
			Type:    "stdio",
			Command: "/usr/local/bin/app2",
			Args:    []string{"mcp", "start"},
		}

		err = m.EnableServer("server1", server1)
		require.NoError(t, err)

		err = m.EnableServer("server2", server2)
		require.NoError(t, err)

		assert.True(t, m.config.HasServer("server1"))
		assert.True(t, m.config.HasServer("server2"))

		// Disable one and verify the other remains
		err = m.DisableServer("server1")
		require.NoError(t, err)

		assert.False(t, m.config.HasServer("server1"))
		assert.True(t, m.config.HasServer("server2"))
	})
}
