# 功能列表

## 實作策略說明

- 🟢 **SDK**：使用官方 SDK 實作
- 🟡 **Custom**：自訂 HTTP 請求實作（重用 SDK 的認證機制）
- 🔴 **Mixed**：部分功能使用 SDK，部分需要自訂實作

### Label 相關功能

針對某個 repo 可以使用的 label

- **列出 Label** 🟢
  - `GET /repos/{owner}/{repo}/labels`
  - SDK: `ListRepoLabels(owner, repo string, opt ListLabelsOptions) ([]*Label, *Response, error)`
- **修改可用的 label 的名稱、說明、顏色** 🟢
  - `PATCH /repos/{owner}/{repo}/labels/{id}`
  - SDK: `EditLabel(owner, repo string, id int64, opt EditLabelOption) (*Label, *Response, error)`
- **新增或刪除 label** 🟢
  - `POST /repos/{owner}/{repo}/labels`
  - SDK: `CreateLabel(owner, repo string, opt CreateLabelOption) (*Label, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/labels/{id}`
  - SDK: `DeleteLabel(owner, repo string, id int64) (*Response, error)`

### Milestone 相關功能 🟢

- **列出 Milestone**
  - `GET /repos/{owner}/{repo}/milestones`
  - SDK: `ListRepoMilestones(owner, repo string, opt ListMilestoneOption) ([]*Milestone, *Response, error)`
- **建立、刪除、修改 milestone (包括標題、到期時間和說明)**
  - `POST /repos/{owner}/{repo}/milestones`
  - SDK: `CreateMilestone(owner, repo string, opt CreateMilestoneOption) (*Milestone, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/milestones/{id}`
  - SDK: `DeleteMilestone(owner, repo string, id int64) (*Response, error)`
  - `PATCH /repos/{owner}/{repo}/milestones/{id}`
  - SDK: `EditMilestone(owner, repo string, id int64, opt EditMilestoneOption) (*Milestone, *Response, error)`

### Issue 相關功能 🔴

- **建立新的 issue** 🟢
  - `POST /repos/{owner}/{repo}/issues`
  - SDK: `CreateIssue(owner, repo string, opt CreateIssueOption) (*Issue, *Response, error)`
- **在現有的 issue 上留言** 🟢
  - `POST /repos/{owner}/{repo}/issues/{index}/comments`
  - SDK: `CreateIssueComment(owner, repo string, index int64, opt CreateIssueCommentOption) (*Comment, *Response, error)`
