// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// UnlinkImpl implements the unlink_gitea tool.
type UnlinkImpl struct {
	Client *tools.Client
}

// Definition describes the unlink_gitea tool with minimal schema.
func (UnlinkImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "unlink_gitea",
		Title: "Unlink Gitea Resources",
		Description: `Remove relationships between resources in Forgejo/Gitea.
Types: issue_label (remove label from issue), issue_dependency (remove dependency), issue_blocking (remove blocking).
Use gitea_manual(action="unlink") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(true),
			IdempotentHint:  true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"type": {
					Type:        "string",
					Description: "Link type to remove",
					Enum:        []any{"issue_label", "issue_dependency", "issue_blocking"},
				},
				"owner": {
					Type:        "string",
					Description: "Repository owner",
				},
				"repo": {
					Type:        "string",
					Description: "Repository name",
				},
				"index": {
					Type:        "integer",
					Description: "Issue number the link applies to",
				},
				"label_id": {
					Type:        "integer",
					Description: "Label ID to remove (for issue_label)",
				},
				"dependency_index": {
					Type:        "integer",
					Description: "Issue number to remove from dependencies (for issue_dependency)",
				},
				"blocked_index": {
					Type:        "integer",
					Description: "Blocked issue number to remove (for issue_blocking)",
				},
			},
			Required:             []string{"type", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate unlink logic based on type.
func (impl UnlinkImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		linkType, _ := args["type"].(string)
		if linkType == "" {
			return nil, nil, fmt.Errorf("type is required. Valid types: %v", ListLinkTypes())
		}

		if _, ok := LookupManual(ActionUnlink, linkType); !ok {
			return nil, nil, fmt.Errorf("unknown type '%s'. Valid types: %v", linkType, ListLinkTypes())
		}

		switch linkType {
		case "issue_label":
			return impl.removeIssueLabel(args)
		case "issue_dependency":
			return impl.removeIssueDependency(args)
		case "issue_blocking":
			return impl.removeIssueBlocking(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionUnlink, linkType, "not implemented"))
		}
	}
}

func (impl UnlinkImpl) removeIssueLabel(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_label", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_label", err.Error()))
	}

	labelID, err := requiredPositiveInt(args, "label_id")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_label", err.Error()))
	}

	_, err = impl.Client.DeleteIssueLabel(owner, repo, index, labelID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to remove label: %w", err)
	}

	return textResult(fmt.Sprintf("Label %d removed from issue #%d", labelID, index)), nil, nil
}

func (impl UnlinkImpl) removeIssueDependency(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_dependency", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_dependency", err.Error()))
	}

	dependencyIndex, err := requiredPositiveInt(args, "dependency_index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_dependency", err.Error()))
	}

	dependency := types.MyIssueMeta{
		Owner: owner,
		Name:  repo,
		Index: dependencyIndex,
	}

	_, err = impl.Client.MyRemoveIssueDependency(owner, repo, index, dependency)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to remove dependency: %w", err)
	}

	return textResult(fmt.Sprintf("Issue #%d no longer depends on issue #%d",
		index, dependencyIndex)), nil, nil
}

func (impl UnlinkImpl) removeIssueBlocking(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_blocking", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_blocking", err.Error()))
	}

	blockedIndex, err := requiredPositiveInt(args, "blocked_index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionUnlink, "issue_blocking", err.Error()))
	}

	blocked := types.MyIssueMeta{
		Owner: owner,
		Name:  repo,
		Index: blockedIndex,
	}

	_, err = impl.Client.MyRemoveIssueBlocking(owner, repo, index, blocked)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to remove blocking relationship: %w", err)
	}

	return textResult(fmt.Sprintf("Issue #%d no longer blocks issue #%d",
		index, blockedIndex)), nil, nil
}
