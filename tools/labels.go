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

// Label related tools

func ListRepoLabels() mcp.Tool {
	return mcp.Tool{
		Name:        "list_repo_labels",
		Title:       "List Repository Labels",
		Description: "List all labels available in a repository. Returns label information including name, description, color, and ID.",
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
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func CreateLabel() mcp.Tool {
	return mcp.Tool{
		Name:        "create_label",
		Title:       "Create Label",
		Description: "Create a new label in a repository. Specify the label name, description, and color.",
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
				"name": {
					Type:        "string",
					Description: "Label name",
				},
				"color": {
					Type:        "string",
					Description: "Label color (hex color code without #, e.g., 'ff0000' for red)",
				},
				"description": {
					Type:        "string",
					Description: "Optional label description",
				},
			},
			Required: []string{"owner", "repo", "name", "color"},
		},
	}
}

func EditLabel() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_label",
		Title:       "Edit Label",
		Description: "Edit an existing label's name, description, or color.",
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
					Description: "Label ID",
				},
				"name": {
					Type:        "string",
					Description: "New label name (optional)",
				},
				"color": {
					Type:        "string",
					Description: "New label color (hex color code without #, e.g., 'ff0000' for red) (optional)",
				},
				"description": {
					Type:        "string",
					Description: "New label description (optional)",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

func DeleteLabel() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_label",
		Title:       "Delete Label",
		Description: "Delete a label from a repository.",
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
					Description: "Label ID to delete",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}
