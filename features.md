# Feature List

## Implementation Strategy

- 🟢 **SDK**: Implementation using official SDK
- 🟡 **Custom**: Custom HTTP request implementation (reusing SDK authentication)
- 🔴 **Mixed**: Some features using SDK, some requiring custom implementation

### Label Features

Labels available for a specific repository

- **List Labels** 🟢
  - `GET /repos/{owner}/{repo}/labels`
  - SDK: `ListRepoLabels(owner, repo string, opt ListLabelsOptions) ([]*Label, *Response, error)`
- **Modify label name, description, and color** 🟢
  - `PATCH /repos/{owner}/{repo}/labels/{id}`
  - SDK: `EditLabel(owner, repo string, id int64, opt EditLabelOption) (*Label, *Response, error)`
- **Create or delete labels** 🟢
  - `POST /repos/{owner}/{repo}/labels`
  - SDK: `CreateLabel(owner, repo string, opt CreateLabelOption) (*Label, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/labels/{id}`
  - SDK: `DeleteLabel(owner, repo string, id int64) (*Response, error)`

### Milestone Features 🟢

- **List Milestones**
  - `GET /repos/{owner}/{repo}/milestones`
  - SDK: `ListRepoMilestones(owner, repo string, opt ListMilestoneOption) ([]*Milestone, *Response, error)`
- **Create, delete, and modify milestones (including title, due date, and description)**
  - `POST /repos/{owner}/{repo}/milestones`
  - SDK: `CreateMilestone(owner, repo string, opt CreateMilestoneOption) (*Milestone, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/milestones/{id}`
  - SDK: `DeleteMilestone(owner, repo string, id int64) (*Response, error)`
  - `PATCH /repos/{owner}/{repo}/milestones/{id}`
  - SDK: `EditMilestone(owner, repo string, id int64, opt EditMilestoneOption) (*Milestone, *Response, error)`

### Issue Features 🔴

- **List Repository Issues** 🟢
  - `GET /repos/{owner}/{repo}/issues`
  - SDK: `ListRepoIssues(owner, repo string, opt ListIssueOption) ([]*Issue, *Response, error)`
  - Supports filters: state, labels, milestones, assignees, search, date filters
- **Get Specific Issue Details** 🟢
  - `GET /repos/{owner}/{repo}/issues/{index}`
  - SDK: `GetIssue(owner, repo string, index int64) (*Issue, *Response, error)`
- **List Issue Comments** 🟢
  - `GET /repos/{owner}/{repo}/issues/{index}/comments`
  - SDK: `ListIssueComments(owner, repo string, index int64, opt ListIssueCommentOptions) ([]*Comment, *Response, error)`
- **Create new issue** 🟢
  - `POST /repos/{owner}/{repo}/issues`
  - SDK: `CreateIssue(owner, repo string, opt CreateIssueOption) (*Issue, *Response, error)`
- **Comment on existing issue** 🟢
  - `POST /repos/{owner}/{repo}/issues/{index}/comments`
  - SDK: `CreateIssueComment(owner, repo string, index int64, opt CreateIssueCommentOption) (*Comment, *Response, error)`
