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

// Milestone related tools

func ListRepoMilestones() mcp.Tool {
	return mcp.Tool{
		Name:        "list_repo_milestones",
		Title:       "List Repository Milestones",
		Description: "List all milestones in a repository. Returns milestone information including title, description, due date, and status.",
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
				"state": {
					Type:        "string",
					Description: "Milestone state filter: 'open', 'closed', or 'all' (optional, defaults to 'open')",
					Enum:        []any{"open", "closed", "all"},
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func CreateMilestone() mcp.Tool {
	return mcp.Tool{
		Name:        "create_milestone",
		Title:       "Create Milestone",
		Description: "Create a new milestone in a repository. Specify the title, description, and optional due date.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false,
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
				"title": {
					Type:        "string",
					Description: "Milestone title",
				},
				"description": {
					Type:        "string",
					Description: "Milestone description (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Milestone due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "title"},
		},
	}
}

func EditMilestone() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_milestone",
		Title:       "Edit Milestone",
		Description: "Edit an existing milestone's title, description, due date, or state.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
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
				"id": {
					Type:        "integer",
					Description: "Milestone ID",
				},
				"title": {
					Type:        "string",
					Description: "New milestone title (optional)",
				},
				"description": {
					Type:        "string",
					Description: "New milestone description (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "New milestone due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
				"state": {
					Type:        "string",
					Description: "New milestone state: 'open' or 'closed' (optional)",
					Enum:        []any{"open", "closed"},
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

func DeleteMilestone() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_milestone",
		Title:       "Delete Milestone",
		Description: "Delete a milestone from a repository.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
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
				"id": {
					Type:        "integer",
					Description: "Milestone ID to delete",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}
