// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Milestone represents a milestone response with embedded SDK milestone
// Used by endpoints:
// - GET /repos/{owner}/{repo}/milestones (list)
// - POST /repos/{owner}/{repo}/milestones (create)
// - PATCH /repos/{owner}/{repo}/milestones/{id} (edit)
type Milestone struct {
	*forgejo.Milestone
}

// ToMarkdown renders milestone with title, state, due date and progress
// Example: **v1.0.0** (open) - Due: 2024-12-31 - Progress: 5/10
// Fix critical bugs before release
func (m *Milestone) ToMarkdown() string {
	if m.Milestone == nil {
		return "*Invalid milestone*"
	}
	markdown := fmt.Sprintf("**%s** #%d", m.Title, m.ID)
	if m.State != "" {
		markdown += " (" + string(m.State) + ")"
	}
	if m.Deadline != nil {
		markdown += " - Due: " + m.Deadline.Format("2006-01-02")
	}
	if m.ClosedIssues > 0 || m.OpenIssues > 0 {
		total := m.ClosedIssues + m.OpenIssues
		markdown += " - Progress: " + fmt.Sprintf("%d/%d", m.ClosedIssues, total)
	}
	if m.Description != "" {
		markdown += "\n" + m.Description
	}
	return markdown
}

// MilestoneList represents a list of milestones response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/milestones
type MilestoneList []*Milestone

// ToMarkdown renders milestones as a numbered list with essential details only
// Description is omitted to reduce memory usage for AI assistants
// Example:
// 1. **v1.0.0** (open) - Due: 2024-12-31 - Progress: 5/10
// 2. **v0.9.0** (closed) - Progress: 10/10
func (ml MilestoneList) ToMarkdown() string {
	if len(ml) == 0 {
		return "*No milestones found*"
	}
	markdown := ""
	for i, milestone := range ml {
		if milestone.Milestone == nil {
			markdown += fmt.Sprintf("%d. *Invalid milestone*\n", i+1)
			continue
		}

		// Format: **Title** (state) - Due: date - Progress: closed/total
		line := fmt.Sprintf("%d. **%s** #%d", i+1, milestone.Title, milestone.ID)
		if milestone.State != "" {
			line += " (" + string(milestone.State) + ")"
		}
		if milestone.Deadline != nil {
			line += " - Due: " + milestone.Deadline.Format("2006-01-02")
		}
		if milestone.ClosedIssues > 0 || milestone.OpenIssues > 0 {
			total := milestone.ClosedIssues + milestone.OpenIssues
			line += " - Progress: " + fmt.Sprintf("%d/%d", milestone.ClosedIssues, total)
		}
		markdown += line + "\n"
	}
	return markdown
}