- **Close issue** 🟢
  - `PATCH /repos/{owner}/{repo}/issues/{index}` (set `state` to `closed`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
- **Modify issue data** 🟢
  - **Description:** `PATCH /repos/{owner}/{repo}/issues/{index}` (modify `body`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **Labels:** 🟢
    - `POST /repos/{owner}/{repo}/issues/{index}/labels` (add)
    - SDK: `AddIssueLabels(owner, repo string, index int64, opt IssueLabelsOption) ([]*Label, *Response, error)`
    - `DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id}` (remove)
    - SDK: `DeleteIssueLabel(owner, repo string, index, label int64) (*Response, error)`
    - `PUT /repos/{owner}/{repo}/issues/{index}/labels` (replace)
    - SDK: `ReplaceIssueLabels(owner, repo string, index int64, opt IssueLabelsOption) ([]*Label, *Response, error)`
  - **Assignees:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (modify `assignees`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **Milestone:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (modify `milestone`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **Due date:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (modify `due_date`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **Dependency management:** 🟡
    - **Dependencies (issues that block this issue):**
      - **Add dependency:** `POST /repos/{owner}/{repo}/issues/{index}/dependencies`
      - Custom: Not supported by SDK, requires custom HTTP request
      - **List dependencies:** `GET /repos/{owner}/{repo}/issues/{index}/dependencies`
      - Custom: Not supported by SDK, requires custom HTTP request
      - **Remove dependency:** `DELETE /repos/{owner}/{repo}/issues/{index}/dependencies` (via request body)
      - Custom: Not supported by SDK, requires custom HTTP request
    - **Blocking (issues blocked by this issue):**
      - **Add blocking:** `POST /repos/{owner}/{repo}/issues/{index}/blocks`
      - Custom: Not supported by SDK, requires custom HTTP request
      - **List blocking:** `GET /repos/{owner}/{repo}/issues/{index}/blocks`
      - Custom: Not supported by SDK, requires custom HTTP request
      - **Remove blocking:** `DELETE /repos/{owner}/{repo}/issues/{index}/blocks` (via request body)
      - Custom: Not supported by SDK, requires custom HTTP request
- **Edit Issue Comments** 🟢
  - `PATCH /repos/{owner}/{repo}/issues/comments/{id}`
  - SDK: `EditIssueComment(owner, repo string, commentID int64, opt EditIssueCommentOption) (*Comment, *Response, error)`
- **Delete Issue Comments** 🟢
  - `DELETE /repos/{owner}/{repo}/issues/comments/{id}`
  - SDK: `DeleteIssueComment(owner, repo string, commentID int64) (*Response, error)`
- **Attachment management** 🟡
  - **List attachments:** `GET /repos/{owner}/{repo}/issues/{index}/assets`
  - Custom: Not supported by SDK, requires custom HTTP request
  - **Add attachment:** `POST /repos/{owner}/{repo}/issues/{index}/assets`
  - Custom: Not supported by SDK, requires custom HTTP request
  - **Delete attachment:** `DELETE /repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}`
  - Custom: Not supported by SDK, requires custom HTTP request
  - **Modify attachment:** `PATCH /repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}`
  - Custom: Not supported by SDK, requires custom HTTP request

### Wiki Features 🟡

- **Query pages**
  - `GET /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: Not supported by SDK, requires custom HTTP request
- **List pages**
  - `GET /repos/{owner}/{repo}/wiki/pages`
  - Custom: Not supported by SDK, requires custom HTTP request
- **Create, delete, and modify pages**
  - `POST /repos/{owner}/{repo}/wiki/new`
  - Custom: Not supported by SDK, requires custom HTTP request
  - `DELETE /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: Not supported by SDK, requires custom HTTP request
  - `PATCH /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: Not supported by SDK, requires custom HTTP request

### Release Management 🟢

- **List Releases**
  - `GET /repos/{owner}/{repo}/releases`
  - SDK: `ListReleases(owner, repo string, opt ListReleasesOptions) ([]*Release, *Response, error)`
- **Create, delete, and modify releases**
  - `POST /repos/{owner}/{repo}/releases`
  - SDK: `CreateRelease(owner, repo string, opt CreateReleaseOption) (*Release, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/releases/{id}`
  - SDK: `DeleteRelease(user, repo string, id int64) (*Response, error)`
  - `PATCH /repos/{owner}/{repo}/releases/{id}`
  - SDK: `EditRelease(owner, repo string, id int64, form EditReleaseOption) (*Release, *Response, error)`
- **Attachment management**
  - **List attachments:** `GET /repos/{owner}/{repo}/releases/{id}/assets`
  - SDK: `ListReleaseAttachments(user, repo string, release int64, opt ListReleaseAttachmentsOptions) ([]*Attachment, *Response, error)`
  - **Add attachment:** `POST /repos/{owner}/{repo}/releases/{id}/assets`
  - SDK: `CreateReleaseAttachment(user, repo string, release int64, file io.Reader, filename string) (*Attachment, *Response, error)`
  - **Delete attachment:** `DELETE /repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}`
  - SDK: `DeleteReleaseAttachment(user, repo string, release, id int64) (*Response, error)`
  - **Modify attachment:** `PATCH /repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}`
  - SDK: `EditReleaseAttachment(user, repo string, release, attachment int64, form EditAttachmentOptions) (*Attachment, *Response, error)`

### PR Management 🟢

- **List and query PRs**
  - `GET /repos/{owner}/{repo}/pulls`
  - SDK: `ListRepoPullRequests(owner, repo string, opt ListPullRequestsOptions) ([]*PullRequest, *Response, error)`
  - `GET /repos/{owner}/{repo}/pulls/{index}`
  - SDK: `GetPullRequest(owner, repo string, index int64) (*PullRequest, *Response, error)`

### Repository Management 🟢

- **List and query repositories**
  - `GET /repos/search`
  - SDK: `SearchRepos(opt SearchRepoOptions) ([]*Repository, *Response, error)`
  - `GET /user/repos`
  - SDK: `ListMyRepos(opt ListReposOptions) ([]*Repository, *Response, error)`
  - `GET /orgs/{org}/repos`
  - SDK: `ListOrgRepos(org string, opt ListOrgReposOptions) ([]*Repository, *Response, error)`
- **Get Specific Repository Information** 🟢
  - `GET /repos/{owner}/{repo}`
  - SDK: `GetRepo(owner, repo string) (*Repository, *Response, error)`

### Forgejo Actions (CI/CD) 🟡

- **List Action execution tasks**
  - `GET /repos/{owner}/{repo}/actions/tasks`
  - Custom: Not supported by SDK, requires custom HTTP request

## Summary

- 🟢 **Fully supported (5/7)**: Label, Milestone, Release, PR, Repository management
- 🔴 **Partially supported (1/7)**: Issue features (attachments and dependencies require custom implementation)
- 🟡 **Requires custom implementation (2/7)**: Wiki, Forgejo Actions

**Recommended Hybrid approach**: Approximately 71% of features can use SDK, remaining features require custom HTTP requests.