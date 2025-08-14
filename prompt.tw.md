# 需求說明

根據一個常見的用例 `幫我看看最近一週的 issue 和 commit，為我分析一下這個專案的進度和開發效率` ，我認為我們必須增加 `list_repo_issues` 接受的參數

- 根據 Forgejo/Gitea API (及 SDK)，可以看到有四個參數我們沒有在 `ListRepoIssuesImpl` 裡提供
  * `since`: Only show items updated after the given time. This is a timestamp in RFC 3339 format
  * `before`: Only show items updated before the given time. This is a timestamp in RFC 3339 format
  * `created_by`: Only show items which were created by the given user (username string)
  * `mentioned_by`: Only show items in which the given user was mentioned (username string)
- 我們這次要加上的是 `since` 及 `before` 兩個參數。另外兩個我還沒遇到用例，所以先不加
- 不追加 commit 搜尋相關的工具
  * Forgejo/Gitea API 提供的 commit 歷史功能並不能完成這個需求
  * 多數情況下用戶應該能輕易透過 `git remote update` 取得伺服器上的 commit，讓 AI 助手可以從本地存取
  * 所以我們不應該做一個搜尋 commit 的工具

# 功能變更說明

- 修改 `tools/issue/crud.go`
  - 把 `Definition` 方法加上新的參數
  - 確認 `Handler` 方法有處理參數 (包括格式檢查)
  
這次的修改，理論上不應該動到任何測試案例

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] **準備工作：開新 branch**: branch name 必須符合 `refactor/short_desc` 的格式
- [x] **第一階段：分析現有 list_repo_issues 工具**
    - [x] 讀取 `tools/issue/crud.go` 檔案，了解現有的 `ListRepoIssuesImpl` 結構
    - [x] 檢查 `ListRepoIssuesParams` 結構中已有的參數
    - [x] 查看 `Definition()` 方法中的 JSON schema 定義
    - [x] 檢查 `Handler()` 方法中的參數處理邏輯
- [x] **第二階段：新增 since 和 before 參數**
    - [x] 在 `ListRepoIssuesParams` 結構中新增 `Since` 和 `Before` 欄位 (型別為 `*string`)
    - [x] 修改 `Definition()` 方法，在 JSON schema 中加入這兩個參數的定義
        - 參數型別：string (RFC 3339 格式)
        - 參數為 optional
        - 加上方便 AI 助手理解的描述文字
    - [x] 修改 `Handler()` 方法，處理新參數的解析和格式驗證
        - 檢查是否為有效的 RFC 3339 時間格式
        - 將參數傳遞給底層的 Forgejo SDK 呼叫
- [x] **第三階段：測試與驗證**
    - [x] 執行 `go build ./...` 確保編譯成功
    - [x] 執行 `go test ./...` 確保沒有破壞現有功能

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成。

# 意外中斷回復指南

如果開發過程中意外中斷，請按以下步驟回復：

1.  **確認當前狀態**: `git status` 以及 `go build ./...`
2.  **查看進度標記**：檢查此檔案 (`prompt.tw.md`) 中的進度標記。
3.  **重新開始的關鍵資訊**:
    - **主要修改檔案**: `tools/issue/list.go`
        - 需要修改 `ListRepoIssuesParams` 結構
        - 需要修改 `ListRepoIssuesImpl` 的 `Definition()` 和 `Handler()` 方法
    - **參數規格**:
        - `Since`: RFC 3339 格式時間字串，optional，只顯示在此時間之後更新的 issue
        - `Before`: RFC 3339 格式時間字串，optional，只顯示在此時間之前更新的 issue
    - **相關 API 文件**: 可用 `jq '.paths["/repos/{owner}/{repo}/issues"].get.parameters' swagger.v1.json` 查看完整參數清單

# 開發時必備知識

- **由編譯錯誤驅動開發**: `go build ./...` 是你的好朋友，它會告訴你修改是否完整、有沒有遺漏。
- **RFC 3339 時間格式**: Go 使用 `time.RFC3339` 常數，格式為 `2006-01-02T15:04:05Z07:00`
- **Forgejo SDK 參數對應**: 
  - SDK 中的 `ListIssueOption` 結構有 `Since` 和 `Before` 欄位 (型別為 `time.Time`)
  - 需要將字串參數解析為 `time.Time` 再傳遞給 SDK
- **JSON Schema 格式定義**: 時間參數在 schema 中應定義為 `"type": "string", "format": "date-time"`
- **錯誤處理**: 時間格式解析失敗時應回傳明確的錯誤訊息給使用者
