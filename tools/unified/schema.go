// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package unified

import "fmt"

// Action represents the type of operation being performed.
type Action string

const (
	ActionCreate Action = "create"
	ActionGet    Action = "get"
	ActionList   Action = "list"
	ActionEdit   Action = "edit"
	ActionDelete Action = "delete"
	ActionLink   Action = "link"
	ActionUnlink Action = "unlink"
)

// Resource represents the type of Forgejo resource being operated on.
type Resource string

const (
	ResourceIssue             Resource = "issue"
	ResourceIssueComment      Resource = "issue_comment"
	ResourceIssueAttachment   Resource = "issue_attachment"
	ResourceLabel             Resource = "label"
	ResourceMilestone         Resource = "milestone"
	ResourceRelease           Resource = "release"
	ResourceReleaseAttachment Resource = "release_attachment"
	ResourceWikiPage          Resource = "wiki_page"
	ResourcePullRequest       Resource = "pull_request"
	ResourceRepository        Resource = "repository"
	ResourceActionTask        Resource = "action_task"
)

// LinkType represents the type of relationship between resources.
type LinkType string

const (
	LinkIssueLabel      LinkType = "issue_label"
	LinkIssueDependency LinkType = "issue_dependency"
	LinkIssueBlocking   LinkType = "issue_blocking"
)

// ParamSpec describes a parameter for documentation purposes.
type ParamSpec struct {
	Name        string
	Type        string // "string", "integer", "boolean", "array"
	Required    bool
	Description string
	Enum        []string // Valid values for enum types
}

// ActionDoc contains documentation for a specific action+resource combination.
type ActionDoc struct {
	Description string
	Params      []ParamSpec
	Example     string
}

// ManualEntry is the documentation for a specific action+resource or action+linktype.
type ManualEntry struct {
	Action      Action
	Resource    Resource // For non-link actions
	LinkType    LinkType // For link/unlink actions
	Description string
	Params      []ParamSpec
	Example     string
}

// commonRepoParams returns the common owner/repo parameters.
func commonRepoParams() []ParamSpec {
	return []ParamSpec{
		{Name: "owner", Type: "string", Required: true, Description: "Repository owner (username or organization)"},
		{Name: "repo", Type: "string", Required: true, Description: "Repository name"},
	}
}

