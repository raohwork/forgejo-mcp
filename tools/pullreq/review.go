// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// reviewPageSize is the page size used when walking all reviews of a pull request.
const reviewPageSize = 50

// maxReviewPages caps how many review pages are fetched when aggregating comments.
const maxReviewPages = 20

func prSchemaProps() map[string]*jsonschema.Schema {
	return map[string]*jsonschema.Schema{
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
	}
}

func textResult(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}

// ----------------------------------------------------------------------------
// list_pull_reviews

// ListPullReviewsParams defines the parameters for the list_pull_reviews tool.
type ListPullReviewsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// Page is the page number for pagination.
	Page int `json:"page,omitempty"`
	// Limit is the number of reviews to return per page.
	Limit int `json:"limit,omitempty"`
}

// ListPullReviewsImpl implements the read-only MCP tool for listing the
// reviews of a pull request.
type ListPullReviewsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_pull_reviews` tool.
func (ListPullReviewsImpl) Definition() *mcp.Tool {
	props := prSchemaProps()
	props["page"] = &jsonschema.Schema{
		Type:        "integer",
		Description: "Page number for pagination (optional, defaults to 1)",
		Minimum:     tools.Float64Ptr(1),
	}
	props["limit"] = &jsonschema.Schema{
		Type:        "integer",
		Description: "Number of reviews per page (optional)",
		Minimum:     tools.Float64Ptr(1),
		Maximum:     tools.Float64Ptr(50),
	}
	return &mcp.Tool{
		Name:        "list_pull_reviews",
		Title:       "List Pull Request Reviews",
		Description: "List reviews of a pull request, including reviewer, state (APPROVED, REQUEST_CHANGES, COMMENT, PENDING), body and number of inline comments.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type:       "object",
			Properties: props,
			Required:   []string{"owner", "repo", "index"},
		},
	}
}

// Handler lists the reviews of a pull request via the Forgejo SDK.
func (impl ListPullReviewsImpl) Handler() mcp.ToolHandlerFor[ListPullReviewsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListPullReviewsParams) (*mcp.CallToolResult, any, error) {
		p := args
		opt := forgejo.ListPullReviewsOptions{}
		if p.Page > 0 {
			opt.Page = p.Page
		}
		if p.Limit > 0 {
			opt.PageSize = p.Limit
		}

		reviews, _, err := impl.Client.ListPullReviews(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list pull request reviews: %w", err)
		}

		if len(reviews) == 0 {
			return textResult(fmt.Sprintf("No reviews found for pull request #%d.", p.Index)), nil, nil
		}

		list := make(types.PullReviewList, len(reviews))
		for i, r := range reviews {
			list[i] = &types.PullReview{PullReview: r}
		}
		return textResult(fmt.Sprintf("Found %d reviews on pull request #%d\n\n%s",
			len(reviews), p.Index, list.ToMarkdown())), nil, nil
	}
}

// ----------------------------------------------------------------------------
// list_pull_review_comments

// ListPullReviewCommentsParams defines the parameters for the
// list_pull_review_comments tool.
type ListPullReviewCommentsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// Path optionally restricts results to a single file.
	Path string `json:"path,omitempty"`
}

// ListPullReviewCommentsImpl implements the read-only MCP tool that collects
// the inline comments of every review on a pull request and groups them into
// conversations by file and line.
type ListPullReviewCommentsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_pull_review_comments` tool.
func (ListPullReviewCommentsImpl) Definition() *mcp.Tool {
	props := prSchemaProps()
	props["path"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Only show comments on this file path (optional)",
	}
	return &mcp.Tool{
		Name:        "list_pull_review_comments",
		Title:       "List Pull Request Review Comments",
		Description: "List inline code review comments on a pull request, grouped into conversations by file and line. Each comment shows its comment_id and review_id, which can be passed to reply_pull_review_comment.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type:       "object",
			Properties: props,
			Required:   []string{"owner", "repo", "index"},
		},
	}
}

// Handler walks all reviews of the pull request, fetches the inline comments of
// each review that has any, and renders them grouped into threads.
func (impl ListPullReviewCommentsImpl) Handler() mcp.ToolHandlerFor[ListPullReviewCommentsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListPullReviewCommentsParams) (*mcp.CallToolResult, any, error) {
		p := args
		owner, repo, index := p.Owner, p.Repo, int64(p.Index)

		var all []*forgejo.PullReviewComment
		for page := 1; page <= maxReviewPages; page++ {
			opt := forgejo.ListPullReviewsOptions{
				ListOptions: forgejo.ListOptions{Page: page, PageSize: reviewPageSize},
			}
			reviews, _, err := impl.Client.ListPullReviews(owner, repo, index, opt)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to list pull request reviews: %w", err)
			}
			for _, r := range reviews {
				if r == nil || r.CodeCommentsCount == 0 {
					continue
				}
				comments, _, err := impl.Client.ListPullReviewComments(owner, repo, index, r.ID)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to list comments of review %d: %w", r.ID, err)
				}
				for _, c := range comments {
					if c == nil || (p.Path != "" && c.Path != p.Path) {
						continue
					}
					all = append(all, c)
				}
			}
			if len(reviews) < reviewPageSize {
				break
			}
		}

		threads := types.GroupReviewThreads(all)
		if len(threads) == 0 {
			msg := fmt.Sprintf("No inline review comments found on pull request #%d", p.Index)
			if p.Path != "" {
				msg += " for " + p.Path
			}
			return textResult(msg + "."), nil, nil
		}

		return textResult(fmt.Sprintf("Found %d inline comments in %d conversations on pull request #%d\n\n%s",
			len(all), len(threads), p.Index, threads.ToMarkdown())), nil, nil
	}
}

