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

// ActionTask represents an action task response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - GET /repos/{owner}/{repo}/actions/tasks
type ActionTask struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	WorkflowID  int64            `json:"workflow_id"`
	RunID       int64            `json:"run_id"`
	JobID       int64            `json:"job_id"`
	Steps       []ActionTaskStep `json:"steps,omitempty"`
}

// ToMarkdown renders action task with name, status and execution time
// Example: **Build and Test** `success`
// Run ID: 123 | Job ID: 456
// Started: 2024-01-15 14:30 | Duration: 2m15s
// Steps:
//   - Setup Go `success`
//   - Run Tests `success`
//   - Build Binary `success`
func (at *ActionTask) ToMarkdown() string {
	markdown := fmt.Sprintf("**%s** `%s`\n", at.Name, at.Status)
	markdown += fmt.Sprintf("Run ID: %d | Job ID: %d\n", at.RunID, at.JobID)
	if at.StartedAt != nil {
		markdown += "Started: " + at.StartedAt.Format("2006-01-02 15:04")
		if at.CompletedAt != nil {
			duration := at.CompletedAt.Sub(*at.StartedAt)
			markdown += " | Duration: " + duration.String()
		}
		markdown += "\n"
	}
	if len(at.Steps) > 0 {
		markdown += "Steps:\n"
		for _, step := range at.Steps {
			markdown += fmt.Sprintf("  - %s `%s`\n", step.Name, step.Status)
		}
	}
	return markdown
}

// ActionTaskStep represents a step in an action task
type ActionTaskStep struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Conclusion  string     `json:"conclusion,omitempty"`
}

// ToMarkdown renders action task step with name, status and timing
// Example: **Run Tests** `success` (passed) - 1m30s
func (ats *ActionTaskStep) ToMarkdown() string {
	markdown := fmt.Sprintf("**%s** `%s`", ats.Name, ats.Status)
	if ats.Conclusion != "" && ats.Conclusion != ats.Status {
		markdown += " (" + ats.Conclusion + ")"
	}
	if ats.StartedAt != nil && ats.CompletedAt != nil {
		duration := ats.CompletedAt.Sub(*ats.StartedAt)
		markdown += " - " + duration.String()
	}
	return markdown
}

// ActionTaskList represents a list of action tasks response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/actions/tasks
type ActionTaskList []*ActionTask

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
	if len(atl) == 0 {
		return "*No action tasks found*"
	}
	markdown := ""
	for i, task := range atl {
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

// MyActionTaskResponse represents the response for listing action tasks.
type MyActionTaskResponse struct {
	TotalCount   int64           `json:"total_count"`
	WorkflowRuns []*MyActionTask `json:"workflow_runs"`
}
