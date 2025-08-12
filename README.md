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
- Manage release attachments

### Other Features
- View Pull Requests
- Manage Wiki pages
- View Forgejo/Gitea Actions tasks

## 📦 Installation

### Method 1: Use docker (Recommended)

For STDIO mode, you can skip to **Usage** section.

For SSE/Streamable HTTP mode, you should run `forgejo-mcp` as server before configuring your MCP client.

```bash
docker run -p 8080:8080 -e FORGEJOMCP_TOKEN="my-forgejo-api-token" ronmi/forgejo-mcp http --address :8080 --server https://git.example.com
```

### Method 2: Install from source

```bash
go install github.com/raohwork/forgejo-mcp@latest
```

### Method 3: Download Pre-compiled Binaries

Download the appropriate version for your operating system from the [Releases page](https://github.com/raohwork/forgejo-mcp/releases).

## 🖥️ Usage

This tool provides two primary modes of operation: `stdio` for local integration and `http` for remote access.

Before actually setup you MCP client, you have to create an access token on the Forgejo/Gitea server.

1. Log in to your Forgejo/Gitea instance
2. Go to **Settings** → **Applications** → **Access Tokens**
3. Click **Generate New Token**
4. Select appropriate permission scopes (recommend at least `repository` and `issue` write permissions)
5. Copy the generated token

💡 **Tip**: For security, consider setting environment variables instead of using tokens directly in config:
```bash
export FORGEJOMCP_SERVER="https://your-forgejo-instance.com"
export FORGEJOMCP_TOKEN="your_access_token"
```

### Stdio Mode (for Local Clients)

This is the recommended mode for integrating with local AI assistant clients like Claude Desktop or Gemini CLI. It uses standard input/output for direct communication.

#### Configure Your AI Client

Using docker:

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

Installed from source or pre-built binary:

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

You might want to take a look at **Security Recommendations** section for best practice.

### HTTP Server Mode (for Remote Access)

This mode starts a web server, allowing remote clients to connect via HTTP. It's ideal for web-based services or setting up a central gateway for multiple users.

Run the following command to start the server:
```bash
# with local binary
/path/to/forgejo-mcp http --address :8080 --server https://your-forgejo-instance.com

# with docker
docker run -p 8080:8080 -d --rm ronmi/forgejo-mcp http --address :8080 --server https://your-forgejo-instance.com
```

The server supports two operational modes:
- **Single-user mode**: If you provide a `--token` (or environment variable `FORGEJOMCP_TOKEN`) at startup, all operations will use that token.
  ```bash
  forgejo-mcp http --address :8080 --server https://git.example.com --token your_token
  ```
- **Multi-user mode**: If no token is provided, the server requires clients to send an `Authorization: Bearer <token>` header with each request, allowing it to serve multiple users securely.

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

or `http` type (for Streamable HTTP, use different path in URL)

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

1. **Use environment variables**: Set `FORGEJOMCP_SERVER` and `FORGEJOMCP_TOKEN`, then remove `--server` and `--token` from your configuration
2. **Limit token permissions**: Only grant necessary permission scopes
3. **Rotate tokens regularly**: Update access tokens periodically

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
