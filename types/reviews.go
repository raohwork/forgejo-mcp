// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"
	"sort"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// PullReview represents a pull request review response with embedded SDK review
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews (list)
// - POST /repos/{owner}/{repo}/pulls/{index}/reviews (create)
type PullReview struct {
	*forgejo.PullReview
}

// ToMarkdown renders a review with its ID, state, reviewer and body
// Example: **Review #7** (REQUEST_CHANGES) by johndoe
// Submitted: 2024-01-15 | Commit: abcdef1 | 3 inline comments
//
// Please fix the error handling
func (r *PullReview) ToMarkdown() string {
	if r.PullReview == nil {
		return "*Invalid review*"
	}
	markdown := fmt.Sprintf("**Review #%d** (%s)", r.ID, r.State)
	if r.Reviewer != nil {
		markdown += " by " + r.Reviewer.UserName
	} else if r.ReviewerTeam != nil {
		markdown += " by team " + r.ReviewerTeam.Name
	}
	markdown += "\n"

	meta := make([]string, 0, 5)
	if !r.Submitted.IsZero() {
		meta = append(meta, "Submitted: "+r.Submitted.Format("2006-01-02 15:04"))
	}
	if r.CommitID != "" {
		meta = append(meta, "Commit: "+shortSHA(r.CommitID))
	}
	if r.CodeCommentsCount > 0 {
		meta = append(meta, fmt.Sprintf("%d inline comments", r.CodeCommentsCount))
	}
	if r.Stale {
		meta = append(meta, "stale")
	}
	if r.Dismissed {
		meta = append(meta, "dismissed")
	}
	if len(meta) > 0 {
		markdown += strings.Join(meta, " | ") + "\n"
	}
	if r.Body != "" {
		markdown += "\n" + r.Body + "\n"
	}
	return markdown
}

// PullReviewList represents a list of pull request reviews
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews
type PullReviewList []*PullReview

// ToMarkdown renders reviews as a numbered list
func (l PullReviewList) ToMarkdown() string {
	if len(l) == 0 {
		return "*No reviews found*"
	}
	markdown := ""
	for i, r := range l {
		markdown += fmt.Sprintf("%d. %s\n", i+1, r.ToMarkdown())
	}
	return markdown
}

// PullReviewComment represents an inline (code) review comment
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments
type PullReviewComment struct {
	*forgejo.PullReviewComment
}

// ToMarkdown renders a single inline comment with the IDs needed to reply to it
// Example: - **johndoe** (2024-01-15 14:30) [comment_id: 42, review_id: 7]
//
//	Why not use a constant here?
func (c *PullReviewComment) ToMarkdown() string {
	if c.PullReviewComment == nil {
		return "*Invalid review comment*"
	}
	author := "unknown"
	if c.Reviewer != nil {
		author = c.Reviewer.UserName
	}
	markdown := fmt.Sprintf("- **%s**", author)
	if !c.Created.IsZero() {
		markdown += " (" + c.Created.Format("2006-01-02 15:04") + ")"
	}
	markdown += fmt.Sprintf(" [comment_id: %d, review_id: %d]", c.ID, c.ReviewID)
	if c.Resolver != nil {
		markdown += " (resolved by " + c.Resolver.UserName + ")"
	}
	markdown += "\n"
	for _, line := range strings.Split(strings.TrimRight(c.Body, "\n"), "\n") {
		markdown += "  " + line + "\n"
	}
	return markdown
}

// Diff side constants for ReviewThread.Side
const (
	// SideOld means the comment is attached to a line of the original file
	SideOld = "old"
	// SideNew means the comment is attached to a line of the changed file
	SideNew = "new"
)

// ReviewThread is a conversation of inline comments attached to the same
// file, diff side and line. Forgejo has no explicit reply relation between
// review comments; the web UI groups them into a conversation by position.
type ReviewThread struct {
	Path     string
	Side     string
	Line     uint64
	Comments []*forgejo.PullReviewComment
}

// ToMarkdown renders the thread location, a short diff excerpt and all comments
func (t *ReviewThread) ToMarkdown() string {
	markdown := fmt.Sprintf("#### %s line %d\n", t.Side, t.Line)
	if len(t.Comments) > 0 && t.Comments[0].DiffHunk != "" {
		markdown += "```diff\n" + lastLines(t.Comments[0].DiffHunk, 8) + "\n```\n"
	}
	for _, c := range t.Comments {
		markdown += (&PullReviewComment{PullReviewComment: c}).ToMarkdown()
	}
	return markdown
}

// ReviewThreadList is a list of inline comment threads
// Used by tools that aggregate:
// - GET /repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments
type ReviewThreadList []*ReviewThread

// ToMarkdown renders threads grouped by file
// Example:
// ### main.go
// #### new line 10
// ```diff
// +foo
// ```
// - **johndoe** (2024-01-15 14:30) [comment_id: 1, review_id: 1]
// root
func (l ReviewThreadList) ToMarkdown() string {
	if len(l) == 0 {
		return "*No inline review comments found*"
	}
	markdown := ""
	lastPath := ""
	for i, t := range l {
		if i == 0 || t.Path != lastPath {
			if i > 0 {
				markdown += "\n"
			}
			markdown += "### " + t.Path + "\n"
			lastPath = t.Path
		}
		markdown += t.ToMarkdown()
	}
	return markdown
}

// CommentSide returns the diff side and line number a review comment is
// attached to. Forgejo reports new-side lines in `position` and old-side
// lines in `original_position`.
func CommentSide(c *forgejo.PullReviewComment) (string, uint64) {
	if c.LineNum > 0 {
		return SideNew, c.LineNum
	}
	return SideOld, c.OldLineNum
}

// GroupReviewThreads groups inline review comments into conversations by
// (path, side, line). Threads are sorted by path, then line (old side first),
// and comments within a thread are sorted by creation time.
func GroupReviewThreads(comments []*forgejo.PullReviewComment) ReviewThreadList {
	type key struct {
		path string
		side string
		line uint64
	}
	index := map[key]*ReviewThread{}
	threads := ReviewThreadList{}
	for _, c := range comments {
		if c == nil {
			continue
		}
		side, line := CommentSide(c)
		k := key{c.Path, side, line}
		t, ok := index[k]
		if !ok {
			t = &ReviewThread{Path: c.Path, Side: side, Line: line}
			index[k] = t
			threads = append(threads, t)
		}
		t.Comments = append(t.Comments, c)
	}

	for _, t := range threads {
		sort.SliceStable(t.Comments, func(i, j int) bool {
			a, b := t.Comments[i], t.Comments[j]
			if !a.Created.Equal(b.Created) {
				return a.Created.Before(b.Created)
			}
			return a.ID < b.ID
		})
	}
	sort.SliceStable(threads, func(i, j int) bool {
		a, b := threads[i], threads[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.Side == SideOld && b.Side != SideOld
	})
	return threads
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

func lastLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
