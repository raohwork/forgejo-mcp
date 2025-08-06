# 需求說明

我們要進行重構，目標是 wiki 相關的 mcp tool，包括了

- CreateWikiPageImpl
- EditWikiPageImpl
- GetWikiPageImpl
- ListWikiPagesImpl

# 重構原因

- 我們預期的格式 `types/wiki.go` 與實際上 forgejo sdk 提供的略有出入
- 重複的格式轉換程式碼太多，應該與其他 tool 一樣，把 sdk 回傳的用 wrapper 提供轉換

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] 開新的 branch
- [x] 把 `tools.MyWikiPage` 及相為的附屬型別，搬移到 `types` package
- [x] 修正編譯錯誤 (具體來說是配合引用路徑的變更進行修改)
- [ ] 人工審查
- [ ] 把 `types.WikiPage` 改成 `type WikiPage struct { *MyWikiPage }`
- [ ] 修正對應的 ToMarkdown 方法以維持原本的輸出格式，確保通過測試
- [ ] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成

