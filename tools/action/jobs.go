// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package action

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListActionRunJobsParams defines the parameters for the list_action_run_jobs tool.
type ListActionRunJobsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// RunID is the identifier of the workflow run.
	RunID int64 `json:"run_id"`
}

// ListActionRunJobsImpl implements the read-only MCP tool for listing jobs of a
// Forgejo Actions workflow run. This is a safe, idempotent operation. Note: This
// feature is not supported by the official Forgejo SDK and requires a custom HTTP
// implementation.
type ListActionRunJobsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_action_run_jobs` tool. It requires `owner`,
// `repo` and `run_id`. It is marked as a safe, read-only operation.
func (ListActionRunJobsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_action_run_jobs",
		Title:       "List Action Run Jobs",
		Description: "List the jobs of a Forgejo Actions workflow run. Use the returned job IDs with get_action_job_logs to read job logs. Requires a Forgejo version that provides the run jobs API.",
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
				"run_id": {
					Type:        "integer",
					Description: "ID of the workflow run",
					Minimum:     tools.Float64Ptr(1),
				},
			},
			Required: []string{"owner", "repo", "run_id"},
		},
	}
}

// Handler implements the logic for listing workflow run jobs. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/actions/runs/{run_id}/jobs`
// endpoint and formats the results as markdown.
func (impl ListActionRunJobsImpl) Handler() mcp.ToolHandlerFor[ListActionRunJobsParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListActionRunJobsParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Call custom client method
		jobs, err := impl.Client.MyListActionRunJobs(p.Owner, p.Repo, p.RunID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list run jobs (the run may not exist, or this Forgejo version may not provide the run jobs API yet): %w", err)
		}

		var content string
		if len(jobs) == 0 {
			content = "No jobs found for this workflow run."
		} else {
			jobList := types.ActionRunJobList(jobs)
			content = fmt.Sprintf("Found %d jobs in run %d\n\n%s",
				len(jobs), p.RunID, jobList.ToMarkdown())
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil, nil
	}
}
