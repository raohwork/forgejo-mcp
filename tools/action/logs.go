// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// Log tails are capped so that a multi-megabyte CI log cannot blow up the
// context window of the calling LLM.
const (
	defaultLogTailLines = 500
	maxLogTailLines     = 5000
)

// GetActionJobLogsParams defines the parameters for the get_action_job_logs tool.
type GetActionJobLogsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// JobID is the identifier of the workflow job (see list_action_run_jobs).
	JobID int64 `json:"job_id"`
	// Attempt is the 1-based attempt number; 0 means the latest attempt.
	Attempt int64 `json:"attempt,omitempty"`
	// Tail limits the output to the last N lines of the log.
	Tail int `json:"tail,omitempty"`
}

// GetActionJobLogsImpl implements the read-only MCP tool for downloading the
// plaintext logs of a Forgejo Actions job. This is a safe, idempotent operation.
// Note: This feature is not supported by the official Forgejo SDK and requires a
// custom HTTP implementation.
type GetActionJobLogsImpl struct {
	Client *tools.Client
}

// Definition describes the `get_action_job_logs` tool. It requires `owner`,
// `repo` and `job_id`, and optionally accepts `attempt` and `tail`. It is
// marked as a safe, read-only operation.
func (GetActionJobLogsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_action_job_logs",
		Title:       "Get Action Job Logs",
		Description: "Download the plaintext logs of a Forgejo Actions job. Get job IDs from list_action_run_jobs. Returns the last 500 lines by default; use `tail` to adjust. Requires a Forgejo version that provides the job logs API.",
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
				"job_id": {
					Type:        "integer",
					Description: "ID of the workflow job",
					Minimum:     tools.Float64Ptr(1),
				},
				"attempt": {
					Type:        "integer",
					Description: "1-based attempt number as reported by list_action_run_jobs (optional, defaults to the latest attempt)",
					Minimum:     tools.Float64Ptr(1),
				},
				"tail": {
					Type:        "integer",
					Description: "Return only the last N lines of the log (optional, defaults to 500, max 5000)",
					Minimum:     tools.Float64Ptr(1),
					Maximum:     tools.Float64Ptr(maxLogTailLines),
				},
			},
			Required: []string{"owner", "repo", "job_id"},
		},
	}
}

// Handler implements the logic for downloading job logs. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/actions/jobs/{job_id}/logs`
// endpoint and returns the (optionally tailed) plaintext log.
func (impl GetActionJobLogsImpl) Handler() mcp.ToolHandlerFor[GetActionJobLogsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args GetActionJobLogsParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Call custom client method
		logText, err := impl.Client.MyGetActionJobLogs(p.Owner, p.Repo, p.JobID, p.Attempt)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get job logs (the job may not exist, its log may have expired, or this Forgejo version may not provide the job logs API yet): %w", err)
		}

		tail := p.Tail
		if tail <= 0 {
			tail = defaultLogTailLines
		}
		if tail > maxLogTailLines {
			tail = maxLogTailLines
		}

		content := formatJobLogs(logText, p.JobID, tail)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil, nil
	}
}

// formatJobLogs renders the log text with a short header, keeping only the
// last `tail` lines. The header states whether the log was truncated so the
// caller knows to re-request with a larger tail if needed.
func formatJobLogs(logText string, jobID int64, tail int) string {
	logText = strings.TrimRight(logText, "\n")
	if logText == "" {
		return fmt.Sprintf("Logs of job %d are empty.", jobID)
	}

	lines := strings.Split(logText, "\n")
	total := len(lines)
	header := fmt.Sprintf("Logs of job %d (%d lines)", jobID, total)
	if total > tail {
		lines = lines[total-tail:]
		header = fmt.Sprintf("Logs of job %d (last %d of %d lines, increase `tail` for more)", jobID, tail, total)
	}

	return header + "\n\n" + strings.Join(lines, "\n")
}
