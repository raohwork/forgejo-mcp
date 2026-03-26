// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/raohwork/forgejo-mcp/types"
)

// wikiPageNameToSlug converts a wiki page title to the URL slug format
// expected by the Gitea API. Gitea's wiki API endpoints require the
// "sub_url" slug rather than the display title. The slug format:
//   - Spaces are replaced with hyphens
//   - Forward slashes are URL-encoded as %2F
//   - Pages with slashes get a ".-" suffix appended
//   - If the input already looks like a sub_url (contains %2F), it is
//     returned as-is to avoid double-encoding
func wikiPageNameToSlug(pageName string) string {
	// If it already contains %2F, assume it's already a sub_url slug
	if strings.Contains(pageName, "%2F") || strings.Contains(pageName, "%2f") {
		return pageName
	}

	// Replace spaces with hyphens (Gitea wiki convention)
	slug := strings.ReplaceAll(pageName, " ", "-")

	// If no slashes, return as-is (flat page like "Home")
	if !strings.Contains(slug, "/") {
		return url.PathEscape(slug)
	}

	// URL-encode each path segment and join with %2F
	parts := strings.Split(slug, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	result := strings.Join(parts, "%2F")

	// Gitea appends ".-" to slugs for pages with path separators
	if !strings.HasSuffix(result, ".-") {
		result += ".-"
	}

	return result
}

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
// The pageName can be either the display title (e.g. "architecture/overview")
// or the sub_url slug — it will be converted automatically.
// GET /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyGetWikiPage(owner, repo, pageName string) (*types.MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, wikiPageNameToSlug(pageName))

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
// The pageName can be either the display title or the sub_url slug.
// DELETE /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyDeleteWikiPage(owner, repo, pageName string) error {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, wikiPageNameToSlug(pageName))

	// DELETE returns 204 No Content on success
	var result interface{}
	err := c.sendSimpleRequest("DELETE", endpoint, nil, &result)
	if err != nil {
		return err
	}

	return nil
}

// MyEditWikiPage edits an existing wiki page.
// The pageName can be either the display title or the sub_url slug.
// PATCH /repos/{owner}/{repo}/wiki/page/{pageName}
func (c *Client) MyEditWikiPage(owner, repo, pageName string, options types.MyCreateWikiPageOptions) (*types.MyWikiPage, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, wikiPageNameToSlug(pageName))

	var result types.MyWikiPage
	err := c.sendSimpleRequest("PATCH", endpoint, options, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
