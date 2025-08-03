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

// Issue related tools

func CreateIssue() mcp.Tool {
	return mcp.Tool{
		Name:        "create_issue",
		Title:       "Create Issue",
		Description: "Create a new issue in a repository. Specify title, body, and optional metadata like labels, assignees, milestone.",
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
					Description: "Issue title",
				},
				"body": {
					Type:        "string",
					Description: "Issue body content (markdown supported)",
				},
				"assignees": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Array of usernames to assign to this issue (optional)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID to assign to this issue (optional)",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to assign to this issue (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Issue due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "title", "body"},
		},
	}
}

func EditIssue() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_issue",
		Title:       "Edit Issue",
		Description: "Edit an existing issue's title, body, state, assignees, milestone, or due date.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"title": {
					Type:        "string",
					Description: "New issue title (optional)",
				},
				"body": {
					Type:        "string",
					Description: "New issue body content (markdown supported) (optional)",
				},
				"state": {
					Type:        "string",
					Description: "New issue state: 'open' or 'closed' (optional)",
					Enum:        []any{"open", "closed"},
				},
				"assignees": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Array of usernames to assign to this issue (optional)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID to assign to this issue (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Issue due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

func CreateIssueComment() mcp.Tool {
	return mcp.Tool{
		Name:        "create_issue_comment",
		Title:       "Add Issue Comment",
		Description: "Add a comment to an existing issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"body": {
					Type:        "string",
					Description: "Comment body content (markdown supported)",
				},
			},
			Required: []string{"owner", "repo", "index", "body"},
		},
	}
}

func AddIssueLabels() mcp.Tool {
	return mcp.Tool{
		Name:        "add_issue_labels",
		Title:       "Add Issue Labels",
		Description: "Add labels to an existing issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to add to this issue",
					MinItems:    intPtr(1),
				},
			},
			Required: []string{"owner", "repo", "index", "labels"},
		},
	}
}

func RemoveIssueLabel() mcp.Tool {
	return mcp.Tool{
		Name:        "remove_issue_label",
		Title:       "Remove Issue Label",
		Description: "Remove a specific label from an issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"label": {
					Type:        "integer",
					Description: "Label ID to remove from this issue",
				},
			},
			Required: []string{"owner", "repo", "index", "label"},
		},
	}
}

func ReplaceIssueLabels() mcp.Tool {
	return mcp.Tool{
		Name:        "replace_issue_labels",
		Title:       "Replace Issue Labels",
		Description: "Replace all labels on an issue with a new set of labels.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to set on this issue (replaces all existing labels)",
				},
			},
			Required: []string{"owner", "repo", "index", "labels"},
		},
	}
}

// Issue attachment management tools (Custom HTTP implementation needed)

func ListIssueAttachments() mcp.Tool {
	return mcp.Tool{
		Name:        "list_issue_attachments",
		Title:       "List Issue Attachments",
		Description: "List all attachments on an issue. Returns attachment information including names, sizes, and download URLs.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

func CreateIssueAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "create_issue_attachment",
		Title:       "Add Issue Attachment",
		Description: "Add a file attachment to an issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"file_path": {
					Type:        "string",
					Description: "Absolute path to the file to attach",
				},
				"name": {
					Type:        "string",
					Description: "Optional display name for the attachment (defaults to filename)",
				},
			},
			Required: []string{"owner", "repo", "index", "file_path"},
		},
	}
}

func DeleteIssueAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_issue_attachment",
		Title:       "Delete Issue Attachment",
		Description: "Delete a specific attachment from an issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"attachment_id": {
					Type:        "string",
					Description: "Attachment ID to delete",
				},
			},
			Required: []string{"owner", "repo", "index", "attachment_id"},
		},
	}
}

func EditIssueAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_issue_attachment",
		Title:       "Edit Issue Attachment",
		Description: "Edit an attachment's metadata such as display name.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"attachment_id": {
					Type:        "string",
					Description: "Attachment ID to edit",
				},
				"name": {
					Type:        "string",
					Description: "New display name for the attachment",
				},
			},
			Required: []string{"owner", "repo", "index", "attachment_id", "name"},
		},
	}
}

// Issue dependency management tools (Custom HTTP implementation needed)

func ListIssueDependencies() mcp.Tool {
	return mcp.Tool{
		Name:        "list_issue_dependencies",
		Title:       "List Issue Dependencies",
		Description: "List all dependency relationships for an issue, showing which issues this issue depends on and which issues depend on this issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

func RemoveIssueDependency() mcp.Tool {
	return mcp.Tool{
		Name:        "remove_issue_dependency",
		Title:       "Remove Issue Dependency",
		Description: "Remove a specific dependency relationship between two issues.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"dependency_index": {
					Type:        "integer",
					Description: "Index of the dependency issue to remove",
				},
			},
			Required: []string{"owner", "repo", "index", "dependency_index"},
		},
	}
}

func AddIssueDependency() mcp.Tool {
	return mcp.Tool{
		Name:        "add_issue_dependency",
		Title:       "Add Issue Dependency",
		Description: "Add a dependency relationship between two issues, where one issue depends on another.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"dependency_index": {
					Type:        "integer",
					Description: "Index of the issue this issue depends on",
				},
			},
			Required: []string{"owner", "repo", "index", "dependency_index"},
		},
	}
}
