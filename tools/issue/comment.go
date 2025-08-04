// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"fmt"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListIssueCommentsParams defines the parameters for the list_issue_comments tool.
// It specifies the issue and includes optional filters for pagination and time range.
type ListIssueCommentsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Since filters for comments updated after the given time.
	Since time.Time `json:"since,omitempty"`
	// Before filters for comments updated before the given time.
	Before time.Time `json:"before,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of comments to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListIssueCommentsImpl implements the read-only MCP tool for listing issue comments.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a list
// of comments for a specific issue.
type ListIssueCommentsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_issue_comments` tool. It requires `owner`, `repo`,
// and the issue `index`. It supports time-based filtering and pagination and is
// marked as a safe, read-only operation.
func (ListIssueCommentsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of comments per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

// Handler implements the logic for listing issue comments. It calls the Forgejo SDK's
// `ListIssueComments` function and formats the results into a markdown list.
func (impl ListIssueCommentsImpl) Handler() mcp.ToolHandlerFor[ListIssueCommentsParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListIssueCommentsParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		opt := forgejo.ListIssueCommentOptions{}
		if !p.Since.IsZero() {
			opt.Since = p.Since
		}
		if !p.Before.IsZero() {
			opt.Before = p.Before
		}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		comments, _, err := impl.Client.ListIssueComments(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list comments: %w", err)
		}

		var content string
		if len(comments) == 0 {
			content = "No comments found for this issue."
		} else {
			var commentsMarkdown string
			for _, comment := range comments {
				commentWrapper := &types.Comment{Comment: comment}
				commentsMarkdown += commentWrapper.ToMarkdown() + "\n\n---\n\n"
			}
			content = fmt.Sprintf("Found %d comments\n\n%s", len(comments), commentsMarkdown)
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

// CreateIssueCommentParams defines the parameters for the create_issue_comment tool.
// It specifies the issue to comment on and the content of the comment.
type CreateIssueCommentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Body is the markdown content of the comment.
	Body string `json:"body"`
}

// CreateIssueCommentImpl implements the MCP tool for creating a new comment on an issue.
// This is a non-idempotent operation that posts a new comment using the Forgejo SDK.
type CreateIssueCommentImpl struct {
	Client *tools.Client
}

// Definition describes the `create_issue_comment` tool. It requires the issue `index`
// and the comment `body`. It is not idempotent, as multiple calls will create
// multiple identical comments.
func (CreateIssueCommentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_issue_comment",
		Title:       "Add Issue Comment",
		Description: "Add a comment to an existing issue.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
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

// Handler implements the logic for creating an issue comment. It calls the Forgejo
// SDK's `CreateIssueComment` function and returns the details of the new comment.
func (impl CreateIssueCommentImpl) Handler() mcp.ToolHandlerFor[CreateIssueCommentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateIssueCommentParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		opt := forgejo.CreateIssueCommentOption{
			Body: p.Body,
		}

		comment, _, err := impl.Client.CreateIssueComment(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to create comment: %w", err)
		}

		commentWrapper := &types.Comment{Comment: comment}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: commentWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// EditIssueCommentParams defines the parameters for the edit_issue_comment tool.
// It specifies the comment to edit by its ID and the new content.
type EditIssueCommentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// CommentID is the unique identifier of the comment to edit.
	CommentID int `json:"comment_id"`
	// Body is the new markdown content for the comment.
	Body string `json:"body"`
}

// EditIssueCommentImpl implements the MCP tool for editing an existing issue comment.
// This is an idempotent operation that modifies a comment's content using the
// Forgejo SDK.
type EditIssueCommentImpl struct {
	Client *tools.Client
}

// Definition describes the `edit_issue_comment` tool. It requires the `comment_id`
// and the new `body`. It is marked as idempotent.
func (EditIssueCommentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_issue_comment",
		Title:       "Edit Issue Comment",
		Description: "Edit an existing comment on an issue. You can modify the comment body content.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
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

// Handler implements the logic for editing an issue comment. It calls the Forgejo
// SDK's `EditIssueComment` function. It will return an error if the comment ID
// is not found.
func (impl EditIssueCommentImpl) Handler() mcp.ToolHandlerFor[EditIssueCommentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditIssueCommentParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		opt := forgejo.EditIssueCommentOption{
			Body: p.Body,
		}

		comment, _, err := impl.Client.EditIssueComment(p.Owner, p.Repo, int64(p.CommentID), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to edit comment: %w", err)
		}

		commentWrapper := &types.Comment{Comment: comment}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: commentWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// DeleteIssueCommentParams defines the parameters for the delete_issue_comment tool.
// It specifies the comment to be deleted by its ID.
type DeleteIssueCommentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// CommentID is the unique identifier of the comment to delete.
	CommentID int `json:"comment_id"`
}

// DeleteIssueCommentImpl implements the destructive MCP tool for deleting an issue comment.
// This is an idempotent but irreversible operation that removes a comment using the
// Forgejo SDK.
type DeleteIssueCommentImpl struct {
	Client *tools.Client
}

// Definition describes the `delete_issue_comment` tool. It requires the `comment_id`
// to be deleted. It is marked as a destructive operation to ensure clients can
// warn the user before execution.
func (DeleteIssueCommentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_issue_comment",
		Title:       "Delete Issue Comment",
		Description: "Delete a specific comment from an issue. This action cannot be undone.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(true),
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

// Handler implements the logic for deleting an issue comment. It calls the Forgejo
// SDK's `DeleteIssueComment` function. On success, it returns a simple text
// confirmation. It will return an error if the comment does not exist.
func (impl DeleteIssueCommentImpl) Handler() mcp.ToolHandlerFor[DeleteIssueCommentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteIssueCommentParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		_, err := impl.Client.DeleteIssueComment(p.Owner, p.Repo, int64(p.CommentID))
		if err != nil {
			return nil, fmt.Errorf("failed to delete comment: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Comment %d successfully deleted.", p.CommentID),
				},
			},
		}, nil
	}
}
