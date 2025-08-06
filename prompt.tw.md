# 需求說明

我們要為專案裡每一個 subpackage 都加上 package document，以便 AI 助手可以透過 `go doc ./package_name` 指令快速理解，減少 AI 助手大量讀取完整檔案的機會

# 執行步驟

已經完成的步驟不需要重複執行

- [x] ./tools/
- [x] ./tools/action/
- [x] ./tools/issue/
- [x] ./tools/label/
- [x] ./tools/milestone/
- [x] ./tools/pullreq/
- [x] ./tools/release/
- [x] ./tools/repo/
- [x] ./tools/wiki/
- [x] ./types/

你必須逐項與我討論每一個 subpackage 的用途
  - 你應該先從檔名和路徑猜測用途
  - 與我討論確認
  - 將確定的討論結果寫入 subpackage 的 doc.go 以提供 package document
  - 更新上方的項目為完成。這是為了記錄我們的進度，確保能回復非預期的中斷
