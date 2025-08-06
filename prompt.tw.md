# 需求說明

我們要進行重構，目標是 issue dependency 相關的 mcp tool，包括了

- AddIssueDependencyImpl
- ListIssueDependenciesImpl
- RemoveIssueDependencyImpl

# 重構原因

- 自訂型別 `MyIssueMeta` 被定義在 `tools` package，應移至 `types` package
- Client 方法回傳的是未經處理的 forgejo sdk 型別，導致 handler 中出現重複的格式轉換程式碼
- `AddIssueDependencyImpl` 和 `RemoveIssueDependencyImpl` 的輸出過於冗餘，應回傳簡潔的成功訊息
- `ListIssueDependenciesImpl` 需要更加精簡的輸出
  - 根據用例，每一個 issue 只要輸出 `#123 **title** (state)` 就夠了
- 這三個工具容易誤解，所以需要更精確的說明 (原始碼的註解文件、features*.md 中英文文件、tool definition)
  - 依照 Forgejo API 的定義 deps 代表 issues that blocks current issue
  - 例如我們針對 issue#123 查詢，若回傳了 #234, #235 ，就代表這兩個 issue 必須先關閉，才能關閉 #123

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] 開新的 branch
- [ ] **第一階段：型別分離**
    - [x] 在 `types` package 中建立一個新檔案 `dependencies.go`
    - [x] 將 `MyIssueMeta` struct 從 `tools/client_issue_dependencies.go` 完整搬移到 `types/dependencies.go`
    - [x] 修正編譯錯誤 (具體來說是配合引用路徑的變更進行修改)
- [ ] 人工審查
- [x] **第二階段：集中格式化邏輯與簡化輸出 (List)**
    - [x] 在 `types/dependencies.go` 中，定義一個新的 wrapper type `IssueDependencyList`，其本質為 `[]*forgejo.Issue`，同時移除 types/issues.go 裡的現有型別 (及附屬型別、方法)
    - [x] 為 `IssueDependencyList` 實作 `ToMarkdown()` 方法，封裝列表格式化的邏輯
    - [x] 重構測試，確保
        - 使用相同的邏輯 (只測資訊，不測格式) 測試
        - 新的測試應該符合我們重構的目的： Add 和 Remove 回應簡潔的成功或失敗訊息，List 回傳必要的 issue 資訊就好
    - [x] 重構 `ListIssueDependenciesImpl` 的 Handler，改為使用 `IssueDependencyList.ToMarkdown()`
- [ ] 人工審查
- [x] **第三階段：簡化輸出 (Add/Remove)**
    - [x] 修改對應的測試，確保輸出只有簡潔的成功或錯誤訊息 (無測試檔案，EmptyResponse 測試已存在)
    - [x] 重構 `AddIssueDependencyImpl` 和 `RemoveIssueDependencyImpl` 的 Handler，使其回傳簡潔的成功訊息字串，而不是完整的 issue 內容
- [ ] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成

