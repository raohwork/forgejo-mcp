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

// PullRequest represents a pull request response with embedded SDK pull request
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls (list)
// - GET /repos/{owner}/{repo}/pulls/{index} (get)
type PullRequest struct {
	*forgejo.PullRequest
}

// ToMarkdown renders pull request with title, state, author and branch info
// Example: **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
func (pr *PullRequest) ToMarkdown() string {
	if pr.PullRequest == nil {
		return "*Invalid pull request*"
	}
	markdown := fmt.Sprintf("**#%d %s** (%s)\n", pr.Index, pr.Title, pr.State)
	if pr.Poster != nil {
		markdown += "Author: " + pr.Poster.UserName + "\n"
	}
	if pr.Head != nil && pr.Base != nil {
		markdown += fmt.Sprintf("Branch: %s → %s\n", pr.Head.Name, pr.Base.Name)
	}
	if pr.Body != "" {
		markdown += "\n" + pr.Body
	}
	return markdown
}

// PullRequestList represents a list of pull requests response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls
type PullRequestList []*PullRequest

// ToMarkdown renders pull requests as a numbered list with basic info
// Example:
// 1. **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
// 2. **#41 Fix database connection** (merged)
// Author: alice
// Branch: bugfix/db → main
func (prl PullRequestList) ToMarkdown() string {
	if len(prl) == 0 {
		return "*No pull requests found*"
	}
	markdown := ""
	for i, pr := range prl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, pr.ToMarkdown())
	}
	return markdown
}