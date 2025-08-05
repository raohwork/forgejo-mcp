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

func TestLabel_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		label    *Label
		required []string
	}{
		{
			name: "complete label with all fields",
			label: &Label{
				Label: testLabel(),
			},
			required: []string{"bug", "ff0000", "Something isn't working"},
		},
		{
			name:     "nil label",
			label:    &Label{Label: nil},
			required: []string{"Invalid label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.label.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestLabelList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		labels   LabelList
		required []string
	}{
		{
			name: "multiple labels with complete information",
			labels: LabelList{
				&Label{Label: testLabel()},
				&Label{
					Label: &forgejo.Label{
						Name:        "enhancement",
						Color:       "a2eeef",
						Description: "New feature or request",
					},
				},
			},
			required: []string{"bug", "ff0000", "Something isn't working", "enhancement", "a2eeef", "New feature or request"},
		},
		{
			name:     "empty label list",
			labels:   LabelList{},
			required: []string{"No labels found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.labels.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
