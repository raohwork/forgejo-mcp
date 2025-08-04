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

func ListRepoIssues() mcp.Tool {
	return mcp.Tool{
		Name:        "list_repo_issues",
		Title:       "List Repository Issues",
		Description: "List issues in a repository with optional filtering by state, labels, milestones, assignees, and search terms.",
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
					Description: "Issue state filter: 'open', 'closed', or 'all' (optional, defaults to 'open')",
					Enum:        []any{"open", "closed", "all"},
				},
				"labels": {
					Type:        "string",
					Description: "Comma-separated list of label names to filter by (optional)",
				},
				"milestones": {
					Type:        "string",
					Description: "Comma-separated list of milestone names or IDs to filter by (optional)",
				},
				"assignees": {
					Type:        "string",
					Description: "Comma-separated list of usernames to filter by assignee (optional)",
				},
				"q": {
					Type:        "string",
					Description: "Search query string to filter issues (optional)",
				},
				"sort": {
					Type:        "string",
					Description: "Sort field: 'created', 'updated', 'comments' (optional, defaults to 'created')",
					Enum:        []any{"created", "updated", "comments"},
				},
				"order": {
					Type:        "string",
					Description: "Sort order: 'asc' or 'desc' (optional, defaults to 'desc')",
					Enum:        []any{"asc", "desc"},
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of issues per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func GetIssue() mcp.Tool {
	return mcp.Tool{
		Name:        "get_issue",
		Title:       "Get Issue Details",
		Description: "Get detailed information about a specific issue, including title, body, state, assignees, labels, and metadata.",
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

func ListIssueComments() mcp.Tool {
	return mcp.Tool{
		Name:        "list_issue_comments",
		Title:       "List Issue Comments",
		Description: "List all comments on a specific issue, including comment body, author, and timestamps.",
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
				"since": {
					Type:        "string",
					Description: "Only show comments updated after this time (ISO 8601 format) (optional)",
					Format:      "date-time",
				},
				"before": {
					Type:        "string",
					Description: "Only show comments updated before this time (ISO 8601 format) (optional)",
					Format:      "date-time",
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of comments per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

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

func EditIssueComment() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_issue_comment",
		Title:       "Edit Issue Comment",
		Description: "Edit an existing comment on an issue. You can modify the comment body content.",
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
				"comment_id": {
					Type:        "integer",
					Description: "Comment ID to edit",
				},
				"body": {
					Type:        "string",
					Description: "New comment body content (markdown supported)",
				},
			},
			Required: []string{"owner", "repo", "comment_id", "body"},
		},
	}
}

func DeleteIssueComment() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_issue_comment",
		Title:       "Delete Issue Comment",
		Description: "Delete a specific comment from an issue. This action cannot be undone.",
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
				"comment_id": {
					Type:        "integer",
					Description: "Comment ID to delete",
				},
			},
			Required: []string{"owner", "repo", "comment_id"},
		},
	}
}
