// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Issue represents an issue response with embedded SDK issue
// Used by endpoints:
// - POST /repos/{owner}/{repo}/issues (create)
// - PATCH /repos/{owner}/{repo}/issues/{index} (edit)
type Issue struct {
	*forgejo.Issue
}

// ToMarkdown renders issue with title, state, assignees, labels and basic info
// Example: **#123 Fix login bug** (open)
// Author: johndoe
// Assignees: [alice bob]
// Labels: [bug priority-high]
// Milestone: v1.0.0
// Due: 2024-12-31
//
// The login page crashes when...
func (i *Issue) ToMarkdown() string {
	if i.Issue == nil {
		return "*Invalid issue*"
	}
	markdown := fmt.Sprintf("**#%d %s** (%s)\n", i.Index, i.Title, i.State)
	if i.Poster != nil {
		markdown += "Author: " + i.Poster.UserName + "\n"
	}
	if len(i.Assignees) > 0 {
		assignees := make([]string, len(i.Assignees))
		for j, assignee := range i.Assignees {
			assignees[j] = assignee.UserName
		}
		markdown += "Assignees: " + fmt.Sprintf("%v", assignees) + "\n"
	}
	if len(i.Labels) > 0 {
		labelNames := make([]string, len(i.Labels))
		for j, label := range i.Labels {
			labelNames[j] = label.Name
		}
		markdown += "Labels: " + fmt.Sprintf("%v", labelNames) + "\n"
	}
	if i.Milestone != nil {
		markdown += "Milestone: " + i.Milestone.Title + "\n"
	}
	if i.Deadline != nil {
		markdown += "Due: " + i.Deadline.Format("2006-01-02") + "\n"
	}
	if i.Body != "" {
		markdown += "\n" + i.Body
	}
	return markdown
}

// Comment represents a comment response with embedded SDK comment
// Used by endpoints:
// - POST /repos/{owner}/{repo}/issues/{index}/comments (create)
type Comment struct {
	*forgejo.Comment
}

// ToMarkdown renders comment with author, timestamp and content
// Example: **alice** (2024-01-15 14:30)
// I think the issue is in the authentication module...
func (c *Comment) ToMarkdown() string {
	if c.Comment == nil {
		return "*Invalid comment*"
	}
	markdown := ""
	if c.Poster != nil {
		markdown += "**" + c.Poster.UserName + "**"
	}
	if !c.Created.IsZero() {
		markdown += " (" + c.Created.Format("2006-01-02 15:04") + ")"
	}
	if c.Body != "" {
		markdown += "\n" + c.Body
	}
	return markdown
}

// IssueDependency represents an issue dependency response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - POST /repos/{owner}/{repo}/issues/{index}/dependencies
type IssueDependency struct {
	ID           int64  `json:"id"`
	IssueID      int64  `json:"issue_id"`
	DependencyID int64  `json:"dependency_id"`
	CreatedUnix  int64  `json:"created_unix"`
	Issue        *Issue `json:"issue,omitempty"`
	Dependency   *Issue `json:"dependency,omitempty"`
}

// ToMarkdown renders issue dependency with issue titles and numbers
// Example: #123 Fix login bug depends on #45 Update authentication library
func (id *IssueDependency) ToMarkdown() string {
	markdown := ""
	if id.Issue != nil {
		markdown += fmt.Sprintf("#%d %s", id.Issue.Index, id.Issue.Title)
	} else {
		markdown += fmt.Sprintf("Issue #%d", id.IssueID)
	}
	markdown += " depends on "
	if id.Dependency != nil {
		markdown += fmt.Sprintf("#%d %s", id.Dependency.Index, id.Dependency.Title)
	} else {
		markdown += fmt.Sprintf("Issue #%d", id.DependencyID)
	}
	return markdown
}

// IssueList represents a list of issues optimized for list display
// Used by list_repo_issues endpoint to show essential information only
type IssueList []*forgejo.Issue

// ToMarkdown renders issues with essential information for quick scanning
// Shows: Index, Title, State, Assignees, Labels, Updated time, Comments count
// Example per issue:
// #123 Fix login bug (open) | [testuser] | [bug priority-high] | 2024-01-15 | 5 comments
func (il IssueList) ToMarkdown() string {
	if len(il) == 0 {
		return "*No issues found*"
	}

	markdown := ""
	for _, issue := range il {
		if issue == nil {
			continue
		}

		// Index, Title, and State
		line := fmt.Sprintf("#%d %s (%s)", issue.Index, issue.Title, issue.State)

		// Assignees
		if len(issue.Assignees) > 0 {
			assigneeNames := make([]string, len(issue.Assignees))
			for i, assignee := range issue.Assignees {
				assigneeNames[i] = assignee.UserName
			}
			line += " | [" + strings.Join(assigneeNames, " ") + "]"
		}

		// Labels
		if len(issue.Labels) > 0 {
			labelNames := make([]string, len(issue.Labels))
			for i, label := range issue.Labels {
				labelNames[i] = label.Name
			}
			line += " | [" + strings.Join(labelNames, " ") + "]"
		}

		// Updated time
		if !issue.Updated.IsZero() {
			line += " | " + issue.Updated.Format("2006-01-02")
		}

		// Comments count
		line += fmt.Sprintf(" | %d", issue.Comments)

		markdown += line + "\n"
	}

	return markdown
}

// IssueDependencyList represents a list of issue dependencies response
type IssueDependencyList []*IssueDependency

// ToMarkdown renders issue dependencies as a bullet list
// Example:
// - #123 Fix login bug depends on #45 Update authentication library
// - #124 Add user profile depends on #46 Database migration
func (idl IssueDependencyList) ToMarkdown() string {
	if len(idl) == 0 {
		return "*No issue dependencies found*"
	}
	markdown := ""
	for _, dep := range idl {
		markdown += "- " + dep.ToMarkdown() + "\n"
	}
	return markdown
}
