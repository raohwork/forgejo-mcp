// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package wiki

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/types"
)

// GetWikiPageParams defines the parameters for the get_wiki_page tool.
// It specifies the owner, repository, and page name to retrieve.
type GetWikiPageParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// PageName is the title of the wiki page to retrieve.
	PageName string `json:"page_name"`
}

// GetWikiPageImpl implements the read-only MCP tool for fetching a single wiki page.
// This operation is safe, idempotent, and does not modify any data. Note: This
// feature is not supported by the official Forgejo SDK and requires a custom
// HTTP implementation.
type GetWikiPageImpl struct {
	Client *tools.Client
}

// Definition describes the `get_wiki_page` tool. It requires `owner`, `repo`,
// and `page_name` as parameters and is marked as a safe, read-only operation.
func (GetWikiPageImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
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

// Handler implements the logic for fetching a wiki page. It performs a custom
// HTTP GET request to the `/repos/{owner}/{repo}/wiki/page/{pageName}` endpoint
// and formats the resulting page content as markdown. Errors will occur if the
// page or repository is not found.
func (impl GetWikiPageImpl) Handler() mcp.ToolHandlerFor[GetWikiPageParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[GetWikiPageParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call custom client method
		page, err := impl.Client.MyGetWikiPage(p.Owner, p.Repo, p.PageName)
		if err != nil {
			return nil, fmt.Errorf("failed to get wiki page: %w", err)
		}

		// Decode base64 content
		content, err := base64.StdEncoding.DecodeString(page.ContentBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode page content: %w", err)
		}

		// Convert to our type and format
		wikiPage := &types.WikiPage{
			Title:          page.Title,
			Content:        string(content),
			HTMLContentURL: page.HTMLURL,
			SubURL:         page.SubURL,
		}
		// Set last modified time if available
		if page.LastCommit != nil && page.LastCommit.Author != nil && page.LastCommit.Author.Date != "" {
			// Try to parse the date string
			if parsedTime, err := time.Parse(time.RFC3339, page.LastCommit.Author.Date); err == nil {
				wikiPage.LastModified = parsedTime
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: wikiPage.ToMarkdown(),
				},
			},
		}, nil
	}
}

// CreateWikiPageParams defines the parameters for the create_wiki_page tool.
// It includes the title and content for the new wiki page.
type CreateWikiPageParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// Title is the title for the new wiki page.
	Title string `json:"title"`
	// Content is the markdown content of the new wiki page.
	Content string `json:"content"`
	// Message is an optional commit message for the creation.
	Message string `json:"message,omitempty"`
}

// CreateWikiPageImpl implements the MCP tool for creating a new wiki page.
// This is a non-idempotent operation that adds a new page to the repository's
// wiki. Note: This feature is not supported by the official Forgejo SDK and
// requires a custom HTTP implementation.
type CreateWikiPageImpl struct {
	Client *tools.Client
}

// Definition describes the `create_wiki_page` tool. It requires `owner`, `repo`,
// `title`, and `content`. It is not idempotent as multiple calls with the same
// parameters will result in multiple pages if the title is not unique.
func (CreateWikiPageImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "create_wiki_page",
		Title:       "Create Wiki Page",
		Description: "Create a new wiki page with specified title and content.",
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

// Handler implements the logic for creating a wiki page. It performs a custom
// HTTP POST request to the `/repos/{owner}/{repo}/wiki/new` endpoint. On success,
// it returns information about the newly created page.
func (impl CreateWikiPageImpl) Handler() mcp.ToolHandlerFor[CreateWikiPageParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateWikiPageParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Prepare options for API call
		options := tools.MyCreateWikiPageOptions{
			Title:         p.Title,
			ContentBase64: base64.StdEncoding.EncodeToString([]byte(p.Content)),
			Message:       p.Message,
		}

		// Call custom client method
		page, err := impl.Client.MyCreateWikiPage(p.Owner, p.Repo, options)
		if err != nil {
			return nil, fmt.Errorf("failed to create wiki page: %w", err)
		}

		// Convert to our type and format
		wikiPage := &types.WikiPage{
			Title:          page.Title,
			Content:        p.Content, // Use original content
			HTMLContentURL: page.HTMLURL,
			SubURL:         page.SubURL,
			LastModified:   time.Now(), // Set to current time for new page
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: wikiPage.ToMarkdown(),
				},
			},
		}, nil
	}
}

