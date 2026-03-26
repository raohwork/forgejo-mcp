// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tools

import "testing"

func TestWikiPageNameToSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "flat page no conversion needed",
			input:    "Home",
			expected: "Home",
		},
		{
			name:     "flat page with spaces",
			input:    "Getting Started",
			expected: "Getting-Started",
		},
		{
			name:     "single level nested page",
			input:    "architecture/overview",
			expected: "architecture%2Foverview.-",
		},
		{
			name:     "deep nested page",
			input:    "architecture/mflow/MFLOW-glossary",
			expected: "architecture%2Fmflow%2FMFLOW-glossary.-",
		},
		{
			name:     "nested page with spaces",
			input:    "getting started/quick reference",
			expected: "getting-started%2Fquick-reference.-",
		},
		{
			name:     "already encoded sub_url passthrough",
			input:    "architecture%2Foverview.-",
			expected: "architecture%2Foverview.-",
		},
		{
			name:     "already encoded lowercase passthrough",
			input:    "architecture%2foverview.-",
			expected: "architecture%2foverview.-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wikiPageNameToSlug(tt.input)
			if got != tt.expected {
				t.Errorf("wikiPageNameToSlug(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
