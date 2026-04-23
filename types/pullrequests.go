// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// PullRequest represents a pull request response with embedded SDK pull request
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls (list)
// - GET /repos/{owner}/{repo}/pulls/{index} (get)
type PullRequest struct {
	*forgejo.PullRequest
}

// ToMarkdown renders pull request with title, state, author and branch info
// Example: **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
func (pr *PullRequest) ToMarkdown() string {
	if pr.PullRequest == nil {
		return "*Invalid pull request*"
	}
	markdown := fmt.Sprintf("**#%d %s** (%s)\n", pr.Index, pr.Title, pr.State)
	if pr.Poster != nil {
		markdown += "Author: " + pr.Poster.UserName + "\n"
	}
	if pr.Head != nil && pr.Base != nil {
		markdown += fmt.Sprintf("Branch: %s → %s\n", pr.Head.Name, pr.Base.Name)
	}
	if pr.Body != "" {
		markdown += "\n" + pr.Body
	}
	return markdown
}

// PullRequestList represents a list of pull requests response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls
type PullRequestList []*PullRequest

// ToMarkdown renders pull requests as a numbered list with basic info
// Example:
// 1. **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
// 2. **#41 Fix database connection** (merged)
// Author: alice
// Branch: bugfix/db → main
func (prl PullRequestList) ToMarkdown() string {
	if len(prl) == 0 {
		return "*No pull requests found*"
	}
	markdown := ""
	for i, pr := range prl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, pr.ToMarkdown())
	}
	return markdown
}

// PullReview represents a pull request review response with embedded SDK review
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews (list)
// - POST /repos/{owner}/{repo}/pulls/{index}/reviews (create)
type PullReview struct {
	*forgejo.PullReview
}

// ToMarkdown renders a review with reviewer, state, commit, submission time and body.
// Also reports how many inline comments belong to the review so a caller can decide
// whether to drill in via list_pull_request_review_comments.
func (r *PullReview) ToMarkdown() string {
	if r.PullReview == nil {
		return "*Invalid review*"
	}
	markdown := fmt.Sprintf("Review#%d", r.ID)
	if r.State != "" {
		markdown += fmt.Sprintf(" (%s)", r.State)
	}
	if r.Reviewer != nil {
		markdown += " by **" + r.Reviewer.UserName + "**"
	}
	if !r.Submitted.IsZero() {
		markdown += " at " + r.Submitted.Format("2006-01-02 15:04")
	}
	markdown += "\n"
	if r.CommitID != "" {
		markdown += "Commit: " + r.CommitID + "\n"
	}
	markdown += fmt.Sprintf("Inline comments: %d\n", r.CodeCommentsCount)
	if r.Stale {
		markdown += "_Stale (PR changed since this review)_\n"
	}
	if r.Dismissed {
		markdown += "_Dismissed_\n"
	}
	if r.Body != "" {
		markdown += "\n" + r.Body
	}
	return markdown
}

// PullReviewList represents a list of reviews for a pull request.
type PullReviewList []*PullReview

// ToMarkdown renders reviews as a numbered list.
func (rl PullReviewList) ToMarkdown() string {
	if len(rl) == 0 {
		return "*No reviews found*"
	}
	markdown := ""
	for i, r := range rl {
		markdown += fmt.Sprintf("%d. %s\n\n", i+1, r.ToMarkdown())
	}
	return strings.TrimRight(markdown, "\n")
}

// PullReviewComment represents an inline review comment on a pull request.
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments
type PullReviewComment struct {
	*forgejo.PullReviewComment
}

// ToMarkdown renders an inline review comment with file path, line number,
// author, and the diff hunk that provides context for where the comment was left.
func (c *PullReviewComment) ToMarkdown() string {
	if c.PullReviewComment == nil {
		return "*Invalid review comment*"
	}
	markdown := fmt.Sprintf("Comment#%d", c.ID)
	if c.Reviewer != nil {
		markdown += " by **" + c.Reviewer.UserName + "**"
	}
	if !c.Created.IsZero() {
		markdown += " at " + c.Created.Format("2006-01-02 15:04")
	}
	markdown += "\n"
	if c.Path != "" {
		line := c.LineNum
		marker := "new"
		if line == 0 && c.OldLineNum != 0 {
			line = c.OldLineNum
			marker = "old"
		}
		if line > 0 {
			markdown += fmt.Sprintf("File: `%s` (%s line %d)\n", c.Path, marker, line)
		} else {
			markdown += fmt.Sprintf("File: `%s`\n", c.Path)
		}
	}
	if c.Resolver != nil {
		markdown += "Resolved by: " + c.Resolver.UserName + "\n"
	}
	if c.DiffHunk != "" {
		markdown += "\n```diff\n" + strings.TrimRight(c.DiffHunk, "\n") + "\n```\n"
	}
	if c.Body != "" {
		markdown += "\n" + c.Body
	}
	return markdown
}

// PullReviewCommentList represents a list of inline review comments.
type PullReviewCommentList []*PullReviewComment

// ToMarkdown renders inline review comments as a numbered list separated by rules.
func (cl PullReviewCommentList) ToMarkdown() string {
	if len(cl) == 0 {
		return "*No inline review comments found*"
	}
	parts := make([]string, 0, len(cl))
	for i, c := range cl {
		parts = append(parts, fmt.Sprintf("%d. %s", i+1, c.ToMarkdown()))
	}
	return strings.Join(parts, "\n\n---\n\n")
}