- **關閉 issue** 🟢
  - `PATCH /repos/{owner}/{repo}/issues/{index}` (設定 `state` 為 `closed`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
- **修改 issue 資料** 🟢
  - **說明:** `PATCH /repos/{owner}/{repo}/issues/{index}` (修改 `body`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **label:** 🟢
    - `POST /repos/{owner}/{repo}/issues/{index}/labels` (新增)
    - SDK: `AddIssueLabels(owner, repo string, index int64, opt IssueLabelsOption) ([]*Label, *Response, error)`
    - `DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id}` (移除)
    - SDK: `DeleteIssueLabel(owner, repo string, index, label int64) (*Response, error)`
    - `PUT /repos/{owner}/{repo}/issues/{index}/labels` (取代)
    - SDK: `ReplaceIssueLabels(owner, repo string, index int64, opt IssueLabelsOption) ([]*Label, *Response, error)`
  - **負責人:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (修改 `assignees`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **milestone:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (修改 `milestone`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **到期時間:** 🟢 `PATCH /repos/{owner}/{repo}/issues/{index}` (修改 `due_date`)
  - SDK: `EditIssue(owner, repo string, index int64, opt EditIssueOption) (*Issue, *Response, error)`
  - **依賴管理:** 🟡
    - **新增依賴:** `POST /repos/{owner}/{repo}/issues/{index}/dependencies`
    - Custom: SDK 無支援，需自訂 HTTP 請求
    - **列出依賴:** `GET /repos/{owner}/{repo}/issues/{index}/dependencies`
    - Custom: SDK 無支援，需自訂 HTTP 請求
    - **移除依賴:** `DELETE /repos/{owner}/{repo}/issues/{index}/dependencies/{dependency_index}`
    - Custom: SDK 無支援，需自訂 HTTP 請求
- **附件管理** 🟡
  - **列出附件:** `GET /repos/{owner}/{repo}/issues/{index}/attachments`
  - Custom: SDK 無支援，需自訂 HTTP 請求
  - **新增附件:** `POST /repos/{owner}/{repo}/issues/{index}/attachments`
  - Custom: SDK 無支援，需自訂 HTTP 請求
  - **刪除附件:** `DELETE /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}`
  - Custom: SDK 無支援，需自訂 HTTP 請求
  - **修改附件:** `PATCH /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}`
  - Custom: SDK 無支援，需自訂 HTTP 請求

### Wiki 相關功能 🟡

- **查詢頁面**
  - `GET /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: SDK 無支援，需自訂 HTTP 請求
- **頁面列表**
  - `GET /repos/{owner}/{repo}/wiki/pages`
  - Custom: SDK 無支援，需自訂 HTTP 請求
- **新增、刪除、修改頁面**
  - `POST /repos/{owner}/{repo}/wiki/new`
  - Custom: SDK 無支援，需自訂 HTTP 請求
  - `DELETE /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: SDK 無支援，需自訂 HTTP 請求
  - `PATCH /repos/{owner}/{repo}/wiki/page/{pageName}`
  - Custom: SDK 無支援，需自訂 HTTP 請求

### Release 管理 🟢

- **列出 Release**
  - `GET /repos/{owner}/{repo}/releases`
  - SDK: `ListReleases(owner, repo string, opt ListReleasesOptions) ([]*Release, *Response, error)`
- **建立、刪除、修改 release**
  - `POST /repos/{owner}/{repo}/releases`
  - SDK: `CreateRelease(owner, repo string, opt CreateReleaseOption) (*Release, *Response, error)`
  - `DELETE /repos/{owner}/{repo}/releases/{id}`
  - SDK: `DeleteRelease(user, repo string, id int64) (*Response, error)`
  - `PATCH /repos/{owner}/{repo}/releases/{id}`
  - SDK: `EditRelease(owner, repo string, id int64, form EditReleaseOption) (*Release, *Response, error)`
- **附件管理**
  - **列出附件:** `GET /repos/{owner}/{repo}/releases/{id}/assets`
  - SDK: `ListReleaseAttachments(user, repo string, release int64, opt ListReleaseAttachmentsOptions) ([]*Attachment, *Response, error)`
  - **新增附件:** `POST /repos/{owner}/{repo}/releases/{id}/assets`
  - SDK: `CreateReleaseAttachment(user, repo string, release int64, file io.Reader, filename string) (*Attachment, *Response, error)`
  - **刪除附件:** `DELETE /repos/{owner}/{repo}/releases/assets/{id}`
  - SDK: `DeleteReleaseAttachment(user, repo string, release, id int64) (*Response, error)`
  - **修改附件:** `PATCH /repos/{owner}/{repo}/releases/assets/{id}`
  - SDK: `EditReleaseAttachment(user, repo string, release, attachment int64, form EditAttachmentOptions) (*Attachment, *Response, error)`

### PR 管理 🟢

- **列出及查詢 PR**
  - `GET /repos/{owner}/{repo}/pulls`
  - SDK: `ListRepoPullRequests(owner, repo string, opt ListPullRequestsOptions) ([]*PullRequest, *Response, error)`
  - `GET /repos/{owner}/{repo}/pulls/{index}`
  - SDK: `GetPullRequest(owner, repo string, index int64) (*PullRequest, *Response, error)`

### Repo 管理 🟢

- **列出及查詢 repo**
  - `GET /repos/search`
  - SDK: `SearchRepos(opt SearchRepoOptions) ([]*Repository, *Response, error)`
  - `GET /user/repos`
  - SDK: `ListMyRepos(opt ListReposOptions) ([]*Repository, *Response, error)`
  - `GET /orgs/{org}/repos`
  - SDK: `ListOrgRepos(org string, opt ListOrgReposOptions) ([]*Repository, *Response, error)`

### Forgejo Actions (CI/CD) 🟡

- **列出 Action 執行任務**
  - `GET /repos/{owner}/{repo}/actions/tasks`
  - Custom: SDK 無支援，需自訂 HTTP 請求

## 總結

- 🟢 **完全支援 (5/7)**：Label、Milestone、Release、PR、Repo 管理
- 🔴 **部分支援 (1/7)**：Issue 功能（附件、依賴需自訂）
- 🟡 **需自訂實作 (2/7)**：Wiki、Forgejo Actions

**建議采用 Hybrid 模式**：約 71% 的功能可使用 SDK，剩餘功能自訂 HTTP 請求。
