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

// CreateParams defines the parameters for the create_gitea tool.
// Uses a generic map to handle varying parameters per resource type.
type CreateParams struct {
	Resource string         `json:"resource"`
	Params   map[string]any `json:"params"`
}

// CreateImpl implements the create_gitea tool.
type CreateImpl struct {
	Client *tools.Client
}

// Definition describes the create_gitea tool with minimal schema.
func (CreateImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "create_gitea",
		Title: "Create Gitea Resource",
		Description: `Create a resource in Forgejo/Gitea.
Resources: issue, issue_comment, label, milestone, release, wiki_page, pull_request.
Use gitea_manual(action="create") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  false,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resource": {
					Type:        "string",
					Description: "Resource type to create",
					Enum:        []any{"issue", "issue_comment", "label", "milestone", "release", "wiki_page", "pull_request"},
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
					Description: "Issue number (required for issue_comment)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID",
				},
				"labels": {
					Type:        "array",
					Items:       &jsonschema.Schema{Type: "integer"},
					Description: "Label IDs (for issue, pull_request)",
				},
				"assignees": {
					Type:        "array",
					Items:       &jsonschema.Schema{Type: "string"},
					Description: "Usernames to assign (for issue, pull_request)",
				},
			},
			Required:             []string{"resource", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate creation logic based on resource type.
func (impl CreateImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		resource, _ := args["resource"].(string)
		if resource == "" {
			resources := ListResourcesForAction(ActionCreate)
			return nil, nil, fmt.Errorf("resource is required. Valid resources: %v", resources)
		}

		// Validate resource type
		if _, ok := LookupManual(ActionCreate, resource); !ok {
			resources := ListResourcesForAction(ActionCreate)
			return nil, nil, fmt.Errorf("unknown resource '%s'. Valid resources: %v", resource, resources)
		}

		switch resource {
		case "issue":
			return impl.createIssue(args)
		case "issue_comment":
			return impl.createIssueComment(args)
		case "label":
			return impl.createLabel(args)
		case "milestone":
			return impl.createMilestone(args)
		case "release":
			return impl.createRelease(args)
		case "wiki_page":
			return impl.createWikiPage(args)
		case "pull_request":
			return impl.createPullRequest(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionCreate, resource, "not implemented"))
		}
	}
}

func (impl CreateImpl) createIssue(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", err.Error()))
	}

	title, _ := args["title"].(string)
	if title == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", "title is required"))
	}

	body, _ := args["body"].(string)
	if body == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", "body is required"))
	}

	opt := forgejo.CreateIssueOption{
		Title: title,
		Body:  body,
	}

	assignees, err := stringList(args, "assignees")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", err.Error()))
	}
	opt.Assignees = assignees

	milestone, ok, err := optionalPositiveInt(args, "milestone")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", err.Error()))
	}
	if ok {
		opt.Milestone = milestone
	}

	labels, err := int64List(args, "labels")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue", err.Error()))
	}
	opt.Labels = labels

	if dueDateStr, ok := args["due_date"].(string); ok && dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid due_date format (expected RFC3339): %w", err)
		}
		opt.Deadline = &dueDate
	}

	issue, _, err := impl.Client.CreateIssue(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return textResult((&types.Issue{Issue: issue}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createIssueComment(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue_comment", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue_comment", err.Error()))
	}

	body, _ := args["body"].(string)
	if body == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "issue_comment", "body is required"))
	}

	opt := forgejo.CreateIssueCommentOption{Body: body}
	comment, _, err := impl.Client.CreateIssueComment(owner, repo, index, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return textResult((&types.Comment{Comment: comment}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createLabel(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "label", err.Error()))
	}

	name, _ := args["name"].(string)
	if name == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "label", "name is required"))
	}

	color, _ := args["color"].(string)
	if color == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "label", "color is required"))
	}

	description, _ := args["description"].(string)

	opt := forgejo.CreateLabelOption{
		Name:        name,
		Color:       color,
		Description: description,
	}

	label, _, err := impl.Client.CreateLabel(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create label: %w", err)
	}

	return textResult((&types.Label{Label: label}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createMilestone(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "milestone", err.Error()))
	}

	title, _ := args["title"].(string)
	if title == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "milestone", "title is required"))
	}

	description, _ := args["description"].(string)

	opt := forgejo.CreateMilestoneOption{
		Title:       title,
		Description: description,
	}

	if dueDateStr, ok := args["due_date"].(string); ok && dueDateStr != "" {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid due_date format (expected RFC3339): %w", err)
		}
		opt.Deadline = &dueDate
	}

	milestone, _, err := impl.Client.CreateMilestone(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create milestone: %w", err)
	}

	return textResult((&types.Milestone{Milestone: milestone}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createRelease(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "release", err.Error()))
	}

	tagName, _ := args["tag_name"].(string)
	if tagName == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "release", "tag_name is required"))
	}

	name, _ := args["name"].(string)
	if name == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "release", "name is required"))
	}

	opt := forgejo.CreateReleaseOption{
		TagName: tagName,
		Title:   name,
	}

	if body, ok := args["body"].(string); ok {
		opt.Note = body
	}
	if target, ok := args["target_commitish"].(string); ok {
		opt.Target = target
	}
	if draft, ok := args["draft"].(bool); ok {
		opt.IsDraft = draft
	}
	if prerelease, ok := args["prerelease"].(bool); ok {
		opt.IsPrerelease = prerelease
	}

	release, _, err := impl.Client.CreateRelease(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create release: %w", err)
	}

	return textResult((&types.Release{Release: release}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createWikiPage(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "wiki_page", err.Error()))
	}

	title, _ := args["title"].(string)
	if title == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "wiki_page", "title is required"))
	}

	content, _ := args["content"].(string)
	if content == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "wiki_page", "content is required"))
	}

	message, _ := args["message"].(string)

	options := types.MyCreateWikiPageOptions{
		Title:         title,
		ContentBase64: base64.StdEncoding.EncodeToString([]byte(content)),
		Message:       message,
	}

	page, err := impl.Client.MyCreateWikiPage(owner, repo, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create wiki page: %w", err)
	}

	return textResult((&types.WikiPage{MyWikiPage: page}).ToMarkdown()), nil, nil
}

func (impl CreateImpl) createPullRequest(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", err.Error()))
	}

	title, _ := args["title"].(string)
	if title == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", "title is required"))
	}

	head, _ := args["head"].(string)
	if head == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", "head is required"))
	}

	base, _ := args["base"].(string)
	if base == "" {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", "base is required"))
	}

	opt := forgejo.CreatePullRequestOption{
		Title: title,
		Head:  head,
		Base:  base,
	}

	if body, ok := args["body"].(string); ok {
		opt.Body = body
	}

	assignees, err := stringList(args, "assignees")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", err.Error()))
	}
	opt.Assignees = assignees

	milestone, ok, err := optionalPositiveInt(args, "milestone")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", err.Error()))
	}
	if ok {
		opt.Milestone = milestone
	}

	labels, err := int64List(args, "labels")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionCreate, "pull_request", err.Error()))
	}
	opt.Labels = labels

	pr, _, err := impl.Client.CreatePullRequest(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return textResult((&types.PullRequest{PullRequest: pr}).ToMarkdown()), nil, nil
}

// Helper functions

func extractOwnerRepo(args map[string]any) (string, string, error) {
	owner, _ := args["owner"].(string)
	repo, _ := args["repo"].(string)
	if owner == "" {
		return "", "", fmt.Errorf("owner is required")
	}
	if repo == "" {
		return "", "", fmt.Errorf("repo is required")
	}
	return owner, repo, nil
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}
