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
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// CreateRepositoryParams defines the parameters for the create_repository tool.
type CreateRepositoryParams struct {
	// Name is the name of the repository to create.
	Name string `json:"name"`
	// Org is the organization to create the repository in; empty means the
	// authenticated user.
	Org string `json:"org,omitempty"`
	// Description is the repository description.
	Description string `json:"description,omitempty"`
	// Private indicates whether the repository is private.
	Private bool `json:"private,omitempty"`
	// AutoInit indicates whether to initialize the repository with a README.
	AutoInit bool `json:"auto_init,omitempty"`
	// DefaultBranch is the default branch name (used when auto-initializing).
	DefaultBranch string `json:"default_branch,omitempty"`
}

// CreateRepositoryImpl implements the MCP tool for creating a new repository.
// This operation creates data and is not idempotent: calling it twice with the
// same name fails on the second call.
type CreateRepositoryImpl struct {
	Client *tools.Client
}

// Definition describes the `create_repository` tool. It requires `name` and
// optionally accepts `org`, `description`, `private`, `auto_init` and
// `default_branch`.
func (CreateRepositoryImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_repository",
		Title:       "Create Repository",
		Description: "Create a new repository for the authenticated user, or in an organization when `org` is given. Note: the access token needs the write:user scope (write:organization for org repositories).",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  false,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Type:        "string",
					Description: "Name of the repository to create",
				},
				"org": {
					Type:        "string",
					Description: "Organization to create the repository in (optional, defaults to the authenticated user)",
				},
				"description": {
					Type:        "string",
					Description: "Repository description (optional)",
				},
				"private": {
					Type:        "boolean",
					Description: "Whether the repository is private (optional, defaults to false)",
				},
				"auto_init": {
					Type:        "boolean",
					Description: "Initialize the repository with a README so it can be cloned immediately (optional, defaults to false; leave false when you plan to push an existing repository)",
				},
				"default_branch": {
					Type:        "string",
					Description: "Default branch name, used when auto-initializing (optional, server default when omitted)",
				},
			},
			Required: []string{"name"},
		},
	}
}

// Handler implements the logic for creating a repository. It calls the Forgejo
// SDK's `CreateRepo` (or `CreateOrgRepo` when `org` is given) and formats the
// created repository as markdown.
func (impl CreateRepositoryImpl) Handler() mcp.ToolHandlerFor[CreateRepositoryParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args CreateRepositoryParams) (*mcp.CallToolResult, any, error) {
		p := args

		opt := forgejo.CreateRepoOption{
			Name:          p.Name,
			Description:   p.Description,
			Private:       p.Private,
			AutoInit:      p.AutoInit,
			DefaultBranch: p.DefaultBranch,
		}

		// Call SDK
		var (
			repo *forgejo.Repository
			err  error
		)
		if p.Org != "" {
			repo, _, err = impl.Client.CreateOrgRepo(p.Org, opt)
		} else {
			repo, _, err = impl.Client.CreateRepo(opt)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create repository (the token may lack the write:user or write:organization scope, or the name may already be taken): %w", err)
		}

		// Convert to our type and format
		repoWrapper := &types.Repository{Repository: repo}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Repository created successfully!\n\n" + repoWrapper.ToMarkdown(),
				},
			},
		}, nil, nil
	}
}
