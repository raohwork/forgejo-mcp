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

func TestIssueDependencyList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		dependencies IssueDependencyList
		required     []string
	}{
		{
			name: "multiple issue dependencies",
			dependencies: IssueDependencyList{
				&forgejo.Issue{
					Index: 123,
					Title: "Fix authentication bug",
					State: "open",
				},
				&forgejo.Issue{
					Index: 45,
					Title: "Update user model",
					State: "closed",
				},
			},
			required: []string{"#123", "**Fix authentication bug**", "(open)", "#45", "**Update user model**", "(closed)"},
		},
		{
			name:         "empty dependency list",
			dependencies: IssueDependencyList{},
			required:     []string{"*No issue dependencies found*"},
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