// Manual is the documentation registry for all action+resource combinations.
// It provides on-demand documentation lookup and powers rich error messages.
var Manual = map[string]ManualEntry{
	// === CREATE ===
	"create:issue": {
		Action:      ActionCreate,
		Resource:    ResourceIssue,
		Description: "Create a new issue in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "title", Type: "string", Required: true, Description: "Issue title"},
			ParamSpec{Name: "body", Type: "string", Required: true, Description: "Issue body (markdown)"},
			ParamSpec{Name: "assignees", Type: "array", Required: false, Description: "Usernames to assign"},
			ParamSpec{Name: "milestone", Type: "integer", Required: false, Description: "Milestone ID"},
			ParamSpec{Name: "labels", Type: "array", Required: false, Description: "Label IDs to attach"},
			ParamSpec{Name: "due_date", Type: "string", Required: false, Description: "Due date (RFC3339 format)"},
		),
		Example: `create_gitea(resource="issue", owner="org", repo="project", title="Bug report", body="Description...")`,
	},
	"create:issue_comment": {
		Action:      ActionCreate,
		Resource:    ResourceIssueComment,
		Description: "Add a comment to an issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "body", Type: "string", Required: true, Description: "Comment body (markdown)"},
		),
		Example: `create_gitea(resource="issue_comment", owner="org", repo="project", index=42, body="Thanks!")`,
	},
	"create:label": {
		Action:      ActionCreate,
		Resource:    ResourceLabel,
		Description: "Create a new label in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "name", Type: "string", Required: true, Description: "Label name"},
			ParamSpec{Name: "color", Type: "string", Required: true, Description: "Hex color (without #, e.g., 'ff0000')"},
			ParamSpec{Name: "description", Type: "string", Required: false, Description: "Label description"},
		),
		Example: `create_gitea(resource="label", owner="org", repo="project", name="bug", color="ff0000")`,
	},
	"create:milestone": {
		Action:      ActionCreate,
		Resource:    ResourceMilestone,
		Description: "Create a new milestone in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "title", Type: "string", Required: true, Description: "Milestone title"},
			ParamSpec{Name: "description", Type: "string", Required: false, Description: "Milestone description"},
			ParamSpec{Name: "due_date", Type: "string", Required: false, Description: "Due date (RFC3339 format)"},
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "State: 'open' or 'closed'", Enum: []string{"open", "closed"}},
		),
		Example: `create_gitea(resource="milestone", owner="org", repo="project", title="v1.0")`,
	},
	"create:release": {
		Action:      ActionCreate,
		Resource:    ResourceRelease,
		Description: "Create a new release in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "tag_name", Type: "string", Required: true, Description: "Git tag name"},
			ParamSpec{Name: "name", Type: "string", Required: true, Description: "Release title"},
			ParamSpec{Name: "body", Type: "string", Required: false, Description: "Release notes (markdown)"},
			ParamSpec{Name: "target_commitish", Type: "string", Required: false, Description: "Target branch or commit SHA"},
			ParamSpec{Name: "draft", Type: "boolean", Required: false, Description: "Is draft release"},
			ParamSpec{Name: "prerelease", Type: "boolean", Required: false, Description: "Is prerelease"},
		),
		Example: `create_gitea(resource="release", owner="org", repo="project", tag_name="v1.0.0", name="Version 1.0")`,
	},
	"create:wiki_page": {
		Action:      ActionCreate,
		Resource:    ResourceWikiPage,
		Description: "Create a new wiki page in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "title", Type: "string", Required: true, Description: "Page title"},
			ParamSpec{Name: "content", Type: "string", Required: true, Description: "Page content (markdown)"},
			ParamSpec{Name: "message", Type: "string", Required: false, Description: "Commit message"},
		),
		Example: `create_gitea(resource="wiki_page", owner="org", repo="project", title="Home", content="# Welcome")`,
	},
	"create:pull_request": {
		Action:      ActionCreate,
		Resource:    ResourcePullRequest,
		Description: "Create a new pull request.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "title", Type: "string", Required: true, Description: "PR title"},
			ParamSpec{Name: "body", Type: "string", Required: false, Description: "PR description (markdown)"},
			ParamSpec{Name: "head", Type: "string", Required: true, Description: "Source branch (or 'user:branch' for forks)"},
			ParamSpec{Name: "base", Type: "string", Required: true, Description: "Target branch"},
			ParamSpec{Name: "assignees", Type: "array", Required: false, Description: "Usernames to assign"},
			ParamSpec{Name: "milestone", Type: "integer", Required: false, Description: "Milestone ID"},
			ParamSpec{Name: "labels", Type: "array", Required: false, Description: "Label IDs"},
		),
		Example: `create_gitea(resource="pull_request", owner="org", repo="project", title="Feature X", head="feature-x", base="main")`,
	},

	// === GET ===
	"get:issue": {
		Action:      ActionGet,
		Resource:    ResourceIssue,
		Description: "Get details of a specific issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
		),
		Example: `get_gitea(resource="issue", owner="org", repo="project", index=42)`,
	},
	"get:wiki_page": {
		Action:      ActionGet,
		Resource:    ResourceWikiPage,
		Description: "Get content of a specific wiki page.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "page_name", Type: "string", Required: true, Description: "Wiki page name"},
		),
		Example: `get_gitea(resource="wiki_page", owner="org", repo="project", page_name="Home")`,
	},
	"get:pull_request": {
		Action:      ActionGet,
		Resource:    ResourcePullRequest,
		Description: "Get details of a specific pull request.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "PR number"},
		),
		Example: `get_gitea(resource="pull_request", owner="org", repo="project", index=42)`,
	},
	"get:repository": {
		Action:      ActionGet,
		Resource:    ResourceRepository,
		Description: "Get details of a specific repository.",
		Params:      commonRepoParams(),
		Example:     `get_gitea(resource="repository", owner="org", repo="project")`,
	},

	// === LIST ===
	"list:issue": {
		Action:      ActionList,
		Resource:    ResourceIssue,
		Description: "List issues in a repository with optional filtering.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "Filter by state", Enum: []string{"open", "closed", "all"}},
			ParamSpec{Name: "labels", Type: "string", Required: false, Description: "Comma-separated label names"},
			ParamSpec{Name: "milestones", Type: "string", Required: false, Description: "Comma-separated milestone names/IDs"},
			ParamSpec{Name: "assignees", Type: "string", Required: false, Description: "Comma-separated usernames"},
			ParamSpec{Name: "q", Type: "string", Required: false, Description: "Search query"},
			ParamSpec{Name: "sort", Type: "string", Required: false, Description: "Sort field", Enum: []string{"created", "updated", "comments"}},
			ParamSpec{Name: "order", Type: "string", Required: false, Description: "Sort order", Enum: []string{"asc", "desc"}},
			ParamSpec{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			ParamSpec{Name: "limit", Type: "integer", Required: false, Description: "Results per page (max 50)"},
			ParamSpec{Name: "since", Type: "string", Required: false, Description: "Only issues updated after (RFC3339)"},
			ParamSpec{Name: "before", Type: "string", Required: false, Description: "Only issues updated before (RFC3339)"},
		),
		Example: `list_gitea(resource="issue", owner="org", repo="project", state="open", labels="bug")`,
	},
	"list:issue_comment": {
		Action:      ActionList,
		Resource:    ResourceIssueComment,
		Description: "List comments on an issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "since", Type: "string", Required: false, Description: "Only comments after (RFC3339)"},
			ParamSpec{Name: "before", Type: "string", Required: false, Description: "Only comments before (RFC3339)"},
		),
		Example: `list_gitea(resource="issue_comment", owner="org", repo="project", index=42)`,
	},
	"list:issue_attachment": {
		Action:      ActionList,
		Resource:    ResourceIssueAttachment,
		Description: "List attachments on an issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
		),
		Example: `list_gitea(resource="issue_attachment", owner="org", repo="project", index=42)`,
	},
	"list:label": {
		Action:      ActionList,
		Resource:    ResourceLabel,
		Description: "List all labels in a repository.",
		Params:      commonRepoParams(),
		Example:     `list_gitea(resource="label", owner="org", repo="project")`,
	},
	"list:milestone": {
		Action:      ActionList,
		Resource:    ResourceMilestone,
		Description: "List milestones in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "Filter by state", Enum: []string{"open", "closed", "all"}},
			ParamSpec{Name: "name", Type: "string", Required: false, Description: "Filter by name"},
			ParamSpec{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			ParamSpec{Name: "limit", Type: "integer", Required: false, Description: "Results per page"},
		),
		Example: `list_gitea(resource="milestone", owner="org", repo="project", state="open")`,
	},
	"list:release": {
		Action:      ActionList,
		Resource:    ResourceRelease,
		Description: "List releases in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			ParamSpec{Name: "limit", Type: "integer", Required: false, Description: "Results per page"},
		),
		Example: `list_gitea(resource="release", owner="org", repo="project")`,
	},
	"list:release_attachment": {
		Action:      ActionList,
		Resource:    ResourceReleaseAttachment,
		Description: "List attachments on a release.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Release ID"},
		),
		Example: `list_gitea(resource="release_attachment", owner="org", repo="project", id=1)`,
	},
	"list:wiki_page": {
		Action:      ActionList,
		Resource:    ResourceWikiPage,
		Description: "List all wiki pages in a repository.",
		Params:      commonRepoParams(),
		Example:     `list_gitea(resource="wiki_page", owner="org", repo="project")`,
	},
	"list:pull_request": {
		Action:      ActionList,
		Resource:    ResourcePullRequest,
		Description: "List pull requests in a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "Filter by state", Enum: []string{"open", "closed", "all"}},
			ParamSpec{Name: "sort", Type: "string", Required: false, Description: "Sort field", Enum: []string{"oldest", "recentupdate", "leastupdate", "mostcomment", "leastcomment", "priority"}},
			ParamSpec{Name: "milestone", Type: "integer", Required: false, Description: "Milestone ID"},
			ParamSpec{Name: "labels", Type: "array", Required: false, Description: "Label IDs"},
			ParamSpec{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			ParamSpec{Name: "limit", Type: "integer", Required: false, Description: "Results per page"},
		),
		Example: `list_gitea(resource="pull_request", owner="org", repo="project", state="open")`,
	},
	"list:repository": {
		Action:      ActionList,
		Resource:    ResourceRepository,
		Description: "List repositories. Use 'scope' to choose: 'my' (authenticated user), 'org' (organization), or 'search' (search all).",
		Params: []ParamSpec{
			{Name: "scope", Type: "string", Required: true, Description: "Listing scope", Enum: []string{"my", "org", "search"}},
			{Name: "org", Type: "string", Required: false, Description: "Organization name (required for scope='org')"},
			{Name: "q", Type: "string", Required: false, Description: "Search query (for scope='search')"},
			{Name: "topic", Type: "boolean", Required: false, Description: "Search in topics (for scope='search')"},
			{Name: "include_desc", Type: "boolean", Required: false, Description: "Search in description (for scope='search')"},
			{Name: "template", Type: "boolean", Required: false, Description: "Only templates (for scope='search')"},
			{Name: "archived", Type: "boolean", Required: false, Description: "Only archived (for scope='search')"},
			{Name: "private", Type: "boolean", Required: false, Description: "Only private (for scope='search')"},
			{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			{Name: "limit", Type: "integer", Required: false, Description: "Results per page"},
		},
		Example: `list_gitea(resource="repository", scope="my")`,
	},
	"list:action_task": {
		Action:      ActionList,
		Resource:    ResourceActionTask,
		Description: "List Forgejo Actions tasks (CI/CD runs).",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "page", Type: "integer", Required: false, Description: "Page number"},
			ParamSpec{Name: "limit", Type: "integer", Required: false, Description: "Results per page"},
		),
		Example: `list_gitea(resource="action_task", owner="org", repo="project")`,
	},
	"list:issue_dependency": {
		Action:      ActionList,
		Resource:    Resource("issue_dependency"),
		Description: "List issues that block this issue (dependencies).",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
		),
		Example: `list_gitea(resource="issue_dependency", owner="org", repo="project", index=42)`,
	},
	"list:issue_blocking": {
		Action:      ActionList,
		Resource:    Resource("issue_blocking"),
		Description: "List issues that are blocked by this issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
		),
		Example: `list_gitea(resource="issue_blocking", owner="org", repo="project", index=42)`,
	},

	// === EDIT ===
	"edit:issue": {
		Action:      ActionEdit,
		Resource:    ResourceIssue,
		Description: "Edit an existing issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "title", Type: "string", Required: false, Description: "New title"},
			ParamSpec{Name: "body", Type: "string", Required: false, Description: "New body"},
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "New state", Enum: []string{"open", "closed"}},
			ParamSpec{Name: "assignees", Type: "array", Required: false, Description: "New assignees"},
			ParamSpec{Name: "milestone", Type: "integer", Required: false, Description: "New milestone ID"},
			ParamSpec{Name: "due_date", Type: "string", Required: false, Description: "New due date (RFC3339)"},
		),
		Example: `edit_gitea(resource="issue", owner="org", repo="project", index=42, state="closed")`,
	},
	"edit:issue_comment": {
		Action:      ActionEdit,
		Resource:    ResourceIssueComment,
		Description: "Edit an issue comment.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Comment ID"},
			ParamSpec{Name: "body", Type: "string", Required: true, Description: "New comment body"},
		),
		Example: `edit_gitea(resource="issue_comment", owner="org", repo="project", id=123, body="Updated comment")`,
	},
	"edit:issue_attachment": {
		Action:      ActionEdit,
		Resource:    ResourceIssueAttachment,
		Description: "Edit an issue attachment's name.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "attachment_id", Type: "integer", Required: true, Description: "Attachment ID"},
			ParamSpec{Name: "name", Type: "string", Required: true, Description: "New filename"},
		),
		Example: `edit_gitea(resource="issue_attachment", owner="org", repo="project", index=42, attachment_id=1, name="new_name.png")`,
	},
	"edit:label": {
		Action:      ActionEdit,
		Resource:    ResourceLabel,
		Description: "Edit a label.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Label ID"},
			ParamSpec{Name: "name", Type: "string", Required: false, Description: "New name"},
			ParamSpec{Name: "color", Type: "string", Required: false, Description: "New color (hex without #)"},
			ParamSpec{Name: "description", Type: "string", Required: false, Description: "New description"},
		),
		Example: `edit_gitea(resource="label", owner="org", repo="project", id=1, color="00ff00")`,
	},
	"edit:milestone": {
		Action:      ActionEdit,
		Resource:    ResourceMilestone,
		Description: "Edit a milestone.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Milestone ID"},
			ParamSpec{Name: "title", Type: "string", Required: false, Description: "New title"},
			ParamSpec{Name: "description", Type: "string", Required: false, Description: "New description"},
			ParamSpec{Name: "due_date", Type: "string", Required: false, Description: "New due date (RFC3339)"},
			ParamSpec{Name: "state", Type: "string", Required: false, Description: "New state", Enum: []string{"open", "closed"}},
		),
		Example: `edit_gitea(resource="milestone", owner="org", repo="project", id=1, state="closed")`,
	},
	"edit:release": {
		Action:      ActionEdit,
		Resource:    ResourceRelease,
		Description: "Edit a release.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Release ID"},
			ParamSpec{Name: "tag_name", Type: "string", Required: false, Description: "New tag name"},
			ParamSpec{Name: "name", Type: "string", Required: false, Description: "New title"},
			ParamSpec{Name: "body", Type: "string", Required: false, Description: "New release notes"},
			ParamSpec{Name: "target_commitish", Type: "string", Required: false, Description: "New target"},
			ParamSpec{Name: "draft", Type: "boolean", Required: false, Description: "Is draft"},
			ParamSpec{Name: "prerelease", Type: "boolean", Required: false, Description: "Is prerelease"},
		),
		Example: `edit_gitea(resource="release", owner="org", repo="project", id=1, prerelease=false)`,
	},
	"edit:release_attachment": {
		Action:      ActionEdit,
		Resource:    ResourceReleaseAttachment,
		Description: "Edit a release attachment's name.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Release ID"},
			ParamSpec{Name: "attachment_id", Type: "integer", Required: true, Description: "Attachment ID"},
			ParamSpec{Name: "name", Type: "string", Required: true, Description: "New filename"},
		),
		Example: `edit_gitea(resource="release_attachment", owner="org", repo="project", id=1, attachment_id=2, name="app-v1.0.zip")`,
	},
	"edit:wiki_page": {
		Action:      ActionEdit,
		Resource:    ResourceWikiPage,
		Description: "Edit a wiki page.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "page_name", Type: "string", Required: true, Description: "Current page name"},
			ParamSpec{Name: "title", Type: "string", Required: false, Description: "New title"},
			ParamSpec{Name: "content", Type: "string", Required: true, Description: "New content"},
			ParamSpec{Name: "message", Type: "string", Required: false, Description: "Commit message"},
		),
		Example: `edit_gitea(resource="wiki_page", owner="org", repo="project", page_name="Home", content="# Updated")`,
	},

	// === DELETE ===
	"delete:issue_comment": {
		Action:      ActionDelete,
		Resource:    ResourceIssueComment,
		Description: "Delete an issue comment.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Comment ID"},
		),
		Example: `delete_gitea(resource="issue_comment", owner="org", repo="project", id=123)`,
	},
	"delete:issue_attachment": {
		Action:      ActionDelete,
		Resource:    ResourceIssueAttachment,
		Description: "Delete an issue attachment.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "attachment_id", Type: "integer", Required: true, Description: "Attachment ID"},
		),
		Example: `delete_gitea(resource="issue_attachment", owner="org", repo="project", index=42, attachment_id=1)`,
	},
	"delete:label": {
		Action:      ActionDelete,
		Resource:    ResourceLabel,
		Description: "Delete a label from a repository.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Label ID"},
		),
		Example: `delete_gitea(resource="label", owner="org", repo="project", id=1)`,
	},
	"delete:milestone": {
		Action:      ActionDelete,
		Resource:    ResourceMilestone,
		Description: "Delete a milestone.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Milestone ID"},
		),
		Example: `delete_gitea(resource="milestone", owner="org", repo="project", id=1)`,
	},
	"delete:release": {
		Action:      ActionDelete,
		Resource:    ResourceRelease,
		Description: "Delete a release.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Release ID"},
		),
		Example: `delete_gitea(resource="release", owner="org", repo="project", id=1)`,
	},
	"delete:release_attachment": {
		Action:      ActionDelete,
		Resource:    ResourceReleaseAttachment,
		Description: "Delete a release attachment.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "id", Type: "integer", Required: true, Description: "Release ID"},
			ParamSpec{Name: "attachment_id", Type: "integer", Required: true, Description: "Attachment ID"},
		),
		Example: `delete_gitea(resource="release_attachment", owner="org", repo="project", id=1, attachment_id=2)`,
	},
	"delete:wiki_page": {
		Action:      ActionDelete,
		Resource:    ResourceWikiPage,
		Description: "Delete a wiki page.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "page_name", Type: "string", Required: true, Description: "Page name to delete"},
		),
		Example: `delete_gitea(resource="wiki_page", owner="org", repo="project", page_name="OldPage")`,
	},

	// === LINK ===
	"link:issue_label": {
		Action:      ActionLink,
		LinkType:    LinkIssueLabel,
		Description: "Add labels to an issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "labels", Type: "array", Required: true, Description: "Label IDs to add"},
		),
		Example: `link_gitea(type="issue_label", owner="org", repo="project", index=42, labels=[1, 2])`,
	},
	"link:issue_dependency": {
		Action:      ActionLink,
		LinkType:    LinkIssueDependency,
		Description: "Add a dependency: issue cannot be closed until dependency_index is closed.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Dependent issue number"},
			ParamSpec{Name: "dependency_index", Type: "integer", Required: true, Description: "Issue that blocks this one"},
		),
		Example: `link_gitea(type="issue_dependency", owner="org", repo="project", index=42, dependency_index=10)`,
	},
	"link:issue_blocking": {
		Action:      ActionLink,
		LinkType:    LinkIssueBlocking,
		Description: "Add a blocking relationship: blocked_index cannot be closed until index is closed.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Blocking issue number"},
			ParamSpec{Name: "blocked_index", Type: "integer", Required: true, Description: "Issue that will be blocked"},
		),
		Example: `link_gitea(type="issue_blocking", owner="org", repo="project", index=42, blocked_index=50)`,
	},

	// === UNLINK ===
	"unlink:issue_label": {
		Action:      ActionUnlink,
		LinkType:    LinkIssueLabel,
		Description: "Remove a label from an issue.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Issue number"},
			ParamSpec{Name: "label_id", Type: "integer", Required: true, Description: "Label ID to remove"},
		),
		Example: `unlink_gitea(type="issue_label", owner="org", repo="project", index=42, label_id=1)`,
	},
	"unlink:issue_dependency": {
		Action:      ActionUnlink,
		LinkType:    LinkIssueDependency,
		Description: "Remove a dependency relationship.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Dependent issue number"},
			ParamSpec{Name: "dependency_index", Type: "integer", Required: true, Description: "Dependency to remove"},
		),
		Example: `unlink_gitea(type="issue_dependency", owner="org", repo="project", index=42, dependency_index=10)`,
	},
	"unlink:issue_blocking": {
		Action:      ActionUnlink,
		LinkType:    LinkIssueBlocking,
		Description: "Remove a blocking relationship.",
		Params: append(commonRepoParams(),
			ParamSpec{Name: "index", Type: "integer", Required: true, Description: "Blocking issue number"},
			ParamSpec{Name: "blocked_index", Type: "integer", Required: true, Description: "Issue to unblock"},
		),
		Example: `unlink_gitea(type="issue_blocking", owner="org", repo="project", index=42, blocked_index=50)`,
	},
}

