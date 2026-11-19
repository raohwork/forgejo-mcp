// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"strings"
	"testing"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

func TestPullReview_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		review   *PullReview
		required []string
	}{
		{
			name: "full review",
			review: &PullReview{PullReview: &forgejo.PullReview{
				ID:                7,
				Reviewer:          testUser(),
				State:             forgejo.ReviewStateRequestChanges,
				Body:              "Please fix the error handling",
				CommitID:          "abcdef1234567890",
				CodeCommentsCount: 3,
				Stale:             true,
				Submitted:         testTime(),
			}},
			required: []string{"Review #7", "REQUEST_CHANGES", "testuser", "3 inline comments", "abcdef1", "stale", "2024-01-15", "Please fix the error handling"},
		},
		{
			name:     "nil review",
			review:   &PullReview{},
			required: []string{"*Invalid review*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertContains(t, tt.review.ToMarkdown(), tt.required)
		})
	}
}

func TestPullReviewList_ToMarkdown(t *testing.T) {
	list := PullReviewList{
		{PullReview: &forgejo.PullReview{ID: 1, State: forgejo.ReviewStateApproved, Reviewer: testUser()}},
		{PullReview: &forgejo.PullReview{ID: 2, State: forgejo.ReviewStateComment, Body: "nit"}},
	}
	assertContains(t, list.ToMarkdown(), []string{"1. **Review #1**", "APPROVED", "2. **Review #2**", "COMMENT", "nit"})
	assertContains(t, PullReviewList{}.ToMarkdown(), []string{"*No reviews found*"})
}

func TestPullReviewComment_ToMarkdown(t *testing.T) {
	c := &PullReviewComment{PullReviewComment: &forgejo.PullReviewComment{
		ID:       42,
		ReviewID: 7,
		Reviewer: testUser(),
		Body:     "Why not use a constant here?",
		Path:     "main.go",
		LineNum:  12,
		Created:  testTime(),
		Resolver: &forgejo.User{UserName: "alice"},
	}}
	assertContains(t, c.ToMarkdown(), []string{"comment_id: 42", "review_id: 7", "testuser", "2024-01-15", "resolved by alice", "Why not use a constant here?"})
	assertContains(t, (&PullReviewComment{}).ToMarkdown(), []string{"*Invalid review comment*"})
}

func TestGroupReviewThreads(t *testing.T) {
	t1 := testTime()
	t2 := t1.Add(time.Hour)
	t3 := t1.Add(2 * time.Hour)

	comments := []*forgejo.PullReviewComment{
		// reply on main.go new line 10, created later
		{ID: 3, ReviewID: 2, Path: "main.go", LineNum: 10, Body: "reply", Created: t2},
		// different file
		{ID: 4, ReviewID: 1, Path: "a.go", LineNum: 5, Body: "a comment", Created: t3},
		// root on main.go new line 10
		{ID: 1, ReviewID: 1, Path: "main.go", LineNum: 10, Body: "root", Created: t1, DiffHunk: "@@ -1,3 +1,4 @@\n+foo"},
		// same file, old side line 10: must be a separate thread
		{ID: 2, ReviewID: 1, Path: "main.go", OldLineNum: 10, Body: "old side", Created: t1},
		nil,
	}

	threads := GroupReviewThreads(comments)
	if len(threads) != 3 {
		t.Fatalf("expected 3 threads, got %d", len(threads))
	}

	// sorted by path, then line, old side before new side
	if threads[0].Path != "a.go" {
		t.Errorf("expected first thread on a.go, got %s", threads[0].Path)
	}
	if threads[1].Path != "main.go" || threads[1].Side != SideOld || threads[1].Line != 10 {
		t.Errorf("unexpected second thread: %+v", threads[1])
	}
	if threads[2].Path != "main.go" || threads[2].Side != SideNew || threads[2].Line != 10 {
		t.Errorf("unexpected third thread: %+v", threads[2])
	}
	if len(threads[2].Comments) != 2 || threads[2].Comments[0].ID != 1 || threads[2].Comments[1].ID != 3 {
		t.Errorf("expected thread comments ordered [1 3], got %+v", threads[2].Comments)
	}

	md := threads.ToMarkdown()
	assertContains(t, md, []string{
		"### a.go",
		"### main.go",
		"new line 10",
		"old line 10",
		"comment_id: 1",
		"comment_id: 3",
		"```diff",
		"+foo",
	})
	// the file header for main.go must appear exactly once
	if strings.Count(md, "### main.go") != 1 {
		t.Errorf("expected one main.go header, output: %s", md)
	}

	assertContains(t, ReviewThreadList{}.ToMarkdown(), []string{"*No inline review comments found*"})
}
