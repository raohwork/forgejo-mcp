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

func TestRelease_ToMarkdown(t *testing.T) {
	created := testTime()
	tests := []struct {
		name     string
		release  *Release
		required []string
	}{
		{
			name: "complete release with all fields",
			release: &Release{
				Release: &forgejo.Release{
					TagName:      "v1.0.0",
					Title:        "Major Release",
					Note:         "This release includes new authentication system and bug fixes",
					IsDraft:      false,
					IsPrerelease: true,
					CreatedAt:    created,
				},
			},
			required: []string{"v1.0.0", "Major Release", "PRERELEASE", "2024-01-15", "This release includes new authentication system"},
		},
		{
			name:     "nil release",
			release:  &Release{Release: nil},
			required: []string{"Invalid release"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.release.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestReleaseList_ToMarkdown(t *testing.T) {
	created := testTime()
	tests := []struct {
		name     string
		releases ReleaseList
		required []string
	}{
		{
			name: "multiple releases with complete information",
			releases: ReleaseList{
				&Release{
					Release: &forgejo.Release{
						TagName:      "v1.0.0",
						Title:        "Major Release",
						Note:         "New features",
						IsPrerelease: true,
						CreatedAt:    created,
					},
				},
				&Release{
					Release: &forgejo.Release{
						TagName:   "v0.9.0",
						Title:     "Beta Release",
						CreatedAt: created,
					},
				},
			},
			required: []string{"1.", "v1.0.0", "Major Release", "PRERELEASE", "2.", "v0.9.0", "Beta Release"},
		},
		{
			name:     "empty release list",
			releases: ReleaseList{},
			required: []string{"No releases found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.releases.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestAttachment_ToMarkdown(t *testing.T) {
	tests := []struct {
		name       string
		attachment *Attachment
		required   []string
	}{
		{
			name: "complete attachment with all fields",
			attachment: &Attachment{
				Attachment: &forgejo.Attachment{
					Name:        "document.pdf",
					Size:        1024,
					DownloadURL: "https://git.example.com/attachments/123",
				},
			},
			required: []string{"document.pdf", "1024 bytes", "Download", "https://git.example.com/attachments/123"},
		},
		{
			name:       "nil attachment",
			attachment: &Attachment{Attachment: nil},
			required:   []string{"Invalid attachment"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.attachment.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}

func TestAttachmentList_ToMarkdown(t *testing.T) {
	tests := []struct {
		name        string
		attachments AttachmentList
		required    []string
	}{
		{
			name: "multiple attachments with complete information",
			attachments: AttachmentList{
				&Attachment{
					Attachment: &forgejo.Attachment{
						Name:        "document.pdf",
						Size:        1024,
						DownloadURL: "https://git.example.com/attachments/123",
					},
				},
				&Attachment{
					Attachment: &forgejo.Attachment{
						Name:        "screenshot.png",
						Size:        2048,
						DownloadURL: "https://git.example.com/attachments/124",
					},
				},
			},
			required: []string{"document.pdf", "1024 bytes", "screenshot.png", "2048 bytes", "Download"},
		},
		{
			name:        "empty attachment list",
			attachments: AttachmentList{},
			required:    []string{"No attachments found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.attachments.ToMarkdown()
			assertContains(t, output, tt.required)
		})
	}
}
