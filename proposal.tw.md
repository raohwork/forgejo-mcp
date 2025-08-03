# Forgejo MCP Server 需求規格

## 專案目標

製作一個 MCP (Model Context Protocol) server，讓 Claude Code 可以協助管理 Gitea/Forgejo server 上的 repository。

### 核心需求

1. **通訊協定**
   - 支援 stdio 模式與 Claude Code 整合
   - 使用官方 Go SDK (github.com/modelcontextprotocol/go-sdk)

2. **功能範圍**
   - 參見 `features.tw.md`

3. **配置管理**
   - CLI 參數：`--server` (Forgejo server URL)、`--token` (access token)
   - 環境變數：`FORGEJOMCP_SERVER`、`FORGEJOMCP_TOKEN`
   - 配置優先順序：CLI args > 環境變數

### 技術規格

- **程式語言**：Go
- **主要依賴**：
  * 官方 MCP Go SDK (github.com/modelcontextprotocol/go-sdk)
  * Forgejo SDK (codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2)
  * CLI 框架 (github.com/spf13/cobra)
  * 配置管理 (github.com/spf13/viper)
- **回傳格式**:
  * 錯誤以 plaintext 格式回傳
  * 正常的回應以 markdown 格式回傳

## 開發原則

- 測試驅動開發 (TDD)
- 敏捷開發，從 MVP 開始
- 每個功能獨立開發和測試
- 遵循 pair programming 精神

## 文件、註解與訊息

考慮到這是開源專案，註解及訊息 (如 logging) 必須使用英文。

文件檔案如果沒有標註語言，則使用英文。如果有標註 tw (例如 xxx.tw.md) 則使用台灣中文。
