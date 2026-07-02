// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

// Package unified provides a consolidated MCP toolset for Forgejo/Gitea operations.
//
// This package implements an action-based tool architecture that reduces token
// consumption by ~80% compared to resource-based tools. Instead of separate tools
// for each operation on each resource type (e.g., create_issue, create_label,
// create_milestone), this package provides unified action tools:
//
//   - create_gitea: Create resources (issues, labels, milestones, releases, etc.)
//   - get_gitea: Get single resources by ID/name
//   - list_gitea: List resources with filtering
//   - edit_gitea: Edit existing resources
//   - delete_gitea: Delete resources
//   - link_gitea: Create relationships (issue↔label, issue↔issue dependencies)
//   - unlink_gitea: Remove relationships
//   - gitea_manual: On-demand documentation lookup
//
// Design Philosophy:
//
// The toolset relies on LLM implicit knowledge about git forge concepts (issues,
// labels, milestones, etc.) to keep tool definitions minimal. Detailed documentation
// is provided through:
//
//  1. The gitea_manual tool for proactive lookup
//  2. Rich error messages when validation fails
//
// This "lazy-loaded documentation" approach means tokens are only spent on
// documentation the LLM actually needs.
package unified
