// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListRepoIssuesParams defines the parameters for the list_repo_issues tool.
// It includes extensive options for filtering, sorting, and paginating issues.
type ListRepoIssuesParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// State filters issues by their state (e.g., 'open', 'closed').
	State string `json:"state,omitempty"`
	// Labels is a comma-separated list of label names to filter by.
	Labels string `json:"labels,omitempty"`
	// Milestones is a comma-separated list of milestone names to filter by.
	Milestones string `json:"milestones,omitempty"`
	// Assignees is a comma-separated list of usernames to filter by assignee.
	Assignees string `json:"assignees,omitempty"`
	// Q is a search query string to filter issues by.
	Q string `json:"q,omitempty"`
	// Sort specifies the sort order for the results.
	Sort string `json:"sort,omitempty"`
	// Order specifies the sort direction (asc or desc).
	Order string `json:"order,omitempty"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of issues to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListRepoIssuesImpl implements the read-only MCP tool for listing repository issues.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a list
// of issues with powerful filtering and sorting capabilities.
type ListRepoIssuesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_repo_issues` tool. It requires `owner` and `repo`
// and supports a rich set of optional parameters for filtering and sorting.
// It is marked as a safe, read-only operation.
func (ListRepoIssuesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_repo_issues",
		Title:       "List Repository Issues",
		Description: "List issues in a repository with optional filtering by state, labels, milestones, assignees, and search terms.",
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
				"state": {
					Type:        "string",
					Description: "Issue state filter: 'open', 'closed', or 'all' (optional, defaults to 'open')",
					Enum:        []any{"open", "closed", "all"},
				},
				"labels": {
					Type:        "string",
					Description: "Comma-separated list of label names to filter by (optional)",
				},
				"milestones": {
					Type:        "string",
					Description: "Comma-separated list of milestone names or IDs to filter by (optional)",
				},
				"assignees": {
					Type:        "string",
					Description: "Comma-separated list of usernames to filter by assignee (optional)",
				},
				"q": {
					Type:        "string",
					Description: "Search query string to filter issues (optional)",
				},
				"sort": {
					Type:        "string",
					Description: "Sort field: 'created', 'updated', 'comments' (optional, defaults to 'created')",
					Enum:        []any{"created", "updated", "comments"},
				},
				"order": {
					Type:        "string",
					Description: "Sort order: 'asc' or 'desc' (optional, defaults to 'desc')",
					Enum:        []any{"asc", "desc"},
				},
				"page": {
					Type:        "integer",
					Description: "Page number for pagination (optional, defaults to 1)",
					Minimum:     tools.Float64Ptr(1),
				},
				"limit": {
					Type:        "integer",
					Description: "Number of issues per page (optional, defaults to 20, max 50)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(50),
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing issues. It calls the Forgejo SDK's
// `ListRepoIssues` function with the provided filters and formats the results
// into a markdown table.
func (impl ListRepoIssuesImpl) Handler() mcp.ToolHandlerFor[ListRepoIssuesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListRepoIssuesParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.ListIssueOption{}
		if p.State != "" {
			opt.State = forgejo.StateType(p.State)
		}
		if p.Labels != "" {
			opt.Labels = strings.Split(p.Labels, ",")
		}
		if p.Milestones != "" {
			opt.Milestones = strings.Split(p.Milestones, ",")
		}
		if p.Q != "" {
			opt.KeyWord = p.Q
		}
		if p.Assignees != "" {
			opt.AssignedBy = p.Assignees
		}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.ListOptions.PageSize = p.Limit
		}

		// Call SDK
		issues, _, err := impl.Client.ListRepoIssues(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		// Convert to our types and format
		issueList := types.IssueList(issues)
		content := fmt.Sprintf("Found %d issues\n\n%s", len(issues), issueList.ToMarkdown())

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}

// GetIssueParams defines the parameters for the get_issue tool.
// It specifies the issue to retrieve by its owner, repository, and index.
type GetIssueParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
}

// GetIssueImpl implements the read-only MCP tool for fetching a single issue.
// This is a safe, idempotent operation that uses the Forgejo SDK to retrieve
// detailed information about a specific issue.
type GetIssueImpl struct {
	Client *tools.Client
}

// Definition describes the `get_issue` tool. It requires `owner`, `repo`, and
// the issue `index`. It is marked as a safe, read-only operation.
func (GetIssueImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_issue",
		Title:       "Get Issue Details",
		Description: "Get detailed information about a specific issue, including title, body, state, assignees, labels, and metadata.",
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

// Handler implements the logic for fetching an issue. It calls the Forgejo SDK's
// `GetIssue` function and formats the result into a detailed markdown view.
// It will return an error if the issue is not found.
func (impl GetIssueImpl) Handler() mcp.ToolHandlerFor[GetIssueParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetIssueParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call SDK
		issue, _, err := impl.Client.GetIssue(p.Owner, p.Repo, int64(p.Index))
		if err != nil {
			return nil, fmt.Errorf("failed to get issue: %w", err)
		}

		// Convert to our type and format
		issueWrapper := &types.Issue{Issue: issue}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: issueWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// CreateIssueParams defines the parameters for the create_issue tool.
// It includes the title, body, and optional metadata for the new issue.
type CreateIssueParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Title is the title of the new issue.
	Title string `json:"title"`
	// Body is the markdown content of the issue.
	Body string `json:"body"`
	// Assignees is a slice of usernames to assign to the issue.
	Assignees []string `json:"assignees,omitempty"`
	// Milestone is the ID of a milestone to assign to the issue.
	Milestone int `json:"milestone,omitempty"`
	// Labels is a slice of label IDs to assign to the issue.
	Labels []int `json:"labels,omitempty"`
	// DueDate is the optional due date for the issue.
	DueDate time.Time `json:"due_date,omitempty"`
}

// CreateIssueImpl implements the MCP tool for creating a new issue.
// This is a non-idempotent operation that creates a new issue using the Forgejo SDK.
type CreateIssueImpl struct {
	Client *tools.Client
}

// Definition describes the `create_issue` tool. It requires `owner`, `repo`,
// a `title`, and a `body`. It is not idempotent, as multiple calls will create
// multiple identical issues.
func (CreateIssueImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_issue",
		Title:       "Create Issue",
		Description: "Create a new issue in a repository. Specify title, body, and optional metadata like labels, assignees, milestone.",
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
				"title": {
					Type:        "string",
					Description: "Issue title",
				},
				"body": {
					Type:        "string",
					Description: "Issue body content (markdown supported)",
				},
				"assignees": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Array of usernames to assign to this issue (optional)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID to assign to this issue (optional)",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to assign to this issue (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Issue due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "title", "body"},
		},
	}
}

