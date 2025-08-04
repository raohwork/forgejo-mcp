// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// ListIssueAttachmentsParams defines the parameters for the list_issue_attachments tool.
// It specifies the issue from which to list attachments.
type ListIssueAttachmentsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
}

// ListIssueAttachmentsImpl implements the read-only MCP tool for listing issue attachments.
// This is a safe, idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type ListIssueAttachmentsImpl struct{}

// Definition describes the `list_issue_attachments` tool. It requires `owner`, `repo`,
// and the issue `index`. It is marked as a safe, read-only operation.
func (ListIssueAttachmentsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_issue_attachments",
		Title:       "List Issue Attachments",
		Description: "List all attachments on an issue. Returns attachment information including names, sizes, and download URLs.",
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

// Handler implements the logic for listing issue attachments. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/issues/{index}/attachments`
// endpoint and formats the results into a markdown list.
func (ListIssueAttachmentsImpl) Handler() mcp.ToolHandlerFor[ListIssueAttachmentsParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListIssueAttachmentsParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// CreateIssueAttachmentParams defines the parameters for creating an issue attachment.
// It specifies the issue and the local file path of the asset to upload.
type CreateIssueAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number to attach the file to.
	Index int `json:"index"`
	// FilePath is the local path to the file to be uploaded as an attachment.
	FilePath string `json:"file_path"`
	// Name is an optional display name for the attachment.
	Name string `json:"name,omitempty"`
}

// CreateIssueAttachmentImpl implements the MCP tool for adding an attachment to an issue.
// This is a non-idempotent operation that uploads a file from the local filesystem.
// Note: This feature is not supported by the official Forgejo SDK and requires a
// custom HTTP implementation.
type CreateIssueAttachmentImpl struct{}

// Definition describes the `create_issue_attachment` tool. It requires the issue `index`
// and a local `file_path`. It is not idempotent, as multiple calls will upload
// multiple files.
func (CreateIssueAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_issue_attachment",
		Title:       "Add Issue Attachment",
		Description: "Add a file attachment to an issue.",
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
				"index": {
					Type:        "integer",
					Description: "Issue index number",
				},
				"file_path": {
					Type:        "string",
					Description: "Absolute path to the file to attach",
				},
				"name": {
					Type:        "string",
					Description: "Optional display name for the attachment (defaults to filename)",
				},
			},
			Required: []string{"owner", "repo", "index", "file_path"},
		},
	}
}

// Handler implements the logic for creating an issue attachment. It reads the file
// from the specified `file_path` and performs a custom HTTP POST request to upload
// it. It will return an error if the file cannot be read.
func (CreateIssueAttachmentImpl) Handler() mcp.ToolHandlerFor[CreateIssueAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateIssueAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// DeleteIssueAttachmentParams defines the parameters for deleting an issue attachment.
// It specifies the attachment to be deleted by its ID.
type DeleteIssueAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number containing the attachment.
	Index int `json:"index"`
	// AttachmentID is the unique identifier of the attachment to delete.
	AttachmentID string `json:"attachment_id"`
}

// DeleteIssueAttachmentImpl implements the destructive MCP tool for deleting an issue attachment.
// This is an idempotent but irreversible operation. Note: This feature is not supported
// by the official Forgejo SDK and requires a custom HTTP implementation.
type DeleteIssueAttachmentImpl struct{}

// Definition describes the `delete_issue_attachment` tool. It requires the issue `index`
// and `attachment_id`. It is marked as a destructive operation.
func (DeleteIssueAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_issue_attachment",
		Title:       "Delete Issue Attachment",
		Description: "Delete a specific attachment from an issue.",
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
				"attachment_id": {
					Type:        "string",
					Description: "Attachment ID to delete",
				},
			},
			Required: []string{"owner", "repo", "index", "attachment_id"},
		},
	}
}

// Handler implements the logic for deleting an issue attachment. It performs a custom
// HTTP DELETE request to the `/repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}`
// endpoint. On success, it returns a simple text confirmation.
func (DeleteIssueAttachmentImpl) Handler() mcp.ToolHandlerFor[DeleteIssueAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteIssueAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// EditIssueAttachmentParams defines the parameters for editing an issue attachment.
// It specifies the attachment to edit and its new name.
type EditIssueAttachmentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number containing the attachment.
	Index int `json:"index"`
	// AttachmentID is the unique identifier of the attachment to edit.
	AttachmentID string `json:"attachment_id"`
	// Name is the new display name for the attachment.
	Name string `json:"name"`
}

// EditIssueAttachmentImpl implements the MCP tool for editing an issue attachment.
// This is an idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type EditIssueAttachmentImpl struct{}

// Definition describes the `edit_issue_attachment` tool. It requires the issue `index`,
// `attachment_id`, and a new `name`. It is marked as idempotent.
func (EditIssueAttachmentImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_issue_attachment",
		Title:       "Edit Issue Attachment",
		Description: "Edit an attachment's metadata such as display name.",
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
				"attachment_id": {
					Type:        "string",
					Description: "Attachment ID to edit",
				},
				"name": {
					Type:        "string",
					Description: "New display name for the attachment",
				},
			},
			Required: []string{"owner", "repo", "index", "attachment_id", "name"},
		},
	}
}

// Handler implements the logic for editing an issue attachment. It performs a custom
// HTTP PATCH request to the `/repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}`
// endpoint. It will return an error if the attachment is not found.
func (EditIssueAttachmentImpl) Handler() mcp.ToolHandlerFor[EditIssueAttachmentParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditIssueAttachmentParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}
