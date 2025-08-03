# Forgejo MCP Server Requirements Specification

## Project Objective

Create an MCP (Model Context Protocol) server that enables Claude Code to manage repositories on Gitea/Forgejo servers.

### Core Requirements

1. **Communication Protocol**
   - Support stdio mode for Claude Code integration
   - Use official Go SDK (github.com/modelcontextprotocol/go-sdk)

2. **Feature Scope**
   - See `features.md`

3. **Configuration Management**
   - CLI arguments: `--server` (Forgejo server URL), `--token` (access token)
   - Environment variables: `FORGEJOMCP_SERVER`, `FORGEJOMCP_TOKEN`
   - Configuration priority: CLI args > Environment variables

### Technical Specifications

- **Programming Language**: Go
- **Main Dependencies**:
  * Official MCP Go SDK (github.com/modelcontextprotocol/go-sdk)
  * Forgejo SDK (codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2)
  * CLI Framework (github.com/spf13/cobra)
  * Configuration Management (github.com/spf13/viper)
- **Response Format**:
  * Errors returned in plaintext format
  * Normal responses returned in markdown format

## Development Principles

- Test-Driven Development (TDD)
- Agile development, starting with MVP
- Independent development and testing for each feature
- Following pair programming principles

## Documentation, Comments, and Messages

Considering this is an open source project, comments and messages (such as logging) must be in English.

Document files without language annotation should use English. Files with tw annotation (e.g., xxx.tw.md) should use Traditional Chinese.