// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolImpl defines the interface that every tool implementation must satisfy.
// This interface standardizes how tools are defined and handled, ensuring they
// can be registered with the MCP server consistently.
//
// The generic types In and Out represent the tool's input and output data structures.
type ToolImpl[In, Out any] interface {
	// Definition returns the formal MCP tool definition, including its name,
	// description, and input schema.
	Definition() *mcp.Tool

	// Handler returns the function that contains the core logic of the tool.
	// This function is executed when the tool is called by an MCP client.
	Handler() mcp.ToolHandlerFor[In, Out]
}

// Register is a helper function that registers a tool implementation with the MCP server.
// It retrieves the tool's definition and handler through the ToolImpl interface
// and adds them to the server's tool registry.
func Register[I, O any](s *mcp.Server, i ToolImpl[I, O]) {
	mcp.AddTool(s, i.Definition(), i.Handler())
}
