// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raohwork/forgejo-mcp/tools"
)

// RegisterAll registers all unified tools with the MCP server.
// This replaces the 47 individual tool registrations with 8 consolidated tools:
// - gitea_manual: On-demand documentation lookup
// - create_gitea: Create resources
// - get_gitea: Get single resources
// - list_gitea: List resources
// - edit_gitea: Edit resources
// - delete_gitea: Delete resources
// - link_gitea: Create relationships
// - unlink_gitea: Remove relationships
func RegisterAll(s *mcp.Server, cl *tools.Client) {
	tools.Register(s, &ManualImpl{Client: cl})
	tools.Register(s, &CreateImpl{Client: cl})
	tools.Register(s, &GetImpl{Client: cl})
	tools.Register(s, &ListImpl{Client: cl})
	tools.Register(s, &EditImpl{Client: cl})
	tools.Register(s, &DeleteImpl{Client: cl})
	tools.Register(s, &LinkImpl{Client: cl})
	tools.Register(s, &UnlinkImpl{Client: cl})
}
