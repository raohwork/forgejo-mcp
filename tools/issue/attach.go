// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
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
type ListIssueAttachmentsImpl struct {
	Client *tools.Client
}

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
// HTTP GET request to the `/repos/{owner}/{repo}/issues/{index}/assets`
// endpoint and formats the results into a markdown list.
func (impl ListIssueAttachmentsImpl) Handler() mcp.ToolHandlerFor[ListIssueAttachmentsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListIssueAttachmentsParams) (*mcp.CallToolResult, any, error) {
		p := args

		// List issue attachments using the custom client method
		attachments, err := impl.Client.MyListIssueAttachments(p.Owner, p.Repo, int64(p.Index))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list issue attachments: %w", err)
		}

		// Convert to types.AttachmentList for consistent formatting
		attachmentList := make(types.AttachmentList, len(attachments))
		for i, attachment := range attachments {
			attachmentList[i] = &types.Attachment{Attachment: attachment}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("# Issue #%d Attachments\n\n%s", p.Index, attachmentList.ToMarkdown()),
				},
			},
		}, nil, nil
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
type DeleteIssueAttachmentImpl struct {
	Client *tools.Client
}

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
// HTTP DELETE request to the `/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}`
// endpoint. On success, it returns a simple text confirmation.
func (impl DeleteIssueAttachmentImpl) Handler() mcp.ToolHandlerFor[DeleteIssueAttachmentParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args DeleteIssueAttachmentParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Convert attachment ID from string to int64
		attachmentID, err := strconv.ParseInt(p.AttachmentID, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid attachment ID: %w", err)
		}

		// Delete the attachment using the custom client method
		err = impl.Client.MyDeleteIssueAttachment(p.Owner, p.Repo, int64(p.Index), attachmentID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete issue attachment: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Issue attachment %s deleted successfully from issue #%d", p.AttachmentID, p.Index),
				},
			},
		}, nil, nil
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
type EditIssueAttachmentImpl struct {
	Client *tools.Client
}

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
// HTTP PATCH request to the `/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}`
// endpoint. It will return an error if the attachment is not found.
func (impl EditIssueAttachmentImpl) Handler() mcp.ToolHandlerFor[EditIssueAttachmentParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args EditIssueAttachmentParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Convert attachment ID from string to int64
		attachmentID, err := strconv.ParseInt(p.AttachmentID, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid attachment ID: %w", err)
		}

		// Create options struct from parameters
		options := tools.MyEditAttachmentOptions{
			Name: p.Name,
		}

		// Edit the attachment using the custom client method
		attachment, err := impl.Client.MyEditIssueAttachment(p.Owner, p.Repo, int64(p.Index), attachmentID, options)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to edit issue attachment: %w", err)
		}

		// Convert to types.Attachment for consistent formatting
		result := &types.Attachment{Attachment: attachment}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("# Issue Attachment Updated\n\n%s", result.ToMarkdown()),
				},
			},
		}, nil, nil
	}
}
