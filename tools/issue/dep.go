// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListIssueDependenciesParams defines the parameters for the list_issue_dependencies tool.
// It specifies the issue for which to list dependencies.
type ListIssueDependenciesParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
}

// ListIssueDependenciesImpl implements the read-only MCP tool for listing issue dependencies.
// This is a safe, idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type ListIssueDependenciesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_issue_dependencies` tool. It requires `owner`, `repo`,
// and the issue `index`. It is marked as a safe, read-only operation.
func (ListIssueDependenciesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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

// Handler implements the logic for listing issue dependencies. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/issues/{index}/dependencies`
// endpoint and formats the results into a markdown list.
func (impl ListIssueDependenciesImpl) Handler() mcp.ToolHandlerFor[ListIssueDependenciesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListIssueDependenciesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		issues, err := impl.Client.MyListIssueDependencies(p.Owner, p.Repo, int64(p.Index))
		if err != nil {
			return nil, fmt.Errorf("failed to list dependencies: %w", err)
		}

		dependencies := types.IssueDependencyList(issues)
		content := fmt.Sprintf("## Issues that block #%d\n\n%s", p.Index, dependencies.ToMarkdown())

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}

// AddIssueDependencyParams defines the parameters for the add_issue_dependency tool.
// It specifies the two issues to link in a dependency relationship.
type AddIssueDependencyParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number of the dependent issue.
	Index int `json:"index"`
	// DependencyIndex is the issue number of the issue that `Index` will depend on.
	DependencyIndex int `json:"dependency_index"`
}

// AddIssueDependencyImpl implements the MCP tool for adding a dependency to an issue.
// This is an idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type AddIssueDependencyImpl struct {
	Client *tools.Client
}

// Definition describes the `add_issue_dependency` tool. It requires the `index` of
// the dependent issue and the `dependency_index` of the issue it depends on.
// It is marked as idempotent.
func (AddIssueDependencyImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "add_issue_dependency",
		Title:       "Add Issue Dependency",
		Description: "Add a dependency relationship between two issues, where one issue depends on another.",
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

// Handler implements the logic for adding an issue dependency. It performs a custom
// HTTP POST request to the `/repos/{owner}/{repo}/issues/{index}/dependencies`
// endpoint. It will return an error if either issue cannot be found.
func (impl AddIssueDependencyImpl) Handler() mcp.ToolHandlerFor[AddIssueDependencyParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[AddIssueDependencyParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		dependency := types.MyIssueMeta{
			Owner: p.Owner,
			Name:  p.Repo,
			Index: int64(p.DependencyIndex),
		}

		issue, err := impl.Client.MyAddIssueDependency(p.Owner, p.Repo, int64(p.Index), dependency)
		if err != nil {
			return nil, fmt.Errorf("failed to add dependency: %w", err)
		}

		issueWrapper := &types.Issue{Issue: issue}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Dependency added to issue #%d\n\n%s", p.Index, issueWrapper.ToMarkdown()),
				},
			},
		}, nil
	}
}

// RemoveIssueDependencyParams defines the parameters for the remove_issue_dependency tool.
// It specifies the two issues for which to remove the dependency relationship.
type RemoveIssueDependencyParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number of the dependent issue.
	Index int `json:"index"`
	// DependencyIndex is the issue number of the dependency to be removed.
	DependencyIndex int `json:"dependency_index"`
}

// RemoveIssueDependencyImpl implements the destructive MCP tool for removing an issue dependency.
// This is an idempotent but destructive operation. Note: This feature is not supported
// by the official Forgejo SDK and requires a custom HTTP implementation.
type RemoveIssueDependencyImpl struct {
	Client *tools.Client
}

// Definition describes the `remove_issue_dependency` tool. It requires the `index` of
// the dependent issue and the `dependency_index` to remove. It is marked as a
// destructive operation.
func (RemoveIssueDependencyImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "remove_issue_dependency",
		Title:       "Remove Issue Dependency",
		Description: "Remove a specific dependency relationship between two issues.",
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

// Handler implements the logic for removing an issue dependency. It performs a custom
// HTTP DELETE request to the `/repos/{owner}/{repo}/issues/{index}/dependencies/{dependency_index}`
// endpoint. On success, it returns a simple text confirmation.
func (impl RemoveIssueDependencyImpl) Handler() mcp.ToolHandlerFor[RemoveIssueDependencyParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[RemoveIssueDependencyParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		dependency := types.MyIssueMeta{
			Owner: p.Owner,
			Name:  p.Repo,
			Index: int64(p.DependencyIndex),
		}

		issue, err := impl.Client.MyRemoveIssueDependency(p.Owner, p.Repo, int64(p.Index), dependency)
		if err != nil {
			return nil, fmt.Errorf("failed to remove dependency: %w", err)
		}

		issueWrapper := &types.Issue{Issue: issue}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Dependency removed from issue #%d\n\n%s", p.Index, issueWrapper.ToMarkdown()),
				},
			},
		}, nil
	}
}
