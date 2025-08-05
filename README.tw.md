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
- 上傳和管理發布附件

### 其他功能
- 查看 Pull Request
- 管理 Wiki 頁面
- 查看 Forgejo Actions 任務

## 📦 安裝

### 方法一：使用 Go 安裝（推薦）

```bash
go install github.com/raohwork/forgejo-mcp@latest
```

### 方法二：下載預編譯版本

從 [Releases 頁面](https://github.com/raohwork/forgejo-mcp/releases) 下載適合你作業系統的版本。

## ⚙️ 設定

### 1. 取得 Forgejo/Gitea 存取權杖

1. 登入你的 Forgejo/Gitea 實例
2. 進入 **設定** → **應用程式** → **存取權杖**
3. 點擊 **產生新權杖**
4. 選擇適當的權限範圍（建議至少包含 `repository` 和 `issue` 的寫入權限）
5. 複製產生的權杖

### 2. 設定 MCP 客戶端

## 🖥️ Claude Desktop

### Windows

編輯 `%APPDATA%\Claude\claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "forgejo": {
      "command": "forgejo-mcp",
      "args": [
        "stdio",
        "--server", "https://your-forgejo-instance.com",
        "--token", "your_access_token"
      ]
    }
  }
}
```

### macOS

編輯 `~/Library/Application Support/Claude/claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "forgejo": {
      "command": "forgejo-mcp",
      "args": [
        "stdio", 
        "--server", "https://your-forgejo-instance.com",
        "--token", "your_access_token"
      ]
    }
  }
}
```

### Linux

編輯 `~/.config/claude/claude_desktop_config.json`：

```json
{
  "mcpServers": {
    "forgejo": {
      "command": "forgejo-mcp",
      "args": [
        "stdio",
        "--server", "https://your-forgejo-instance.com", 
        "--token", "your_access_token"
      ]
    }
  }
}
```

## 💎 Gemini CLI

如果你使用 [Gemini CLI](https://github.com/google-gemini/gemini-cli)，請在你的設定檔中新增：

```yaml
mcp_servers:
  forgejo:
    command: forgejo-mcp
    args:
      - stdio  
      - --server
      - https://your-forgejo-instance.com
      - --token
      - your_access_token
```

## 🛡️ 安全性建議

1. **使用環境變數**：避免在設定檔中直接寫入權杖
   ```bash
   export FORGEJOMCP_SERVER="https://your-forgejo-instance.com"
   export FORGEJOMCP_TOKEN="your_access_token"
   ```
   
   然後在設定中移除 `--server` 和 `--token` 參數。

2. **限制權杖權限**：只給予必要的權限範圍

3. **定期輪換權杖**：建議定期更新存取權杖

## 📋 使用範例

設定完成後，你就可以在 AI 助手中使用自然語言來管理你的倉庫了：

```
「幫我查看 user/myproject 倉庫中所有標記為 'bug' 且狀態為 'open' 的議題」

「在 user/myproject 倉庫建立一個新議題，標題是『修正登入問題』，內容包含重現步驟」

「幫我將議題 #123 標記為 'urgent' 和 'backend'」

「列出 user/myproject 倉庫的所有里程碑，並告訴我哪些即將到期」

「分析最近的 pull request，告訴我哪些需要優先審查」
```

## 🤝 支援與貢獻

- **問題回報**：[GitHub Issues](https://github.com/raohwork/forgejo-mcp/issues)
- **貢獻程式碼**：歡迎提交 Pull Request！

## 📄 授權

本專案採用 [Mozilla Public License 2.0](LICENSE) 授權。
