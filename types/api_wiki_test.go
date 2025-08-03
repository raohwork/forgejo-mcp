// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"testing"
)

func TestWikiPage_ToMarkdown(t *testing.T) {
	modified := testTime()
	tests := []struct {
		name     string
		wiki     *WikiPage
		required []string
	}{
		{
			name: "complete wiki page with all fields",
			wiki: &WikiPage{
				Title:        "Getting Started",
				Content:      "Welcome to our project wiki. This guide will help you get started with the project.",
				LastModified: modified,
			},
			required: []string{"# Getting Started", "Last modified: 2024-01-15 14:30", "Welcome to our project wiki"},
		},
		{
			name: "wiki page without content",
			wiki: &WikiPage{
				Title: "Empty Page",
			},
			required: []string{"# Empty Page"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.wiki.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestWikiPageList_ToMarkdown(t *testing.T) {
	modified := testTime()
	tests := []struct {
		name     string
		wikis    WikiPageList
		required []string
	}{
		{
			name: "multiple wiki pages with complete information",
			wikis: WikiPageList{
				&WikiPage{
					Title:        "Getting Started",
					LastModified: modified,
				},
				&WikiPage{
					Title:        "API Documentation",
					LastModified: modified,
				},
				&WikiPage{
					Title: "Contributing Guide",
				},
			},
			required: []string{"## Wiki Pages", "Getting Started", "API Documentation", "Contributing Guide", "2024-01-15"},
		},
		{
			name:     "empty wiki page list",
			wikis:    WikiPageList{},
			required: []string{"No wiki pages found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.wikis.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
