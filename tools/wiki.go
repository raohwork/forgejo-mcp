// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Wiki related tools (Custom HTTP implementation needed)

func ListWikiPages() mcp.Tool {
	return mcp.Tool{
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

func GetWikiPage() mcp.Tool {
	return mcp.Tool{
		Name:        "get_wiki_page",
		Title:       "Get Wiki Page",
		Description: "Get the content and metadata of a specific wiki page.",
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
				"page_name": {
					Type:        "string",
					Description: "Wiki page name",
				},
			},
			Required: []string{"owner", "repo", "page_name"},
		},
	}
}

func CreateWikiPage() mcp.Tool {
	return mcp.Tool{
		Name:        "create_wiki_page",
		Title:       "Create Wiki Page",
		Description: "Create a new wiki page with specified title and content.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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
					Description: "Wiki page title",
				},
				"content": {
					Type:        "string",
					Description: "Wiki page content (markdown supported)",
				},
				"message": {
					Type:        "string",
					Description: "Optional commit message (defaults to 'Create page {title}')",
				},
			},
			Required: []string{"owner", "repo", "title", "content"},
		},
	}
}

func EditWikiPage() mcp.Tool {
	return mcp.Tool{
		Name:        "edit_wiki_page",
		Title:       "Edit Wiki Page",
		Description: "Edit an existing wiki page's title and content.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
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
				"page_name": {
					Type:        "string",
					Description: "Wiki page name to edit",
				},
				"title": {
					Type:        "string",
					Description: "New wiki page title (optional, defaults to current title)",
				},
				"content": {
					Type:        "string",
					Description: "New wiki page content (markdown supported)",
				},
				"message": {
					Type:        "string",
					Description: "Optional commit message (defaults to 'Update page {page_name}')",
				},
			},
			Required: []string{"owner", "repo", "page_name", "content"},
		},
	}
}

func DeleteWikiPage() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_wiki_page",
		Title:       "Delete Wiki Page",
		Description: "Delete a wiki page from the repository.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(true),
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
				"page_name": {
					Type:        "string",
					Description: "Wiki page name to delete",
				},
			},
			Required: []string{"owner", "repo", "page_name"},
		},
	}
}
