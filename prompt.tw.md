# 需求說明

我們要進行重構，目標是 issue dependency 相關的 mcp tool，包括了

- AddIssueDependencyImpl
- ListIssueDependenciesImpl
- RemoveIssueDependencyImpl

# 重構原因

- 自訂型別 `MyIssueMeta` 被定義在 `tools` package，應移至 `types` package
- Client 方法回傳的是未經處理的 forgejo sdk 型別，導致 handler 中出現重複的格式轉換程式碼
- `AddIssueDependencyImpl` 和 `RemoveIssueDependencyImpl` 的輸出過於冗餘，應回傳簡潔的成功訊息

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] 開新的 branch
- [ ] **第一階段：型別分離**
    - [x] 在 `types` package 中建立一個新檔案 `dependencies.go`
    - [x] 將 `MyIssueMeta` struct 從 `tools/client_issue_dependencies.go` 完整搬移到 `types/dependencies.go`
    - [x] 修正編譯錯誤 (具體來說是配合引用路徑的變更進行修改)
- [ ] 人工審查
- [ ] **第二階段：集中格式化邏輯與簡化輸出**
    - [ ] 在 `types/dependencies.go` 中，定義一個新的 wrapper type `IssueDependencyList`，其本質為 `[]*forgejo.Issue`
    - [ ] 為 `IssueDependencyList` 實作 `ToMarkdown()` 方法，封裝列表格式化的邏輯
    - [ ] 重構 `ListIssueDependenciesImpl` 的 Handler，改為使用 `IssueDependencyList.ToMarkdown()`
    - [ ] 重構 `AddIssueDependencyImpl` 和 `RemoveIssueDependencyImpl` 的 Handler，使其回傳簡潔的成功訊息字串，而不是完整的 issue 內容
- [ ] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成

