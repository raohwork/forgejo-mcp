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

// GetPullRequestParams defines the parameters for the get_pull_request tool.
// It specifies the pull request to retrieve by its owner, repository, and index.
type GetPullRequestParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
}

// GetPullRequestImpl implements the read-only MCP tool for fetching a single pull request.
// This is a safe, idempotent operation that uses the Forgejo SDK to retrieve
// detailed information about a specific pull request.
type GetPullRequestImpl struct {
	Client *tools.Client
}

// Definition describes the `get_pull_request` tool. It requires `owner`, `repo`,
// and the pull request `index`. It is marked as a safe, read-only operation.
func (GetPullRequestImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_pull_request",
		Title:       "Get Pull Request",
		Description: "Get detailed information about a specific pull request including diff, commits, and review status.",
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
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

// Handler implements the logic for fetching a pull request. It calls the Forgejo
// SDK's `GetPullRequest` function and formats the result into a detailed markdown
// view. It will return an error if the pull request is not found.
func (impl GetPullRequestImpl) Handler() mcp.ToolHandlerFor[GetPullRequestParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetPullRequestParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Call SDK
		pr, _, err := impl.Client.GetPullRequest(p.Owner, p.Repo, int64(p.Index))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get pull request: %w", err)
		}

		// Convert to our type and format
		prWrapper := &types.PullRequest{PullRequest: pr}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: prWrapper.ToMarkdown(),
				},
			},
		}, nil, nil
	}
}
