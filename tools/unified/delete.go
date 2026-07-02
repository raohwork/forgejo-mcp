// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// DeleteImpl implements the delete_gitea tool.
type DeleteImpl struct {
	Client *tools.Client
}

// Definition describes the delete_gitea tool with minimal schema.
func (DeleteImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "delete_gitea",
		Title: "Delete Gitea Resource",
		Description: `Delete a resource from Forgejo/Gitea. This action cannot be undone.
Resources: issue_comment, issue_attachment, label, milestone, release, release_attachment, wiki_page.
Use gitea_manual(action="delete") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(true),
			IdempotentHint:  true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resource": {
					Type:        "string",
					Description: "Resource type to delete",
					Enum: []any{
						"issue_comment", "issue_attachment", "label",
						"milestone", "release", "release_attachment", "wiki_page",
					},
				},
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"id": {
					Type:        "integer",
					Description: "Resource ID (for issue_comment, label, milestone, release, release_attachment)",
				},
				"index": {
					Type:        "integer",
					Description: "Issue number (for issue_attachment)",
				},
				"attachment_id": {
					Type:        "integer",
					Description: "Attachment ID (for issue_attachment, release_attachment)",
				},
			},
			Required:             []string{"resource", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate delete logic based on resource type.
func (impl DeleteImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		resource, _ := args["resource"].(string)
		if resource == "" {
			resources := ListResourcesForAction(ActionDelete)
			return nil, nil, fmt.Errorf("resource is required. Valid resources: %v", resources)
		}

		if _, ok := LookupManual(ActionDelete, resource); !ok {
			resources := ListResourcesForAction(ActionDelete)
			return nil, nil, fmt.Errorf("unknown resource '%s'. Valid resources: %v", resource, resources)
		}

		switch resource {
		case "issue_comment":
			return impl.deleteIssueComment(args)
		case "issue_attachment":
			return impl.deleteIssueAttachment(args)
		case "label":
			return impl.deleteLabel(args)
		case "milestone":
			return impl.deleteMilestone(args)
		case "release":
			return impl.deleteRelease(args)
		case "release_attachment":
			return impl.deleteReleaseAttachment(args)
		case "wiki_page":
			return impl.deleteWikiPage(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionDelete, resource, "not implemented"))
		}
	}
}

func (impl DeleteImpl) deleteIssueComment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "issue_comment", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "issue_comment", err.Error()))
	}

	_, err = impl.Client.DeleteIssueComment(owner, repo, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete comment: %w", err)
	}

	return textResult(fmt.Sprintf("Comment %d successfully deleted.", id)), nil, nil
}

func (impl DeleteImpl) deleteIssueAttachment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "issue_attachment", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "issue_attachment", err.Error()))
	}

	attachmentID, err := requiredPositiveInt(args, "attachment_id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "issue_attachment", err.Error()))
	}

	err = impl.Client.MyDeleteIssueAttachment(owner, repo, index, attachmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete attachment: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}

func (impl DeleteImpl) deleteLabel(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "label", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "label", err.Error()))
	}

	_, err = impl.Client.DeleteLabel(owner, repo, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete label: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}

func (impl DeleteImpl) deleteMilestone(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "milestone", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "milestone", err.Error()))
	}

	_, err = impl.Client.DeleteMilestone(owner, repo, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete milestone: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}

func (impl DeleteImpl) deleteRelease(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "release", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "release", err.Error()))
	}

	_, err = impl.Client.DeleteRelease(owner, repo, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete release: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}

func (impl DeleteImpl) deleteReleaseAttachment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "release_attachment", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "release_attachment", err.Error()))
	}

	attachmentID, err := requiredPositiveInt(args, "attachment_id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "release_attachment", err.Error()))
	}

	_, err = impl.Client.DeleteReleaseAttachment(owner, repo, id, attachmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete release attachment: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}

func (impl DeleteImpl) deleteWikiPage(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "wiki_page", err.Error()))
	}

	pageName, _ := args["page_name"].(string)
	if pageName == "" {
		return nil, nil, errors.New(FormatValidationError(ActionDelete, "wiki_page", "page_name is required"))
	}

	err = impl.Client.MyDeleteWikiPage(owner, repo, pageName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete wiki page: %w", err)
	}

	return textResult(types.EmptyResponse{}.ToMarkdown()), nil, nil
}
