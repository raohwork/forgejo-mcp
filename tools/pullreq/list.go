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

// ListPullRequestsParams defines the parameters for the list_pull_requests tool.
// It includes options for filtering, sorting, and paginating pull requests.
type ListPullRequestsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// State filters pull requests by their state (e.g., 'open', 'closed').
	State string `json:"state,omitempty"`
	// Sort specifies the sort order for the results.
	Sort string `json:"sort,omitempty"`
	// Direction specifies the sort direction (asc or desc).
	Direction string `json:"direction,omitempty"`
	// Milestone filters pull requests by a milestone title.
	Milestone string `json:"milestone,omitempty"`
	// Labels filters pull requests by a list of label names.
	Labels []string `json:"labels,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of pull requests to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListPullRequestsImpl implements the read-only MCP tool for listing pull requests.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a list
// of pull requests with powerful filtering and sorting capabilities.
type ListPullRequestsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_pull_requests` tool. It requires `owner` and `repo`
// and supports a rich set of optional parameters for filtering and sorting.
// It is marked as a safe, read-only operation.
func (ListPullRequestsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_pull_requests",
		Title:       "List Pull Requests",
		Description: "List pull requests in a repository. Returns PR information including title, state, author, and metadata.",
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
					Description: "PR state filter: 'open', 'closed', or 'all' (optional, defaults to 'open')",
					Enum:        []any{"open", "closed", "all"},
				},
				"sort": {
					Type:        "string",
					Description: "Sort order: 'created', 'updated', 'popularity' (comment count), or 'long-running' (age, filtering by time updated) (optional, defaults to 'created')",
					Enum:        []any{"created", "updated", "popularity", "long-running"},
				},
				"direction": {
					Type:        "string",
					Description: "Sort direction: 'asc' or 'desc' (optional, defaults to 'desc')",
					Enum:        []any{"asc", "desc"},
				},
				"milestone": {
					Type:        "string",
					Description: "Filter by milestone title (optional)",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Filter by label names (optional)",
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of pull requests per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing pull requests. It calls the Forgejo SDK's
// `ListRepoPullRequests` function with the provided filters and formats the results
// into a markdown table.
func (impl ListPullRequestsImpl) Handler() mcp.ToolHandlerFor[ListPullRequestsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListPullRequestsParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Build options for SDK call
		opt := forgejo.ListPullRequestsOptions{}
		if p.State != "" {
			opt.State = forgejo.StateType(p.State)
		}
		if p.Sort != "" {
			opt.Sort = p.Sort
		}
		// Note: Direction is not supported in the SDK options
		// Note: Labels filtering is not directly supported in SDK
		// Note: Milestone string name is not supported, only ID
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		// Call SDK
		prs, _, err := impl.Client.ListRepoPullRequests(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list pull requests: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(prs) == 0 {
			content = "No pull requests found matching the criteria."
		} else {
			// Convert PRs to our type
			prList := make(types.PullRequestList, len(prs))
			for i, pr := range prs {
				prList[i] = &types.PullRequest{PullRequest: pr}
			}

			content = fmt.Sprintf("Found %d pull requests\n\n%s",
				len(prs), prList.ToMarkdown())
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
