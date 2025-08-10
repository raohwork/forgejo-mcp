# 需求說明

根據 issue #1 的決議，我們需要 **移除** issue 和 release 的附件上傳工具 (`create_issue_attachment`, `create_release_attachment`)。

# 功能變更說明

## 原有流程 (即將移除)

AI 助手透過 `create_issue_attachment` 或 `create_release_attachment` 工具直接上傳檔案。這需要 AI 讀取完整的檔案內容，或者依賴 MCP server 存取本地檔案系統。

## 全新流程

當使用者需要上傳附件時，AI 助手將 **不再呼叫工具**。取而代之的是，AI 助手會根據當前的上下文 (issue 或 release)，動態生成一段可執行的 `curl` 命令，並提供給使用者。使用者只需複製這段命令到自己的終端機中執行，即可完成上傳。

這個變更將 AI 的職責限定在「生成指令」，而將「執行指令」的責任交還給使用者，從而解決了底層的技術複雜性。

# 技術實作原因

- **高昂的 Token 消耗**：舊方法若讓 AI 讀取檔案內容，會大量消耗寶貴的 context token，尤其是在處理大型檔案時。
- **檔案系統存取限制**：讓 MCP server 透過檔案路徑存取使用者本地檔案，在不同執行模式下會遇到權限和路徑問題。
- **職責分離**：新流程明確區分了 AI 和使用者的職責，讓架構更清晰、更健壯。

# API 端點分析 (供 AI 生成 curl 命令參考)

AI 助手在 **未來* 需要利用以下資訊來組合出給使用者的 `curl` 指令。

## Issue 附件上傳

- **Endpoint**: `POST /api/v1/repos/{owner}/{repo}/issues/{index}/assets`
- **Content-Type**: `multipart/form-data`
- **curl 命令樣板**:
  ```bash
  curl -X POST "${FORGEJO_SERVER}/api/v1/repos/{owner}/{repo}/issues/{index}/assets" \
    -H "Authorization: token ${FORGEJOMCP_TOKEN}" \
    -F "attachment=@/absolute/path/to/file" \
    -F "name=optional-display-name"
  ```

## Release 附件上傳

- **Endpoint**: `POST /api/v1/repos/{owner}/{repo}/releases/{id}/assets`
- **Content-Type**: `multipart/form-data` 或 `application/octet-stream`
- **curl 命令樣板**:
  ```bash
  curl -X POST "${FORGEJO_SERVER}/api/v1/repos/{owner}/{repo}/releases/{id}/assets" \
    -H "Authorization: token ${FORGEJOMCP_TOKEN}" \
    -F "attachment=@/absolute/path/to/file" \
    -F "name=optional-display-name"
  ```

# 執行步驟

具體的作業步驟如下: (已完成的步驟不需重複確認)

- [ ] **準備工作：建立新的 branch**
- [x] **第一階段：移除 Issue Attachment 建立工具**
    - [x] 讀取 `tools/issue/attach.go`。
    - [x] 從檔案中移除 `CreateIssueAttachmentParams` struct, `CreateIssueAttachmentImpl` struct, 以及它對應的 `Definition()` 和 `Handler()` 方法。
    - [x] **注意**：請勿移除檔案中其他 (如 `List`, `Delete`, `Edit`) 的工具實作。
- [ ] 人工審查
- [x] **第二階段：移除 Release Attachment 建立工具**
    - [x] 讀取 `tools/release/attach.go`。
    - [x] 從檔案中移除 `CreateReleaseAttachmentParams` struct, `CreateReleaseAttachmentImpl` struct, 以及它對應的 `Definition()` 和 `Handler()` 方法。
    - [x] **注意**：請勿移除檔案中其他 (如 `List`, `Delete`, `Edit`) 的工具實作。
- [ ] 人工審查
- [x] **第三階段：移除 Client 的附件上傳方法**
    - [x] 讀取 `tools/client_issue_attachments.go`。
    - [x] 從檔案中移除 `MyCreateIssueAttachment` 函式。
    - [x] **注意**：請勿移除檔案中其他 (如 `List`, `Delete`, `Edit`) 的 client 方法。
- [ ] 人工審查
- [ ] **第四階段：移除工具註冊**
    - [ ] 讀取 `cmd/lib.go` 檔案。
    - [ ] 在工具註冊列表中，找到並移除 `issue.CreateIssueAttachmentImpl` 和 `release.CreateReleaseAttachmentImpl` 的註冊程式碼。
    - [ ] 執行 `go build ./...`，預期此時應該能成功編譯。
- [ ] 人工審查
- [ ] **第五階段：清理測試案例與最終確認**
    - [ ] 檢查 `tools/client_test.go`，移除與 `MyCreateIssueAttachment` 相關的測試程式碼。
    - [ ] 執行 `go test ./...` 確保所有測試仍然通過。
    - [ ] 執行 `go mod tidy` 整理相依性。
- [ ] 人工審查

你 **必須** 在完成每一個步驟之後，進入下一個步驟之前，先更新這個檔案，把對應的步驟標記為完成。

# 意外中斷回復指南

如果開發過程中意外中斷，請按以下步驟回復：

1.  **確認當前狀態**: `git status` 以及 `go build ./...`
2.  **查看進度標記**：檢查此檔案 (`prompt.tw.md`) 中的進度標記。
3.  **重新開始的關鍵資訊**:
    - **要修改的檔案**:
        - `tools/issue/attach.go` (移除 Create...)
        - `tools/release/attach.go` (移除 Create...)
        - `tools/client_issue_attachments.go` (移除 MyCreateIssueAttachment)
        - `cmd/lib.go` (移除工具註冊)
        - `tools/client_test.go` (移除相關測試)
    - **不應刪除的檔案**: `types/attachments.go` 仍然被其他功能使用。

# 開發時必備知識

- **精確修改**: 本次任務是外科手術式的程式碼移除，不是刪除整個檔案。需仔細辨識要移除的 struct 和 function。
- **由編譯錯誤驅動開發**: `go build ./...` 是你的好朋友，它會告訴你修改是否完整、有沒有遺漏。
