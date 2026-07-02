package unified

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// fakeForgejo records requests and returns canned responses.
type fakeForgejo struct {
	mu       sync.Mutex
	requests []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Body   string
}

func (f *fakeForgejo) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.requests = append(f.requests, recordedRequest{r.Method, r.URL.Path, string(body)})
		f.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/version":
			w.Write([]byte(`{"version":"11.0.0"}`))
		case r.URL.Path == "/api/v1/repos/o/r/issues" && r.Method == "POST":
			w.WriteHeader(201)
			w.Write([]byte(`{"id":1,"number":326,"title":"t","body":"b","state":"open","labels":[]}`))
		case r.URL.Path == "/api/v1/repos/o/r/issues/326/labels" && r.Method == "POST":
			w.Write([]byte(`[{"id":1,"name":"bug","color":"ff0000"}]`))
		default:
			w.Write([]byte(`{}`))
		}
	})
}

// startSession spins up the MCP server over an in-memory transport and
// returns a connected client session.
func startSession(t *testing.T, backend string) *mcp.ClientSession {
	t.Helper()

	cl, err := tools.NewClient(backend, "dummy-token", "11.0.0", nil)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{Title: "test", Version: "0.0.0"}, nil)
	RegisterAll(server, cl)

	st, ct := mcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(context.Background(), ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

func callTool(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool(%s): %v", name, err)
	}
	return res
}

func resultText(res *mcp.CallToolResult) string {
	out := ""
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			out += tc.Text
		}
	}
	return out
}

func TestLinkIssueLabelIntegerIndex(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	res := callTool(t, cs, "link_gitea", map[string]any{
		"type":   "issue_label",
		"owner":  "o",
		"repo":   "r",
		"index":  326,
		"labels": []any{1, 2},
	})
	if res.IsError {
		t.Fatalf("tool errored: %s", resultText(res))
	}

	// Verify the backend actually received the label IDs.
	found := false
	for _, r := range fake.requests {
		if r.Path == "/api/v1/repos/o/r/issues/326/labels" {
			found = true
			var body struct {
				Labels []int64 `json:"labels"`
			}
			if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
				t.Fatalf("bad body %q: %v", r.Body, err)
			}
			if len(body.Labels) != 2 || body.Labels[0] != 1 || body.Labels[1] != 2 {
				t.Errorf("backend got labels %v, want [1 2]", body.Labels)
			}
		}
	}
	if !found {
		t.Errorf("no label request reached backend; requests: %+v", fake.requests)
	}
}

// callToolExpectTypeError calls a tool with intentionally mistyped arguments
// and asserts the call fails loudly — either as a schema-validation protocol
// error or as a tool error — with a message that names the offending
// parameter. Returns the error text.
func callToolExpectTypeError(t *testing.T, cs *mcp.ClientSession, name string, args map[string]any, param string) string {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	var msg string
	switch {
	case err != nil:
		msg = err.Error()
	case res.IsError:
		msg = resultText(res)
	default:
		t.Fatalf("CallTool(%s) with mistyped %q succeeded, want loud error; result: %s", name, param, resultText(res))
	}
	if !strings.Contains(msg, param) {
		t.Errorf("error should name the offending parameter %q, got: %s", param, msg)
	}
	return msg
}

// TestLinkIssueLabelStringArguments covers the failure reported from the
// McRogueFace project: LLM clients sent numbers as JSON strings because the
// minimal schema gave no type information. The schema now declares integer
// types, so mistyped values are rejected up front with a message naming the
// parameter and the expected type — not the old misleading "index is
// required" — and nothing reaches the backend.
func TestLinkIssueLabelStringArguments(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	msg := callToolExpectTypeError(t, cs, "link_gitea", map[string]any{
		"type":   "issue_label",
		"owner":  "o",
		"repo":   "r",
		"index":  "326",
		"labels": []any{"1", "2"},
	}, "index")
	if !strings.Contains(msg, "integer") {
		t.Errorf("error should state the expected type, got: %s", msg)
	}
	for _, r := range fake.requests {
		if strings.Contains(r.Path, "/labels") {
			t.Errorf("mistyped call must not reach the backend; got %s %s", r.Method, r.Path)
		}
	}
}

// TestCreateIssueStringLabelsNotSilentlyDropped covers the silent label
// drop: create_gitea used to ignore label IDs it could not type-assert,
// creating the issue unlabeled with no error. Now the mistyped array is
// rejected before the issue is created.
func TestCreateIssueStringLabelsNotSilentlyDropped(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	callToolExpectTypeError(t, cs, "create_gitea", map[string]any{
		"resource": "issue",
		"owner":    "o",
		"repo":     "r",
		"title":    "t",
		"body":     "b",
		"labels":   []any{"1", "2", "3"},
	}, "labels")

	for _, r := range fake.requests {
		if r.Path == "/api/v1/repos/o/r/issues" && r.Method == "POST" {
			t.Errorf("issue must not be created when labels are mistyped; body: %s", r.Body)
		}
	}
}

// TestCreateIssueGarbageLabelsFailLoudly: unparseable label IDs must produce
// an explicit error, never a silently-unlabeled issue.
func TestCreateIssueGarbageLabelsFailLoudly(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	callToolExpectTypeError(t, cs, "create_gitea", map[string]any{
		"resource": "issue",
		"owner":    "o",
		"repo":     "r",
		"title":    "t",
		"body":     "b",
		"labels":   []any{"bug", "enhancement"},
	}, "labels")

	for _, r := range fake.requests {
		if r.Path == "/api/v1/repos/o/r/issues" && r.Method == "POST" {
			t.Errorf("issue was created despite invalid labels; body: %s", r.Body)
		}
	}
}

// TestLinkIssueLabelMissingIndex: the error must clearly name the missing
// parameter.
func TestLinkIssueLabelMissingIndex(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	res := callTool(t, cs, "link_gitea", map[string]any{
		"type":   "issue_label",
		"owner":  "o",
		"repo":   "r",
		"labels": []any{1},
	})
	if !res.IsError {
		t.Fatalf("expected error for missing index, got: %s", resultText(res))
	}
}

// TestLinkIssueLabelBadIndexType: a non-numeric index must produce an error
// that names the parameter and expected type, not a bare "index is required".
func TestLinkIssueLabelBadIndexType(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	msg := callToolExpectTypeError(t, cs, "link_gitea", map[string]any{
		"type":   "issue_label",
		"owner":  "o",
		"repo":   "r",
		"index":  "no. 326",
		"labels": []any{1},
	}, "index")
	if !strings.Contains(msg, "integer") {
		t.Errorf("error should state the expected type, got: %s", msg)
	}
}

func TestCreateIssueWithLabels(t *testing.T) {
	fake := &fakeForgejo{}
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()
	cs := startSession(t, srv.URL)

	res := callTool(t, cs, "create_gitea", map[string]any{
		"resource": "issue",
		"owner":    "o",
		"repo":     "r",
		"title":    "t",
		"body":     "b",
		"labels":   []any{1, 2, 3},
	})
	if res.IsError {
		t.Fatalf("tool errored: %s", resultText(res))
	}

	for _, r := range fake.requests {
		if r.Path == "/api/v1/repos/o/r/issues" && r.Method == "POST" {
			var body struct {
				Labels []int64 `json:"labels"`
			}
			if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
				t.Fatalf("bad body %q: %v", r.Body, err)
			}
			if len(body.Labels) != 3 {
				t.Errorf("backend got labels %v, want [1 2 3]; raw body: %s", body.Labels, r.Body)
			}
			return
		}
	}
	t.Errorf("no issue creation request reached backend; requests: %+v", fake.requests)
}
