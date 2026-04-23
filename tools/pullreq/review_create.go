// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// CreatePullRequestReviewComment is the per-file inline comment payload accepted
// by the create_pull_request_review tool. Either NewLine or OldLine must be set
// to indicate the side of the diff the comment applies to.
type CreatePullRequestReviewComment struct {
	// Path is the file path within the repository (relative to the repo root).
	Path string `json:"path"`
	// Body is the markdown body of the inline comment.
	Body string `json:"body"`
	// NewLine is the 1-based line number on the new (post-change) side of the diff.
	NewLine int `json:"new_line,omitempty"`
	// OldLine is the 1-based line number on the old (pre-change) side of the diff.
	OldLine int `json:"old_line,omitempty"`
}

// CreatePullRequestReviewParams defines the parameters for the
// create_pull_request_review tool.
type CreatePullRequestReviewParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// Event is the review state: APPROVED, REQUEST_CHANGES, COMMENT, or PENDING.
	Event string `json:"event,omitempty"`
	// Body is the overall review body (markdown).
	Body string `json:"body,omitempty"`
	// CommitID pins the review to a specific commit SHA (optional).
	CommitID string `json:"commit_id,omitempty"`
	// Comments is the list of inline comments to include in the review.
	Comments []CreatePullRequestReviewComment `json:"comments,omitempty"`
}

// CreatePullRequestReviewImpl implements the MCP tool for submitting a new pull
// request review, optionally with inline comments. This is how to "reply" to
// existing inline review comments — Forgejo threads inline replies by posting a
// new review containing comments at the same path/line.
type CreatePullRequestReviewImpl struct {
	Client *tools.Client
}

func (CreatePullRequestReviewImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_pull_request_review",
		Title:       "Create Pull Request Review",
		Description: "Submit a new review on a pull request, optionally including inline comments at specific file paths and line numbers. Use this to reply to existing inline review comments by posting a new review with comments at the same path and line.",
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
					Description: "Pull request index number",
				},
				"event": {
					Type:        "string",
					Description: "Review state (optional, defaults to COMMENT): APPROVED, REQUEST_CHANGES, COMMENT, or PENDING",
					Enum:        []any{"APPROVED", "REQUEST_CHANGES", "COMMENT", "PENDING"},
				},
				"body": {
					Type:        "string",
					Description: "Overall review body in markdown (optional, but required when event is not APPROVED and there are no inline comments)",
				},
				"commit_id": {
					Type:        "string",
					Description: "Commit SHA to pin the review to (optional)",
				},
				"comments": {
					Type:        "array",
					Description: "Inline comments to include in the review (optional)",
					Items: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"path": {
								Type:        "string",
								Description: "File path relative to the repository root",
							},
							"body": {
								Type:        "string",
								Description: "Inline comment body (markdown supported)",
							},
							"new_line": {
								Type:        "integer",
								Description: "1-based line number on the new (post-change) side of the diff. Set this OR old_line.",
								Minimum:     tools.Float64Ptr(1),
							},
							"old_line": {
								Type:        "integer",
								Description: "1-based line number on the old (pre-change) side of the diff. Set this OR new_line.",
								Minimum:     tools.Float64Ptr(1),
							},
						},
						Required: []string{"path", "body"},
					},
				},
			},
			Required: []string{"owner", "repo", "index"},
		},
	}
}

func (impl CreatePullRequestReviewImpl) Handler() mcp.ToolHandlerFor[CreatePullRequestReviewParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args CreatePullRequestReviewParams) (*mcp.CallToolResult, any, error) {
		p := args

		opt := forgejo.CreatePullReviewOptions{
			Body:     p.Body,
			CommitID: p.CommitID,
		}
		if p.Event != "" {
			opt.State = forgejo.ReviewStateType(p.Event)
		}
		if len(p.Comments) > 0 {
			opt.Comments = make([]forgejo.CreatePullReviewComment, len(p.Comments))
			for i, c := range p.Comments {
				opt.Comments[i] = forgejo.CreatePullReviewComment{
					Path:       c.Path,
					Body:       c.Body,
					NewLineNum: int64(c.NewLine),
					OldLineNum: int64(c.OldLine),
				}
			}
		}

		review, _, err := impl.Client.CreatePullReview(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create review: %w", err)
		}

		wrapper := &types.PullReview{PullReview: review}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: wrapper.ToMarkdown(),
				},
			},
		}, nil, nil
	}
}
