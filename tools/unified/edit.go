// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// EditImpl implements the edit_gitea tool.
type EditImpl struct {
	Client *tools.Client
}

// Definition describes the edit_gitea tool with minimal schema.
func (EditImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "edit_gitea",
		Title: "Edit Gitea Resource",
		Description: `Edit an existing resource in Forgejo/Gitea.
Resources: issue, issue_comment, issue_attachment, label, milestone, release, release_attachment, wiki_page.
Use gitea_manual(action="edit") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resource": {
					Type:        "string",
					Description: "Resource type to edit",
					Enum: []any{
						"issue", "issue_comment", "issue_attachment", "label",
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
				"index": {
					Type:        "integer",
					Description: "Issue number (for issue, issue_attachment)",
				},
				"id": {
					Type:        "integer",
					Description: "Resource ID (for issue_comment, label, milestone, release, release_attachment)",
				},
				"attachment_id": {
					Type:        "integer",
					Description: "Attachment ID (for issue_attachment, release_attachment)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID (for issue)",
				},
				"assignees": {
					Type:        "array",
					Items:       &jsonschema.Schema{Type: "string"},
					Description: "Usernames to assign (for issue)",
				},
			},
			Required:             []string{"resource", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate edit logic based on resource type.
func (impl EditImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		resource, _ := args["resource"].(string)
		if resource == "" {
			resources := ListResourcesForAction(ActionEdit)
			return nil, nil, fmt.Errorf("resource is required. Valid resources: %v", resources)
		}

		if _, ok := LookupManual(ActionEdit, resource); !ok {
			resources := ListResourcesForAction(ActionEdit)
			return nil, nil, fmt.Errorf("unknown resource '%s'. Valid resources: %v", resource, resources)
		}

		switch resource {
		case "issue":
			return impl.editIssue(args)
		case "issue_comment":
			return impl.editIssueComment(args)
		case "issue_attachment":
			return impl.editIssueAttachment(args)
		case "label":
			return impl.editLabel(args)
		case "milestone":
			return impl.editMilestone(args)
		case "release":
			return impl.editRelease(args)
		case "release_attachment":
			return impl.editReleaseAttachment(args)
		case "wiki_page":
			return impl.editWikiPage(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionEdit, resource, "not implemented"))
		}
	}
}

func (impl EditImpl) editIssue(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue", err.Error()))
	}

	opt := forgejo.EditIssueOption{}

	if title, ok := args["title"].(string); ok && title != "" {
		opt.Title = title
	}
	if body, ok := args["body"].(string); ok && body != "" {
		opt.Body = &body
	}
	if state, ok := args["state"].(string); ok && state != "" {
		s := forgejo.StateType(state)
		opt.State = &s
	}
	assignees, err := stringList(args, "assignees")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue", err.Error()))
	}
	if assignees != nil {
		opt.Assignees = assignees
	}
	milestone, ok, err := optionalPositiveInt(args, "milestone")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue", err.Error()))
	}
	if ok {
		m := milestone
		opt.Milestone = &m
	}
	if dueDateStr, ok := args["due_date"].(string); ok && dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid due_date format (expected RFC3339): %w", err)
		}
		opt.Deadline = &dueDate
	}

	issue, _, err := impl.Client.EditIssue(owner, repo, index, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit issue: %w", err)
	}

	return textResult((&types.Issue{Issue: issue}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editIssueComment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_comment", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_comment", err.Error()))
	}

	body, _ := args["body"].(string)
	if body == "" {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_comment", "body is required"))
	}

	opt := forgejo.EditIssueCommentOption{Body: body}
	comment, _, err := impl.Client.EditIssueComment(owner, repo, id, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit comment: %w", err)
	}

	return textResult((&types.Comment{Comment: comment}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editIssueAttachment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_attachment", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_attachment", err.Error()))
	}

	attachmentID, err := requiredPositiveInt(args, "attachment_id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_attachment", err.Error()))
	}

	name, _ := args["name"].(string)
	if name == "" {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "issue_attachment", "name is required"))
	}

	options := tools.MyEditAttachmentOptions{Name: name}
	attachment, err := impl.Client.MyEditIssueAttachment(owner, repo, index, attachmentID, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit attachment: %w", err)
	}

	return textResult((&types.Attachment{Attachment: attachment}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editLabel(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "label", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "label", err.Error()))
	}

	opt := forgejo.EditLabelOption{}
	if name, ok := args["name"].(string); ok && name != "" {
		opt.Name = &name
	}
	if color, ok := args["color"].(string); ok && color != "" {
		opt.Color = &color
	}
	if description, ok := args["description"].(string); ok && description != "" {
		opt.Description = &description
	}

	label, _, err := impl.Client.EditLabel(owner, repo, id, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit label: %w", err)
	}

	return textResult((&types.Label{Label: label}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editMilestone(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "milestone", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "milestone", err.Error()))
	}

	opt := forgejo.EditMilestoneOption{}
	if title, ok := args["title"].(string); ok && title != "" {
		opt.Title = title
	}
	if description, ok := args["description"].(string); ok && description != "" {
		opt.Description = &description
	}
	if state, ok := args["state"].(string); ok && state != "" {
		s := forgejo.StateType(state)
		opt.State = &s
	}
	if dueDateStr, ok := args["due_date"].(string); ok && dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid due_date format (expected RFC3339): %w", err)
		}
		opt.Deadline = &dueDate
	}

	milestone, _, err := impl.Client.EditMilestone(owner, repo, id, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit milestone: %w", err)
	}

	return textResult((&types.Milestone{Milestone: milestone}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editRelease(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release", err.Error()))
	}

	opt := forgejo.EditReleaseOption{}
	if tagName, ok := args["tag_name"].(string); ok && tagName != "" {
		opt.TagName = tagName
	}
	if target, ok := args["target_commitish"].(string); ok && target != "" {
		opt.Target = target
	}
	if name, ok := args["name"].(string); ok && name != "" {
		opt.Title = name
	}
	if body, ok := args["body"].(string); ok && body != "" {
		opt.Note = body
	}
	if draft, ok := args["draft"].(bool); ok {
		opt.IsDraft = &draft
	}
	if prerelease, ok := args["prerelease"].(bool); ok {
		opt.IsPrerelease = &prerelease
	}

	release, _, err := impl.Client.EditRelease(owner, repo, id, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit release: %w", err)
	}

	return textResult((&types.Release{Release: release}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editReleaseAttachment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release_attachment", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release_attachment", err.Error()))
	}

	attachmentID, err := requiredPositiveInt(args, "attachment_id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release_attachment", err.Error()))
	}

	name, _ := args["name"].(string)
	if name == "" {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "release_attachment", "name is required"))
	}

	opt := forgejo.EditAttachmentOptions{Name: name}
	attachment, _, err := impl.Client.EditReleaseAttachment(owner, repo, id, attachmentID, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit release attachment: %w", err)
	}

	return textResult((&types.Attachment{Attachment: attachment}).ToMarkdown()), nil, nil
}

func (impl EditImpl) editWikiPage(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "wiki_page", err.Error()))
	}

	pageName, _ := args["page_name"].(string)
	if pageName == "" {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "wiki_page", "page_name is required"))
	}

	content, _ := args["content"].(string)
	if content == "" {
		return nil, nil, errors.New(FormatValidationError(ActionEdit, "wiki_page", "content is required"))
	}

	title, _ := args["title"].(string)
	if title == "" {
		title = pageName
	}
	message, _ := args["message"].(string)

	options := types.MyCreateWikiPageOptions{
		Title:         title,
		ContentBase64: base64.StdEncoding.EncodeToString([]byte(content)),
		Message:       message,
	}

	page, err := impl.Client.MyEditWikiPage(owner, repo, pageName, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to edit wiki page: %w", err)
	}

	return textResult((&types.WikiPage{MyWikiPage: page}).ToMarkdown()), nil, nil
}
