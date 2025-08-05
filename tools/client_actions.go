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

// MyListActionTasks lists all Forgejo Actions tasks in a repository.
// GET /repos/{owner}/{repo}/actions/tasks
func (c *Client) MyListActionTasks(owner, repo string) (*types.MyActionTaskResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/actions/tasks", owner, repo)

	var result types.MyActionTaskResponse
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
