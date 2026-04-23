// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListPullRequestReviewsParams defines the parameters for the
// list_pull_request_reviews tool.
type ListPullRequestReviewsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of reviews to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListPullRequestReviewsImpl implements the read-only MCP tool for listing all
// reviews on a pull request. Each review summary includes the inline-comment
// count so the caller can then drill into specific reviews.
type ListPullRequestReviewsImpl struct {
	Client *tools.Client
}

func (ListPullRequestReviewsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_pull_request_reviews",
		Title:       "List Pull Request Reviews",
		Description: "List all reviews on a pull request, including reviewer, state (approved, changes requested, commented), overall body, and inline comment counts.",
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
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of reviews per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

func (impl ListPullRequestReviewsImpl) Handler() mcp.ToolHandlerFor[ListPullRequestReviewsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListPullRequestReviewsParams) (*mcp.CallToolResult, any, error) {
		p := args

		opt := forgejo.ListPullReviewsOptions{}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		reviews, _, err := impl.Client.ListPullReviews(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list reviews: %w", err)
		}

		var content string
		if len(reviews) == 0 {
			content = "No reviews found for this pull request."
		} else {
			list := make(types.PullReviewList, len(reviews))
			for i, r := range reviews {
				list[i] = &types.PullReview{PullReview: r}
			}
			content = fmt.Sprintf("Found %d reviews\n\n%s", len(reviews), list.ToMarkdown())
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
