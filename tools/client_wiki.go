// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"

	"github.com/raohwork/forgejo-mcp/types"
)

// MyListWikiPages lists all wiki pages in a repository.
// GET /repos/{owner}/{repo}/wiki/pages
func (c *Client) MyListWikiPages(owner, repo string) ([]*types.MyWikiPageMetaData, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/pages", owner, repo)

	var result []*types.MyWikiPageMetaData
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MyGetWikiPage gets a single wiki page by name.
// GET /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyGetWikiPage(owner, repo, pageName string) (*types.MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)

	var result types.MyWikiPage
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MyCreateWikiPage creates a new wiki page.
// POST /repos/{owner}/{repo}/wiki/new
func (c *Client) MyCreateWikiPage(owner, repo string, options types.MyCreateWikiPageOptions) (*types.MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/new", owner, repo)

	var result types.MyWikiPage
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
func (c *Client) MyEditWikiPage(owner, repo, pageName string, options types.MyCreateWikiPageOptions) (*types.MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)

	var result types.MyWikiPage
	err := c.sendSimpleRequest("PATCH", endpoint, options, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
