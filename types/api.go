// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Label represents a label response with embedded SDK label
// Used by endpoints:
// - GET /repos/{owner}/{repo}/labels (list)
// - POST /repos/{owner}/{repo}/labels (create)
// - PATCH /repos/{owner}/{repo}/labels/{id} (edit)
// - POST /repos/{owner}/{repo}/issues/{index}/labels (add to issue)
// - PUT /repos/{owner}/{repo}/issues/{index}/labels (replace issue labels)
// - DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id} (remove from issue)
type Label struct {
	*forgejo.Label
}

// ToMarkdown renders a label as a colored badge with name and description
// Example: **bug** `#ff0000` - Something isn't working
func (l *Label) ToMarkdown() string {
	if l.Label == nil {
		return "*Invalid label*"
	}
	markdown := "**" + l.Name + "**"
	if l.Color != "" {
		markdown += " `#" + l.Color + "`"
	}
	if l.Description != "" {
		markdown += " - " + l.Description
	}
	return markdown
}

// LabelList represents a list of labels response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/labels
// - POST /repos/{owner}/{repo}/issues/{index}/labels
// - PUT /repos/{owner}/{repo}/issues/{index}/labels
type LabelList []*Label

// ToMarkdown renders labels as a bullet list of colored badges
// Example:
// - **bug** `#ff0000` - Something isn't working
// - **enhancement** `#a2eeef` - New feature or request
func (ll LabelList) ToMarkdown() string {
	if len(ll) == 0 {
		return "*No labels found*"
	}
	markdown := ""
	for _, label := range ll {
		markdown += "- " + label.ToMarkdown() + "\n"
	}
	return markdown
}

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
	markdown := "**" + m.Title + "**"
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

// ToMarkdown renders milestones as a numbered list with details
// Example:
// 1. **v1.0.0** (open) - Due: 2024-12-31 - Progress: 5/10
// Fix critical bugs before release
// 2. **v0.9.0** (closed) - Progress: 10/10
// Beta release with new features
func (ml MilestoneList) ToMarkdown() string {
	if len(ml) == 0 {
		return "*No milestones found*"
	}
	markdown := ""
	for i, milestone := range ml {
		markdown += fmt.Sprintf("%d. %s\n", i+1, milestone.ToMarkdown())
	}
	return markdown
}

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

// Release represents a release response with embedded SDK release
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases (list)
// - POST /repos/{owner}/{repo}/releases (create)
// - PATCH /repos/{owner}/{repo}/releases/{id} (edit)
type Release struct {
	*forgejo.Release
}

// ToMarkdown renders release with tag, name, draft/prerelease status and description
// Example: **v1.0.0** - Major Release `PRERELEASE` (2024-01-15)
// This release includes new authentication system and bug fixes...
func (r *Release) ToMarkdown() string {
	if r.Release == nil {
		return "*Invalid release*"
	}
	markdown := "**" + r.TagName + "**"
	if r.Title != "" {
		markdown += " - " + r.Title
	}
	if r.IsDraft {
		markdown += " `DRAFT`"
	}
	if r.IsPrerelease {
		markdown += " `PRERELEASE`"
	}
	if !r.CreatedAt.IsZero() {
		markdown += " (" + r.CreatedAt.Format("2006-01-02") + ")"
	}
	if r.Note != "" {
		markdown += "\n" + r.Note
	}
	return markdown
}

// ReleaseList represents a list of releases response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases
type ReleaseList []*Release

// ToMarkdown renders releases as a numbered list with details
// Example:
// 1. **v1.0.0** - Major Release `PRERELEASE` (2024-01-15)
// This release includes new authentication system...
// 2. **v0.9.0** - Beta Release (2024-01-01)
// Initial beta version with core features
func (rl ReleaseList) ToMarkdown() string {
	if len(rl) == 0 {
		return "*No releases found*"
	}
	markdown := ""
	for i, release := range rl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, release.ToMarkdown())
	}
	return markdown
}

