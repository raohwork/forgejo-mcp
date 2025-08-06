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
// GET /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyListIssueBlocking(owner, repo string, index int64) ([]*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issues []*forgejo.Issue
	err := c.sendSimpleRequest("GET", endpoint, nil, &issues)
	return issues, err
}

// MyAddIssueBlocking blocks the issue given in the body by the issue in path.
// POST /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyAddIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issue *forgejo.Issue
	err := c.sendSimpleRequest("POST", endpoint, blocked, &issue)
	return issue, err
}

// MyRemoveIssueBlocking unblocks the issue given in the body by the issue in path.
// DELETE /repos/{owner}/{repo}/issues/{index}/blocks
func (c *Client) MyRemoveIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)

	var issue *forgejo.Issue
	err := c.sendSimpleRequest("DELETE", endpoint, blocked, &issue)
	return issue, err
}
