// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

func TestIssue_ToMarkdown(t *testing.T) {
	deadline := testTime()
	tests := []struct {
		name     string
		issue    *Issue
		required []string
	}{
		{
			name: "complete issue with all fields",
			issue: &Issue{
				Issue: &forgejo.Issue{
					Index:     123,
					Title:     "Fix login bug",
					Body:      "The login page crashes when using special characters",
					State:     "open",
					Poster:    testUser(),
					Assignees: []*forgejo.User{testUser()},
					Labels: []*forgejo.Label{
						{Name: "bug"},
						{Name: "priority-high"},
					},
					Milestone: testMilestone(),
					Deadline:  &deadline,
				},
			},
			required: []string{"#123", "Fix login bug", "open", "testuser", "[testuser]", "[bug priority-high]", "v1.0.0", "2024-01-15", "The login page crashes"},
		},
		{
			name:     "nil issue",
			issue:    &Issue{Issue: nil},
			required: []string{"Invalid issue"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.issue.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestComment_ToMarkdown(t *testing.T) {
	created := testTime()
	tests := []struct {
		name     string
		comment  *Comment
		required []string
	}{
		{
			name: "complete comment with all fields",
			comment: &Comment{
				Comment: &forgejo.Comment{
					Poster:  testUser(),
					Body:    "I think the issue is in the authentication module",
					Created: created,
				},
			},
			required: []string{"testuser", "2024-01-15 14:30", "I think the issue is in the authentication module"},
		},
		{
			name:     "nil comment",
			comment:  &Comment{Comment: nil},
			required: []string{"Invalid comment"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.comment.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