// Handler implements the logic for creating an issue. It calls the Forgejo SDK's
// `CreateIssue` function and returns the details of the newly created issue.
func (impl CreateIssueImpl) Handler() mcp.ToolHandlerFor[CreateIssueParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateIssueParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.CreateIssueOption{
			Title:     p.Title,
			Body:      p.Body,
			Assignees: p.Assignees,
		}

		// Set milestone if provided
		if p.Milestone > 0 {
			opt.Milestone = int64(p.Milestone)
		}

		// Convert label IDs from int to int64
		if len(p.Labels) > 0 {
			opt.Labels = make([]int64, len(p.Labels))
			for i, label := range p.Labels {
				opt.Labels[i] = int64(label)
			}
		}

		// Set due date if provided
		if !p.DueDate.IsZero() {
			opt.Deadline = &p.DueDate
		}

		// Call SDK
		issue, _, err := impl.Client.CreateIssue(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to create issue: %w", err)
		}

		// Convert to our type and format
		issueWrapper := &types.Issue{Issue: issue}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: issueWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// EditIssueParams defines the parameters for the edit_issue tool.
// It specifies the issue to edit by ID and the fields to update.
type EditIssueParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Title is the new title for the issue.
	Title string `json:"title,omitempty"`
	// Body is the new markdown content for the issue.
	Body string `json:"body,omitempty"`
	// State is the new state for the issue (e.g., 'open', 'closed').
	State string `json:"state,omitempty"`
	// Assignees is the new list of usernames to assign to the issue.
	Assignees []string `json:"assignees,omitempty"`
	// Milestone is the new milestone ID to assign to the issue.
	Milestone int `json:"milestone,omitempty"`
	// DueDate is the new optional due date for the issue.
	DueDate time.Time `json:"due_date,omitempty"`
}

// EditIssueImpl implements the MCP tool for editing an existing issue.
// This is an idempotent operation that modifies an issue's metadata using the
// Forgejo SDK.
type EditIssueImpl struct {
	Client *tools.Client
}

// Definition describes the `edit_issue` tool. It requires `owner`, `repo`, and the
// issue `index`. It is marked as idempotent.
func (EditIssueImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_issue",
		Title:       "Edit Issue",
		Description: "Edit an existing issue's title, body, state, assignees, milestone, or due date.",
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
				"title": {
					Type:        "string",
					Description: "New issue title (optional)",
				},
				"body": {
					Type:        "string",
					Description: "New issue body content (markdown supported) (optional)",
				},
				"state": {
					Type:        "string",
					Description: "New issue state: 'open' or 'closed' (optional)",
					Enum:        []any{"open", "closed"},
				},
				"assignees": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Array of usernames to assign to this issue (optional)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID to assign to this issue (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Issue due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

// Handler implements the logic for editing an issue. It calls the Forgejo SDK's
// `EditIssue` function. It will return an error if the issue is not found.
func (impl EditIssueImpl) Handler() mcp.ToolHandlerFor[EditIssueParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditIssueParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Build options for SDK call
		opt := forgejo.EditIssueOption{
			Assignees: p.Assignees,
		}

		// Set title if provided
		if p.Title != "" {
			opt.Title = p.Title
		}

		// Set body if provided
		if p.Body != "" {
			opt.Body = &p.Body
		}

		// Set state if provided
		if p.State != "" {
			state := forgejo.StateType(p.State)
			opt.State = &state
		}

		// Set milestone if provided
		if p.Milestone > 0 {
			milestone := int64(p.Milestone)
			opt.Milestone = &milestone
		}

		// Set due date if provided
		if !p.DueDate.IsZero() {
			opt.Deadline = &p.DueDate
		}

		// Call SDK
		issue, _, err := impl.Client.EditIssue(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to edit issue: %w", err)
		}

		// Convert to our type and format
		issueWrapper := &types.Issue{Issue: issue}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: issueWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}
