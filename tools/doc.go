// Package tools provides a framework for implementing and registering Model Context Protocol (MCP)
// tools, and an extended Forgejo API client.
//
// This package defines the `ToolImpl` interface and `Register` helper function for MCP tool development.
// Its subpackages contain concrete implementations of these tools.
//
// It also offers `tools.Client`, an extension of the standard Forgejo SDK client, providing
// additional functionalities for interacting with Forgejo API endpoints not fully covered by the SDK.
package tools
