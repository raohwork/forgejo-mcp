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

func TestActionTask_ToMarkdown(t *testing.T) {
	started := testTime()
	completed := started.Add(2 * time.Minute)
	tests := []struct {
		name     string
		task     *ActionTask
		required []string
	}{
		{
			name: "complete action task with all fields",
			task: &ActionTask{
				Name:        "Build and Test",
				Status:      "success",
				RunID:       123,
				JobID:       456,
				StartedAt:   &started,
				CompletedAt: &completed,
				Steps: []ActionTaskStep{
					{
						Name:   "Setup Go",
						Status: "success",
					},
					{
						Name:   "Run Tests",
						Status: "success",
					},
				},
			},
			required: []string{"Build and Test", "success", "Run ID: 123", "Job ID: 456", "Started: 2024-01-15 14:30", "Duration: 2m0s", "Setup Go", "Run Tests"},
		},
		{
			name: "minimal action task",
			task: &ActionTask{
				Name:   "Simple Task",
				Status: "running",
				RunID:  789,
				JobID:  101,
			},
			required: []string{"Simple Task", "running", "Run ID: 789", "Job ID: 101"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.task.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestActionTaskStep_ToMarkdown(t *testing.T) {
	started := testTime()
	completed := started.Add(30 * time.Second)
	tests := []struct {
		name     string
		step     *ActionTaskStep
		required []string
	}{
		{
			name: "complete action task step with all fields",
			step: &ActionTaskStep{
				Name:        "Run Tests",
				Status:      "success",
				Conclusion:  "passed",
				StartedAt:   &started,
				CompletedAt: &completed,
			},
			required: []string{"Run Tests", "success", "passed", "30s"},
		},
		{
			name: "minimal action task step",
			step: &ActionTaskStep{
				Name:   "Simple Step",
				Status: "running",
			},
			required: []string{"Simple Step", "running"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.step.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestActionTaskList_ToMarkdown(t *testing.T) {
	started := testTime()
	tests := []struct {
		name     string
		tasks    ActionTaskList
		required []string
	}{
		{
			name: "multiple action tasks with complete information",
			tasks: ActionTaskList{
				&ActionTask{
					Name:      "Build and Test",
					Status:    "success",
					RunID:     123,
					JobID:     456,
					StartedAt: &started,
				},
				&ActionTask{
					Name:   "Deploy",
					Status: "running",
					RunID:  124,
					JobID:  457,
				},
			},
			required: []string{"1.", "Build and Test", "success", "2.", "Deploy", "running"},
		},
		{
			name:     "empty action task list",
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
