// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"fmt"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// CreatePullRequestParams defines the parameters for the create_pull_request tool.
// It captures the branches, title, description, and optional metadata for the new pull request.
type CreatePullRequestParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Head identifies the source branch (optionally `owner:branch`) to merge from.
	Head string `json:"head"`
	// Base identifies the target branch to merge into.
	Base string `json:"base"`
	// Title is the pull request title.
	Title string `json:"title"`
	// Body is the pull request description in Markdown.
	Body string `json:"body,omitempty"`
	// Assignee assigns the pull request to a single user.
	Assignee string `json:"assignee,omitempty"`
	// Assignees assigns multiple users to the pull request.
	Assignees []string `json:"assignees,omitempty"`
	// Milestone is the milestone ID associated with the pull request.
	Milestone int `json:"milestone,omitempty"`
	// Labels is a slice of label IDs to attach to the pull request.
	Labels []int `json:"labels,omitempty"`
	// DueDate is the optional deadline for the pull request.
	DueDate time.Time `json:"due_date"`
}

// CreatePullRequestImpl implements the MCP tool for creating a new pull request.
// This is a non-idempotent operation that opens a new pull request via the Forgejo SDK.
type CreatePullRequestImpl struct {
	Client *tools.Client
}

// Definition describes the `create_pull_request` tool. It requires `owner`, `repo`,
// `head`, `base`, and a `title`. Repeated calls will create multiple pull requests.
func (CreatePullRequestImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_pull_request",
		Title:       "Create Pull Request",
		Description: "Create a new pull request in a repository. Provide the source and target branches, title, body, and optional metadata.",
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
				"head": {
					Type:        "string",
					Description: "Source branch to merge from (optionally include owner as owner:branch)",
				},
				"base": {
					Type:        "string",
					Description: "Target branch to merge into",
				},
				"title": {
					Type:        "string",
					Description: "Pull request title",
				},
				"body": {
					Type:        "string",
					Description: "Pull request description (markdown supported) (optional)",
				},
				"assignee": {
					Type:        "string",
					Description: "Username of the primary assignee (optional)",
				},
				"assignees": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
					Description: "Additional usernames to assign to this pull request (optional)",
				},
				"milestone": {
					Type:        "integer",
					Description: "Milestone ID to associate with this pull request (optional)",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to attach to this pull request (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Pull request due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "head", "base", "title"},
		},
	}
}

// Handler implements the logic for creating a pull request. It calls the Forgejo SDK's
// `CreatePullRequest` function and returns the details of the newly created pull request.
func (impl CreatePullRequestImpl) Handler() mcp.ToolHandlerFor[CreatePullRequestParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args CreatePullRequestParams) (*mcp.CallToolResult, any, error) {
		p := args

		opt := forgejo.CreatePullRequestOption{
			Head:      p.Head,
			Base:      p.Base,
			Title:     p.Title,
			Body:      p.Body,
			Assignee:  p.Assignee,
			Assignees: p.Assignees,
		}

		if p.Milestone > 0 {
			opt.Milestone = int64(p.Milestone)
		}

		if len(p.Labels) > 0 {
			opt.Labels = make([]int64, len(p.Labels))
			for i, label := range p.Labels {
				opt.Labels[i] = int64(label)
			}
		}

		if !p.DueDate.IsZero() {
			opt.Deadline = &p.DueDate
		}

		pr, _, err := impl.Client.CreatePullRequest(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create pull request: %w", err)
		}

		prWrapper := &types.PullRequest{PullRequest: pr}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: prWrapper.ToMarkdown(),
				},
			},
		}, nil, nil
	}
}
