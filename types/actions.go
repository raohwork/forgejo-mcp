// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"
	"time"
)


// ActionTaskList represents a list of action tasks response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/actions/tasks
type ActionTaskList struct {
	*MyActionTaskResponse
}

// ToMarkdown renders action tasks as a numbered list with status
// Example:
// 1. **Build and Test** `success`
// Run ID: 123 | Job ID: 456
// Started: 2024-01-15 14:30 | Duration: 2m15s
// Steps:
//   - Setup Go `success`
//   - Run Tests `success`
//
// 2. **Deploy** `running`
// Run ID: 124 | Job ID: 457
// Started: 2024-01-15 14:35
func (atl ActionTaskList) ToMarkdown() string {
	if atl.MyActionTaskResponse == nil || len(atl.WorkflowRuns) == 0 {
		return "*No action tasks found*"
	}
	markdown := ""
	for i, task := range atl.WorkflowRuns {
		markdown += fmt.Sprintf("%d. %s\n", i+1, task.ToMarkdown())
	}
	return markdown
}

// MyActionTask represents a Forgejo Actions task.
type MyActionTask struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"` // title of last commit
	Status       string    `json:"status"`
	Event        string    `json:"event"`       //push, pull_request, etc.
	WorkflowID   string    `json:"workflow_id"` // filename of yaml
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	RunNumber    int64     `json:"run_number"` // run#N
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
}

// ToMarkdown renders action task with name, status, execution info and timing
func (at *MyActionTask) ToMarkdown() string {
	markdown := fmt.Sprintf("**%s** `%s` - Run #%d", at.DisplayTitle, at.Status, at.RunNumber)
	
	// Add timing information for statistical analysis
	if !at.CreatedAt.IsZero() {
		markdown += fmt.Sprintf(" | Created: %s", at.CreatedAt.Format("2006-01-02 15:04"))
	}
	
	// Add duration if both start and update times are available
	if !at.RunStartedAt.IsZero() && !at.UpdatedAt.IsZero() {
		duration := at.UpdatedAt.Sub(at.RunStartedAt)
		markdown += fmt.Sprintf(" | Duration: %s", duration.String())
	}
	
	return markdown
}

// MyActionTaskResponse represents the response for listing action tasks.
type MyActionTaskResponse struct {
	TotalCount   int64           `json:"total_count"`
	WorkflowRuns []*MyActionTask `json:"workflow_runs"`
}
