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

// ListPullRequestReviewCommentsParams defines the parameters for the
// list_pull_request_review_comments tool.
type ListPullRequestReviewCommentsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// ReviewID is the ID of the review whose inline comments should be listed.
	ReviewID int `json:"review_id"`
}

// ListPullRequestReviewCommentsImpl implements the read-only MCP tool for listing
// inline comments belonging to a specific pull request review. Each comment
// includes the file path, line number, surrounding diff hunk, author, and body.
type ListPullRequestReviewCommentsImpl struct {
	Client *tools.Client
}

func (ListPullRequestReviewCommentsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_pull_request_review_comments",
		Title:       "List Pull Request Review Comments",
		Description: "List inline review comments for a specific review on a pull request. Each comment includes the file path, line number, diff hunk for context, author, and body.",
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
					Description: "Pull request index number",
				},
				"review_id": {
					Type:        "integer",
					Description: "Review ID (from list_pull_request_reviews) whose inline comments to fetch",
				},
			},
			Required: []string{"owner", "repo", "index", "review_id"},
		},
	}
}

func (impl ListPullRequestReviewCommentsImpl) Handler() mcp.ToolHandlerFor[ListPullRequestReviewCommentsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListPullRequestReviewCommentsParams) (*mcp.CallToolResult, any, error) {
		p := args

		comments, _, err := impl.Client.ListPullReviewComments(p.Owner, p.Repo, int64(p.Index), int64(p.ReviewID))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list review comments: %w", err)
		}

		var content string
		if len(comments) == 0 {
			content = "No inline comments found for this review."
		} else {
			list := make(types.PullReviewCommentList, len(comments))
			for i, c := range comments {
				list[i] = &types.PullReviewComment{PullReviewComment: c}
			}
			content = fmt.Sprintf("Found %d inline comments\n\n%s", len(comments), list.ToMarkdown())
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil, nil
	}
}