// Attachment represents an attachment response with embedded SDK attachment
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases/{id}/assets (list release attachments)
// - POST /repos/{owner}/{repo}/releases/{id}/assets (create release attachment)
// - PATCH /repos/{owner}/{repo}/releases/assets/{id} (edit release attachment)
// - GET /repos/{owner}/{repo}/issues/{index}/attachments (list issue attachments)
// - POST /repos/{owner}/{repo}/issues/{index}/attachments (create issue attachment)
// - PATCH /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id} (edit issue attachment)
type Attachment struct {
	*forgejo.Attachment
}

// ToMarkdown renders attachment with name, size and download link
// Example: **document.pdf** (1024 bytes) [Download](https://git.example.com/attachments/123)
func (a *Attachment) ToMarkdown() string {
	if a.Attachment == nil {
		return "*Invalid attachment*"
	}
	markdown := "**" + a.Name + "**"
	if a.Size > 0 {
		markdown += fmt.Sprintf(" (%d bytes)", a.Size)
	}
	if a.DownloadURL != "" {
		markdown += " [Download](" + a.DownloadURL + ")"
	}
	return markdown
}

// AttachmentList represents a list of attachments response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases/{id}/assets
// - GET /repos/{owner}/{repo}/issues/{index}/attachments
type AttachmentList []*Attachment

// ToMarkdown renders attachments as a bullet list with download info
// Example:
// - **document.pdf** (1024 bytes) [Download](https://git.example.com/attachments/123)
// - **screenshot.png** (2048 bytes) [Download](https://git.example.com/attachments/124)
func (al AttachmentList) ToMarkdown() string {
	if len(al) == 0 {
		return "*No attachments found*"
	}
	markdown := ""
	for _, attachment := range al {
		markdown += "- " + attachment.ToMarkdown() + "\n"
	}
	return markdown
}

// PullRequest represents a pull request response with embedded SDK pull request
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls (list)
// - GET /repos/{owner}/{repo}/pulls/{index} (get)
type PullRequest struct {
	*forgejo.PullRequest
}

// ToMarkdown renders pull request with title, state, author and branch info
// Example: **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
func (pr *PullRequest) ToMarkdown() string {
	if pr.PullRequest == nil {
		return "*Invalid pull request*"
	}
	markdown := fmt.Sprintf("**#%d %s** (%s)\n", pr.Index, pr.Title, pr.State)
	if pr.Poster != nil {
		markdown += "Author: " + pr.Poster.UserName + "\n"
	}
	if pr.Head != nil && pr.Base != nil {
		markdown += fmt.Sprintf("Branch: %s → %s\n", pr.Head.Name, pr.Base.Name)
	}
	if pr.Body != "" {
		markdown += "\n" + pr.Body
	}
	return markdown
}

// PullRequestList represents a list of pull requests response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/pulls
type PullRequestList []*PullRequest

// ToMarkdown renders pull requests as a numbered list with basic info
// Example:
// 1. **#42 Add user authentication** (open)
// Author: johndoe
// Branch: feature/auth → main
//
// This PR implements OAuth2 authentication...
// 2. **#41 Fix database connection** (merged)
// Author: alice
// Branch: bugfix/db → main
func (prl PullRequestList) ToMarkdown() string {
	if len(prl) == 0 {
		return "*No pull requests found*"
	}
	markdown := ""
	for i, pr := range prl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, pr.ToMarkdown())
	}
	return markdown
}

// Repository represents a repository response with embedded SDK repository
// Used by endpoints:
// - GET /repos/search
// - GET /user/repos
// - GET /orgs/{org}/repos
type Repository struct {
	*forgejo.Repository
}

