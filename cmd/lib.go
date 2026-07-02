// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package cmd

import (
	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/tools/unified"
	"github.com/raohwork/forgejo-mcp/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCommands(s *mcp.Server, cl *tools.Client) {
	// Use unified tools (8 tools instead of 47)
	// This reduces token consumption by ~80% while maintaining full functionality.
	// The unified tools use action-based organization:
	// - gitea_manual: On-demand documentation lookup
	// - create_gitea: Create resources (issue, label, milestone, release, wiki_page, pull_request, issue_comment)
	// - get_gitea: Get single resources (issue, wiki_page, pull_request, repository)
	// - list_gitea: List resources (issues, labels, milestones, releases, wiki_pages, pull_requests, repositories, etc.)
	// - edit_gitea: Edit resources (issue, label, milestone, release, wiki_page, etc.)
	// - delete_gitea: Delete resources (label, milestone, release, wiki_page, issue_comment, etc.)
	// - link_gitea: Create relationships (issue↔label, issue dependencies, issue blocking)
	// - unlink_gitea: Remove relationships
	unified.RegisterAll(s, cl)
}

func createServer(cl *tools.Client) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Title:   "Forgejo MCP Server",
		Version: types.VERSION[1:], // strip leading 'v'
	}, &mcp.ServerOptions{
		PageSize:     50,
		Instructions: "An MCP server to interact with repositories on a Forgejo/Gitea instance.",
	})
	registerCommands(server, cl)

	return server
}
