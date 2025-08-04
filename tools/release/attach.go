// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package release

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// ListReleaseAttachmentsParams defines the parameters for the list_release_attachments tool.
// It specifies the release to list attachments from.
type ListReleaseAttachmentsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ReleaseID is the unique identifier of the release.
	ReleaseID int `json:"release_id"`
}

// ListReleaseAttachmentsImpl implements the read-only MCP tool for listing release attachments.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a list
// of all attachments for a specific release.
type ListReleaseAttachmentsImpl struct{}

// Definition describes the `list_release_attachments` tool. It requires `owner`,
// `repo`, and `release_id`. It is marked as a safe, read-only operation.
func (ListReleaseAttachmentsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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

// Handler implements the logic for listing release attachments. It calls the Forgejo
// SDK's `ListReleaseAttachments` function and formats the results into a markdown list.
func (ListReleaseAttachmentsImpl) Handler() mcp.ToolHandlerFor[ListReleaseAttachmentsParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListReleaseAttachmentsParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// CreateReleaseAttachmentParams defines the parameters for creating a release attachment.
// It specifies the release and the local file path of the asset to upload.
type CreateReleaseAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ReleaseID is the unique identifier of the release to attach the file to.
	ReleaseID int `json:"release_id"`
	// FilePath is the local path to the file to be uploaded as an attachment.
	FilePath string `json:"file_path"`
	// Name is an optional display name for the attachment.
	Name string `json:"name,omitempty"`
}

// CreateReleaseAttachmentImpl implements the MCP tool for adding an attachment to a release.
// This is a non-idempotent operation that uploads a file from the local filesystem
// to the specified release using the Forgejo SDK.
type CreateReleaseAttachmentImpl struct{}

// Definition describes the `create_release_attachment` tool. It requires `release_id`
// and a local `file_path`. It is not idempotent, as multiple calls will upload
// multiple files.
func (CreateReleaseAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_release_attachment",
		Title:       "Add Release Attachment",
		Description: "Add an attachment to a release.",
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

// Handler implements the logic for creating a release attachment. It reads the file
// from the specified `file_path`, then calls the Forgejo SDK's `CreateReleaseAttachment`
// function to upload it. It will return an error if the file cannot be read.
func (CreateReleaseAttachmentImpl) Handler() mcp.ToolHandlerFor[CreateReleaseAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateReleaseAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// EditReleaseAttachmentParams defines the parameters for editing a release attachment.
// It specifies the attachment to edit and its new name.
type EditReleaseAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ReleaseID is the unique identifier of the release containing the attachment.
	ReleaseID int `json:"release_id"`
	// AttachmentID is the unique identifier of the attachment to edit.
	AttachmentID int `json:"attachment_id"`
	// Name is the new display name for the attachment.
	Name string `json:"name"`
}

// EditReleaseAttachmentImpl implements the MCP tool for editing a release attachment.
// This is an idempotent operation that renames an existing attachment using the
// Forgejo SDK.
type EditReleaseAttachmentImpl struct{}

// Definition describes the `edit_release_attachment` tool. It requires `release_id`,
// `attachment_id`, and a new `name`. It is marked as idempotent.
func (EditReleaseAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_release_attachment",
		Title:       "Edit Release Attachment",
		Description: "Edit a release attachment's metadata (like display name).",
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

// Handler implements the logic for editing a release attachment. It calls the Forgejo
// SDK's `EditReleaseAttachment` function. It will return an error if the attachment
// ID is not found.
func (EditReleaseAttachmentImpl) Handler() mcp.ToolHandlerFor[EditReleaseAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditReleaseAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// DeleteReleaseAttachmentParams defines the parameters for deleting a release attachment.
// It specifies the attachment to be deleted by its ID.
type DeleteReleaseAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ReleaseID is the unique identifier of the release containing the attachment.
	ReleaseID int `json:"release_id"`
	// AttachmentID is the unique identifier of the attachment to delete.
	AttachmentID int `json:"attachment_id"`
}

// DeleteReleaseAttachmentImpl implements the destructive MCP tool for deleting a release attachment.
// This is an idempotent but irreversible operation that removes an attachment from a
// release using the Forgejo SDK.
type DeleteReleaseAttachmentImpl struct{}

// Definition describes the `delete_release_attachment` tool. It requires `release_id`
// and `attachment_id`. It is marked as a destructive operation to ensure clients
// can warn the user before execution.
func (DeleteReleaseAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_release_attachment",
		Title:       "Delete Release Attachment",
		Description: "Delete an attachment from a release.",
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

// Handler implements the logic for deleting a release attachment. It calls the Forgejo
// SDK's `DeleteReleaseAttachment` function. On success, it returns a simple text
// confirmation. It will return an error if the attachment does not exist.
func (DeleteReleaseAttachmentImpl) Handler() mcp.ToolHandlerFor[DeleteReleaseAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteReleaseAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}
