// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"

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

// MyListWikiPages lists all wiki pages in a repository.
// GET /repos/{owner}/{repo}/wiki/pages
func (c *Client) MyListWikiPages(owner, repo string) ([]*MyWikiPageMetaData, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/pages", owner, repo)

	var result []*MyWikiPageMetaData
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MyGetWikiPage gets a single wiki page by name.
// GET /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyGetWikiPage(owner, repo, pageName string) (*MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)

	var result MyWikiPage
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MyCreateWikiPage creates a new wiki page.
// POST /repos/{owner}/{repo}/wiki/new
func (c *Client) MyCreateWikiPage(owner, repo string, options MyCreateWikiPageOptions) (*MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/new", owner, repo)

	var result MyWikiPage
	err := c.sendSimpleRequest("POST", endpoint, options, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MyDeleteWikiPage deletes a wiki page.
// DELETE /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyDeleteWikiPage(owner, repo, pageName string) error {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)

	// DELETE returns 204 No Content on success
	var result interface{}
	err := c.sendSimpleRequest("DELETE", endpoint, nil, &result)
	if err != nil {
		return err
	}

	return nil
}

// MyEditWikiPage edits an existing wiki page.
// PATCH /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyEditWikiPage(owner, repo, pageName string, options MyCreateWikiPageOptions) (*MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)

	var result MyWikiPage
	err := c.sendSimpleRequest("PATCH", endpoint, options, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
