# 需求說明

我們要進行重構，目標是 `tools.ListActionTasksImpl`。

# 重構原因

- API resp `tools.MyActionTaskResponse` 跟原本規劃用來產生回應的型別 `types.ActionTask` (以及它的附屬型別) 有不少落差，導致轉換變得困難、多了無用欄位、對欄位的理解錯誤
- `types.ActionTask` 原本的設計方式，在現在的實作來說顯得落伍了
- `tools.ListActionTasksImpl` 目前的回應格式也因此不是十分有效，甚至會產生誤導

# 重構目標

- `types.ActionTask` 的定義改為與其他使用 SDK 的型別類似: 直接 embed API 回應的型別。
- 針對預估的用例 (用戶想要知道最近 CI 的成功率、最近哪次失敗、產生統計報表)，確保 `ToMarkdown` 方法揭露有用且必要的資訊

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [x] 開新 branch
- [x] 把 `tools.MyActionTaskResponse` (及會有連帶影響的型別) 移到 `types` 裡，避免後續修改造成 cyclic import
- [x] 修正移動定義造成的編譯錯誤 (修正引用)
- [x] reformat code
- [x] 人工審查
- [x] 把 `types.ActionTaskList` 改成 `type ActionTaskList struct { *MyActionTaskResponse }`
- [x] 拿掉不再使用的附屬型別
- [x] 決定必要的資訊有哪些，並配合修改測試案例 (測試方式維持原狀，只測資訊不測格式)
- [x] 實作 `ToMarkdown` 方法以通過測試
- [x] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成
