// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"testing"
)

func TestMilestone_ToMarkdown(t *testing.T) {
	tests := []struct {
		name      string
		milestone *Milestone
		required  []string
	}{
		{
			name: "complete milestone with all fields",
			milestone: &Milestone{
				Milestone: testMilestone(),
			},
			required: []string{"v1.0.0", "open", "2024-01-15", "10/15", "Major release"},
		},
		{
			name:      "nil milestone",
			milestone: &Milestone{Milestone: nil},
			required:  []string{"Invalid milestone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.milestone.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestMilestoneList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name       string
		milestones MilestoneList
		required   []string
	}{
		{
			name: "multiple milestones with complete information",
			milestones: MilestoneList{
				&Milestone{Milestone: testMilestone()},
				&Milestone{
					Milestone: testMilestone(),
				},
			},
			required: []string{"1.", "v1.0.0", "Major release", "2.", "v1.0.0"},
		},
		{
			name:       "empty milestone list",
			milestones: MilestoneList{},
			required:   []string{"No milestones found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.milestones.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
