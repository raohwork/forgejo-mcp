// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package milestone

import (
	"context"
	"errors"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// ListRepoMilestonesParams defines the parameters for the list_repo_milestones tool.
// It specifies the repository and an optional state filter.
type ListRepoMilestonesParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// State filters milestones by their state (e.g., 'open', 'closed').
	State string `json:"state,omitempty"`
}

// ListRepoMilestonesImpl implements the read-only MCP tool for listing repository milestones.
// This is a safe, idempotent operation that uses the Forgejo SDK to fetch a list
// of milestones, optionally filtered by their state.
type ListRepoMilestonesImpl struct{}

// Definition describes the `list_repo_milestones` tool. It requires `owner` and `repo`,
// supports an optional `state` filter, and is marked as a safe, read-only operation.
func (ListRepoMilestonesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_repo_milestones",
		Title:       "List Repository Milestones",
		Description: "List all milestones in a repository. Returns milestone information including title, description, due date, and status.",
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
				"state": {
					Type:        "string",
					Description: "Milestone state filter: 'open', 'closed', or 'all' (optional, defaults to 'open')",
					Enum:        []any{"open", "closed", "all"},
				},
			},
			Required: []string{"owner", "repo"},
		},
	}
}

// Handler implements the logic for listing milestones. It calls the Forgejo SDK's
// `ListRepoMilestones` function and formats the results into a markdown list.
func (ListRepoMilestonesImpl) Handler() mcp.ToolHandlerFor[ListRepoMilestonesParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[ListRepoMilestonesParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// CreateMilestoneParams defines the parameters for the create_milestone tool.
// It includes the title and optional description and due date.
type CreateMilestoneParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Title is the title of the new milestone.
	Title string `json:"title"`
	// Description is the markdown description of the milestone.
	Description string `json:"description,omitempty"`
	// DueDate is the optional due date for the milestone.
	DueDate time.Time `json:"due_date,omitempty"`
}

// CreateMilestoneImpl implements the MCP tool for creating a new milestone.
// This is a non-idempotent operation that creates a new milestone object
// using the Forgejo SDK.
type CreateMilestoneImpl struct{}

// Definition describes the `create_milestone` tool. It requires `owner`, `repo`,
// and a `title`. It is not idempotent, as multiple calls with the same title
// will create multiple milestones.
func (CreateMilestoneImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_milestone",
		Title:       "Create Milestone",
		Description: "Create a new milestone in a repository. Specify the title, description, and optional due date.",
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
				"title": {
					Type:        "string",
					Description: "Milestone title",
				},
				"description": {
					Type:        "string",
					Description: "Milestone description (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "Milestone due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
			},
			Required: []string{"owner", "repo", "title"},
		},
	}
}

// Handler implements the logic for creating a milestone. It calls the Forgejo SDK's
// `CreateMilestone` function and returns the details of the newly created milestone.
func (CreateMilestoneImpl) Handler() mcp.ToolHandlerFor[CreateMilestoneParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateMilestoneParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// EditMilestoneParams defines the parameters for the edit_milestone tool.
// It specifies the milestone to edit by ID and the fields to update.
type EditMilestoneParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the milestone to edit.
	ID int `json:"id"`
	// Title is the new title for the milestone.
	Title string `json:"title,omitempty"`
	// Description is the new markdown description for the milestone.
	Description string `json:"description,omitempty"`
	// DueDate is the new optional due date for the milestone.
	DueDate time.Time `json:"due_date,omitempty"`
	// State is the new state for the milestone (e.g., 'open', 'closed').
	State string `json:"state,omitempty"`
}

// EditMilestoneImpl implements the MCP tool for editing an existing milestone.
// This is an idempotent operation that modifies a milestone's metadata using
// the Forgejo SDK.
type EditMilestoneImpl struct{}

// Definition describes the `edit_milestone` tool. It requires `owner`, `repo`, and
// the milestone `id`. It is marked as idempotent.
func (EditMilestoneImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_milestone",
		Title:       "Edit Milestone",
		Description: "Edit an existing milestone's title, description, due date, or state.",
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
					Description: "Milestone ID",
				},
				"title": {
					Type:        "string",
					Description: "New milestone title (optional)",
				},
				"description": {
					Type:        "string",
					Description: "New milestone description (optional)",
				},
				"due_date": {
					Type:        "string",
					Description: "New milestone due date in ISO 8601 format (e.g., '2024-12-31T23:59:59Z') (optional)",
					Format:      "date-time",
				},
				"state": {
					Type:        "string",
					Description: "New milestone state: 'open' or 'closed' (optional)",
					Enum:        []any{"open", "closed"},
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for editing a milestone. It calls the Forgejo SDK's
// `EditMilestone` function. It will return an error if the milestone ID is not found.
func (EditMilestoneImpl) Handler() mcp.ToolHandlerFor[EditMilestoneParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditMilestoneParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}

// DeleteMilestoneParams defines the parameters for the delete_milestone tool.
// It specifies the milestone to be deleted by its ID.
type DeleteMilestoneParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// ID is the unique identifier of the milestone to delete.
	ID int `json:"id"`
}

// DeleteMilestoneImpl implements the destructive MCP tool for deleting a milestone.
// This is an idempotent but irreversible operation that removes a milestone from
// a repository using the Forgejo SDK.
type DeleteMilestoneImpl struct{}

// Definition describes the `delete_milestone` tool. It requires `owner`, `repo`,
// and the milestone `id`. It is marked as a destructive operation to ensure
// clients can warn the user before execution.
func (DeleteMilestoneImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_milestone",
		Title:       "Delete Milestone",
		Description: "Delete a milestone from a repository.",
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
					Description: "Milestone ID to delete",
				},
			},
			Required: []string{"owner", "repo", "id"},
		},
	}
}

// Handler implements the logic for deleting a milestone. It calls the Forgejo SDK's
// `DeleteMilestone` function. On success, it returns a simple text confirmation.
// It will return an error if the milestone does not exist.
func (DeleteMilestoneImpl) Handler() mcp.ToolHandlerFor[DeleteMilestoneParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteMilestoneParams]) (*mcp.CallToolResult, error) {
		// TODO: Implement handler logic
		return nil, errors.New("not implemented yet")
	}
}
