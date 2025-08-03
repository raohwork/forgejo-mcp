// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Forgejo Actions (CI/CD) related tools (Custom HTTP implementation needed)

func ListActionTasks() mcp.Tool {
	return mcp.Tool{
		Name:        "list_action_tasks",
		Title:       "List Action Tasks",
		Description: "List Forgejo Actions execution tasks in a repository. Returns task information including status, workflow, and execution details.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"owner": {
					Type:        "string",
					Description: "Repository owner (username or organization name)",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"status": {
					Type:        "string",
					Description: "Filter by task status: 'success', 'failure', 'cancelled', 'skipped', 'running', 'waiting', or 'blocked' (optional)",
					Enum:        []any{"success", "failure", "cancelled", "skipped", "running", "waiting", "blocked"},
				},
				"workflow": {
					Type:        "string",
					Description: "Filter by workflow name (optional)",
				},
				"actor": {
					Type:        "string",
					Description: "Filter by the user who triggered the action (optional)",
				},
				"branch": {
					Type:        "string",
					Description: "Filter by branch name (optional)",
				},
				"event": {
					Type:        "string",
					Description: "Filter by trigger event: 'push', 'pull_request', 'schedule', etc. (optional)",
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of tasks per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}
