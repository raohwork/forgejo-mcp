// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package issue

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// AddIssueLabelsParams defines the parameters for the add_issue_labels tool.
// It specifies the issue and the label IDs to be added.
type AddIssueLabelsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Labels is a slice of label IDs to add to the issue.
	Labels []int `json:"labels"`
}

// AddIssueLabelsImpl implements the MCP tool for adding labels to an issue.
// This is an idempotent operation that uses the Forgejo SDK to associate one
// or more existing labels with an issue.
type AddIssueLabelsImpl struct {
	Client *tools.Client
}

// Definition describes the `add_issue_labels` tool. It requires the issue's `index`
// and an array of `labels` (IDs). It is marked as idempotent.
func (AddIssueLabelsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "add_issue_labels",
		Title:       "Add Issue Labels",
		Description: "Add labels to an existing issue.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  true,
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
					Description: "Issue index number",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to add to this issue",
					MinItems:    tools.IntPtr(1),
				},
			},
			Required: []string{"owner", "repo", "index", "labels"},
		},
	}
}

// Handler implements the logic for adding labels to an issue. It calls the
// Forgejo SDK's `AddIssueLabels` function. It will return an error if the issue
// or any of the label IDs are not found.
func (impl AddIssueLabelsImpl) Handler() mcp.ToolHandlerFor[AddIssueLabelsParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[AddIssueLabelsParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Convert int labels to int64
		labelIDs := make([]int64, len(p.Labels))
		for i, label := range p.Labels {
			labelIDs[i] = int64(label)
		}

		opt := forgejo.IssueLabelsOption{
			Labels: labelIDs,
		}

		labels, _, err := impl.Client.AddIssueLabels(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to add labels: %w", err)
		}

		// Convert to our types
		var labelsMarkdown string
		for _, label := range labels {
			labelWrapper := &types.Label{Label: label}
			labelsMarkdown += labelWrapper.ToMarkdown() + "\n"
		}

		content := fmt.Sprintf("Added %d labels to issue #%d\n\n%s", len(labels), p.Index, labelsMarkdown)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}

// RemoveIssueLabelParams defines the parameters for the remove_issue_label tool.
// It specifies the issue and the single label ID to be removed.
type RemoveIssueLabelParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Label is the ID of the label to remove from the issue.
	Label int `json:"label"`
}

// RemoveIssueLabelImpl implements the MCP tool for removing a label from an issue.
// This is an idempotent operation that uses the Forgejo SDK to disassociate a
// label from an issue.
type RemoveIssueLabelImpl struct {
	Client *tools.Client
}

// Definition describes the `remove_issue_label` tool. It requires the issue's
// `index` and a single `label` ID to remove. It is marked as idempotent.
func (RemoveIssueLabelImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "remove_issue_label",
		Title:       "Remove Issue Label",
		Description: "Remove a specific label from an issue.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  true,
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
					Description: "Issue index number",
				},
				"label": {
					Type:        "integer",
					Description: "Label ID to remove from this issue",
				},
			},
			Required: []string{"owner", "repo", "index", "label"},
		},
	}
}

// Handler implements the logic for removing a label from an issue. It calls the
// Forgejo SDK's `DeleteIssueLabel` function. On success, it returns a simple
// text confirmation. It will return an error if the issue or label is not found.
func (impl RemoveIssueLabelImpl) Handler() mcp.ToolHandlerFor[RemoveIssueLabelParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[RemoveIssueLabelParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		_, err := impl.Client.DeleteIssueLabel(p.Owner, p.Repo, int64(p.Index), int64(p.Label))
		if err != nil {
			return nil, fmt.Errorf("failed to remove label: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Label %d successfully removed from issue #%d.", p.Label, p.Index),
				},
			},
		}, nil
	}
}

// ReplaceIssueLabelsParams defines the parameters for the replace_issue_labels tool.
// It specifies the issue and the new set of label IDs.
type ReplaceIssueLabelsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Index is the issue number.
	Index int `json:"index"`
	// Labels is a slice of label IDs that will replace all existing labels on the issue.
	Labels []int `json:"labels"`
}

// ReplaceIssueLabelsImpl implements the MCP tool for replacing all labels on an issue.
// This is an idempotent operation that uses the Forgejo SDK to set the definitive
// list of labels for an issue.
type ReplaceIssueLabelsImpl struct {
	Client *tools.Client
}

// Definition describes the `replace_issue_labels` tool. It requires the issue's
// `index` and an array of `labels` (IDs) to apply. It is marked as idempotent.
func (ReplaceIssueLabelsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "replace_issue_labels",
		Title:       "Replace Issue Labels",
		Description: "Replace all labels on an issue with a new set of labels.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  true,
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
					Description: "Issue index number",
				},
				"labels": {
					Type: "array",
					Items: &jsonschema.Schema{
						Type: "integer",
					},
					Description: "Array of label IDs to set on this issue (replaces all existing labels)",
				},
			},
			Required: []string{"owner", "repo", "index", "labels"},
		},
	}
}

// Handler implements the logic for replacing issue labels. It calls the Forgejo
// SDK's `ReplaceIssueLabels` function. It will return an error if the issue or
// any of the label IDs are not found.
func (impl ReplaceIssueLabelsImpl) Handler() mcp.ToolHandlerFor[ReplaceIssueLabelsParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[ReplaceIssueLabelsParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Convert int labels to int64
		labelIDs := make([]int64, len(p.Labels))
		for i, label := range p.Labels {
			labelIDs[i] = int64(label)
		}

		opt := forgejo.IssueLabelsOption{
			Labels: labelIDs,
		}

		labels, _, err := impl.Client.ReplaceIssueLabels(p.Owner, p.Repo, int64(p.Index), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to replace labels: %w", err)
		}

		// Convert to our types
		var labelsMarkdown string
		for _, label := range labels {
			labelWrapper := &types.Label{Label: label}
			labelsMarkdown += labelWrapper.ToMarkdown() + "\n"
		}

		content := fmt.Sprintf("Replaced labels for issue #%d with %d labels\n\n%s", p.Index, len(labels), labelsMarkdown)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}
