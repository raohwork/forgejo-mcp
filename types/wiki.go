// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"time"
)

// WikiPage represents a wiki page response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/page/{pageName}
// - POST /repos/{owner}/{repo}/wiki/new
// - PATCH /repos/{owner}/{repo}/wiki/page/{pageName}
type WikiPage struct {
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CommitMessage  string    `json:"commit_message,omitempty"`
	LastCommitSHA  string    `json:"last_commit_sha,omitempty"`
	LastModified   time.Time `json:"last_modified,omitempty"`
	HTMLContentURL string    `json:"html_content_url,omitempty"`
	SubURL         string    `json:"sub_url,omitempty"`
}

// ToMarkdown renders wiki page with title, last modified date and content
// Example: # Getting Started
// *Last modified: 2024-01-15 14:30*
//
// Welcome to our project wiki...
func (w *WikiPage) ToMarkdown() string {
	markdown := "# " + w.Title + "\n"
	if !w.LastModified.IsZero() {
		markdown += "*Last modified: " + w.LastModified.Format("2006-01-02 15:04") + "*\n\n"
	}
	if w.Content != "" {
		markdown += w.Content
	}
	return markdown
}

// WikiPageList represents a list of wiki pages response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/pages
type WikiPageList []*WikiPage

// ToMarkdown renders wiki pages as a table of contents with links
// Example: ## Wiki Pages
// - **Getting Started** (2024-01-15)
// - **API Documentation** (2024-01-10)
// - **Contributing Guide** (2024-01-05)
func (wl WikiPageList) ToMarkdown() string {
	if len(wl) == 0 {
		return "*No wiki pages found*"
	}
	markdown := "## Wiki Pages\n"
	for _, page := range wl {
		markdown += "- **" + page.Title + "**"
		if !page.LastModified.IsZero() {
			markdown += " (" + page.LastModified.Format("2006-01-02") + ")"
		}
		markdown += "\n"
	}
	return markdown
}
