# Forgejo MCP Server

A Model Context Protocol (MCP) server that enables MCP clients to interact with Gitea/Forgejo repositories.

## Project Overview

This project provides MCP integration for managing Forgejo/Gitea repositories through MCP-compatible clients such as Claude Desktop, Continue, and other LLM applications.

### Supported Operations
- Issues (create, edit, comment, close)
- Labels (list, create, edit, delete)
- Milestones (list, create, edit, delete)
- Releases (list, create, edit, delete, manage assets)
- Pull requests (list, view, create, reviews, reply to inline review comments)
- Repository search and listing
- Wiki pages (create, edit, delete)
- Forgejo Actions tasks (view)

### Transport Modes
- **stdio**: Standard input/output (best for local integration)
- **sse**: Server-Sent Events over HTTP (best for web apps) - *planned*
- **http**: HTTP POST requests (best for simple integrations) - *planned*

## Architecture

### Core Components

- **cmd/**: CLI application using Cobra framework
  - `root.go`: Main command with global configuration
  - `stdio.go`: Stdio transport mode implementation
- **types/**: Data structures and response types
  - `api.go`: MCP response types wrapping Forgejo SDK types
- **main.go**: Application entry point

### Key Dependencies

- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk`
  * There are 3 important sub packages: `mcp`, `jsonschema` and `jsonrpc`
- **Forgejo SDK**: `codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2`
- **CLI Framework**: `github.com/spf13/cobra`
- **Configuration**: `github.com/spf13/viper`

If you need to get documentation of these packages, `go doc` should be the best bet.

### Implementation Strategy

- **🟢 SDK-based (71%)**: Use official Forgejo SDK where available
- **🟡 Custom HTTP (29%)**: Implement custom requests for unsupported features (Wiki, Actions, Issue dependencies)

## Configuration

### CLI Arguments
```bash
forgejo-mcp stdio --server https://git.example.com --token your_token
```

### Environment Variables
- `FORGEJOMCP_SERVER`: Forgejo server URL
- `FORGEJOMCP_TOKEN`: Access token

### Config File
Default location: `~/.forgejo-mcp.yaml`

Priority: CLI args > Environment variables > Config file

## Development

### Language and Style
- **Code/Comments**: English (open source project)
- **Documentation**: English (unless `.tw.md` suffix for Traditional Chinese)

### Key Files
- `proposal.tw.md`: Project requirements (Traditional Chinese)
- `features.tw.md`: Feature specifications (Traditional Chinese)
- `design.tw.md`: Architecture documentation (Traditional Chinese)
- `swagger.v1.json`: API documentation

### Response Format
- **Error responses**: Plain text
- **Success responses**: Markdown format in MCP TextContent.Text field
- **Structured data**: Dual format (markdown + JSON) for MCP compatibility

### Architecture Decisions
- Use **endpoint-based markdown formatting** rather than type-based to avoid LLM confusion
- Each data type implements `ToMarkdown()` method for reusability
- Endpoint handlers add context-specific headers and descriptions
- **Tool responses**: Use single TextContent with markdown formatting (follows MCP best practices)
- **Tool organization**: Mixed modular architecture with sub-packages for logical grouping

### Tool Architecture

The project uses a **mixed modular architecture** for organizing MCP tools:

#### Structure
```
tools/
├── module.go           # ToolImpl interface definition
├── registry.go         # Unified tool registration
├── helpers.go          # Shared utility functions
├── issue/              # Issue-related operations
│   ├── crud.go         # create/edit/delete issue
│   ├── list.go         # list_issues
│   ├── comment.go      # issue comments
│   ├── label.go        # issue label operations
│   ├── attach.go       # issue attachments
│   └── dep.go          # issue dependencies
├── label/              # Label management
│   └── crud.go         # all label operations
├── milestone/          # Milestone management
│   └── crud.go         # all milestone operations
├── release/            # Release management
│   ├── crud.go         # create/edit/delete release
│   └── attach.go       # release attachments
├── pullreq/            # Pull Request operations
│   ├── list.go         # list pull requests
│   ├── view.go         # get pull request details
│   └── review.go       # reviews, inline comments, replies
├── action/             # Forgejo Actions (CI/CD)
│   └── list.go         # list_action_tasks
├── wiki/               # Wiki pages
│   ├── crud.go         # create/edit/delete wiki pages
│   └── list.go         # list wiki pages
└── repo/               # Repository operations
    └── list.go         # repository listing and search
```

#### Design Principles
- **Logical grouping**: Related tools grouped in sub-packages
- **Single responsibility**: Each file handles closely related operations
- **Interface-driven**: All tools implement the `ToolImpl` interface
- **Extensibility**: Easy to add new tools and categories
- **Testability**: Each module can be independently tested

## Development Workflow
1. Test-Driven Development (TDD)
2. Red-Green-Refactor cycle
3. Human review at each step
4. Git commit after each milestone

## Useful Commands

### Go Commands
```bash
# Build the project for release
go build -o forgejo-mcp

# Check compilation error
go build ./...

# Run tests
go test ./...

# Format code
goimports -w .

# Get dependencies
go mod tidy
```

### Swagger

```
# Find definition of specified api endpoint
#   .paths["/path/to/endpoint"].http_method
jq '.paths["/repos/{owner}/{repo}/labels/{id}"].patch' swagger.v1.json

# Get summary string of specified api endpoint
#   .paths["/path/to/endpoint"].http_method.summary
#
# Change "summary" to "parameters" or "responses" to get params/responses
jq '.paths["/repos/{owner}/{repo}/labels/{id}"].patch.summary' swagger.v1.json

# Get referenced type definition
#
# for "$refs": "#/abc/def"
jq '.abc.def' swagger.v1.json
```

## Additional tools

- `my-git-server`: Tools to access remote git server. The repository is `ronmi/forgejo-mcp`.
- `gopls`: Go language server to query for difinitions, references, and more.
