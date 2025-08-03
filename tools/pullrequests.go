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

// Pull Request related tools

func ListPullRequests() mcp.Tool {
	return mcp.Tool{
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
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of pull requests per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func GetPullRequest() mcp.Tool {
	return mcp.Tool{
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
