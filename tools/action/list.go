// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package action

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListActionTasksParams defines the parameters for the list_action_tasks tool.
// It includes basic pagination options for Forgejo Actions tasks.
type ListActionTasksParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of tasks to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListActionTasksImpl implements the read-only MCP tool for listing Forgejo Actions tasks.
// This is a safe, idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type ListActionTasksImpl struct {
	Client *tools.Client
}

// Definition describes the `list_action_tasks` tool. It requires `owner` and `repo`
// and supports various optional parameters for filtering. It is marked as a safe,
// read-only operation.
func (ListActionTasksImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_action_tasks",
		Title:       "List Action Tasks",
		Description: "List Forgejo Actions execution tasks in a repository with basic pagination support.",
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
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of tasks per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing action tasks. It performs a custom HTTP
// GET request to the `/repos/{owner}/{repo}/actions/tasks` endpoint and formats
// the results into a markdown table.
func (impl ListActionTasksImpl) Handler() mcp.ToolHandlerFor[ListActionTasksParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListActionTasksParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call custom client method
		response, err := impl.Client.MyListActionTasks(p.Owner, p.Repo)
		if err != nil {
			return nil, fmt.Errorf("failed to list action tasks: %w", err)
		}

		// Convert to our types and format
		var content string
		if response.TotalCount == 0 || len(response.WorkflowRuns) == 0 {
			content = "No action tasks found in this repository."
		} else {
			// Convert response to our type
			taskList := types.ActionTaskList{
				MyActionTaskResponse: response,
			}

			content = fmt.Sprintf("Found %d action tasks\n\n%s",
				response.TotalCount, taskList.ToMarkdown())
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}
