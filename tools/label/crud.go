// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package label

import (
	"context"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListRepoLabelsParams defines the parameters for the list_repo_labels tool.
// It specifies the owner and repository name to list labels from.
type ListRepoLabelsParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
}

// ListRepoLabelsImpl implements the read-only MCP tool for listing repository labels.
// This operation is safe, idempotent, and does not modify any data. It fetches
// all available labels for a specified repository using the Forgejo SDK.
type ListRepoLabelsImpl struct {
	Client *tools.Client
}

// Definition describes the `list_repo_labels` tool. It requires `owner` and `repo`
// as parameters and is marked as a safe, read-only operation.
func (ListRepoLabelsImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_repo_labels",
		Title:       "List Repository Labels",
		Description: "List all labels available in a repository. Returns label information including name, description, color, and ID.",
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
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing labels. It calls the Forgejo SDK's
// `ListRepoLabels` function and formats the resulting slice of labels into
// a markdown list. Errors will occur if the repository is not found or
// authentication fails.
func (impl ListRepoLabelsImpl) Handler() mcp.ToolHandlerFor[ListRepoLabelsParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[ListRepoLabelsParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Call SDK
		labels, _, err := impl.Client.ListRepoLabels(p.Owner, p.Repo, forgejo.ListLabelsOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list labels: %w", err)
		}

		// Convert to our types and format
		var content string
		if len(labels) == 0 {
			content = "No labels found in this repository."
		} else {
			// Convert labels to our type
			labelList := make(types.LabelList, len(labels))
			for i, label := range labels {
				labelList[i] = &types.Label{Label: label}
			}

			content = fmt.Sprintf("Found %d labels\n\n%s",
				len(labels), labelList.ToMarkdown())
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: content,
				},
			},
		}, nil
	}
}

// CreateLabelParams defines the parameters for the create_label tool.
// It includes the label's name, color, and optional description.
type CreateLabelParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Name is the name of the new label.
	Name string `json:"name"`
	// Color is the hex color code for the label (without the '#').
	Color string `json:"color"`
	// Description is the optional markdown description of the label.
	Description string `json:"description,omitempty"`
}

// CreateLabelImpl implements the MCP tool for creating a new repository label.
// This is a non-idempotent operation that creates a new label using the Forgejo SDK.
type CreateLabelImpl struct {
	Client *tools.Client
}

// Definition describes the `create_label` tool. It requires `owner`, `repo`,
// a `name`, and a `color`. It is not idempotent, as multiple calls with the
// same name will fail once the first label is created.
func (CreateLabelImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_label",
		Title:       "Create Label",
		Description: "Create a new label in a repository. Specify the label name, description, and color.",
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
				"name": {
					Type:        "string",
					Description: "Label name",
				},
				"color": {
					Type:        "string",
					Description: "Label color (hex color code without #, e.g., 'ff0000' for red)",
				},
				"description": {
					Type:        "string",
					Description: "Optional label description",
				},
			},
			Required: []string{"owner", "repo", "name", "color"},
		},
	}
}

// Handler implements the logic for creating a label. It calls the Forgejo SDK's
// `CreateLabel` function and returns the details of the newly created label.
func (impl CreateLabelImpl) Handler() mcp.ToolHandlerFor[CreateLabelParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[CreateLabelParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Build options for SDK call
		opt := forgejo.CreateLabelOption{
			Name:        p.Name,
			Color:       p.Color,
			Description: p.Description,
		}

		// Call SDK
		label, _, err := impl.Client.CreateLabel(p.Owner, p.Repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to create label: %w", err)
		}

		// Convert to our type and format
		labelWrapper := &types.Label{Label: label}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: labelWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// EditLabelParams defines the parameters for the edit_label tool.
// It specifies the label to edit by ID and the fields to update.
type EditLabelParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the label to edit.
	ID int `json:"id"`
	// Name is the new name for the label.
	Name string `json:"name,omitempty"`
	// Color is the new hex color code for the label (without the '#').
	Color string `json:"color,omitempty"`
	// Description is the new optional markdown description for the label.
	Description string `json:"description,omitempty"`
}

// EditLabelImpl implements the MCP tool for editing an existing repository label.
// This is an idempotent operation that modifies a label's metadata using the
// Forgejo SDK.
type EditLabelImpl struct {
	Client *tools.Client
}

// Definition describes the `edit_label` tool. It requires `owner`, `repo`, and the
// label `id`. It is marked as idempotent.
func (EditLabelImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_label",
		Title:       "Edit Label",
		Description: "Edit an existing label's name, description, or color.",
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
				"id": {
					Type:        "integer",
					Description: "Label ID",
				},
				"name": {
					Type:        "string",
					Description: "New label name (optional)",
				},
				"color": {
					Type:        "string",
					Description: "New label color (hex color code without #, e.g., 'ff0000' for red) (optional)",
				},
				"description": {
					Type:        "string",
					Description: "New label description (optional)",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for editing a label. It calls the Forgejo SDK's
// `EditLabel` function. It will return an error if the label ID is not found.
func (impl EditLabelImpl) Handler() mcp.ToolHandlerFor[EditLabelParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[EditLabelParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Build options for SDK call
		opt := forgejo.EditLabelOption{}
		if p.Name != "" {
			opt.Name = &p.Name
		}
		if p.Color != "" {
			opt.Color = &p.Color
		}
		if p.Description != "" {
			opt.Description = &p.Description
		}

		// Call SDK
		label, _, err := impl.Client.EditLabel(p.Owner, p.Repo, int64(p.ID), opt)
		if err != nil {
			return nil, fmt.Errorf("failed to edit label: %w", err)
		}

		// Convert to our type and format
		labelWrapper := &types.Label{Label: label}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: labelWrapper.ToMarkdown(),
				},
			},
		}, nil
	}
}

// DeleteLabelParams defines the parameters for the delete_label tool.
// It specifies the label to be deleted by its ID.
type DeleteLabelParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the label to delete.
	ID int `json:"id"`
}

// DeleteLabelImpl implements the destructive MCP tool for deleting a repository label.
// This is an idempotent but irreversible operation that removes a label from a
// repository using the Forgejo SDK.
type DeleteLabelImpl struct {
	Client *tools.Client
}

// Definition describes the `delete_label` tool. It requires `owner`, `repo`, and
// the label `id`. It is marked as a destructive operation to ensure clients
// can warn the user before execution.
func (DeleteLabelImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_label",
		Title:       "Delete Label",
		Description: "Delete a label from a repository.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(true),
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
				"id": {
					Type:        "integer",
					Description: "Label ID to delete",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for deleting a label. It calls the Forgejo SDK's
// `DeleteLabel` function. On success, it returns a simple text confirmation.
// It will return an error if the label does not exist.
func (impl DeleteLabelImpl) Handler() mcp.ToolHandlerFor[DeleteLabelParams, any] {
	return func(ctx context.Context, req *mcp.ServerRequest[*mcp.CallToolParamsFor[DeleteLabelParams]]) (*mcp.CallToolResult, error) {
		p := req.Params.Arguments

		// Call SDK
		_, err := impl.Client.DeleteLabel(p.Owner, p.Repo, int64(p.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to delete label: %w", err)
		}

		// Return success message
		emptyResponse := types.EmptyResponse{}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: emptyResponse.ToMarkdown(),
				},
			},
		}, nil
	}
}
