# forgejo-mcp - Inventory Report (2026-04-05)

## Summary

forgejo-mcp is an open-source Model Context Protocol (MCP) server written in Go that enables AI assistants (Claude, Gemini, GitHub Copilot, etc.) to manage Gitea/Forgejo code repositories through natural language. It exposes repository operations — issue management, labels, milestones, releases, pull requests, wiki pages, and CI/CD action tasks — as MCP tools that LLM clients can invoke. The project supports both stdio (local) and HTTP transport modes, and uses the official Forgejo SDK for ~71% of operations with custom HTTP implementations for features the SDK doesn't cover (wiki, Actions, issue dependencies). It includes a compiled binary (`forgejo-mcp`, 13MB), Traditional Chinese documentation (`.tw.md` files), and a pre-built glama.ai integration spec. The project was developed with a rapid iteration cadence (72 commits since January 2025) and reached a major refactor milestone in November 2025 that consolidated 47 tools into 8 unified action-based tools.

## Technologies & Libraries

- **Language**: Go 1.24.5
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk` v0.4.0 (upgraded from v0.2.x during development)
- **Forgejo SDK**: `codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2` v2.2.0
- **CLI framework**: Cobra (`github.com/spf13/cobra` v1.9.1)
- **Configuration**: Viper (`github.com/spf13/viper` v1.20.1)
- **JSON Schema**: `github.com/google/jsonschema-go`
- **Auth**: HTTP signatures (`github.com/go-fed/httpsig`, `github.com/42wim/httpsig`)
- **Containerization**: Dockerfile present (minimal, single-stage)
- **CI**: Forgejo Actions (`.forgejo/` directory)
- **Discovery**: `glama.json` spec for glama.ai MCP registry auto-discovery
- **Documentation**: English README + Traditional Chinese variants (`README.tw.md`, `features.tw.md`, `proposal.tw.md`, `prompt.tw.md`)
- **API reference**: `swagger.v1.json` (763KB) — full Forgejo/Gitea Swagger spec for reference

## Timeline of Work

- **Started**: August 3, 2025 — `init` commit (Taiwanese timezone, +0800)
- **Phase 1** (Aug 3–4, 2025): Basic CLI entry, API response types with Markdown rendering, CI setup, issue dependency management specs
- **Active development** (Aug–Oct 2025): 70 commits adding features — create PR tool, time filters on issues, improved response formatting, SDK upgrade, attachment upload then removal, logo, glama.ai spec
- **Major refactor** (November 27, 2025): Consolidated 47 tools into 8 unified action-based tools (the largest single commit)
- **Most recent commit**: November 27, 2025
- **Total commits since Jan 2025**: 72
- **Last repo activity (git ops)**: April 5, 2026

## Current State

- **Runnable**: Yes — a pre-compiled binary `forgejo-mcp` (13.3MB) is present in the repository root, ready to run without a build step. The `.env` and `.env.local` files are present with actual credentials (server URL, token).
- **Branch**: `master`, up to date with `origin/master`
- **Uncommitted changes**: Two untracked files — `.mcp.json` (MCP client configuration for Claude Desktop, 293 bytes) and the compiled `forgejo-mcp` binary. Neither is harmful to leave untracked.
- **Last thing worked on**: The November 27, 2025 refactor consolidating 47 tools into 8 — a significant architectural improvement for LLM usability (fewer, more generalized tools are easier for models to use effectively)

## Successes

- Fully functional MCP server binary with a complete feature set across issues, labels, milestones, releases, PRs, wiki, and Actions
- The tool consolidation from 47 to 8 tools is a meaningful UX improvement — MCP clients (Claude, etc.) work better with fewer, broader tools
- MCP SDK upgraded to v0.4.0, keeping up with a rapidly evolving protocol
- Registered/discoverable on glama.ai (an MCP server directory) via `glama.json`
- Traditional Chinese documentation alongside English, suggesting the developer community or target users are partly Taiwanese/Chinese-speaking
- 72 commits in ~4 months demonstrates sustained, active development
- Clean separation of SDK-based vs. custom HTTP implementations with explicit documentation of which is which

## Open Items / Failures

- SSE and streamable HTTP transport modes are listed as "planned" in CLAUDE.md but not yet implemented — only stdio is available
- The compiled binary is committed to the repository (13MB), which is unusual practice for a Go project and may cause issues for contributors who need platform-specific builds
- Attachment upload functionality was added and then removed (three commits in October 2025), suggesting API surface instability
- No test files are visible in the repository structure
- The `.mcp.json` config file is untracked, suggesting local usage configuration hasn't been committed (possibly intentional for security reasons)
- The project appears to be a personal fork/adaptation — the `go.mod` module path is `github.com/raohwork/forgejo-mcp`, suggesting the original author is "raohwork" and this may be a local working copy rather than the developer's own fork

## Long-term Vision

The project aims to make Forgejo/Gitea repositories AI-manageable through the MCP protocol, enabling AI assistants to perform repository management tasks (triaging issues, managing releases, organizing wikis) without requiring the user to switch context to the web UI. The glama.ai registration suggests intent to publish this as a community tool. The Traditional Chinese documentation parallels suggest a goal of reaching the Taiwanese/Chinese open-source community specifically. The long-term vision is likely a fully-featured, multi-transport MCP server that covers the entire Forgejo API surface.

## Portfolio Fit

This project demonstrates Go proficiency, protocol-level work (MCP and REST API integration), and practical tool-building for the AI developer ecosystem. MCP servers are an emerging category as of 2025-2026, and building one in Go (rather than the more common TypeScript/Python implementations) shows differentiation. The tool consolidation refactor demonstrates architectural thinking about developer experience and LLM usability. The project is well-suited for a portfolio targeting roles in developer tooling, AI infrastructure, or open-source Go development. The glama.ai registration also shows awareness of distribution and discoverability for developer tools.
