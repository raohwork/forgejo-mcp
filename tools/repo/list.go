// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package repo

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// SearchRepositoriesParams defines the parameters for the search_repositories tool.
// It includes the search query and various options for filtering and sorting.
type SearchRepositoriesParams struct {
	// Q is the search query string.
	Q string `json:"q"`
	// Topic indicates whether to search in repository topics.
	Topic bool `json:"topic,omitempty"`
	// IncludeDesc indicates whether to include repository descriptions in the search.
	IncludeDesc bool `json:"include_desc,omitempty"`
	// Sort specifies the sort order for the results.
	Sort string `json:"sort,omitempty"`
	// Order specifies the sort direction (asc or desc).
	Order string `json:"order,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of repositories to return per page.
	Limit int `json:"limit,omitempty"`
}

// SearchRepositoriesImpl implements the read-only MCP tool for searching repositories.
// This operation is safe, idempotent, and uses the Forgejo SDK to find repositories
// across the entire Forgejo instance based on a query string.
type SearchRepositoriesImpl struct {
	Client *tools.Client
}

// Definition describes the `search_repositories` tool. It requires a search query `q`
// and supports various optional parameters for sorting and pagination. It is
// marked as a safe, read-only operation.
func (SearchRepositoriesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 10, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"q"},
		},
	}
}

// Handler implements the logic for searching repositories. It calls the Forgejo SDK's
// `SearchRepos` function and formats the results into a markdown list.
func (impl SearchRepositoriesImpl) Handler() mcp.ToolHandlerFor[SearchRepositoriesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[SearchRepositoriesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.SearchRepoOptions{
			Keyword: p.Q,
		}
		if p.Topic {
			opt.KeywordIsTopic = p.Topic
		}
		if p.IncludeDesc {
			opt.KeywordInDescription = p.IncludeDesc
		}
		if p.Sort != "" {
			opt.Sort = p.Sort
		}
		if p.Order != "" {
			opt.Order = p.Order
		}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		// Call SDK
		repos, _, err := impl.Client.SearchRepos(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to search repositories: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(repos) == 0 {
			content = "No repositories found matching the search criteria."
		} else {
			// Convert repos to our type
			repoList := make(types.RepositoryList, len(repos))
			for i, repo := range repos {
				repoList[i] = &types.Repository{Repository: repo}
			}

			content = fmt.Sprintf("Found %d repositories\n\n%s",
				len(repos), repoList.ToMarkdown())
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

// ListMyRepositoriesParams defines the parameters for the list_my_repositories tool.
// It allows filtering and sorting of the authenticated user's repositories.
type ListMyRepositoriesParams struct {
	// Affiliation filters repositories by the user's role.
	Affiliation string `json:"affiliation,omitempty"`
	// Visibility filters repositories by their public or private status.
	Visibility string `json:"visibility,omitempty"`
	// Sort specifies the sort order for the results.
	Sort string `json:"sort,omitempty"`
	// Direction specifies the sort direction (asc or desc).
	Direction string `json:"direction,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of repositories to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListMyRepositoriesImpl implements the read-only MCP tool for listing the
// authenticated user's repositories. This is a safe, idempotent operation that
// uses the Forgejo SDK.
type ListMyRepositoriesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_my_repositories` tool. It supports optional
// parameters for filtering and sorting, and is marked as a safe, read-only operation.
func (ListMyRepositoriesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{},
		},
	}
}

// Handler implements the logic for listing the user's repositories. It calls the
// Forgejo SDK's `ListMyRepos` function and formats the results into a markdown list.
func (impl ListMyRepositoriesImpl) Handler() mcp.ToolHandlerFor[ListMyRepositoriesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListMyRepositoriesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.ListReposOptions{}
		// Note: ListReposOptions is quite limited in the SDK
		// Many filtering options are not available
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		// Call SDK
		repos, _, err := impl.Client.ListMyRepos(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list my repositories: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(repos) == 0 {
			content = "No repositories found for the authenticated user."
		} else {
			// Convert repos to our type
			repoList := make(types.RepositoryList, len(repos))
			for i, repo := range repos {
				repoList[i] = &types.Repository{Repository: repo}
			}

			content = fmt.Sprintf("Found %d repositories\n\n%s",
				len(repos), repoList.ToMarkdown())
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

// ListOrgRepositoriesParams defines the parameters for the list_org_repositories tool.
// It specifies the organization and allows for filtering and sorting.
type ListOrgRepositoriesParams struct {
	// Org is the name of the organization.
	Org string `json:"org"`
	// Type filters repositories by their type (e.g., forks, sources).
	Type string `json:"type,omitempty"`
	// Sort specifies the sort order for the results.
	Sort string `json:"sort,omitempty"`
	// Direction specifies the sort direction (asc or desc).
	Direction string `json:"direction,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of repositories to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListOrgRepositoriesImpl implements the read-only MCP tool for listing an
// organization's repositories. This is a safe, idempotent operation that uses
// the Forgejo SDK.
type ListOrgRepositoriesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_org_repositories` tool. It requires an `org` name
// and supports optional parameters for filtering and sorting. It is marked as a
// safe, read-only operation.
func (ListOrgRepositoriesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of repositories per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"org"},
		},
	}
}

// Handler implements the logic for listing organization repositories. It calls the
// Forgejo SDK's `ListOrgRepos` function and formats the results into a markdown list.
func (impl ListOrgRepositoriesImpl) Handler() mcp.ToolHandlerFor[ListOrgRepositoriesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListOrgRepositoriesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.ListOrgReposOptions{}
		// Note: ListOrgReposOptions is quite limited in the SDK
		// Type filtering is not available
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		// Call SDK
		repos, _, err := impl.Client.ListOrgRepos(p.Org, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list organization repositories: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(repos) == 0 {
			content = fmt.Sprintf("No repositories found for organization '%s'.", p.Org)
		} else {
			// Convert repos to our type
			repoList := make(types.RepositoryList, len(repos))
			for i, repo := range repos {
				repoList[i] = &types.Repository{Repository: repo}
			}

			content = fmt.Sprintf("Found %d repositories for organization '%s'\n\n%s",
				len(repos), p.Org, repoList.ToMarkdown())
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

// GetRepositoryParams defines the parameters for the get_repository tool.
// It specifies the owner and repository name to retrieve.
type GetRepositoryParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
}

// GetRepositoryImpl implements the read-only MCP tool for fetching detailed
// information about a single repository. This is a safe, idempotent operation
// that uses the Forgejo SDK.
type GetRepositoryImpl struct {
	Client *tools.Client
}

// Definition describes the `get_repository` tool. It requires `owner` and `repo`
// as parameters and is marked as a safe, read-only operation.
func (GetRepositoryImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_repository",
		Title:       "Get Repository Information",
		Description: "Get detailed information about a specific repository, including description, stats, permissions, and metadata.",
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
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for fetching repository details. It calls the
// Forgejo SDK's `GetRepo` function and formats the full repository object into
// a detailed markdown view.
func (impl GetRepositoryImpl) Handler() mcp.ToolHandlerFor[GetRepositoryParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetRepositoryParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call SDK
		repo, _, err := impl.Client.GetRepo(p.Owner, p.Repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}

		// Convert to our type and format
		repoWrapper := &types.Repository{Repository: repo}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: repoWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}
