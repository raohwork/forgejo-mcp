// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"encoding/base64"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// MyWikiCommit represents wiki page commit/revision information.
type MyWikiCommit struct {
	ID        string              `json:"sha"`
	Author    *forgejo.CommitUser `json:"author"`
	Committer *forgejo.CommitUser `json:"commiter"` // Note: API has typo "commiter"
	Message   string              `json:"message"`
}

// MyWikiPageMetaData represents wiki page meta information.
type MyWikiPageMetaData struct {
	Title      string        `json:"title"`
	HTMLURL    string        `json:"html_url"`
	SubURL     string        `json:"sub_url"`
	LastCommit *MyWikiCommit `json:"last_commit"`
}

// MyWikiPage represents a complete wiki page with content.
type MyWikiPage struct {
	Title         string        `json:"title"`
	HTMLURL       string        `json:"html_url"`
	SubURL        string        `json:"sub_url"`
	LastCommit    *MyWikiCommit `json:"last_commit"`
	ContentBase64 string        `json:"content_base64"`
	CommitCount   int64         `json:"commit_count"`
	Sidebar       string        `json:"sidebar"`
	Footer        string        `json:"footer"`
}

// MyCreateWikiPageOptions represents options for creating a wiki page.
type MyCreateWikiPageOptions struct {
	Title         string `json:"title"`
	ContentBase64 string `json:"content_base64"`
	Message       string `json:"message,omitempty"`
}

// WikiPage represents a wiki page response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/page/{pageName}
// - POST /repos/{owner}/{repo}/wiki/new
// - PATCH /repos/{owner}/{repo}/wiki/page/{pageName}
type WikiPage struct {
	*MyWikiPage
}

// ToMarkdown renders wiki page with title, last modified date and content
// Example: # Getting Started
// *Last modified: 2024-01-15 14:30*
//
// Welcome to our project wiki...
func (w *WikiPage) ToMarkdown() string {
	markdown := "# " + w.Title + "\n"
	if w.LastCommit != nil && w.LastCommit.Author != nil && w.LastCommit.Author.Date != "" {
		if t, err := time.Parse(time.RFC3339, w.LastCommit.Author.Date); err == nil {
			markdown += "*Last modified: " + t.Format("2006-01-02 15:04") + "*\n\n"
		}
	}
	if w.ContentBase64 != "" {
		if content, err := base64.StdEncoding.DecodeString(w.ContentBase64); err == nil {
			markdown += string(content)
		}
	}
	return markdown
}

// WikiPageList represents a list of wiki pages response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/pages
type WikiPageList []*MyWikiPageMetaData

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
		if page.LastCommit != nil && page.LastCommit.Author != nil && page.LastCommit.Author.Date != "" {
			if t, err := time.Parse(time.RFC3339, page.LastCommit.Author.Date); err == nil {
				markdown += " (" + t.Format("2006-01-02") + ")"
			}
		}
		markdown += "\n"
	}
	return markdown
}
