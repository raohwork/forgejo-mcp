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

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// LinkImpl implements the link_gitea tool.
type LinkImpl struct {
	Client *tools.Client
}

// Definition describes the link_gitea tool with minimal schema.
func (LinkImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "link_gitea",
		Title: "Link Gitea Resources",
		Description: `Create relationships between resources in Forgejo/Gitea.
Types: issue_label (add labels to issue), issue_dependency (issue depends on another), issue_blocking (issue blocks another).
Use gitea_manual(action="link") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: tools.BoolPtr(false),
			IdempotentHint:  true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"type": {
					Type:        "string",
					Description: "Link type",
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
				"labels": {
					Type:        "array",
					Items:       &jsonschema.Schema{Type: "integer"},
					Description: "Label IDs (for issue_label)",
				},
				"dependency_index": {
					Type:        "integer",
					Description: "Issue number this issue depends on (for issue_dependency)",
				},
				"blocked_index": {
					Type:        "integer",
					Description: "Issue number blocked by this issue (for issue_blocking)",
				},
			},
			Required:             []string{"type", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate link logic based on type.
func (impl LinkImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		linkType, _ := args["type"].(string)
		if linkType == "" {
			return nil, nil, fmt.Errorf("type is required. Valid types: %v", ListLinkTypes())
		}

		if _, ok := LookupManual(ActionLink, linkType); !ok {
			return nil, nil, fmt.Errorf("unknown type '%s'. Valid types: %v", linkType, ListLinkTypes())
		}

		switch linkType {
		case "issue_label":
			return impl.addIssueLabels(args)
		case "issue_dependency":
			return impl.addIssueDependency(args)
		case "issue_blocking":
			return impl.addIssueBlocking(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionLink, linkType, "not implemented"))
		}
	}
}

func (impl LinkImpl) addIssueLabels(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_label", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_label", err.Error()))
	}

	labelIDs, err := int64List(args, "labels")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_label", err.Error()))
	}
	if len(labelIDs) == 0 {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_label", "labels is required (array of label IDs)"))
	}

	opt := forgejo.IssueLabelsOption{Labels: labelIDs}

	labels, _, err := impl.Client.AddIssueLabels(owner, repo, index, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add labels: %w", err)
	}

	labelList := make(types.LabelList, len(labels))
	for i, l := range labels {
		labelList[i] = &types.Label{Label: l}
	}
	return textResult(fmt.Sprintf("Labels added to issue #%d\n\n%s", index, labelList.ToMarkdown())), nil, nil
}

func (impl LinkImpl) addIssueDependency(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_dependency", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_dependency", err.Error()))
	}

	dependencyIndex, err := requiredPositiveInt(args, "dependency_index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_dependency", err.Error()))
	}

	dependency := types.MyIssueMeta{
		Owner: owner,
		Name:  repo,
		Index: dependencyIndex,
	}

	_, err = impl.Client.MyAddIssueDependency(owner, repo, index, dependency)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add dependency: %w", err)
	}

	return textResult(fmt.Sprintf("Issue #%d now depends on issue #%d (must close #%d first)",
		index, dependencyIndex, dependencyIndex)), nil, nil
}

func (impl LinkImpl) addIssueBlocking(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_blocking", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_blocking", err.Error()))
	}

	blockedIndex, err := requiredPositiveInt(args, "blocked_index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionLink, "issue_blocking", err.Error()))
	}

	blocked := types.MyIssueMeta{
		Owner: owner,
		Name:  repo,
		Index: blockedIndex,
	}

	_, err = impl.Client.MyAddIssueBlocking(owner, repo, index, blocked)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add blocking relationship: %w", err)
	}

	return textResult(fmt.Sprintf("Issue #%d now blocks issue #%d (must close #%d first)",
		index, blockedIndex, index)), nil, nil
}