// LookupManual retrieves documentation for an action+resource or action+linktype combination.
func LookupManual(action Action, resourceOrType string) (ManualEntry, bool) {
	key := fmt.Sprintf("%s:%s", action, resourceOrType)
	entry, ok := Manual[key]
	return entry, ok
}

// FormatManualEntry renders a ManualEntry as human-readable documentation.
func FormatManualEntry(entry ManualEntry) string {
	var result string

	// Header
	if entry.Resource != "" {
		result = fmt.Sprintf("## %s %s\n\n", entry.Action, entry.Resource)
	} else {
		result = fmt.Sprintf("## %s %s\n\n", entry.Action, entry.LinkType)
	}

	// Description
	result += entry.Description + "\n\n"

	// Parameters
	result += "### Parameters\n\n"
	result += "| Name | Type | Required | Description |\n"
	result += "|------|------|----------|-------------|\n"
	for _, p := range entry.Params {
		req := "no"
		if p.Required {
			req = "yes"
		}
		desc := p.Description
		if len(p.Enum) > 0 {
			desc += fmt.Sprintf(" (values: %v)", p.Enum)
		}
		result += fmt.Sprintf("| %s | %s | %s | %s |\n", p.Name, p.Type, req, desc)
	}

	// Example
	result += fmt.Sprintf("\n### Example\n\n```\n%s\n```\n", entry.Example)

	return result
}

// FormatValidationError creates a helpful error message with the relevant schema.
func FormatValidationError(action Action, resourceOrType, message string) string {
	entry, ok := LookupManual(action, resourceOrType)
	if !ok {
		return fmt.Sprintf("Error: %s\n\nNo documentation found for %s:%s", message, action, resourceOrType)
	}

	result := fmt.Sprintf("Error: %s\n\n", message)
	result += "Here's how to use this operation:\n\n"
	result += FormatManualEntry(entry)
	return result
}

// ListResourcesForAction returns all valid resources for a given action.
func ListResourcesForAction(action Action) []string {
	var resources []string
	prefix := fmt.Sprintf("%s:", action)
	for key := range Manual {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			resources = append(resources, key[len(prefix):])
		}
	}
	return resources
}

// ListLinkTypes returns all valid link types.
func ListLinkTypes() []string {
	return []string{
		string(LinkIssueLabel),
		string(LinkIssueDependency),
		string(LinkIssueBlocking),
	}
}
