// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// MyIssueMeta represents basic issue information for dependency operations.
// This type is not available in the Forgejo SDK.
type MyIssueMeta struct {
	Index int64  `json:"index"`
	Owner string `json:"owner,omitempty"`
	Name  string `json:"repo,omitempty"`
}

// IssueDependencyList represents a list of issues that block the current issue.
// According to Forgejo API definition, these are issues that must be closed
// before the current issue can be closed.
// Used by list_issue_dependencies endpoint.
type IssueDependencyList []*forgejo.Issue

// ToMarkdown renders issue dependencies with essential information for quick scanning
// Shows: #Index **Title** (state)
// Example per issue:
// #123 **Fix authentication bug** (open)
// #45 **Update user model** (closed)
func (idl IssueDependencyList) ToMarkdown() string {
	if len(idl) == 0 {
		return "*No issue dependencies found*"
	}

	markdown := ""
	for _, issue := range idl {
		if issue == nil {
			continue
		}
		markdown += fmt.Sprintf("#%d **%s** (%s)\n", issue.Index, issue.Title, issue.State)
	}

	return markdown
}