// ToMarkdown renders repository with name, description, stats and key info
// Example: **owner/repo-name** `PRIVATE` `FORK`
// A sample repository for testing purposes
// Stars: 42 | Forks: 7 | Issues: 3 | PRs: 1
// [View Repository](https://git.example.com/owner/repo-name)
func (r *Repository) ToMarkdown() string {
	if r.Repository == nil {
		return "*Invalid repository*"
	}
	markdown := "**" + r.FullName + "**"
	if r.Private {
		markdown += " `PRIVATE`"
	}
	if r.Fork {
		markdown += " `FORK`"
	}
	if r.Template {
		markdown += " `TEMPLATE`"
	}
	markdown += "\n"
	if r.Description != "" {
		markdown += r.Description + "\n"
	}
	markdown += fmt.Sprintf("Stars: %d | Forks: %d | Issues: %d | PRs: %d\n", r.Stars, r.Forks, r.OpenIssues, r.OpenPulls)
	if r.HTMLURL != "" {
		markdown += "[View Repository](" + r.HTMLURL + ")"
	}
	return markdown
}

// RepositoryList represents a list of repositories response
// Used by endpoints:
// - GET /repos/search
// - GET /user/repos
// - GET /orgs/{org}/repos
type RepositoryList []*Repository

// ToMarkdown renders repositories as a numbered list with basic stats
// Example:
// 1. **owner/repo-name** `PRIVATE` `FORK`
// A sample repository for testing purposes
// Stars: 42 | Forks: 7 | Issues: 3 | PRs: 1
// [View Repository](https://git.example.com/owner/repo-name)
// 2. **owner/another-repo**
// Another repository
// Stars: 15 | Forks: 2 | Issues: 0 | PRs: 0
func (rl RepositoryList) ToMarkdown() string {
	if len(rl) == 0 {
		return "*No repositories found*"
	}
	markdown := ""
	for i, repo := range rl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, repo.ToMarkdown())
	}
	return markdown
}

// WikiPage represents a wiki page response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/page/{pageName}
// - POST /repos/{owner}/{repo}/wiki/new
// - PATCH /repos/{owner}/{repo}/wiki/page/{pageName}
type WikiPage struct {
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CommitMessage  string    `json:"commit_message,omitempty"`
	LastCommitSHA  string    `json:"last_commit_sha,omitempty"`
	LastModified   time.Time `json:"last_modified,omitempty"`
	HTMLContentURL string    `json:"html_content_url,omitempty"`
	SubURL         string    `json:"sub_url,omitempty"`
}

// ToMarkdown renders wiki page with title, last modified date and content
// Example: # Getting Started
// *Last modified: 2024-01-15 14:30*
//
// Welcome to our project wiki...
func (w *WikiPage) ToMarkdown() string {
	markdown := "# " + w.Title + "\n"
	if !w.LastModified.IsZero() {
		markdown += "*Last modified: " + w.LastModified.Format("2006-01-02 15:04") + "*\n\n"
	}
	if w.Content != "" {
		markdown += w.Content
	}
	return markdown
}

// WikiPageList represents a list of wiki pages response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/wiki/pages
type WikiPageList []*WikiPage

// ToMarkdown renders wiki pages as a table of contents with links
// Example: ## Wiki Pages
// - **Getting Started** (2024-01-15)
// - **API Documentation** (2024-01-10)
// - **Contributing Guide** (2024-01-05)
func (wl WikiPageList) ToMarkdown() string {
	if len(wl) == 0 {
		return "*No wiki pages found*"
	}
	markdown := "## Wiki Pages\n"
	for _, page := range wl {
		markdown += "- **" + page.Title + "**"
		if !page.LastModified.IsZero() {
			markdown += " (" + page.LastModified.Format("2006-01-02") + ")"
		}
		markdown += "\n"
	}
	return markdown
}

// ActionTask represents an action task response (custom implementation as SDK doesn't support)
// Used by endpoints:
// - GET /repos/{owner}/{repo}/actions/tasks
type ActionTask struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	WorkflowID  int64            `json:"workflow_id"`
	RunID       int64            `json:"run_id"`
	JobID       int64            `json:"job_id"`
	Steps       []ActionTaskStep `json:"steps,omitempty"`
}

