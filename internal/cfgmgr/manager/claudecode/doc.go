// Package claudecode provides configuration management for Claude Code MCP servers.
//
// This package handles:
//   - Claude Code project configuration (.mcp.json in the current working directory)
//   - MCP server entry management
//
// Claude Code reads project-level MCP servers from .mcp.json in the current working
// directory. No per-OS variants are needed since the path is cwd-relative.
//
// This is an internal package and should not be imported by users of the ophis library.
package claudecode
