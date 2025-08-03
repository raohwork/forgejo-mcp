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

// Repository related tools

func SearchRepositories() mcp.Tool {
	return mcp.Tool{
		Name:        "search_repositories",
		Title:       "Search Repositories",
		Description: "Search for repositories across the Forgejo instance. Returns repository information including name, description, and metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"q": {
					Type:        "string",
					Description: "Search query string",
				},
				"topic": {
					Type:        "boolean",
					Description: "Whether to search in repository topics (optional, defaults to true)",
				},
				"include_desc": {
					Type:        "boolean",
					Description: "Whether to include repository descriptions in search (optional, defaults to true)",
				},
				"sort": {
					Type:        "string",
					Description: "Sort order: 'alpha', 'created', 'updated', 'size', 'id' (optional, defaults to 'alpha')",
					Enum:        []any{"alpha", "created", "updated", "size", "id"},
				},
				"order": {
					Type:        "string",
					Description: "Sort direction: 'asc' or 'desc' (optional, defaults to 'asc')",
					Enum:        []any{"asc", "desc"},
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 10, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"q"},
		},
	}
}

func ListMyRepositories() mcp.Tool {
	return mcp.Tool{
		Name:        "list_my_repositories",
		Title:       "List My Repositories",
		Description: "List repositories owned by the authenticated user. Returns repository information including name, description, and metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"affiliation": {
					Type:        "string",
					Description: "Repository affiliation filter: 'owner', 'collaborator', 'organization_member', or 'all' (optional, defaults to 'all')",
					Enum:        []any{"owner", "collaborator", "organization_member", "all"},
				},
				"visibility": {
					Type:        "string",
					Description: "Repository visibility filter: 'all', 'public', 'private' (optional, defaults to 'all')",
					Enum:        []any{"all", "public", "private"},
				},
				"sort": {
					Type:        "string",
					Description: "Sort order: 'created', 'updated', 'pushed', 'full_name' (optional, defaults to 'full_name')",
					Enum:        []any{"created", "updated", "pushed", "full_name"},
				},
				"direction": {
					Type:        "string",
					Description: "Sort direction: 'asc' or 'desc' (optional, defaults to 'asc')",
					Enum:        []any{"asc", "desc"},
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{},
		},
	}
}

func ListOrgRepositories() mcp.Tool {
	return mcp.Tool{
		Name:        "list_org_repositories",
		Title:       "List Organization Repositories",
		Description: "List repositories owned by a specific organization. Returns repository information including name, description, and metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"org": {
					Type:        "string",
					Description: "Organization name",
				},
				"type": {
					Type:        "string",
					Description: "Repository type filter: 'all', 'public', 'private', 'forks', 'sources', 'member' (optional, defaults to 'all')",
					Enum:        []any{"all", "public", "private", "forks", "sources", "member"},
				},
				"sort": {
					Type:        "string",
					Description: "Sort order: 'created', 'updated', 'pushed', 'full_name' (optional, defaults to 'created')",
					Enum:        []any{"created", "updated", "pushed", "full_name"},
				},
				"direction": {
					Type:        "string",
					Description: "Sort direction: 'asc' or 'desc' (optional, defaults to 'desc')",
					Enum:        []any{"asc", "desc"},
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 20, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"org"},
		},
	}
}
