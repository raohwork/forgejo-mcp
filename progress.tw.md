# Tool Implementation Progress

## Issue 基本控制 (CRUD) ✅
- [x] ListRepoIssuesImpl: List repository issues
  * [x] 使用 API: Client.ListRepoIssues
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] GetIssueImpl: Get specific issue details
  * [x] 使用 API: Client.GetIssue
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateIssueImpl: Create new issue
  * [x] 使用 API: Client.CreateIssue
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditIssueImpl: Edit existing issue
  * [x] 使用 API: Client.EditIssue
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Issue Comments ✅
- [x] ListIssueCommentsImpl: List issue comments
  * [x] 使用 API: Client.ListIssueComments
  * [x] 對應的 types 格式: types.Comment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateIssueCommentImpl: Create issue comment
  * [x] 使用 API: Client.CreateIssueComment
  * [x] 對應的 types 格式: types.Comment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditIssueCommentImpl: Edit issue comment
  * [x] 使用 API: Client.EditIssueComment
  * [x] 對應的 types 格式: types.Comment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteIssueCommentImpl: Delete issue comment
  * [x] 使用 API: Client.DeleteIssueComment
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Issue Labels + Attachments ✅
- [x] AddIssueLabelsImpl: Add labels to issue
  * [x] 使用 API: Client.AddIssueLabels
  * [x] 對應的 types 格式: types.Label
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] RemoveIssueLabelImpl: Remove label from issue
  * [x] 使用 API: Client.DeleteIssueLabel
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ReplaceIssueLabelsImpl: Replace issue labels
  * [x] 使用 API: Client.ReplaceIssueLabels
  * [x] 對應的 types 格式: types.Label
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ListIssueAttachmentsImpl: List issue attachments
  * [x] 使用 API: Client.MyListIssueAttachments
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateIssueAttachmentImpl: Create issue attachment
  * [x] 使用 API: Client.MyCreateIssueAttachment
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditIssueAttachmentImpl: Edit issue attachment
  * [x] 使用 API: Client.MyEditIssueAttachment
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteIssueAttachmentImpl: Delete issue attachment
  * [x] 使用 API: Client.MyDeleteIssueAttachment
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Issue Dependencies ✅
- [x] ListIssueDependenciesImpl: List issue dependencies
  * [x] 使用 API: Client.MyListIssueDependencies
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] AddIssueDependencyImpl: Add issue dependency
  * [x] 使用 API: Client.MyAddIssueDependency
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] RemoveIssueDependencyImpl: Remove issue dependency
  * [x] 使用 API: Client.MyRemoveIssueDependency
  * [x] 對應的 types 格式: types.Issue
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Repo 基本控制 ✅
- [x] SearchRepositoriesImpl: Search repositories
  * [x] 使用 API: Client.SearchRepos
  * [x] 對應的 types 格式: types.Repository
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ListMyRepositoriesImpl: List user repositories
  * [x] 使用 API: Client.ListMyRepos
  * [x] 對應的 types 格式: types.Repository
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ListOrgRepositoriesImpl: List organization repositories
  * [x] 使用 API: Client.ListOrgRepos
  * [x] 對應的 types 格式: types.Repository
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] GetRepositoryImpl: Get repository details
  * [x] 使用 API: Client.GetRepo
  * [x] 對應的 types 格式: types.Repository
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Repo Label 管理 ✅
- [x] ListRepoLabelsImpl: List repository labels
  * [x] 使用 API: Client.ListRepoLabels
  * [x] 對應的 types 格式: types.Label
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateLabelImpl: Create repository label
  * [x] 使用 API: Client.CreateLabel
  * [x] 對應的 types 格式: types.Label
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditLabelImpl: Edit repository label
  * [x] 使用 API: Client.EditLabel
  * [x] 對應的 types 格式: types.Label
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteLabelImpl: Delete repository label
  * [x] 使用 API: Client.DeleteLabel
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Milestones 管理 ✅
- [x] ListRepoMilestonesImpl: List repository milestones
  * [x] 使用 API: Client.ListRepoMilestones
  * [x] 對應的 types 格式: types.Milestone
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateMilestoneImpl: Create milestone
  * [x] 使用 API: Client.CreateMilestone
  * [x] 對應的 types 格式: types.Milestone
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditMilestoneImpl: Edit milestone
  * [x] 使用 API: Client.EditMilestone
  * [x] 對應的 types 格式: types.Milestone
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteMilestoneImpl: Delete milestone
  * [x] 使用 API: Client.DeleteMilestone
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Releases CRUD + Attachments ✅
- [x] ListReleasesImpl: List repository releases
  * [x] 使用 API: Client.ListReleases
  * [x] 對應的 types 格式: types.Release
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateReleaseImpl: Create release
  * [x] 使用 API: Client.CreateRelease
  * [x] 對應的 types 格式: types.Release
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditReleaseImpl: Edit release
  * [x] 使用 API: Client.EditRelease
  * [x] 對應的 types 格式: types.Release
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteReleaseImpl: Delete release
  * [x] 使用 API: Client.DeleteRelease
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ListReleaseAttachmentsImpl: List release attachments
  * [x] 使用 API: Client.ListReleaseAttachments
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateReleaseAttachmentImpl: Create release attachment
  * [x] 使用 API: Client.CreateReleaseAttachment
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditReleaseAttachmentImpl: Edit release attachment
  * [x] 使用 API: Client.EditReleaseAttachment
  * [x] 對應的 types 格式: types.Attachment
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteReleaseAttachmentImpl: Delete release attachment
  * [x] 使用 API: Client.DeleteReleaseAttachment
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Pull Requests 管理 + Actions ✅
- [x] ListPullRequestsImpl: List pull requests
  * [x] 使用 API: Client.ListRepoPullRequests
  * [x] 對應的 types 格式: types.PullRequest
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] GetPullRequestImpl: Get pull request details
  * [x] 使用 API: Client.GetPullRequest
  * [x] 對應的 types 格式: types.PullRequest
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] ListActionTasksImpl: List action tasks
  * [x] 使用 API: Client.MyListActionTasks
  * [x] 對應的 types 格式: types.ActionTask
  * [x] struct 加上 API Client
  * [x] 實作 handler

## Wiki 管理 ✅
- [x] ListWikiPagesImpl: List wiki pages
  * [x] 使用 API: Client.MyListWikiPages
  * [x] 對應的 types 格式: types.WikiPage
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] GetWikiPageImpl: Get wiki page content
  * [x] 使用 API: Client.MyGetWikiPage
  * [x] 對應的 types 格式: types.WikiPage
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] CreateWikiPageImpl: Create wiki page
  * [x] 使用 API: Client.MyCreateWikiPage
  * [x] 小應的 types 格式: types.WikiPage
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] EditWikiPageImpl: Edit wiki page
  * [x] 使用 API: Client.MyEditWikiPage
  * [x] 對應的 types 格式: types.WikiPage
  * [x] struct 加上 API Client
  * [x] 實作 handler
- [x] DeleteWikiPageImpl: Delete wiki page
  * [x] 使用 API: Client.MyDeleteWikiPage
  * [x] 對應的 types 格式: none (deletion)
  * [x] struct 加上 API Client
  * [x] 實作 handler