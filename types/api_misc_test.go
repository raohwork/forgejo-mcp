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

func TestIssueDependency_ToMarkdown(t *testing.T) {
	tests := []struct {
		name       string
		dependency *IssueDependency
		required   []string
	}{
		{
			name: "complete issue dependency with both issues",
			dependency: &IssueDependency{
				IssueID:      123,
				DependencyID: 45,
				Issue: &Issue{
					Issue: &forgejo.Issue{
						Index: 123,
						Title: "Fix login bug",
					},
				},
				Dependency: &Issue{
					Issue: &forgejo.Issue{
						Index: 45,
						Title: "Update authentication library",
					},
				},
			},
			required: []string{"#123", "Fix login bug", "depends on", "#45", "Update authentication library"},
		},
		{
			name: "issue dependency with only IDs",
			dependency: &IssueDependency{
				IssueID:      123,
				DependencyID: 45,
			},
			required: []string{"Issue #123", "depends on", "Issue #45"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.dependency.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestIssueDependencyList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		dependencies IssueDependencyList
		required     []string
	}{
		{
			name: "multiple issue dependencies",
			dependencies: IssueDependencyList{
				&IssueDependency{
					IssueID:      123,
					DependencyID: 45,
				},
				&IssueDependency{
					IssueID:      124,
					DependencyID: 46,
				},
			},
			required: []string{"Issue #123", "depends on", "Issue #45", "Issue #124", "Issue #46"},
		},
		{
			name:         "empty dependency list",
			dependencies: IssueDependencyList{},
			required:     []string{"No issue dependencies found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.dependencies.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestEmptyResponse_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		response EmptyResponse
		required []string
	}{
		{
			name:     "empty response",
			response: EmptyResponse{},
			required: []string{"Operation completed successfully"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.response.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
