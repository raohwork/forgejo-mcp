// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package release

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListReleasesParams defines the parameters for the list_releases tool.
// It specifies the owner and repository, with optional pagination.
type ListReleasesParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of releases to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListReleasesImpl implements the read-only MCP tool for listing repository releases.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a
// paginated list of releases.
type ListReleasesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_releases` tool. It requires `owner` and `repo`
// and supports pagination. It is marked as a safe, read-only operation.
func (ListReleasesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_releases",
		Title:       "List Releases",
		Description: "List all releases in a repository. Returns release information including tag name, title, description, and assets.",
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
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of releases per page (optional, defaults to 10, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing releases. It calls the Forgejo SDK's
// `ListReleases` function and formats the results into a markdown list.
func (impl ListReleasesImpl) Handler() mcp.ToolHandlerFor[ListReleasesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListReleasesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.ListReleasesOptions{}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		// Call SDK
		releases, _, err := impl.Client.ListReleases(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(releases) == 0 {
			content = "No releases found in this repository."
		} else {
			// Convert releases to our type
			releaseList := make(types.ReleaseList, len(releases))
			for i, release := range releases {
				releaseList[i] = &types.Release{Release: release}
			}

			content = fmt.Sprintf("Found %d releases\n\n%s",
				len(releases), releaseList.ToMarkdown())
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

// CreateReleaseParams defines the parameters for the create_release tool.
// It includes all necessary details for creating a new release.
type CreateReleaseParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// TagName is the name of the git tag for the release.
	TagName string `json:"tag_name"`
	// TargetCommitish is the branch or commit SHA to create the release from.
	TargetCommitish string `json:"target_commitish,omitempty"`
	// Name is the title of the release.
	Name string `json:"name"`
	// Body is the markdown description of the release.
	Body string `json:"body,omitempty"`
	// Draft indicates whether the release is a draft.
	Draft bool `json:"draft,omitempty"`
	// Prerelease indicates whether the release is a pre-release.
	Prerelease bool `json:"prerelease,omitempty"`
}

// CreateReleaseImpl implements the MCP tool for creating a new release.
// This is a non-idempotent operation that creates a new git tag and a
// corresponding release object via the Forgejo SDK.
type CreateReleaseImpl struct {
	Client *tools.Client
}

// Definition describes the `create_release` tool. It requires `owner`, `repo`,
// `tag_name`, and a `name`. It is not idempotent, as multiple calls will fail
// once the tag is created.
func (CreateReleaseImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_release",
		Title:       "Create Release",
		Description: "Create a new release in a repository. Specify tag name, title, description, and other metadata.",
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
				"tag_name": {
					Type:        "string",
					Description: "Git tag name for this release",
				},
				"target_commitish": {
					Type:        "string",
					Description: "Target branch or commit SHA (optional, defaults to default branch)",
				},
				"name": {
					Type:        "string",
					Description: "Release title",
				},
				"body": {
					Type:        "string",
					Description: "Release description (markdown supported) (optional)",
				},
				"draft": {
					Type:        "boolean",
					Description: "Whether this is a draft release (optional, defaults to false)",
				},
				"prerelease": {
					Type:        "boolean",
					Description: "Whether this is a prerelease (optional, defaults to false)",
				},
			},
			Required: []string{"owner", "repo", "tag_name", "name"},
		},
	}
}

// Handler implements the logic for creating a release. It calls the Forgejo SDK's
// `CreateRelease` function. On success, it returns details of the new release.
func (impl CreateReleaseImpl) Handler() mcp.ToolHandlerFor[CreateReleaseParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateReleaseParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.CreateReleaseOption{
			TagName:      p.TagName,
			Target:       p.TargetCommitish,
			Title:        p.Name,
			Note:         p.Body,
			IsDraft:      p.Draft,
			IsPrerelease: p.Prerelease,
		}

		// Call SDK
		release, _, err := impl.Client.CreateRelease(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to create release: %w", err)
		}

		// Convert to our type and format
		releaseWrapper := &types.Release{Release: release}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: releaseWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// EditReleaseParams defines the parameters for the edit_release tool.
// It specifies the release to edit by ID and the fields to update.
type EditReleaseParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the release to edit.
	ID int `json:"id"`
	// TagName is the new git tag name for the release.
	TagName string `json:"tag_name,omitempty"`
	// TargetCommitish is the new target branch or commit SHA.
	TargetCommitish string `json:"target_commitish,omitempty"`
	// Name is the new title for the release.
	Name string `json:"name,omitempty"`
	// Body is the new markdown description for the release.
	Body string `json:"body,omitempty"`
	// Draft indicates whether the release should be marked as a draft.
	Draft bool `json:"draft,omitempty"`
	// Prerelease indicates whether the release should be marked as a pre-release.
	Prerelease bool `json:"prerelease,omitempty"`
}

// EditReleaseImpl implements the MCP tool for editing an existing release.
// This is an idempotent operation that modifies the metadata of a release
// identified by its ID, using the Forgejo SDK.
type EditReleaseImpl struct {
	Client *tools.Client
}

// Definition describes the `edit_release` tool. It requires `owner`, `repo`, and
// the release `id`. It is marked as idempotent, as multiple identical calls
// will result in the same state.
func (EditReleaseImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_release",
		Title:       "Edit Release",
		Description: "Edit an existing release's title, description, or other metadata.",
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
				"id": {
					Type:        "integer",
					Description: "Release ID",
				},
				"tag_name": {
					Type:        "string",
					Description: "New git tag name for this release (optional)",
				},
				"target_commitish": {
					Type:        "string",
					Description: "New target branch or commit SHA (optional)",
				},
				"name": {
					Type:        "string",
					Description: "New release title (optional)",
				},
				"body": {
					Type:        "string",
					Description: "New release description (markdown supported) (optional)",
				},
				"draft": {
					Type:        "boolean",
					Description: "Whether this is a draft release (optional)",
				},
				"prerelease": {
					Type:        "boolean",
					Description: "Whether this is a prerelease (optional)",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for editing a release. It calls the Forgejo SDK's
// `EditRelease` function. It will return an error if the release ID is not found.
func (impl EditReleaseImpl) Handler() mcp.ToolHandlerFor[EditReleaseParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditReleaseParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.EditReleaseOption{}
		if p.TagName != "" {
			opt.TagName = p.TagName
		}
		if p.TargetCommitish != "" {
			opt.Target = p.TargetCommitish
		}
		if p.Name != "" {
			opt.Title = p.Name
		}
		if p.Body != "" {
			opt.Note = p.Body
		}
		// Note: For boolean fields, we need to check if they were explicitly set
		// For now, we'll pass them directly assuming they're properly handled
		opt.IsDraft = &p.Draft
		opt.IsPrerelease = &p.Prerelease

		// Call SDK
		release, _, err := impl.Client.EditRelease(p.Owner, p.Repo, int64(p.ID), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to edit release: %w", err)
		}

		// Convert to our type and format
		releaseWrapper := &types.Release{Release: release}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: releaseWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// DeleteReleaseParams defines the parameters for the delete_release tool.
// It specifies the release to be deleted by its ID.
type DeleteReleaseParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the release to delete.
	ID int `json:"id"`
}

// DeleteReleaseImpl implements the destructive MCP tool for deleting a release.
// This is an idempotent but irreversible operation that removes both the release
// object and its associated git tag. It uses the Forgejo SDK.
type DeleteReleaseImpl struct {
	Client *tools.Client
}

// Definition describes the `delete_release` tool. It requires `owner`, `repo`,
// and the release `id`. It is marked as a destructive operation to ensure
// clients can warn the user before execution.
func (DeleteReleaseImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_release",
		Title:       "Delete Release",
		Description: "Delete a release from a repository.",
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
				"id": {
					Type:        "integer",
					Description: "Release ID to delete",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for deleting a release. It calls the Forgejo SDK's
// `DeleteRelease` function. On success, it returns a simple text confirmation.
// It will return an error if the release does not exist.
func (impl DeleteReleaseImpl) Handler() mcp.ToolHandlerFor[DeleteReleaseParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteReleaseParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call SDK
		_, err := impl.Client.DeleteRelease(p.Owner, p.Repo, int64(p.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to delete release: %w", err)
		}

		// Return success message
		emptyResponse := types.EmptyResponse{}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: emptyResponse.ToMarkdown(),
				},
			},
		}, nil
	}
}