// ----------------------------------------------------------------------------
// create_pull_review

// ReviewCommentParam is an inline comment attached to a line of a file in a
// pull request review. Exactly one of NewLine or OldLine must be set.
type ReviewCommentParam struct {
	// Path is the file path in the repository.
	Path string `json:"path"`
	// Body is the comment text in Markdown.
	Body string `json:"body"`
	// NewLine is the line number in the changed (new) version of the file.
	NewLine int `json:"new_line,omitempty"`
	// OldLine is the line number in the original (old) version of the file.
	OldLine int `json:"old_line,omitempty"`
}

// toSDK validates the comment and converts it to the SDK representation.
func (c ReviewCommentParam) toSDK() (forgejo.CreatePullReviewComment, error) {
	ret := forgejo.CreatePullReviewComment{
		Path:       c.Path,
		Body:       c.Body,
		NewLineNum: int64(c.NewLine),
		OldLineNum: int64(c.OldLine),
	}
	if strings.TrimSpace(c.Path) == "" {
		return ret, fmt.Errorf("path is required")
	}
	if strings.TrimSpace(c.Body) == "" {
		return ret, fmt.Errorf("body is required")
	}
	if c.NewLine < 0 || c.OldLine < 0 {
		return ret, fmt.Errorf("line numbers must be positive")
	}
	if (c.NewLine > 0) == (c.OldLine > 0) {
		return ret, fmt.Errorf("exactly one of new_line or old_line must be set")
	}
	return ret, nil
}

// CreatePullReviewParams defines the parameters for the create_pull_review tool.
type CreatePullReviewParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// Event is the review state: APPROVED, REQUEST_CHANGES, COMMENT or PENDING.
	Event string `json:"event,omitempty"`
	// Body is the overall review message in Markdown.
	Body string `json:"body,omitempty"`
	// CommitID is the commit the review applies to; defaults to the PR head.
	CommitID string `json:"commit_id,omitempty"`
	// Comments are inline comments attached to file lines.
	Comments []ReviewCommentParam `json:"comments,omitempty"`
}

// CreatePullReviewImpl implements the MCP tool for submitting a review on a
// pull request, optionally with inline comments.
type CreatePullReviewImpl struct {
	Client *tools.Client
}

var validReviewEvents = []string{
	string(forgejo.ReviewStateComment),
	string(forgejo.ReviewStateApproved),
	string(forgejo.ReviewStateRequestChanges),
	string(forgejo.ReviewStatePending),
}

// Definition describes the `create_pull_review` tool.
func (CreatePullReviewImpl) Definition() *mcp.Tool {
	events := make([]any, len(validReviewEvents))
	for i, e := range validReviewEvents {
		events[i] = e
	}

	props := prSchemaProps()
	props["event"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Review state: 'COMMENT', 'APPROVED', 'REQUEST_CHANGES', or 'PENDING' (saved as a draft, not submitted) (optional, defaults to 'COMMENT')",
		Enum:        events,
	}
	props["body"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Overall review message (markdown supported). Required unless event is APPROVED or inline comments are given.",
	}
	props["commit_id"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Commit SHA the review applies to (optional, defaults to the pull request head)",
	}
	props["comments"] = &jsonschema.Schema{
		Type:        "array",
		Description: "Inline comments on specific file lines (optional)",
		Items: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Type:        "string",
					Description: "File path in the repository",
				},
				"body": {
					Type:        "string",
					Description: "Comment text (markdown supported)",
				},
				"new_line": {
					Type:        "integer",
					Description: "Line number in the new version of the file (for added or unchanged lines). Set this or old_line, not both.",
					Minimum:     tools.Float64Ptr(1),
				},
				"old_line": {
					Type:        "integer",
					Description: "Line number in the old version of the file (for removed lines). Set this or new_line, not both.",
					Minimum:     tools.Float64Ptr(1),
				},
			},
			Required: []string{"path", "body"},
		},
	}

	return &mcp.Tool{
		Name:        "create_pull_review",
		Title:       "Create Pull Request Review",
		Description: "Submit a review on a pull request: approve, request changes, or comment, optionally with inline comments on specific file lines.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  false,
		},
		InputSchema: &jsonschema.Schema{
			Type:       "object",
			Properties: props,
			Required:   []string{"owner", "repo", "index"},
		},
	}
}

