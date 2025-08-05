// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

func TestRepository_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		repo     *Repository
		required []string
	}{
		{
			name: "complete repository with all fields",
			repo: &Repository{
				Repository: &forgejo.Repository{
					FullName:    "owner/repo-name",
					Description: "A sample repository for testing purposes",
					Private:     true,
					Fork:        true,
					Template:    false,
					Stars:       42,
					Forks:       7,
					OpenIssues:  3,
					OpenPulls:   1,
					HTMLURL:     "https://git.example.com/owner/repo-name",
				},
			},
			required: []string{"owner/repo-name", "PRIVATE", "FORK", "A sample repository for testing purposes", "Stars: 42", "Forks: 7", "Issues: 3", "PRs: 1", "View Repository", "https://git.example.com/owner/repo-name"},
		},
		{
			name:     "nil repository",
			repo:     &Repository{Repository: nil},
			required: []string{"Invalid repository"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.repo.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestRepositoryList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		repositories RepositoryList
		required     []string
	}{
		{
			name: "multiple repositories with complete information",
			repositories: RepositoryList{
				&Repository{
					Repository: &forgejo.Repository{
						FullName:    "owner/repo-name",
						Description: "A sample repository",
						Private:     true,
						Fork:        true,
						Stars:       42,
						Forks:       7,
						OpenIssues:  3,
						OpenPulls:   1,
					},
				},
				&Repository{
					Repository: &forgejo.Repository{
						FullName:    "owner/another-repo",
						Description: "Another repository",
						Stars:       15,
						Forks:       2,
					},
				},
			},
			required: []string{"1.", "owner/repo-name", "PRIVATE", "FORK", "2.", "owner/another-repo"},
		},
		{
			name:         "empty repository list",
			repositories: RepositoryList{},
			required:     []string{"No repositories found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.repositories.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestPullRequest_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		pr       *PullRequest
		required []string
	}{
		{
			name: "complete pull request with all fields",
			pr: &PullRequest{
				PullRequest: &forgejo.PullRequest{
					Index:  42,
					Title:  "Add user authentication",
					Body:   "This PR implements OAuth2 authentication system",
					State:  "open",
					Poster: testUser(),
					Head: &forgejo.PRBranchInfo{
						Name: "feature/auth",
					},
					Base: &forgejo.PRBranchInfo{
						Name: "main",
					},
				},
			},
			required: []string{"#42", "Add user authentication", "open", "testuser", "feature/auth", "main", "This PR implements OAuth2 authentication"},
		},
		{
			name:     "nil pull request",
			pr:       &PullRequest{PullRequest: nil},
			required: []string{"Invalid pull request"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.pr.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestPullRequestList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		prs      PullRequestList
		required []string
	}{
		{
			name: "multiple pull requests with complete information",
			prs: PullRequestList{
				&PullRequest{
					PullRequest: &forgejo.PullRequest{
						Index:  42,
						Title:  "Add user authentication",
						State:  "open",
						Poster: testUser(),
					},
				},
				&PullRequest{
					PullRequest: &forgejo.PullRequest{
						Index: 41,
						Title: "Fix database connection",
						State: "merged",
					},
				},
			},
			required: []string{"1.", "#42", "Add user authentication", "2.", "#41", "Fix database connection"},
		},
		{
			name:     "empty pull request list",
			prs:      PullRequestList{},
			required: []string{"No pull requests found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.prs.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
