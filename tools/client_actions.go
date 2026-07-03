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

// MyListActionRunJobs lists the jobs of a Forgejo Actions workflow run.
// GET /repos/{owner}/{repo}/actions/runs/{run_id}/jobs
func (c *Client) MyListActionRunJobs(owner, repo string, runID int64) ([]*types.MyActionRunJob, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID)

	var result []*types.MyActionRunJob
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MyGetActionJobLogs downloads the plaintext logs of a Forgejo Actions job.
// Pass attempt > 0 to fetch a specific historical attempt, 0 for the latest.
// GET /repos/{owner}/{repo}/actions/jobs/{job_id}/logs
func (c *Client) MyGetActionJobLogs(owner, repo string, jobID, attempt int64) (string, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobID)
	if attempt > 0 {
		endpoint += fmt.Sprintf("?attempt=%d", attempt)
	}

	return c.sendTextRequest("GET", endpoint)
}

// MyDispatchWorkflow triggers a workflow_dispatch run of the given workflow file.
// It always requests run info so the caller can learn the resulting run number
// and job names; servers that predate that field simply return an empty result.
// POST /repos/{owner}/{repo}/actions/workflows/{workflowfilename}/dispatches
func (c *Client) MyDispatchWorkflow(owner, repo, workflow, ref string, inputs map[string]string) (*types.MyDispatchWorkflowRun, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflow)

	payload := types.MyDispatchWorkflowOption{
		Ref:           ref,
		Inputs:        inputs,
		ReturnRunInfo: true,
	}

	var result types.MyDispatchWorkflowRun
	if err := c.sendJSONRequestAllowEmpty("POST", endpoint, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
