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

## 🖥️ Usage

This tool provides two primary modes of operation: `stdio` for local integration and `http` for remote access.

### Stdio Mode (for Local Clients)

This is the recommended mode for integrating with local AI assistant clients like Claude Desktop or Gemini CLI. It uses standard input/output for direct communication.

#### 1. Get Forgejo/Gitea Access Token

1. Log in to your Forgejo/Gitea instance
2. Go to **Settings** → **Applications** → **Access Tokens**
3. Click **Generate New Token**
4. Select appropriate permission scopes (recommend at least `repository` and `issue` write permissions)
5. Copy the generated token

#### 2. Configure Your AI Client

##### Claude Desktop

- **Windows**: Edit `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: Edit `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: Edit `~/.config/claude/claude_desktop_config.json`

Add the following configuration:
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

##### Gemini CLI

If you're using [Gemini CLI](https://github.com/google-gemini/gemini-cli), add this to your configuration file:

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

### HTTP Server Mode (for Remote Access)

This mode starts a web server, allowing remote clients to connect via HTTP. It's ideal for web-based services or setting up a central gateway for multiple users.

Run the following command to start the server:
```bash
forgejo-mcp http --address :8080 --server https://your-forgejo-instance.com
```

The server supports two operational modes:
- **Single-user mode**: If you provide a `--token` at startup, all operations will use that token.
  ```bash
  forgejo-mcp http --address :8080 --server https://git.example.com --token your_token
  ```
- **Multi-user mode**: If no `--token` is provided, the server requires clients to send an `Authorization: Bearer <token>` header with each request, allowing it to serve multiple users securely.

#### Client Configuration

For clients that support connecting to a remote MCP server via HTTP, you can add a configuration like this. This example shows how to connect to a server running in multi-user mode:

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

or `http` type (different URL)

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

If connecting to a server in single-user mode, you can omit the `headers` field.

## 🛡️ Security Recommendations

1. **Use environment variables**: Avoid writing tokens directly in configuration files
   ```bash
   export FORGEJOMCP_SERVER="https://your-forgejo-instance.com"
   export FORGEJOMCP_TOKEN="your_access_token"
   ```
   
   Then remove the `--server` and `--token` parameters from your configuration.
   
   For sse/http type, update your config:

   ```json
   {
     "headers": {
       "Authorization": "Bearer ${FORGEJOMCP_TOKEN}"
     }
   }
   ```

2. **Limit token permissions**: Only grant necessary permission scopes

3. **Rotate tokens regularly**: Recommended to update access tokens periodically

## 📋 Usage Examples

After configuration, you can use natural language in your AI assistant to manage your repositories:

```
"Show me critical bug reports of this repo on my gitea server"

"According to our discussion above, create a detailed issue about this bug, then leave a comment on the issue to describe how we will fix it."

"Give me a report about current milestone. Recent progression in particular."

"Analyze recent pull requests and tell me which ones need priority review"
```

## 🤝 Support & Contributing

- **Bug Reports**: [GitHub Issues](https://github.com/raohwork/forgejo-mcp/issues)
- **Code Contributions**: Pull Requests are welcome!

## 📄 License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

---

**Start making AI your code repository management partner!** 🚀