// ToMarkdown renders action task with name, status and execution time
// Example: **Build and Test** `success`
// Run ID: 123 | Job ID: 456
// Started: 2024-01-15 14:30 | Duration: 2m15s
// Steps:
//   - Setup Go `success`
//   - Run Tests `success`
//   - Build Binary `success`
func (at *ActionTask) ToMarkdown() string {
	markdown := fmt.Sprintf("**%s** `%s`\n", at.Name, at.Status)
	markdown += fmt.Sprintf("Run ID: %d | Job ID: %d\n", at.RunID, at.JobID)
	if at.StartedAt != nil {
		markdown += "Started: " + at.StartedAt.Format("2006-01-02 15:04")
		if at.CompletedAt != nil {
			duration := at.CompletedAt.Sub(*at.StartedAt)
			markdown += " | Duration: " + duration.String()
		}
		markdown += "\n"
	}
	if len(at.Steps) > 0 {
		markdown += "Steps:\n"
		for _, step := range at.Steps {
			markdown += fmt.Sprintf("  - %s `%s`\n", step.Name, step.Status)
		}
	}
	return markdown
}

// ActionTaskStep represents a step in an action task
type ActionTaskStep struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Conclusion  string     `json:"conclusion,omitempty"`
}

// ToMarkdown renders action task step with name, status and timing
// Example: **Run Tests** `success` (passed) - 1m30s
func (ats *ActionTaskStep) ToMarkdown() string {
	markdown := fmt.Sprintf("**%s** `%s`", ats.Name, ats.Status)
	if ats.Conclusion != "" && ats.Conclusion != ats.Status {
		markdown += " (" + ats.Conclusion + ")"
	}
	if ats.StartedAt != nil && ats.CompletedAt != nil {
		duration := ats.CompletedAt.Sub(*ats.StartedAt)
		markdown += " - " + duration.String()
	}
	return markdown
}

// ActionTaskList represents a list of action tasks response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/actions/tasks
type ActionTaskList []*ActionTask

// ToMarkdown renders action tasks as a numbered list with status
// Example:
// 1. **Build and Test** `success`
// Run ID: 123 | Job ID: 456
// Started: 2024-01-15 14:30 | Duration: 2m15s
// Steps:
//   - Setup Go `success`
//   - Run Tests `success`
//
// 2. **Deploy** `running`
// Run ID: 124 | Job ID: 457
// Started: 2024-01-15 14:35
func (atl ActionTaskList) ToMarkdown() string {
	if len(atl) == 0 {
		return "*No action tasks found*"
	}
	markdown := ""
	for i, task := range atl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, task.ToMarkdown())
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

// MyActionTask represents a Forgejo Actions task.
type MyActionTask struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`
	Event        string    `json:"event"`
	WorkflowID   string    `json:"workflow_id"`
	HeadBranch   string    `json:"head_branch"`
	HeadSHA      string    `json:"head_sha"`
	RunNumber    int64     `json:"run_number"`
	URL          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
}

// MyActionTaskResponse represents the response for listing action tasks.
type MyActionTaskResponse struct {
	TotalCount   int64           `json:"total_count"`
	WorkflowRuns []*MyActionTask `json:"workflow_runs"`
}

// EmptyResponse represents an empty response for endpoints that don't return data
// Used by endpoints that only return status codes:
// - DELETE /repos/{owner}/{repo}/labels/{id}
// - DELETE /repos/{owner}/{repo}/milestones/{id}
// - DELETE /repos/{owner}/{repo}/releases/{id}
// - DELETE /repos/{owner}/{repo}/releases/assets/{id}
// - DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id}
// - DELETE /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}
// - DELETE /repos/{owner}/{repo}/wiki/page/{pageName}
type EmptyResponse struct{}

// ToMarkdown renders simple success message for empty responses
// Example: *Operation completed successfully*
func (er EmptyResponse) ToMarkdown() string {
	return "*Operation completed successfully*"
}
