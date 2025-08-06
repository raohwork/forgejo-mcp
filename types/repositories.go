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

// Repository represents a repository response with embedded SDK repository
// Used by endpoints:
// - GET /repos/search
// - GET /user/repos
// - GET /orgs/{org}/repos
type Repository struct {
	*forgejo.Repository
}

// ToMarkdown renders repository with name, description, stats and key info
// Example: **owner/repo-name** `PRIVATE` `FORK`
// A sample repository for testing purposes
// Stars: 42 | Forks: 7 | Issues: 3 | PRs: 1
// [View Repository](https://git.example.com/owner/repo-name)
func (r *Repository) ToMarkdown() string {
	if r.Repository == nil {
		return "*Invalid repository*"
	}
	markdown := "**" + r.FullName + "**"
	if r.Private {
		markdown += " `PRIVATE`"
	}
	if r.Fork {
		markdown += " `FORK`"
	}
	if r.Template {
		markdown += " `TEMPLATE`"
	}
	markdown += "\n"
	if r.Description != "" {
		markdown += r.Description + "\n"
	}
	markdown += fmt.Sprintf("Stars: %d | Forks: %d | Issues: %d | PRs: %d\n", r.Stars, r.Forks, r.OpenIssues, r.OpenPulls)
	if r.HTMLURL != "" {
		markdown += "[View Repository](" + r.HTMLURL + ")"
	}
	return markdown
}

// RepositoryList represents a list of repositories response
// Used by endpoints:
// - GET /repos/search
// - GET /user/repos
// - GET /orgs/{org}/repos
type RepositoryList []*Repository

// ToMarkdown renders repositories as a numbered list with basic stats
// Example:
// 1. **owner/repo-name** `PRIVATE` `FORK`
// A sample repository for testing purposes
// Stars: 42 | Forks: 7 | Issues: 3 | PRs: 1
// [View Repository](https://git.example.com/owner/repo-name)
// 2. **owner/another-repo**
// Another repository
// Stars: 15 | Forks: 2 | Issues: 0 | PRs: 0
func (rl RepositoryList) ToMarkdown() string {
	if len(rl) == 0 {
		return "*No repositories found*"
	}
	markdown := ""
	for i, repo := range rl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, repo.ToMarkdown())
	}
	return markdown
}
