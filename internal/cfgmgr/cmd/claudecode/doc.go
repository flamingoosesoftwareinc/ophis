// Package claudecode provides CLI commands for managing Claude Code MCP servers.
//
// This package implements the 'mcp claude-code' subcommands:
//   - enable: Add MCP server to Claude Code configuration
//   - disable: Remove MCP server from Claude Code configuration
//   - list: Show all configured MCP servers
//
// Claude Code reads project-level MCP servers from .mcp.json in the current
// working directory.
//
// This is an internal package and should not be imported directly by users of the ophis library.
// These commands are automatically available when using ophis.Command() in your application.
package claudecode
