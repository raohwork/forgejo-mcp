// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package wiki

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// ListWikiPagesParams defines the parameters for the list_wiki_pages tool.
// It specifies the owner and repository name to list wiki pages from.
type ListWikiPagesParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
}

// ListWikiPagesImpl implements the read-only MCP tool for listing repository wiki pages.
// This operation is safe, idempotent, and does not modify any data. It fetches
// all available wiki pages for a specified repository. Note: This feature is not
// supported by the official Forgejo SDK and requires a custom HTTP implementation.
type ListWikiPagesImpl struct {
	Client *tools.Client
}

// Definition describes the `list_wiki_pages` tool. It requires `owner` and `repo`
// as parameters and is marked as a safe, read-only operation.
func (ListWikiPagesImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_wiki_pages",
		Title:       "List Wiki Pages",
		Description: "List all wiki pages in a repository. Returns page information including titles and metadata.",
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

// Handler implements the logic for listing wiki pages. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/wiki/pages` endpoint and
// formats the resulting list of pages into a markdown table. Errors will occur
// if the repository is not found or authentication fails.
func (impl ListWikiPagesImpl) Handler() mcp.ToolHandlerFor[ListWikiPagesParams, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, args ListWikiPagesParams) (*mcp.CallToolResult, any, error) {
		p := args

		// Call custom client method
		pages, err := impl.Client.MyListWikiPages(p.Owner, p.Repo)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list wiki pages: %w", err)
		}

		// Convert to our types and format
		pageList := types.WikiPageList(pages)

		var content string
		if len(pages) == 0 {
			content = "No wiki pages found in this repository."
		} else {
			content = fmt.Sprintf("Found %d wiki pages\n\n%s",
				len(pages), pageList.ToMarkdown())
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
