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

// GetImpl implements the get_gitea tool.
type GetImpl struct {
	Client *tools.Client
}

// Definition describes the get_gitea tool with minimal schema.
func (GetImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:  "get_gitea",
		Title: "Get Gitea Resource",
		Description: `Get details of a single resource from Forgejo/Gitea.
Resources: issue, wiki_page, pull_request, repository.
Use gitea_manual(action="get") for details.`,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resource": {
					Type:        "string",
					Description: "Resource type to get",
					Enum:        []any{"issue", "wiki_page", "pull_request", "repository"},
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
					Description: "Issue or pull request number (for issue, pull_request)",
				},
			},
			Required:             []string{"resource", "owner", "repo"},
			AdditionalProperties: &jsonschema.Schema{},
		},
	}
}

// Handler dispatches to the appropriate get logic based on resource type.
func (impl GetImpl) Handler() mcp.ToolHandlerFor[map[string]any, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args map[string]any) (*mcp.CallToolResult, any, error) {
		resource, _ := args["resource"].(string)
		if resource == "" {
			resources := ListResourcesForAction(ActionGet)
			return nil, nil, fmt.Errorf("resource is required. Valid resources: %v", resources)
		}

		if _, ok := LookupManual(ActionGet, resource); !ok {
			resources := ListResourcesForAction(ActionGet)
			return nil, nil, fmt.Errorf("unknown resource '%s'. Valid resources: %v", resource, resources)
		}

		switch resource {
		case "issue":
			return impl.getIssue(args)
		case "wiki_page":
			return impl.getWikiPage(args)
		case "pull_request":
			return impl.getPullRequest(args)
		case "repository":
			return impl.getRepository(args)
		default:
			return nil, nil, errors.New(FormatValidationError(ActionGet, resource, "not implemented"))
		}
	}
}

func (impl GetImpl) getIssue(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "issue", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "issue", err.Error()))
	}

	issue, _, err := impl.Client.GetIssue(owner, repo, index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return textResult((&types.Issue{Issue: issue}).ToMarkdown()), nil, nil
}

func (impl GetImpl) getWikiPage(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "wiki_page", err.Error()))
	}

	pageName, _ := args["page_name"].(string)
	if pageName == "" {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "wiki_page", "page_name is required"))
	}

	page, err := impl.Client.MyGetWikiPage(owner, repo, pageName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get wiki page: %w", err)
	}

	return textResult((&types.WikiPage{MyWikiPage: page}).ToMarkdown()), nil, nil
}

func (impl GetImpl) getPullRequest(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "pull_request", err.Error()))
	}

	index, err := requiredPositiveInt(args, "index")
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "pull_request", err.Error()))
	}

	pr, _, err := impl.Client.GetPullRequest(owner, repo, index)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return textResult((&types.PullRequest{PullRequest: pr}).ToMarkdown()), nil, nil
}

func (impl GetImpl) getRepository(args map[string]any) (*mcp.CallToolResult, any, error) {
	owner, repo, err := extractOwnerRepo(args)
	if err != nil {
		return nil, nil, errors.New(FormatValidationError(ActionGet, "repository", err.Error()))
	}

	repository, _, err := impl.Client.GetRepo(owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return textResult((&types.Repository{Repository: repository}).ToMarkdown()), nil, nil
}
