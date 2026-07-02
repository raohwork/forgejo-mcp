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
	"strings"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListImpl implements the list_gitea tool.
type ListImpl struct {
	Client *tools.Client
}

// Definition describes the list_gitea tool with minimal schema.
func (ListImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "list_gitea",
		Title: "List Gitea Resources",
		Description: `List resources from Forgejo/Gitea with filtering.
Resources: issue, issue_comment, issue_attachment, label, milestone, release, release_attachment, wiki_page, pull_request, repository, action_task, issue_dependency, issue_blocking.
Use gitea_manual(action="list") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resource": {
					Type:        "string",
					Description: "Resource type to list",
					Enum: []any{
						"issue", "issue_comment", "issue_attachment", "label",
						"milestone", "release", "release_attachment", "wiki_page",
						"pull_request", "repository", "action_task",
						"issue_dependency", "issue_blocking",
					},
				},
				"owner": {
					Type:        "string",
					Description: "Repository owner (not required for repository with scope='my')",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name (not required for repository listing)",
				},
				"index": {
					Type:        "integer",
					Description: "Issue number (for issue_comments, issue_attachments, issue_dependencies, issue_blocking)",
				},
				"id": {
					Type:        "integer",
					Description: "Release ID (for release_attachments)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Filter by milestone ID (for issues)",
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination",
				},
				"limit": {
					Type:        "integer",
					Description: "Results per page",
				},
			},
			Required:             []string{"resource"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate list logic based on resource type.
func (impl ListImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		resource, _ := args["resource"].(string)
		if resource == "" {
			resources := ListResourcesForAction(ActionList)
			return nil, nil, fmt.Errorf("resource is required. Valid resources: %v", resources)
		}

		if _, ok := LookupManual(ActionList, resource); !ok {
			resources := ListResourcesForAction(ActionList)
			return nil, nil, fmt.Errorf("unknown resource '%s'. Valid resources: %v", resource, resources)
		}

		switch resource {
		case "issue":
			return impl.listIssues(args)
		case "issue_comment":
			return impl.listIssueComments(args)
		case "issue_attachment":
			return impl.listIssueAttachments(args)
		case "label":
			return impl.listLabels(args)
		case "milestone":
			return impl.listMilestones(args)
		case "release":
			return impl.listReleases(args)
		case "release_attachment":
			return impl.listReleaseAttachments(args)
		case "wiki_page":
			return impl.listWikiPages(args)
		case "pull_request":
			return impl.listPullRequests(args)
		case "repository":
			return impl.listRepositories(args)
		case "action_task":
			return impl.listActionTasks(args)
		case "issue_dependency":
			return impl.listIssueDependencies(args)
		case "issue_blocking":
			return impl.listIssueBlocking(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionList, resource, "not implemented"))
		}
	}
}

func (impl ListImpl) listIssues(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue", err.Error()))
	}

	opt := forgejo.ListIssueOption{}

	if state, ok := args["state"].(string); ok && state != "" {
		opt.State = forgejo.StateType(state)
	}
	if labels, ok := args["labels"].(string); ok && labels != "" {
		opt.Labels = strings.Split(labels, ",")
	}
	if milestones, ok := args["milestones"].(string); ok && milestones != "" {
		opt.Milestones = strings.Split(milestones, ",")
	}
	if q, ok := args["q"].(string); ok && q != "" {
		opt.KeyWord = q
	}
	if assignees, ok := args["assignees"].(string); ok && assignees != "" {
		opt.AssignedBy = assignees
	}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}
	if sinceStr, ok := args["since"].(string); ok && sinceStr != "" {
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid since format (expected RFC3339): %w", err)
		}
		opt.Since = since
	}
	if beforeStr, ok := args["before"].(string); ok && beforeStr != "" {
		before, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid before format (expected RFC3339): %w", err)
		}
		opt.Before = before
	}

	issues, _, err := impl.Client.ListRepoIssues(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list issues: %w", err)
	}

	issueList := types.IssueList(issues)
	content := fmt.Sprintf("Found %d issues\n\n%s", len(issues), issueList.ToMarkdown())
	return textResult(content), nil, nil
}

func (impl ListImpl) listIssueComments(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_comment", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_comment", err.Error()))
	}

	opt := forgejo.ListIssueCommentOptions{}
	if sinceStr, ok := args["since"].(string); ok && sinceStr != "" {
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid since format: %w", err)
		}
		opt.Since = since
	}
	if beforeStr, ok := args["before"].(string); ok && beforeStr != "" {
		before, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid before format: %w", err)
		}
		opt.Before = before
	}

	comments, _, err := impl.Client.ListIssueComments(owner, repo, index, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list comments: %w", err)
	}

	if len(comments) == 0 {
		return textResult("No comments found for this issue."), nil, nil
	}

	var sb strings.Builder
	for _, comment := range comments {
		sb.WriteString((&types.Comment{Comment: comment}).ToMarkdown())
		sb.WriteString("\n\n---\n\n")
	}
	return textResult(fmt.Sprintf("Found %d comments\n\n%s", len(comments), sb.String())), nil, nil
}

func (impl ListImpl) listIssueAttachments(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_attachment", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_attachment", err.Error()))
	}

	attachments, err := impl.Client.MyListIssueAttachments(owner, repo, index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list attachments: %w", err)
	}

	if len(attachments) == 0 {
		return textResult("No attachments found for this issue."), nil, nil
	}

	list := make(types.AttachmentList, len(attachments))
	for i, a := range attachments {
		list[i] = &types.Attachment{Attachment: a}
	}
	return textResult(fmt.Sprintf("Found %d attachments\n\n%s", len(attachments), list.ToMarkdown())), nil, nil
}

func (impl ListImpl) listLabels(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "label", err.Error()))
	}

	labels, _, err := impl.Client.ListRepoLabels(owner, repo, forgejo.ListLabelsOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list labels: %w", err)
	}

	if len(labels) == 0 {
		return textResult("No labels found in this repository."), nil, nil
	}

	labelList := make(types.LabelList, len(labels))
	for i, label := range labels {
		labelList[i] = &types.Label{Label: label}
	}
	return textResult(fmt.Sprintf("Found %d labels\n\n%s", len(labels), labelList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listMilestones(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "milestone", err.Error()))
	}

	opt := forgejo.ListMilestoneOption{}
	if state, ok := args["state"].(string); ok && state != "" {
		opt.State = forgejo.StateType(state)
	}
	if name, ok := args["name"].(string); ok && name != "" {
		opt.Name = name
	}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	milestones, _, err := impl.Client.ListRepoMilestones(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list milestones: %w", err)
	}

	if len(milestones) == 0 {
		return textResult("No milestones found in this repository."), nil, nil
	}

	milestoneList := make(types.MilestoneList, len(milestones))
	for i, m := range milestones {
		milestoneList[i] = &types.Milestone{Milestone: m}
	}
	return textResult(fmt.Sprintf("Found %d milestones\n\n%s", len(milestones), milestoneList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listReleases(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "release", err.Error()))
	}

	opt := forgejo.ListReleasesOptions{}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	releases, _, err := impl.Client.ListReleases(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list releases: %w", err)
	}

	if len(releases) == 0 {
		return textResult("No releases found in this repository."), nil, nil
	}

	releaseList := make(types.ReleaseList, len(releases))
	for i, r := range releases {
		releaseList[i] = &types.Release{Release: r}
	}
	return textResult(fmt.Sprintf("Found %d releases\n\n%s", len(releases), releaseList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listReleaseAttachments(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "release_attachment", err.Error()))
	}

	id, err := requiredPositiveInt(args, "id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "release_attachment", err.Error()))
	}

	attachments, _, err := impl.Client.ListReleaseAttachments(owner, repo, id, forgejo.ListReleaseAttachmentsOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list release attachments: %w", err)
	}

	if len(attachments) == 0 {
		return textResult("No attachments found for this release."), nil, nil
	}

	list := make(types.AttachmentList, len(attachments))
	for i, a := range attachments {
		list[i] = &types.Attachment{Attachment: a}
	}
	return textResult(fmt.Sprintf("Found %d attachments\n\n%s", len(attachments), list.ToMarkdown())), nil, nil
}

func (impl ListImpl) listWikiPages(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "wiki_page", err.Error()))
	}

	pages, err := impl.Client.MyListWikiPages(owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list wiki pages: %w", err)
	}

	if len(pages) == 0 {
		return textResult("No wiki pages found in this repository."), nil, nil
	}

	list := types.WikiPageList(pages)
	return textResult(fmt.Sprintf("Found %d wiki pages\n\n%s", len(pages), list.ToMarkdown())), nil, nil
}

func (impl ListImpl) listPullRequests(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "pull_request", err.Error()))
	}

	opt := forgejo.ListPullRequestsOptions{}
	if state, ok := args["state"].(string); ok && state != "" {
		opt.State = forgejo.StateType(state)
	}
	if sort, ok := args["sort"].(string); ok && sort != "" {
		opt.Sort = sort
	}
	milestone, ok, err := optionalPositiveInt(args, "milestone")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Milestone = milestone
	}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	prs, _, err := impl.Client.ListRepoPullRequests(owner, repo, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	if len(prs) == 0 {
		return textResult("No pull requests found in this repository."), nil, nil
	}

	prList := make(types.PullRequestList, len(prs))
	for i, pr := range prs {
		prList[i] = &types.PullRequest{PullRequest: pr}
	}
	return textResult(fmt.Sprintf("Found %d pull requests\n\n%s", len(prs), prList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listRepositories(args map[string]any) (*mcp.CallToolResult, any, error) {
	scope, _ := args["scope"].(string)
	if scope == "" {
		return nil, nil, errors.New(FormatValidationError(ActionList, "repository", "scope is required ('my', 'org', or 'search')"))
	}

	switch scope {
	case "my":
		return impl.listMyRepositories(args)
	case "org":
		return impl.listOrgRepositories(args)
	case "search":
		return impl.searchRepositories(args)
	default:
		return nil, nil, errors.New(FormatValidationError(ActionList, "repository", "scope must be 'my', 'org', or 'search'"))
	}
}

func (impl ListImpl) listMyRepositories(args map[string]any) (*mcp.CallToolResult, any, error) {
	opt := forgejo.ListReposOptions{}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	repos, _, err := impl.Client.ListMyRepos(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list my repositories: %w", err)
	}

	if len(repos) == 0 {
		return textResult("No repositories found for the authenticated user."), nil, nil
	}

	repoList := make(types.RepositoryList, len(repos))
	for i, r := range repos {
		repoList[i] = &types.Repository{Repository: r}
	}
	return textResult(fmt.Sprintf("Found %d repositories\n\n%s", len(repos), repoList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listOrgRepositories(args map[string]any) (*mcp.CallToolResult, any, error) {
	org, _ := args["org"].(string)
	if org == "" {
		return nil, nil, errors.New(FormatValidationError(ActionList, "repository", "org is required for scope='org'"))
	}

	opt := forgejo.ListOrgReposOptions{}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	repos, _, err := impl.Client.ListOrgRepos(org, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list organization repositories: %w", err)
	}

	if len(repos) == 0 {
		return textResult(fmt.Sprintf("No repositories found for organization '%s'.", org)), nil, nil
	}

	repoList := make(types.RepositoryList, len(repos))
	for i, r := range repos {
		repoList[i] = &types.Repository{Repository: r}
	}
	return textResult(fmt.Sprintf("Found %d repositories for '%s'\n\n%s", len(repos), org, repoList.ToMarkdown())), nil, nil
}

func (impl ListImpl) searchRepositories(args map[string]any) (*mcp.CallToolResult, any, error) {
	q, _ := args["q"].(string)

	opt := forgejo.SearchRepoOptions{Keyword: q}
	if topic, ok := args["topic"].(bool); ok {
		opt.KeywordIsTopic = topic
	}
	if includeDesc, ok := args["include_desc"].(bool); ok {
		opt.KeywordInDescription = includeDesc
	}
	page, ok, err := optionalPositiveInt(args, "page")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.Page = int(page)
	}
	limit, ok, err := optionalPositiveInt(args, "limit")
	if err != nil {
		return nil, nil, err
	}
	if ok {
		opt.PageSize = int(limit)
	}

	repos, _, err := impl.Client.SearchRepos(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	if len(repos) == 0 {
		return textResult("No repositories found matching the search criteria."), nil, nil
	}

	repoList := make(types.RepositoryList, len(repos))
	for i, r := range repos {
		repoList[i] = &types.Repository{Repository: r}
	}
	return textResult(fmt.Sprintf("Found %d repositories\n\n%s", len(repos), repoList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listActionTasks(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "action_task", err.Error()))
	}

	response, err := impl.Client.MyListActionTasks(owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list action tasks: %w", err)
	}

	if response.TotalCount == 0 || len(response.WorkflowRuns) == 0 {
		return textResult("No action tasks found in this repository."), nil, nil
	}

	taskList := types.ActionTaskList{MyActionTaskResponse: response}
	return textResult(fmt.Sprintf("Found %d action tasks\n\n%s", response.TotalCount, taskList.ToMarkdown())), nil, nil
}

func (impl ListImpl) listIssueDependencies(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_dependency", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_dependency", err.Error()))
	}

	issues, err := impl.Client.MyListIssueDependencies(owner, repo, index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list dependencies: %w", err)
	}

	deps := types.IssueDependencyList(issues)
	return textResult(fmt.Sprintf("## Issues that block #%d\n\n%s", index, deps.ToMarkdown())), nil, nil
}

func (impl ListImpl) listIssueBlocking(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_blocking", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionList, "issue_blocking", err.Error()))
	}

	issues, err := impl.Client.MyListIssueBlocking(owner, repo, index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list blocking issues: %w", err)
	}

	blocking := types.IssueBlockingList(issues)
	return textResult(fmt.Sprintf("## Issues blocked by #%d\n\n%s", index, blocking.ToMarkdown())), nil, nil
}
