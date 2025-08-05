// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"testing"
	"time"
)

func TestMyActionTask_ToMarkdown(t *testing.T) {
	createdTime := testTime()
	startedTime := createdTime.Add(1 * time.Minute)
	updatedTime := startedTime.Add(5 * time.Minute)
	
	tests := []struct {
		name     string
		task     *MyActionTask
		required []string
	}{
		{
			name: "complete action task with all fields",
			task: &MyActionTask{
				DisplayTitle: "Add new feature",
				Status:       "success",
				RunNumber:    123,
				WorkflowID:   "ci.yml",
				HeadBranch:   "main",
				Event:        "push",
				CreatedAt:    createdTime,
				RunStartedAt: startedTime,
				UpdatedAt:    updatedTime,
			},
			required: []string{"Add new feature", "success", "Run #123", "Created: 2024-01-15 14:30", "Duration: 5m0s"},
		},
		{
			name: "failed action task with minimal timing",
			task: &MyActionTask{
				DisplayTitle: "Fix bug in authentication",
				Status:       "failure",
				RunNumber:    456,
				WorkflowID:   "test.yml",
				HeadBranch:   "feature/auth-fix",
				Event:        "pull_request",
				CreatedAt:    createdTime,
			},
			required: []string{"Fix bug in authentication", "failure", "Run #456", "Created: 2024-01-15 14:30"},
		},
		{
			name: "running task without timing",
			task: &MyActionTask{
				DisplayTitle: "Deploy to staging",
				Status:       "running",
				RunNumber:    789,
			},
			required: []string{"Deploy to staging", "running", "Run #789"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.task.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}


func TestActionTaskList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		tasks    ActionTaskList
		required []string
	}{
		{
			name: "multiple action tasks with complete information",
			tasks: ActionTaskList{
				MyActionTaskResponse: &MyActionTaskResponse{
					TotalCount: 2,
					WorkflowRuns: []*MyActionTask{
						{
							DisplayTitle: "Add new feature",
							Status:       "success",
							RunNumber:    123,
							WorkflowID:   "ci.yml",
						},
						{
							DisplayTitle: "Fix authentication bug",
							Status:       "failure",
							RunNumber:    124,
							WorkflowID:   "test.yml",
						},
					},
				},
			},
			required: []string{"1.", "Add new feature", "success", "Run #123", "2.", "Fix authentication bug", "failure", "Run #124"},
		},
		{
			name: "empty action task list",
			tasks: ActionTaskList{
				MyActionTaskResponse: &MyActionTaskResponse{
					TotalCount:   0,
					WorkflowRuns: []*MyActionTask{},
				},
			},
			required: []string{"No action tasks found"},
		},
		{
			name:     "nil response",
			tasks:    ActionTaskList{},
			required: []string{"No action tasks found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.tasks.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
