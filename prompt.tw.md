# 需求說明

我們要新增 issue blocking 相關的 MCP tool，包括了

- ListIssueBlockingImpl
- AddIssueBlockingImpl
- RemoveIssueBlockingImpl

這組工具與現有的 issue dependency 工具形成完整的雙向關係。

# 功能說明與語義差異

## Issue Dependencies vs Issue Blocking

**Issue Dependencies (現有)**：
- 查詢會阻塞當前 issue 的其他 issues
- 語義：`#123 depends on #234` → 必須先關閉 #234 才能關閉 #123
- API: `/repos/{owner}/{repo}/issues/{index}/dependencies`

**Issue Blocking (新增)**：
- 查詢被當前 issue 阻塞的其他 issues
- 語義：`#123 blocks #234` → 必須先關閉 #123 才能關閉 #234
- API: `/repos/{owner}/{repo}/issues/{index}/blocks`

## 實際用例對比

假設 issue #123 "修復登入 bug"，issue #234 "更新用戶介面"：

**Dependencies 查詢 #234**：
- 回傳：[#123]
- 含義：#234 被 #123 阻塞，必須先修復登入 bug

**Blocking 查詢 #123**：
- 回傳：[#234]
- 含義：#123 阻塞了 #234，修復登入後才能更新介面

# 技術實作原因

- 現有的 dependency 工具只能單向查詢，無法反向查詢哪些 issues 被當前 issue 阻塞
- Forgejo API 提供了對稱的 `/blocks` 端點，我們應該完整支援
- 使用者在專案管理時需要雙向視角：既要知道自己被誰阻塞，也要知道自己阻塞了誰
- 與現有 dependency 工具共用大部分程式碼，開發成本低但價值高

# API 端點分析

根據 `swagger.v1.json` 的定義：

## GET `/repos/{owner}/{repo}/issues/{index}/blocks`
- **功能**：List issues that are blocked by this issue
- **回傳**：`[]*forgejo.Issue` (與 dependencies 相同格式)
- **參數**：支援 page, limit 分頁參數

## POST `/repos/{owner}/{repo}/issues/{index}/blocks`
- **功能**：Block the issue given in the body by the issue in path
- **請求體**：`IssueMeta` (與 MyIssueMeta 相同結構)
- **回傳**：`*forgejo.Issue`

## DELETE `/repos/{owner}/{repo}/issues/{index}/blocks`
- **功能**：Unblock the issue given in the body by the issue in path
- **請求體**：`IssueMeta`
- **回傳**：`*forgejo.Issue`

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] 開新的 branch `feature/issue-blocking`
- [x] **第一階段：擴展類型定義**
    - [x] 在 `types/dependencies.go` 中新增 `IssueBlockingList` 類型 (基於 `[]*forgejo.Issue`)
    - [x] 為 `IssueBlockingList` 實作 `ToMarkdown()` 方法，格式與 `IssueDependencyList` 相同但空值訊息不同
    - [x] 檢查編譯，確保新類型定義正確
- [x] 人工審查
- [x] **第二階段：實作 Client 方法**
    - [x] 在 `tools/client_issue_dependencies.go` 新增三個方法：
      - `MyListIssueBlocking(owner, repo string, index int64) ([]*forgejo.Issue, error)`
      - `MyAddIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error)`
      - `MyRemoveIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error)`
    - [x] 檢查編譯，確保 client 方法實作正確
- [x] 人工審查
- [ ] **第三階段：實作 MCP 工具 - List**
    - [ ] 在 `tools/issue/dep.go` 新增 `ListIssueBlockingImpl` 結構及相關類型
    - [ ] 實作 `Definition()` 方法，工具名稱為 `list_issue_blocking`
    - [ ] 實作 `Handler()` 方法，輸出標頭為 `"## Issues blocked by #%d"`
    - [ ] 使用 `IssueBlockingList.ToMarkdown()` 進行格式化
    - [ ] 檢查編譯和基本功能測試
- [ ] 人工審查
- [ ] **第四階段：實作 MCP 工具 - Add/Remove**
    - [ ] 在 `tools/issue/dep.go` 新增 `AddIssueBlockingImpl` 和 `RemoveIssueBlockingImpl`
    - [ ] 實作對應的 `Definition()` 方法，工具名稱為 `add_issue_blocking` 和 `remove_issue_blocking`
    - [ ] 實作 `Handler()` 方法，使用 `EmptyResponse` 回傳簡潔訊息：
      - Add: `"Issue #%d now blocks issue #%d"`
      - Remove: `"Issue #%d no longer blocks issue #%d"`
    - [ ] 注意參數順序：blocking issue 在前，blocked issue 在後
    - [ ] 檢查編譯和基本功能測試
- [ ] 人工審查
- [ ] **第五階段：註冊工具到 MCP**
    - [ ] 在 `cmd/lib.go` 的工具註冊列表中新增三個新工具
    - [ ] 檢查完整編譯和工具可用性
- [ ] 人工審查
- [ ] **第六階段：測試與驗證**
    - [ ] 新增對應測試到 `types/misc_test.go` 測試 `IssueBlockingList.ToMarkdown()`
    - [ ] 把 `types/misc_test.go` 改名為 `types/dependencies_test.go`
    - [ ] 執行 `go test ./...` 確保所有測試通過
    - [ ] 執行 `go build ./...` 確保完整編譯成功
- [ ] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成

# 意外中斷回復指南

如果開發過程中意外中斷，請按以下步驟回復：

## 1. 確認當前狀態
```bash
git status                    # 檢查工作目錄狀態
git branch                    # 確認當前分支
git log --oneline -5          # 檢查最近提交記錄
```

## 2. 檢查編譯狀態
```bash
go build ./...                # 確保程式碼可編譯
go test ./...                 # 執行測試檢查狀態
```

## 3. 查看進度標記
檢查此檔案 (prompt.tw.md) 中的進度標記，確認已完成的步驟

## 4. 重新開始的關鍵資訊

### 類型定義參考 (types/dependencies.go)
```go
// 需新增的類型
type IssueBlockingList []*forgejo.Issue

func (ibl IssueBlockingList) ToMarkdown() string {
    if len(ibl) == 0 {
        return "*This issue is not blocking any other issues*"
    }

    markdown := ""
    for _, issue := range ibl {
        if issue == nil {
            continue
        }
        markdown += fmt.Sprintf("#%d **%s** (%s)\n", issue.Index, issue.Title, issue.State)
    }

    return markdown
}
```

### Client 方法簽名 (`tools/client_issue_dependencies.go`)
```go
func (c *Client) MyListIssueBlocking(owner, repo string, index int64) ([]*forgejo.Issue, error) {
    endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)
    var issues []*forgejo.Issue
    err := c.sendSimpleRequest("GET", endpoint, nil, &issues)
    return issues, err
}

func (c *Client) MyAddIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
    endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)
    var issue *forgejo.Issue
    err := c.sendSimpleRequest("POST", endpoint, blocked, &issue)
    return issue, err
}

func (c *Client) MyRemoveIssueBlocking(owner, repo string, index int64, blocked types.MyIssueMeta) (*forgejo.Issue, error) {
    endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/blocks", owner, repo, index)
    var issue *forgejo.Issue
    err := c.sendSimpleRequest("DELETE", endpoint, blocked, &issue)
    return issue, err
}
```

### MCP 工具註冊位置 (cmd/lib.go)
在 `tools.RegisterTool()` 呼叫列表中新增：
```go
tools.RegisterTool(toolRegistry, &issue.ListIssueBlockingImpl{Client: client})
tools.RegisterTool(toolRegistry, &issue.AddIssueBlockingImpl{Client: client})
tools.RegisterTool(toolRegistry, &issue.RemoveIssueBlockingImpl{Client: client})
```

### 輸出訊息格式
- **List**: `fmt.Sprintf("## Issues blocked by #%d\n\n%s", p.Index, blockingList.ToMarkdown())`
- **Add**: `fmt.Sprintf("Issue #%d now blocks issue #%d", p.Index, p.BlockedIndex)` (使用 `EmptyResponse`)
- **Remove**: `fmt.Sprintf("Issue #%d no longer blocks issue #%d", p.Index, p.BlockedIndex)` (使用 `EmptyResponse`)

### 測試檔案參考 (types/misc_test.go)
需新增 `TestIssueBlockingList_ToMarkdown` 測試，格式參考現有 `TestIssueDependencyList_ToMarkdown`

## 5. 常見問題排除
- **編譯錯誤**：檢查 import 路徑，確保 `types` package 正確引用
- **工具未註冊**：檢查 `cmd/lib.go` 中的註冊呼叫
- **API 呼叫失敗**：檢查 endpoint 路徑和 HTTP method 是否正確
- **測試失敗**：檢查輸出格式是否與預期一致

# 開發時必備知識

## Forgejo SDK 類型重用
- `*forgejo.Issue`: 完整的 issue 物件，包含 Index, Title, State 等欄位
- `MyIssueMeta`: 自定義類型，用於 API 請求，包含 Index, Owner, Name 欄位
- 兩種類型可以在同一個檔案中共存，各有不同用途

## MCP 工具模式
- `ToolImpl` interface: 需實作 `Definition()` 和 `Handler()` 方法
- `Definition()`: 定義工具元資訊，包含 name, description, parameters
- `Handler()`: 實際執行邏輯，回傳 `*mcp.CallToolResult`

## HTTP Client 模式
- `sendSimpleRequest()`: 用於 JSON API 呼叫
- GET 請求：paramObj 傳 nil
- POST/DELETE 請求：paramObj 傳請求體物件
- respObj 用於接收回應資料

## 語義對稱性設計原則
- Dependencies: "issues that block this issue" (被動視角)
- Blocking: "issues blocked by this issue" (主動視角)
- 兩者形成完整的雙向關係，在訊息措辭上保持對稱但意義相反
