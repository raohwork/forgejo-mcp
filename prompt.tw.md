# 需求說明

我們已經完成 custom forgejo api client 了，現在的目標是

- 把所有的 tool impl 都加上 `*tools.Client` 的依賴

```
type ListActionTasksImpl struct {
    Client *tools.Client
}
```

- 為每一個 tool impl 實作 handler (Handler 方法回傳的 ToolHandlerFor 的內容)

# 工作流程說明

依照之前的決定，這次的目標不需要寫測試，所以我們的流程會調整一下

1. 開新的 branch
2. 建立 progress.tw.md 文件 (內部控管用，不 commit)
  * 掃描 `./tools/*/*.go` 找到所有 tool impl
  * 內容是用 task list 列出每一個 tool impl (需要分類)
  * 分成以下幾類 (顯示名稱由你決定)
    - Issue 基本控制 (CRUD)
    - Issue Comments
    - Issue Labels + Attachments
    - Issue Dependencies
    - Repo 基本控制
    - Repo Label 管理
    - Milestones 管理
    - Releases CRUD + Attachments
    - Pull Requests 管理 + Actions
    - Wiki 管理
  * 每一個 tool impl 都要再列幾個待辦項目，如以下範例

```範例內容
- [ ] issue 基本控制:
  - [ ] tool 1: short tool description
    * [ ] 使用 API:
    * [ ] 對應的 types 格式:
    * [ ] struct 加上 API Client
    * [ ] 實作 handler
```

3. 第一次迴圈處理：檢查及記錄。你必須逐一處理每一個 tool 的前兩項待辦事項
   - 找到實作 handler 所需要的 api method (可能是使用 forgejo go sdk 裡的，例如 `ListRepoIssuesImpl.Client.ListIssues`；也可能是我們自訂的，例如 `ListActionTasksImpl.Client.MyListActionTasks`
   - 把 api method 記錄到第一個待辦事項，並標記成完成 (這需要更新 progress.tw.md)
   - 找到對應的 types 格式:
     * 對於使用 forgejo go sdk 的 impl，它應該會是一個 wrapper type，例如 `types.Issue`
     * 對於使用自訂 api 實作的 impl (API 回應的型別會是 MyXXXResponse)，它對應的應該會是一個 solid type。例如 `tools.MyActionTaskResponse` 是由 WorkflowRuns 欄位去對應到 `types.ActionTaskList` (slice of a solid type)
   - 把對應的 types 格式記錄到第二個待辦事項並標記完成 (這需要更新 progress.tw.md)
   - 這兩個資訊中的任何一個，若是找不到，應該立刻停止執行並詢問我該如何處理
4. 第二次迴圈處理：實作。你必須逐一處理每一個 tool 的後兩項待辦事項
   - 為 impl struct 加上 `Client *tools.Client` 欄位
   - 把第三個待辦事項標記成完成 (這需要更新 progress.tw.md)
   - 依照第一、第二個待辦事項的記錄，實作 ToolHandler (參考下一節 "其他資訊")
   - 把第四個待辦事項，以及這個 tool impl 本身都標記成完成 (這需要更新 progress.tw.md)
   - 如果整個分類裡的所有 impl 都已經完成，就執行以下步驟
     * 執行 goimports 來 format 程式碼
     * 修正編譯錯誤 (如果這需要修改 impl 以外的程式，必須中斷作業並詢問我該如何處理)
     * 再次執行 goimports 來 format 程式碼
     * commit (不 commit 文件)
     * 把分類標記為完成 (這需要更新 progress.tw.md)
   
# 其他指示

- custom client 已經提供 forgejo go sdk 的功能 (透過 embedded struct)，所有 tool impl 都依賴它呼叫 forgejo api
- 你可以使用 `go doc` 查看文件
- 我預期這次的作業不會建立新的程式碼檔案，所以你要建立新的 go 檔案之前，必需中斷作業，詢問我該怎麼辦
- ToolHandler 的輸出會是一個 `[]Content{*TextContent}`，內容可能是一小段總結 (有的 tool 可能不需要總結)，後面接著 types 的 ToMarkdown() 得到的格式。以下是幾個範例，但不要照抄：你應該使用對於用戶及 LLM 更友善的方式呈現，且應該使用英文

```範例1 issue 列表
這是符合條件的 50 個 issue (第 2/5 頁)

/* issuelist.ToMarkdown() 的內容放在這裡 */
```

```範例2 issue 列表
這是符合條件的 5 個 issue

/* issuelist.ToMarkdown() 的內容放在這裡 */
```

```範例3 指定 issue 資料
/* issue.ToMarkdown() 的內容放在這裡，不需要總結 */
```
