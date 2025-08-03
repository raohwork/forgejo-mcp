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

// Release related tools

func ListReleases() mcp.Tool {
	return mcp.Tool{
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
					Minimum:     float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of releases per page (optional, defaults to 10, max 50)",
					Minimum:     float64Ptr(1),
					Maximum:     float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

func CreateRelease() mcp.Tool {
	return mcp.Tool{
		Name:        "create_release",
		Title:       "Create Release",
		Description: "Create a new release in a repository. Specify tag name, title, description, and other metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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

func EditRelease() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_release",
		Title:       "Edit Release",
		Description: "Edit an existing release's title, description, or other metadata.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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

func DeleteRelease() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_release",
		Title:       "Delete Release",
		Description: "Delete a release from a repository.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(true),
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

// Release attachment management tools

func ListReleaseAttachments() mcp.Tool {
	return mcp.Tool{
		Name:        "list_release_attachments",
		Title:       "List Release Attachments",
		Description: "List all attachments for a specific release.",
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
				"release_id": {
					Type:        "integer",
					Description: "Release ID",
				},
			},
			Required: []string{"owner", "repo", "release_id"},
		},
	}
}

func CreateReleaseAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "create_release_attachment",
		Title:       "Add Release Attachment",
		Description: "Add an attachment to a release.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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
				"release_id": {
					Type:        "integer",
					Description: "Release ID",
				},
				"file_path": {
					Type:        "string",
					Description: "Path to the file to attach",
				},
				"name": {
					Type:        "string",
					Description: "Optional display name for the attachment (defaults to filename)",
				},
			},
			Required: []string{"owner", "repo", "release_id", "file_path"},
		},
	}
}

func EditReleaseAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_release_attachment",
		Title:       "Edit Release Attachment",
		Description: "Edit a release attachment's metadata (like display name).",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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
				"release_id": {
					Type:        "integer",
					Description: "Release ID",
				},
				"attachment_id": {
					Type:        "integer",
					Description: "Attachment ID to edit",
				},
				"name": {
					Type:        "string",
					Description: "New display name for the attachment",
				},
			},
			Required: []string{"owner", "repo", "release_id", "attachment_id", "name"},
		},
	}
}

func DeleteReleaseAttachment() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_release_attachment",
		Title:       "Delete Release Attachment",
		Description: "Delete an attachment from a release.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(true),
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
				"release_id": {
					Type:        "integer",
					Description: "Release ID",
				},
				"attachment_id": {
					Type:        "integer",
					Description: "Attachment ID to delete",
				},
			},
			Required: []string{"owner", "repo", "release_id", "attachment_id"},
		},
	}
}
