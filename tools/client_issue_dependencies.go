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

// MyIssueMeta represents basic issue information for dependency operations.
// This type is not available in the Forgejo SDK.
type MyIssueMeta struct {
	Index int64  `json:"index"`
	Owner string `json:"owner,omitempty"`
	Name  string `json:"repo,omitempty"`
}

// MyAddIssueDependency adds a dependency to an issue.
// POST /repos/{owner}/{repo}/issues/{index}/dependencies
func (c *Client) MyAddIssueDependency(owner, repo string, index int64, dependency MyIssueMeta) (*forgejo.Issue, error) {
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
func (c *Client) MyRemoveIssueDependency(owner, repo string, index int64, dependency MyIssueMeta) (*forgejo.Issue, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/dependencies", owner, repo, index)

	var result forgejo.Issue
	err := c.sendSimpleRequest("DELETE", endpoint, dependency, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
