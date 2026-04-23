// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ReplyToReviewCommentParams defines the parameters for the
// reply_to_review_comment tool.
type ReplyToReviewCommentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// ReviewID is the ID of the existing review whose thread the reply is added to.
	ReviewID int `json:"review_id"`
	// Path is the file path of the comment being replied to.
	Path string `json:"path"`
	// Body is the markdown body of the reply.
	Body string `json:"body"`
	// NewLine is the 1-based line number on the new (post-change) side of the diff.
	NewLine int `json:"new_line,omitempty"`
	// OldLine is the 1-based line number on the old (pre-change) side of the diff.
	OldLine int `json:"old_line,omitempty"`
}

// ReplyToReviewCommentImpl implements the MCP tool for adding a reply to an
// existing pull request review thread. Forgejo groups inline comments by file
// path and line number, so posting a new comment to the original review at the
// same path and line continues the conversation in the UI.
type ReplyToReviewCommentImpl struct {
	Client *tools.Client
}

func (ReplyToReviewCommentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "reply_to_review_comment",
		Title:       "Reply to Pull Request Review Comment",
		Description: "Reply to an inline pull request review comment by adding a new comment to the existing review at the same file path and line. Use list_pull_request_review_comments to find the review_id, path, and line of the comment you want to reply to.",
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
					Description: "Pull request index number",
				},
				"review_id": {
					Type:        "integer",
					Description: "ID of the review whose thread to reply to (from list_pull_request_reviews or list_pull_request_review_comments)",
				},
				"path": {
					Type:        "string",
					Description: "File path of the comment being replied to, relative to the repository root",
				},
				"body": {
					Type:        "string",
					Description: "Reply body (markdown supported)",
				},
				"new_line": {
					Type:        "integer",
					Description: "1-based line number on the new (post-change) side of the diff. Set this OR old_line to match the original comment.",
					Minimum:     tools.Float64Ptr(1),
				},
				"old_line": {
					Type:        "integer",
					Description: "1-based line number on the old (pre-change) side of the diff. Set this OR new_line to match the original comment.",
					Minimum:     tools.Float64Ptr(1),
				},
			},
			Required: []string{"owner", "repo", "index", "review_id", "path", "body"},
		},
	}
}

func (impl ReplyToReviewCommentImpl) Handler() mcp.ToolHandlerFor[ReplyToReviewCommentParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ReplyToReviewCommentParams) (*mcp.CallToolResult, any, error) {
		p := args

		if p.NewLine == 0 && p.OldLine == 0 {
			return nil, nil, fmt.Errorf("either new_line or old_line must be provided to locate the comment thread")
		}

		opt := tools.MyCreatePullReviewCommentOptions{
			Body:        p.Body,
			Path:        p.Path,
			NewPosition: int64(p.NewLine),
			OldPosition: int64(p.OldLine),
		}

		comment, err := impl.Client.MyCreatePullReviewComment(p.Owner, p.Repo, int64(p.Index), int64(p.ReviewID), opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to reply to review comment: %w", err)
		}

		wrapper := &types.PullReviewComment{PullReviewComment: comment}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: wrapper.ToMarkdown(),
				},
			},
		}, nil, nil
	}
}
