// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package pullreq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

const testForgejoVersion = "11.0.1+gitea-1.22.0"

// fakeReviewServer serves a minimal subset of the pull review API for one PR
// (owner/repo#5) and records the review creation requests it receives.
type fakeReviewServer struct {
	mu       sync.Mutex
	reviews  []*forgejo.PullReview
	comments map[int64][]*forgejo.PullReviewComment
	created  []map[string]any
}

func (f *fakeReviewServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	const prefix = "/api/v1/repos/owner/repo/pulls/5/reviews"
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == prefix:
		if page := r.URL.Query().Get("page"); page != "" && page != "1" {
			json.NewEncoder(w).Encode([]any{})
			return
		}
		json.NewEncoder(w).Encode(f.reviews)
	case r.Method == http.MethodPost && r.URL.Path == prefix:
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		f.created = append(f.created, body)
		json.NewEncoder(w).Encode(&forgejo.PullReview{
			ID:    100,
			State: forgejo.ReviewStateType(body["event"].(string)),
			Body:  body["body"].(string),
		})
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, prefix+"/") && strings.HasSuffix(r.URL.Path, "/comments"):
		var id int64
		for rid := range f.comments {
			if r.URL.Path == prefix+"/"+itoa(rid)+"/comments" {
				id = rid
			}
		}
		if id == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message":"not found"}`))
			return
		}
		json.NewEncoder(w).Encode(f.comments[id])
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"unexpected request"}`))
	}
}

func itoa(i int64) string {
	b, _ := json.Marshal(i)
	return string(b)
}

func newFakeReviewServer(t *testing.T) (*fakeReviewServer, *tools.Client) {
	t.Helper()
	f := &fakeReviewServer{
		reviews: []*forgejo.PullReview{
			{ID: 1, State: forgejo.ReviewStateRequestChanges, Reviewer: &forgejo.User{UserName: "alice"}, CodeCommentsCount: 2},
			{ID: 2, State: forgejo.ReviewStateApproved, Reviewer: &forgejo.User{UserName: "bob"}},
		},
		comments: map[int64][]*forgejo.PullReviewComment{
			1: {
				{ID: 11, ReviewID: 1, Path: "main.go", LineNum: 10, Body: "use a constant", Reviewer: &forgejo.User{UserName: "alice"}},
				{ID: 12, ReviewID: 1, Path: "util.go", OldLineNum: 4, Body: "why removed?", Reviewer: &forgejo.User{UserName: "alice"}},
			},
		},
	}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)

	cl, err := tools.NewClient(srv.URL, "token", testForgejoVersion, srv.Client())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return f, cl
}

func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if res == nil || len(res.Content) != 1 {
		t.Fatalf("expected a single content item, got %+v", res)
	}
	txt, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return txt.Text
}

func TestListPullReviews(t *testing.T) {
	_, cl := newFakeReviewServer(t)
	res, _, err := ListPullReviewsImpl{Client: cl}.Handler()(context.Background(), nil, ListPullReviewsParams{
		Owner: "owner", Repo: "repo", Index: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := resultText(t, res)
	for _, want := range []string{"Review #1", "REQUEST_CHANGES", "alice", "Review #2", "APPROVED", "bob"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestListPullReviewComments(t *testing.T) {
	_, cl := newFakeReviewServer(t)
	handler := ListPullReviewCommentsImpl{Client: cl}.Handler()

	t.Run("all files", func(t *testing.T) {
		res, _, err := handler(context.Background(), nil, ListPullReviewCommentsParams{
			Owner: "owner", Repo: "repo", Index: 5,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := resultText(t, res)
		for _, want := range []string{"### main.go", "new line 10", "comment_id: 11", "### util.go", "old line 4", "comment_id: 12"} {
			if !strings.Contains(out, want) {
				t.Errorf("output missing %q: %s", want, out)
			}
		}
	})

	t.Run("path filter", func(t *testing.T) {
		res, _, err := handler(context.Background(), nil, ListPullReviewCommentsParams{
			Owner: "owner", Repo: "repo", Index: 5, Path: "util.go",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := resultText(t, res)
		if strings.Contains(out, "main.go") {
			t.Errorf("expected main.go to be filtered out: %s", out)
		}
		if !strings.Contains(out, "comment_id: 12") {
			t.Errorf("expected util.go comment: %s", out)
		}
	})
}

func TestReplyPullReviewComment(t *testing.T) {
	tests := []struct {
		name      string
		commentID int
		wantPath  string
		wantNew   float64
		wantOld   float64
	}{
		{name: "new side", commentID: 11, wantPath: "main.go", wantNew: 10},
		{name: "old side", commentID: 12, wantPath: "util.go", wantOld: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cl := newFakeReviewServer(t)
			res, _, err := ReplyPullReviewCommentImpl{Client: cl}.Handler()(context.Background(), nil, ReplyPullReviewCommentParams{
				Owner: "owner", Repo: "repo", Index: 5, ReviewID: 1, CommentID: tt.commentID, Body: "done, thanks",
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(f.created) != 1 {
				t.Fatalf("expected 1 review created, got %d", len(f.created))
			}
			got := f.created[0]
			if got["event"] != "COMMENT" {
				t.Errorf("expected COMMENT event, got %v", got["event"])
			}
			cs, _ := got["comments"].([]any)
			if len(cs) != 1 {
				t.Fatalf("expected 1 inline comment, got %v", got["comments"])
			}
			c := cs[0].(map[string]any)
			if c["path"] != tt.wantPath || c["new_position"] != tt.wantNew || c["old_position"] != tt.wantOld || c["body"] != "done, thanks" {
				t.Errorf("unexpected inline comment: %v", c)
			}
			out := resultText(t, res)
			if !strings.Contains(out, tt.wantPath) {
				t.Errorf("output should mention %s: %s", tt.wantPath, out)
			}
		})
	}

	t.Run("comment not found", func(t *testing.T) {
		f, cl := newFakeReviewServer(t)
		_, _, err := ReplyPullReviewCommentImpl{Client: cl}.Handler()(context.Background(), nil, ReplyPullReviewCommentParams{
			Owner: "owner", Repo: "repo", Index: 5, ReviewID: 1, CommentID: 999, Body: "x",
		})
		if err == nil || !strings.Contains(err.Error(), "999") {
			t.Errorf("expected not found error mentioning 999, got %v", err)
		}
		if len(f.created) != 0 {
			t.Errorf("no review should be created")
		}
	})

	t.Run("empty body", func(t *testing.T) {
		_, cl := newFakeReviewServer(t)
		_, _, err := ReplyPullReviewCommentImpl{Client: cl}.Handler()(context.Background(), nil, ReplyPullReviewCommentParams{
			Owner: "owner", Repo: "repo", Index: 5, ReviewID: 1, CommentID: 11, Body: "  ",
		})
		if err == nil {
			t.Errorf("expected error for empty body")
		}
	})
}

func TestCreatePullReview(t *testing.T) {
	t.Run("with inline comments", func(t *testing.T) {
		f, cl := newFakeReviewServer(t)
		res, _, err := CreatePullReviewImpl{Client: cl}.Handler()(context.Background(), nil, CreatePullReviewParams{
			Owner: "owner", Repo: "repo", Index: 5,
			Event: "REQUEST_CHANGES",
			Body:  "a few issues",
			Comments: []ReviewCommentParam{
				{Path: "main.go", Body: "typo", NewLine: 3},
				{Path: "old.go", Body: "keep this", OldLine: 8},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := f.created[0]
		if got["event"] != "REQUEST_CHANGES" || got["body"] != "a few issues" {
			t.Errorf("unexpected review: %v", got)
		}
		cs := got["comments"].([]any)
		if len(cs) != 2 {
			t.Fatalf("expected 2 comments, got %v", cs)
		}
		c0, c1 := cs[0].(map[string]any), cs[1].(map[string]any)
		if c0["new_position"] != float64(3) || c0["old_position"] != float64(0) {
			t.Errorf("unexpected first comment: %v", c0)
		}
		if c1["old_position"] != float64(8) || c1["new_position"] != float64(0) {
			t.Errorf("unexpected second comment: %v", c1)
		}
		if !strings.Contains(resultText(t, res), "Review #100") {
			t.Errorf("expected created review in output")
		}
	})

	t.Run("defaults to COMMENT", func(t *testing.T) {
		f, cl := newFakeReviewServer(t)
		_, _, err := CreatePullReviewImpl{Client: cl}.Handler()(context.Background(), nil, CreatePullReviewParams{
			Owner: "owner", Repo: "repo", Index: 5, Body: "looks fine overall",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.created[0]["event"] != "COMMENT" {
			t.Errorf("expected COMMENT, got %v", f.created[0]["event"])
		}
	})

	invalid := []struct {
		name string
		args CreatePullReviewParams
	}{
		{"bad event", CreatePullReviewParams{Event: "MERGE", Body: "x"}},
		{"both lines", CreatePullReviewParams{Comments: []ReviewCommentParam{{Path: "a", Body: "x", NewLine: 1, OldLine: 1}}}},
		{"no line", CreatePullReviewParams{Comments: []ReviewCommentParam{{Path: "a", Body: "x"}}}},
		{"no path", CreatePullReviewParams{Comments: []ReviewCommentParam{{Body: "x", NewLine: 1}}}},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			f, cl := newFakeReviewServer(t)
			args := tt.args
			args.Owner, args.Repo, args.Index = "owner", "repo", 5
			_, _, err := CreatePullReviewImpl{Client: cl}.Handler()(context.Background(), nil, args)
			if err == nil {
				t.Errorf("expected validation error")
			}
			if len(f.created) != 0 {
				t.Errorf("no review should be created on validation error")
			}
		})
	}
}