// EditWikiPageParams defines the parameters for the edit_wiki_page tool.
// It specifies the page to edit and the new content.
type EditWikiPageParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// PageName is the current title of the wiki page to edit.
	PageName string `json:"page_name"`
	// Title is the optional new title for the wiki page.
	Title string `json:"title,omitempty"`
	// Content is the new markdown content for the wiki page.
	Content string `json:"content"`
	// Message is an optional commit message for the update.
	Message string `json:"message,omitempty"`
}

// EditWikiPageImpl implements the MCP tool for editing an existing wiki page.
// This is an idempotent operation. Note: This feature is not supported by the
// official Forgejo SDK and requires a custom HTTP implementation.
type EditWikiPageImpl struct {
	Client *tools.Client
}

// Definition describes the `edit_wiki_page` tool. It requires `owner`, `repo`,
// `page_name`, and new `content`. It is marked as idempotent.
func (EditWikiPageImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "edit_wiki_page",
		Title:       "Edit Wiki Page",
		Description: "Edit an existing wiki page's title and content.",
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

// Handler implements the logic for editing a wiki page. It performs a custom
// HTTP PATCH request to the `/repos/{owner}/{repo}/wiki/page/{pageName}` endpoint.
// It returns an error if the page is not found.
func (impl EditWikiPageImpl) Handler() mcp.ToolHandlerFor[EditWikiPageParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[EditWikiPageParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Prepare options for API call
		title := p.Title
		if title == "" {
			title = p.PageName // Use current name if no new title
		}
		options := tools.MyCreateWikiPageOptions{
			Title:         title,
			ContentBase64: base64.StdEncoding.EncodeToString([]byte(p.Content)),
			Message:       p.Message,
		}

		// Call custom client method
		page, err := impl.Client.MyEditWikiPage(p.Owner, p.Repo, p.PageName, options)
		if err != nil {
			return nil, fmt.Errorf("failed to edit wiki page: %w", err)
		}

		// Convert to our type and format
		wikiPage := &types.WikiPage{
			Title:          page.Title,
			Content:        p.Content, // Use updated content
			HTMLContentURL: page.HTMLURL,
			SubURL:         page.SubURL,
			LastModified:   time.Now(), // Set to current time for updated page
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: wikiPage.ToMarkdown(),
				},
			},
		}, nil
	}
}

// DeleteWikiPageParams defines the parameters for the delete_wiki_page tool.
// It specifies the page to be deleted.
type DeleteWikiPageParams struct {
	// Owner is the username or organization name that owns the repository.
	Owner string `json:"owner"`
	// Repo is the name of the repository.
	Repo string `json:"repo"`
	// PageName is the title of the wiki page to delete.
	PageName string `json:"page_name"`
}

// DeleteWikiPageImpl implements the destructive MCP tool for deleting a wiki page.
// This is an idempotent but irreversible operation. Note: This feature is not
// supported by the official Forgejo SDK and requires a custom HTTP implementation.
type DeleteWikiPageImpl struct {
	Client *tools.Client
}

// Definition describes the `delete_wiki_page` tool. It requires `owner`, `repo`,
// and `page_name`. It is marked as a destructive operation to ensure clients
// can warn the user before execution.
func (DeleteWikiPageImpl) Definition() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_wiki_page",
		Title:       "Delete Wiki Page",
		Description: "Delete a wiki page from the repository.",
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
				"page_name": {
					Type:        "string",
					Description: "Wiki page name to delete",
				},
			},
			Required: []string{"owner", "repo", "page_name"},
		},
	}
}

// Handler implements the logic for deleting a wiki page. It performs a custom
// HTTP DELETE request to the `/repos/{owner}/{repo}/wiki/page/{pageName}` endpoint.
// On success, it returns a simple text confirmation.
func (impl DeleteWikiPageImpl) Handler() mcp.ToolHandlerFor[DeleteWikiPageParams, any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[DeleteWikiPageParams]) (*mcp.CallToolResult, error) {
		p := params.Arguments

		// Call custom client method
		err := impl.Client.MyDeleteWikiPage(p.Owner, p.Repo, p.PageName)
		if err != nil {
			return nil, fmt.Errorf("failed to delete wiki page: %w", err)
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
