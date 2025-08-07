// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package cmd

import (
	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/raohwork/forgejo-mcp/tools/action"
	"github.com/raohwork/forgejo-mcp/tools/issue"
	"github.com/raohwork/forgejo-mcp/tools/label"
	"github.com/raohwork/forgejo-mcp/tools/milestone"
	"github.com/raohwork/forgejo-mcp/tools/pullreq"
	"github.com/raohwork/forgejo-mcp/tools/release"
	"github.com/raohwork/forgejo-mcp/tools/repo"
	"github.com/raohwork/forgejo-mcp/tools/wiki"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCommands(s *mcp.Server, cl *tools.Client) {
	// Issue tools
	tools.Register(s, &issue.ListRepoIssuesImpl{Client: cl})
	tools.Register(s, &issue.GetIssueImpl{Client: cl})
	tools.Register(s, &issue.CreateIssueImpl{Client: cl})
	tools.Register(s, &issue.EditIssueImpl{Client: cl})

	// Issue label tools
	tools.Register(s, &issue.AddIssueLabelsImpl{Client: cl})
	tools.Register(s, &issue.RemoveIssueLabelImpl{Client: cl})
	tools.Register(s, &issue.ReplaceIssueLabelsImpl{Client: cl})

	// Issue comment tools
	tools.Register(s, &issue.ListIssueCommentsImpl{Client: cl})
	tools.Register(s, &issue.CreateIssueCommentImpl{Client: cl})
	tools.Register(s, &issue.EditIssueCommentImpl{Client: cl})
	tools.Register(s, &issue.DeleteIssueCommentImpl{Client: cl})

	// Issue attachment tools
	tools.Register(s, &issue.ListIssueAttachmentsImpl{Client: cl})
	tools.Register(s, &issue.CreateIssueAttachmentImpl{Client: cl})
	tools.Register(s, &issue.DeleteIssueAttachmentImpl{Client: cl})
	tools.Register(s, &issue.EditIssueAttachmentImpl{Client: cl})

	// Issue dependency tools
	tools.Register(s, &issue.ListIssueDependenciesImpl{Client: cl})
	tools.Register(s, &issue.AddIssueDependencyImpl{Client: cl})
	tools.Register(s, &issue.RemoveIssueDependencyImpl{Client: cl})

	// Issue blocking tools
	tools.Register(s, &issue.ListIssueBlockingImpl{Client: cl})
	tools.Register(s, &issue.AddIssueBlockingImpl{Client: cl})
	tools.Register(s, &issue.RemoveIssueBlockingImpl{Client: cl})

	// Label tools
	tools.Register(s, &label.ListRepoLabelsImpl{Client: cl})
	tools.Register(s, &label.CreateLabelImpl{Client: cl})
	tools.Register(s, &label.EditLabelImpl{Client: cl})
	tools.Register(s, &label.DeleteLabelImpl{Client: cl})

	// Milestone tools
	tools.Register(s, &milestone.ListRepoMilestonesImpl{Client: cl})
	tools.Register(s, &milestone.CreateMilestoneImpl{Client: cl})
	tools.Register(s, &milestone.EditMilestoneImpl{Client: cl})
	tools.Register(s, &milestone.DeleteMilestoneImpl{Client: cl})

	// Release tools
	tools.Register(s, &release.ListReleasesImpl{Client: cl})
	tools.Register(s, &release.CreateReleaseImpl{Client: cl})
	tools.Register(s, &release.EditReleaseImpl{Client: cl})
	tools.Register(s, &release.DeleteReleaseImpl{Client: cl})

	// Release attachment tools
	tools.Register(s, &release.ListReleaseAttachmentsImpl{Client: cl})
	tools.Register(s, &release.CreateReleaseAttachmentImpl{Client: cl})
	tools.Register(s, &release.EditReleaseAttachmentImpl{Client: cl})
	tools.Register(s, &release.DeleteReleaseAttachmentImpl{Client: cl})

	// Pull request tools
	tools.Register(s, &pullreq.ListPullRequestsImpl{Client: cl})
	tools.Register(s, &pullreq.GetPullRequestImpl{Client: cl})

	// Repository tools
	tools.Register(s, &repo.SearchRepositoriesImpl{Client: cl})
	tools.Register(s, &repo.ListMyRepositoriesImpl{Client: cl})
	tools.Register(s, &repo.ListOrgRepositoriesImpl{Client: cl})
	tools.Register(s, &repo.GetRepositoryImpl{Client: cl})

	// Wiki tools
	tools.Register(s, &wiki.GetWikiPageImpl{Client: cl})
	tools.Register(s, &wiki.CreateWikiPageImpl{Client: cl})
	tools.Register(s, &wiki.EditWikiPageImpl{Client: cl})
	tools.Register(s, &wiki.DeleteWikiPageImpl{Client: cl})
	tools.Register(s, &wiki.ListWikiPagesImpl{Client: cl})

	// Action tools
	tools.Register(s, &action.ListActionTasksImpl{Client: cl})
}

func createServer(cl *tools.Client) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Title:   "Forgejo MCP Server",
		Version: "0.0.1",
	}, &mcp.ServerOptions{
		PageSize:     50,
		Instructions: "An MCP server to interact with repositories on a Forgejo/Gitea instance.",
	})
	registerCommands(server, cl)

	return server
}
