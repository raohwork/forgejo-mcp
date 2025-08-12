# Gitea/Forgejo MCP Server

> 讓 AI 成為你的程式碼倉庫管理助手

一個 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 伺服器，讓你可以透過 AI 助手（如 Claude、Gemini、Copilot）來管理 Gitea/Forgejo 倉庫。

## 🚀 為什麼需要 Forgejo MCP Server？

如果你想要：
- **智慧進度管控**：讓 AI 幫你追蹤專案進度、分析瓶頸
- **自動議題分類**：根據內容自動標記議題標籤、設定里程碑  
- **緊急程度排序**：讓 AI 分析議題內容，協助排定優先順序
- **程式碼審查協助**：在 Pull Request 中獲得 AI 的建議和洞察
- **專案資料整理**：自動整理 Wiki 文件、發布說明

那麼這個工具就是為你準備的！

## ✨ 支援功能

### 議題管理
- 建立、編輯、查看議題
- 新增、移除、替換標籤  
- 管理議題評論和附件
- 設定議題相依關係

### 專案組織
- 管理標籤（建立、編輯、刪除）
- 管理里程碑（建立、編輯、刪除）
- 倉庫搜尋和列表

### 發布管理
- 管理版本發布
- 管理發布附件

### 其他功能
- 查看 Pull Request
- 管理 Wiki 頁面
- 查看 Forgejo/Gitea Actions 任務

## 📦 安裝

### 方法一：使用 docker（推薦）

對於 STDIO 模式，你可以直接跳到 **使用方式** 部分。

對於 SSE/Streamable HTTP 模式，你應該在設定 MCP 客戶端之前先執行 `forgejo-mcp` 作為伺服器。

```bash
docker run -p 8080:8080 -e FORGEJOMCP_TOKEN="my-forgejo-api-token" ronmi/forgejo-mcp http --address :8080 --server https://git.example.com
```

### 方法二：從原始碼安裝

```bash
go install github.com/raohwork/forgejo-mcp@latest
```

### 方法三：下載預編譯版本

從 [Releases 頁面](https://github.com/raohwork/forgejo-mcp/releases) 下載適合你作業系統的版本。

## 🖥️ 使用方式

此工具提供兩種主要操作模式：`stdio` 用於本機整合，`http` 用於遠端存取。

在實際設定 MCP 客戶端之前，你必須先在 Forgejo/Gitea 伺服器上建立存取權杖。

1. 登入你的 Forgejo/Gitea 實例
2. 前往 **設定** → **應用程式** → **存取權杖**
3. 點擊 **產生新權杖**
4. 選擇適當的權限範圍（建議至少 `repository` 和 `issue` 的寫入權限）
5. 複製產生的權杖

💡 **提示**：為了安全起見，建議使用環境變數而不是直接在設定中使用權杖：
```bash
export FORGEJOMCP_SERVER="https://your-forgejo-instance.com"
export FORGEJOMCP_TOKEN="your_access_token"
```

### Stdio 模式（適用於本機客戶端）

這是與 Claude Desktop 或 Gemini CLI 等本機 AI 助理客戶端整合的建議模式。它使用標準輸入/輸出進行直接通訊。

#### 設定你的 AI 客戶端

使用 docker：

```json
{
  "mcpServers": {
    "forgejo": {
      "command": "docker",
      "args": [
        "--rm",
        "ronmi/forgejo-mcp",
        "stdio",
        "--server", "https://your-forgejo-instance.com",
        "--token", "your_access_token"
      ]
    }
  }
}
```

從原始碼安裝或預編譯版本：

```json
{
  "mcpServers": {
    "forgejo": {
      "command": "/path/to/forgejo-mcp",
      "args": [
        "stdio",
        "--server", "https://your-forgejo-instance.com",
        "--token", "your_access_token"
      ]
    }
  }
}
```

你可能會想要參考 **安全性建議** 部分的最佳實踐。

### HTTP 伺服器模式（適用於遠端存取）

此模式會啟動一個 HTTP 伺服器，允許遠端客戶端透過 HTTP 連線。它非常適合網頁服務或為多個使用者設定中央閘道。

執行以下指令以啟動伺服器：
```bash
# 使用本機執行檔
/path/to/forgejo-mcp http --address :8080 --server https://your-forgejo-instance.com

# 使用 docker
docker run -p 8080:8080 -d --rm ronmi/forgejo-mcp http --address :8080 --server https://your-forgejo-instance.com
```

伺服器支援兩種運作模式：
- **單使用者模式**：如果在啟動時提供 `--token`（或環境變數 `FORGEJOMCP_TOKEN`），所有操作都將使用該權杖。
  ```bash
  forgejo-mcp http --address :8080 --server https://git.example.com --token your_token
  ```
- **多使用者模式**：如果未提供 `--token`，伺服器會要求客戶端在每個請求中發送 `Authorization: Bearer <token>` 標頭，從而安全地為多個使用者提供服務。

#### 客戶端設定

對於支援透過 HTTP 連線到遠端 MCP 伺服器的客戶端，你可以新增如下設定。此範例展示如何連線到以多使用者模式運行的伺服器：

```json
{
  "mcpServers": {
    "forgejo-remote": {
      "type": "sse",
      "url": "http://localhost:8080/sse",
      "headers": {
        "Authorization": "Bearer your_token"
      }
    }
  }
}
```

或 `http` 類型（適用於 Streamable HTTP，使用不同的 URL 路徑）

```json
{
  "mcpServers": {
    "forgejo-remote": {
      "type": "http",
      "url": "http://localhost:8080/",
      "headers": {
        "Authorization": "Bearer your_token"
      }
    }
  }
}
```

如果連線到單使用者模式的伺服器，你可以省略 `headers` 欄位。

## 🛡️ 安全性建議

1. **使用環境變數**：設定 `FORGEJOMCP_SERVER` 和 `FORGEJOMCP_TOKEN`，然後從設定中移除 `--server` 和 `--token`

2. **限制權杖權限**：只給予必要的權限範圍

3. **定期輪換權杖**：定期更新存取權杖

## 📋 使用範例

設定完成後，你就可以在 AI 助手中使用自然語言來管理你的倉庫了：

```
「顯示我的 gitea 伺服器上這個 repo 的嚴重錯誤報告」

「根據我們上面的討論，建立一個關於此錯誤的詳細議題，然後在議題上留言說明我們將如何修復它。」

「給我一份關於當前里程碑的報告，特別是最近的進展。」

「分析最近的 pull request，告訴我哪些需要優先審查」
```

## 🤝 支援與貢獻

- **問題回報**：[GitHub Issues](https://github.com/raohwork/forgejo-mcp/issues)
- **貢獻程式碼**：歡迎提交 Pull Request！

## 📄 授權

本專案採用 [Mozilla Public License 2.0](LICENSE) 授權。
