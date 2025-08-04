// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"
	"time"
)

// MyActionTask represents a Forgejo Actions task.
type MyActionTask struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`
	Event        string    `json:"event"`
	WorkflowID   string    `json:"workflow_id"`
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	RunNumber    int64     `json:"run_number"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
}

// MyActionTaskResponse represents the response for listing action tasks.
type MyActionTaskResponse struct {
	TotalCount   int64           `json:"total_count"`
	WorkflowRuns []*MyActionTask `json:"workflow_runs"`
}

// MyListActionTasks lists all Forgejo Actions tasks in a repository.
// GET /repos/{owner}/{repo}/actions/tasks
func (c *Client) MyListActionTasks(owner, repo string) (*MyActionTaskResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/actions/tasks", owner, repo)

	var result MyActionTaskResponse
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
