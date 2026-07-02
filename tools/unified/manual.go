// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
)

// ManualParams defines the parameters for the gitea_manual tool.
type ManualParams struct {
	// Action is the operation type to get documentation for.
	// If omitted, lists all available actions.
	Action string `json:"action,omitempty"`
	// Resource is the resource type to get documentation for.
	// Required when action is specified (except for link/unlink which use 'type').
	Resource string `json:"resource,omitempty"`
	// Type is the link type for link/unlink actions.
	Type string `json:"type,omitempty"`
}

// ManualImpl implements the gitea_manual tool for on-demand documentation lookup.
type ManualImpl struct {
	Client *tools.Client // Not used but kept for interface consistency
}

// Definition describes the gitea_manual tool.
func (ManualImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "gitea_manual",
		Title: "Gitea Documentation",
		Description: `Look up documentation for Gitea operations.
Call without arguments to see all available actions.
Call with just 'action' to see resources for that action.
Call with 'action' and 'resource' (or 'type' for link/unlink) for full documentation.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"action": {
					Type:        "string",
					Description: "Action to look up: create, get, list, edit, delete, link, unlink",
					Enum:        []any{"create", "get", "list", "edit", "delete", "link", "unlink"},
				},
				"resource": {
					Type:        "string",
					Description: "Resource type (for create/get/list/edit/delete actions)",
				},
				"type": {
					Type:        "string",
					Description: "Link type (for link/unlink actions): issue_label, issue_dependency, issue_blocking",
					Enum:        []any{"issue_label", "issue_dependency", "issue_blocking"},
				},
			},
		},
	}
}

// Handler implements the documentation lookup logic.
func (impl ManualImpl) Handler() mcp.ToolHandlerFor[ManualParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ManualParams) (*mcp.CallToolResult, any, error) {
		var content string

		if args.Action == "" {
			// List all available actions
			content = formatOverview()
		} else if args.Action == "link" || args.Action == "unlink" {
			if args.Type == "" {
				// List all link types for this action
				content = formatLinkTypes(Action(args.Action))
			} else {
				// Show specific link type documentation
				entry, ok := LookupManual(Action(args.Action), args.Type)
				if !ok {
					return nil, nil, fmt.Errorf("unknown link type '%s' for action '%s'. Valid types: %v",
						args.Type, args.Action, ListLinkTypes())
				}
				content = FormatManualEntry(entry)
			}
		} else {
			if args.Resource == "" {
				// List all resources for this action
				content = formatResourcesForAction(Action(args.Action))
			} else {
				// Show specific resource documentation
				entry, ok := LookupManual(Action(args.Action), args.Resource)
				if !ok {
					resources := ListResourcesForAction(Action(args.Action))
					return nil, nil, fmt.Errorf("unknown resource '%s' for action '%s'. Valid resources: %v",
						args.Resource, args.Action, resources)
				}
				content = FormatManualEntry(entry)
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	}
}

// formatOverview returns a summary of all available actions.
func formatOverview() string {
	var sb strings.Builder
	sb.WriteString("# Gitea MCP Tools\n\n")
	sb.WriteString("This server provides unified tools for Forgejo/Gitea operations.\n\n")
	sb.WriteString("## Available Actions\n\n")
	sb.WriteString("| Action | Description |\n")
	sb.WriteString("|--------|-------------|\n")
	sb.WriteString("| `create_gitea` | Create resources (issues, labels, milestones, etc.) |\n")
	sb.WriteString("| `get_gitea` | Get a single resource by ID/name |\n")
	sb.WriteString("| `list_gitea` | List resources with filtering |\n")
	sb.WriteString("| `edit_gitea` | Edit existing resources |\n")
	sb.WriteString("| `delete_gitea` | Delete resources |\n")
	sb.WriteString("| `link_gitea` | Create relationships (labels to issues, dependencies) |\n")
	sb.WriteString("| `unlink_gitea` | Remove relationships |\n")
	sb.WriteString("\n")
	sb.WriteString("## How to Use\n\n")
	sb.WriteString("Each tool takes a `resource` parameter (or `type` for link/unlink) to specify what you're operating on.\n\n")
	sb.WriteString("**Examples:**\n")
	sb.WriteString("```\n")
	sb.WriteString("create_gitea(resource=\"issue\", owner=\"org\", repo=\"project\", title=\"Bug\", body=\"...\")\n")
	sb.WriteString("list_gitea(resource=\"label\", owner=\"org\", repo=\"project\")\n")
	sb.WriteString("link_gitea(type=\"issue_label\", owner=\"org\", repo=\"project\", index=42, labels=[1,2])\n")
	sb.WriteString("```\n\n")
	sb.WriteString("Call `gitea_manual(action=\"create\")` to see resources for a specific action.\n")
	sb.WriteString("Call `gitea_manual(action=\"create\", resource=\"issue\")` for full documentation.\n")

	return sb.String()
}

// formatResourcesForAction lists all resources available for a given action.
func formatResourcesForAction(action Action) string {
	resources := ListResourcesForAction(action)
	sort.Strings(resources)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s_gitea Resources\n\n", action))
	sb.WriteString("| Resource | Description |\n")
	sb.WriteString("|----------|-------------|\n")

	for _, res := range resources {
		entry, _ := LookupManual(action, res)
		// Truncate description for table
		desc := entry.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		sb.WriteString(fmt.Sprintf("| `%s` | %s |\n", res, desc))
	}

	sb.WriteString(fmt.Sprintf("\nCall `gitea_manual(action=\"%s\", resource=\"<resource>\")` for full documentation.\n", action))
	return sb.String()
}

// formatLinkTypes lists all link types for link/unlink actions.
func formatLinkTypes(action Action) string {
	types := ListLinkTypes()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s_gitea Link Types\n\n", action))
	sb.WriteString("| Type | Description |\n")
	sb.WriteString("|------|-------------|\n")

	for _, t := range types {
		entry, _ := LookupManual(action, t)
		desc := entry.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		sb.WriteString(fmt.Sprintf("| `%s` | %s |\n", t, desc))
	}

	sb.WriteString(fmt.Sprintf("\nCall `gitea_manual(action=\"%s\", type=\"<type>\")` for full documentation.\n", action))
	return sb.String()
}