// Handler validates the input and creates the review via the Forgejo SDK.
func (impl CreatePullReviewImpl) Handler() mcp.ToolHandlerFor[CreatePullReviewParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args CreatePullReviewParams) (*mcp.CallToolResult, any, error) {
		p := args

		event := strings.ToUpper(strings.TrimSpace(p.Event))
		if event == "" {
			event = string(forgejo.ReviewStateComment)
		}
		valid := false
		for _, e := range validReviewEvents {
			if e == event {
				valid = true
				break
			}
		}
		if !valid {
			return nil, nil, fmt.Errorf("invalid event %q: must be one of %s", p.Event, strings.Join(validReviewEvents, ", "))
		}

		opt := forgejo.CreatePullReviewOptions{
			State:    forgejo.ReviewStateType(event),
			Body:     p.Body,
			CommitID: p.CommitID,
		}
		for i, c := range p.Comments {
			sc, err := c.toSDK()
			if err != nil {
				return nil, nil, fmt.Errorf("invalid comment #%d: %w", i+1, err)
			}
			opt.Comments = append(opt.Comments, sc)
		}

		review, _, err := impl.Client.CreatePullReview(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create pull request review: %w", err)
		}

		return textResult(fmt.Sprintf("Review created on pull request #%d\n\n%s",
			p.Index, (&types.PullReview{PullReview: review}).ToMarkdown())), nil, nil
	}
}

// ----------------------------------------------------------------------------
// reply_pull_review_comment

// ReplyPullReviewCommentParams defines the parameters for the
// reply_pull_review_comment tool.
type ReplyPullReviewCommentParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the pull request number.
	Index int `json:"index"`
	// ReviewID is the ID of the review containing the comment being replied to.
	ReviewID int `json:"review_id"`
	// CommentID is the ID of the inline comment being replied to.
	CommentID int `json:"comment_id"`
	// Body is the reply text in Markdown.
	Body string `json:"body"`
}

// ReplyPullReviewCommentImpl implements the MCP tool for replying to an inline
// review comment.
//
// Forgejo has no explicit reply relation for review comments: the web UI
// groups inline comments into a conversation by file path and line. A reply is
// therefore posted as a new COMMENT review with one inline comment at the
// same path and line (and diff side) as the target comment.
type ReplyPullReviewCommentImpl struct {
	Client *tools.Client
}

// Definition describes the `reply_pull_review_comment` tool.
func (ReplyPullReviewCommentImpl) Definition() *mcp.Tool {
	props := prSchemaProps()
	props["review_id"] = &jsonschema.Schema{
		Type:        "integer",
		Description: "ID of the review that contains the comment (from list_pull_review_comments)",
	}
	props["comment_id"] = &jsonschema.Schema{
		Type:        "integer",
		Description: "ID of the inline comment to reply to (from list_pull_review_comments)",
	}
	props["body"] = &jsonschema.Schema{
		Type:        "string",
		Description: "Reply text (markdown supported)",
	}
	return &mcp.Tool{
		Name:        "reply_pull_review_comment",
		Title:       "Reply to Pull Request Review Comment",
		Description: "Reply to an inline code review comment on a pull request. The reply is posted on the same file and line, so it appears in the same conversation.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  false,
		},
		InputSchema: &jsonschema.Schema{
			Type:       "object",
			Properties: props,
			Required:   []string{"owner", "repo", "index", "review_id", "comment_id", "body"},
		},
	}
}

// Handler looks up the target comment to learn its position, then posts the
// reply as a COMMENT review with a single inline comment at that position.
func (impl ReplyPullReviewCommentImpl) Handler() mcp.ToolHandlerFor[ReplyPullReviewCommentParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ReplyPullReviewCommentParams) (*mcp.CallToolResult, any, error) {
		p := args
		if strings.TrimSpace(p.Body) == "" {
			return nil, nil, fmt.Errorf("body is required")
		}
		owner, repo, index := p.Owner, p.Repo, int64(p.Index)

		comments, _, err := impl.Client.ListPullReviewComments(owner, repo, index, int64(p.ReviewID))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list comments of review %d: %w", p.ReviewID, err)
		}
		var target *forgejo.PullReviewComment
		for _, c := range comments {
			if c != nil && c.ID == int64(p.CommentID) {
				target = c
				break
			}
		}
		if target == nil {
			return nil, nil, fmt.Errorf("comment %d not found in review %d of pull request #%d", p.CommentID, p.ReviewID, p.Index)
		}

		side, line := types.CommentSide(target)
		if line == 0 {
			return nil, nil, fmt.Errorf("comment %d is not attached to a file line", p.CommentID)
		}
		reply := forgejo.CreatePullReviewComment{Path: target.Path, Body: p.Body}
		if side == types.SideNew {
			reply.NewLineNum = int64(line)
		} else {
			reply.OldLineNum = int64(line)
		}

		review, _, err := impl.Client.CreatePullReview(owner, repo, index, forgejo.CreatePullReviewOptions{
			State:    forgejo.ReviewStateComment,
			Comments: []forgejo.CreatePullReviewComment{reply},
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to post reply: %w", err)
		}

		return textResult(fmt.Sprintf("Replied to comment %d on %s (%s line %d) in pull request #%d\n\n%s",
			p.CommentID, target.Path, side, line, p.Index,
			(&types.PullReview{PullReview: review}).ToMarkdown())), nil, nil
	}
}
