// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"

	"github.com/raohwork/forgejo-mcp/types"
)

// MyAddIssueDependency adds a dependency to an issue.
// Creates a relationship where the current issue (index) depends on another issue (dependency).
// This means the dependency issue must be closed before the current issue can be closed.
// POST /repos/{owner}/{repo}/issues/{index}/dependencies
func (c *Client) MyAddIssueDependency(owner, repo string, index int64, dependency types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/dependencies", owner, repo, index)

	var result forgejo.Issue
	err := c.sendSimpleRequest("POST", endpoint, dependency, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MyListIssueDependencies lists all dependencies of an issue.
// Returns issues that must be closed before the current issue can be closed.
// GET /repos/{owner}/{repo}/issues/{index}/dependencies
func (c *Client) MyListIssueDependencies(owner, repo string, index int64) ([]*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/dependencies", owner, repo, index)

	var result []*forgejo.Issue
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MyRemoveIssueDependency removes a dependency from an issue.
// Removes the relationship where the current issue depends on another issue.
// DELETE /repos/{owner}/{repo}/issues/{index}/dependencies
func (c *Client) MyRemoveIssueDependency(owner, repo string, index int64, dependency types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/dependencies", owner, repo, index)

	var result forgejo.Issue
	err := c.sendSimpleRequest("DELETE", endpoint, dependency, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// MyListIssueBlocking lists all issues blocked by this issue.
// Returns issues that cannot be closed until the current issue is closed.
// GET /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyListIssueBlocking(owner, repo string, index int64) ([]*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issues []*forgejo.Issue
	err := c.sendSimpleRequest("GET", endpoint, nil, &issues)
	return issues, err
}

// MyAddIssueBlocking blocks the issue given in the body by the issue in path.
// Creates a relationship where the current issue (index) blocks another issue (blocked).
// This means the current issue must be closed before the blocked issue can be closed.
// POST /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyAddIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issue *forgejo.Issue
	err := c.sendSimpleRequest("POST", endpoint, blocked, &issue)
	return issue, err
}

// MyRemoveIssueBlocking unblocks the issue given in the body by the issue in path.
// Removes the relationship where the current issue blocks another issue.
// DELETE /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyRemoveIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issue *forgejo.Issue
	err := c.sendSimpleRequest("DELETE", endpoint, blocked, &issue)
	return issue, err
}
