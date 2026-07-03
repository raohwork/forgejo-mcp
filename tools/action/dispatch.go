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
)

// DispatchWorkflowParams defines the parameters for the dispatch_workflow tool.
type DispatchWorkflowParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Workflow is the workflow file name (e.g. "ci.yml").
	Workflow string `json:"workflow"`
	// Ref is the git reference (branch, tag or SHA) to run the workflow against.
	Ref string `json:"ref"`
	// Inputs are the workflow_dispatch input keys and values defined in the workflow file.
	Inputs map[string]string `json:"inputs,omitempty"`
}

// DispatchWorkflowImpl implements the MCP tool for triggering a workflow via the
// workflow_dispatch event. This operation creates a new workflow run and is not
// idempotent. Note: This feature is not supported by the official Forgejo SDK and
// requires a custom HTTP implementation.
type DispatchWorkflowImpl struct {
	Client *tools.Client
}

// Definition describes the `dispatch_workflow` tool. It requires `owner`, `repo`,
// `workflow` and `ref`, and optionally accepts `inputs`.
func (DispatchWorkflowImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "dispatch_workflow",
		Title:       "Dispatch Workflow",
		Description: "Trigger a Forgejo Actions workflow via the workflow_dispatch event. The workflow file must declare a workflow_dispatch trigger. Returns the created run's ID and number.",
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
				"workflow": {
					Type:        "string",
					Description: "Workflow file name, e.g. 'ci.yml' (as stored under .forgejo/workflows or .github/workflows)",
				},
				"ref": {
					Type:        "string",
					Description: "Git reference to run against: a branch name, tag, or commit SHA",
				},
				"inputs": {
					Type:        "object",
					Description: "workflow_dispatch inputs as string key/value pairs, matching the inputs declared in the workflow file (optional)",
					AdditionalProperties: &jsonschema.Schema{
						Type: "string",
					},
				},
			},
			Required: []string{"owner", "repo", "workflow", "ref"},
		},
	}
}

// Handler implements the logic for dispatching a workflow. It performs a custom
// HTTP POST to the `/repos/{owner}/{repo}/actions/workflows/{workflow}/dispatches`
// endpoint and reports the created run.
func (impl DispatchWorkflowImpl) Handler() mcp.ToolHandlerFor[DispatchWorkflowParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args DispatchWorkflowParams) (*mcp.CallToolResult, any, error) {
		p := args

		run, err := impl.Client.MyDispatchWorkflow(p.Owner, p.Repo, p.Workflow, p.Ref, p.Inputs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to dispatch workflow (check that the workflow file exists, declares a workflow_dispatch trigger, and that the ref is valid): %w", err)
		}

		content := fmt.Sprintf("Workflow `%s` dispatched on `%s`.\n\n%s",
			p.Workflow, p.Ref, run.ToMarkdown())

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil, nil
	}
}
