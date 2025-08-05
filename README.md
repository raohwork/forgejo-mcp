# Gitea/Forgejo MCP Server

> Turn AI into your code repository management assistant

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that enables you to manage Gitea/Forgejo repositories through AI assistants like Claude, Gemini, and Copilot.

## 🚀 Why Use Forgejo MCP Server?

If you want to:
- **Smart progress tracking**: Let AI help you track project progress and analyze bottlenecks
- **Automated issue categorization**: Automatically tag issue labels and set milestones based on content
- **Priority sorting**: Let AI analyze issue content to help prioritize tasks
- **Code review assistance**: Get AI suggestions and insights in Pull Requests
- **Project documentation organization**: Automatically organize Wiki documents and release notes

Then this tool is made for you!

## ✨ Supported Features

### Issue Management
- Create, edit, and view issues
- Add, remove, and replace labels
- Manage issue comments and attachments
- Set issue dependencies

### Project Organization
- Manage labels (create, edit, delete)
- Manage milestones (create, edit, delete)
- Repository search and listing

### Release Management
- Manage version releases
- Upload and manage release attachments

### Other Features
- View Pull Requests
- Manage Wiki pages
- View Forgejo Actions tasks

## 📦 Installation

### Method 1: Install with Go (Recommended)

```bash
go install github.com/raohwork/forgejo-mcp@latest
```

### Method 2: Download Pre-compiled Binaries

Download the appropriate version for your operating system from the [Releases page](https://github.com/raohwork/forgejo-mcp/releases).

## ⚙️ Configuration

### 1. Get Forgejo/Gitea Access Token

1. Log in to your Forgejo/Gitea instance
2. Go to **Settings** → **Applications** → **Access Tokens**
3. Click **Generate New Token**
4. Select appropriate permission scopes (recommend at least `repository` and `issue` write permissions)
5. Copy the generated token

### 2. Configure MCP Client

## 🖥️ Claude Desktop

### Windows

Edit `%APPDATA%\Claude\claude_desktop_config.json`:

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

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

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

Edit `~/.config/claude/claude_desktop_config.json`:

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

If you're using [Gemini CLI](https://github.com/google-gemini/gemini-cli), add this to your configuration file:

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

## 🛡️ Security Recommendations

1. **Use environment variables**: Avoid writing tokens directly in configuration files
   ```bash
   export FORGEJOMCP_SERVER="https://your-forgejo-instance.com"
   export FORGEJOMCP_TOKEN="your_access_token"
   ```
   
   Then remove the `--server` and `--token` parameters from your configuration.

2. **Limit token permissions**: Only grant necessary permission scopes

3. **Rotate tokens regularly**: Recommended to update access tokens periodically

## 📋 Usage Examples

After configuration, you can use natural language in your AI assistant to manage your repositories:

```
"Show me all issues labeled 'bug' with status 'open' in the user/myproject repository"

"Create a new issue in the user/myproject repository titled 'Fix login problem' with reproduction steps"

"Help me label issue #123 as 'urgent' and 'backend'"

"List all milestones in the user/myproject repository and tell me which ones are due soon"

"Analyze recent pull requests and tell me which ones need priority review"
```

## 🤝 Support & Contributing

- **Bug Reports**: [GitHub Issues](https://github.com/raohwork/forgejo-mcp/issues)
- **Code Contributions**: Pull Requests are welcome!

## 📄 License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

---

**Start making AI your code repository management partner!** 🚀
