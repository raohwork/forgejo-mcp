// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Attachment represents an attachment response with embedded SDK attachment
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases/{id}/assets (list release attachments)
// - POST /repos/{owner}/{repo}/releases/{id}/assets (create release attachment)
// - PATCH /repos/{owner}/{repo}/releases/assets/{id} (edit release attachment)
// - GET /repos/{owner}/{repo}/issues/{index}/attachments (list issue attachments)
// - POST /repos/{owner}/{repo}/issues/{index}/attachments (create issue attachment)
// - PATCH /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id} (edit issue attachment)
type Attachment struct {
	*forgejo.Attachment
}

// ToMarkdown renders attachment with name, size and download link
// Example: **document.pdf** (1024 bytes) [Download](https://git.example.com/attachments/123)
func (a *Attachment) ToMarkdown() string {
	if a.Attachment == nil {
		return "*Invalid attachment*"
	}
	markdown := "**" + a.Name + "**"
	if a.Size > 0 {
		markdown += fmt.Sprintf(" (%d bytes)", a.Size)
	}
	if a.DownloadURL != "" {
		markdown += " [Download](" + a.DownloadURL + ")"
	}
	return markdown
}

// AttachmentList represents a list of attachments response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases/{id}/assets
// - GET /repos/{owner}/{repo}/issues/{index}/attachments
type AttachmentList []*Attachment

// ToMarkdown renders attachments as a bullet list with download info
// Example:
// - **document.pdf** (1024 bytes) [Download](https://git.example.com/attachments/123)
// - **screenshot.png** (2048 bytes) [Download](https://git.example.com/attachments/124)
func (al AttachmentList) ToMarkdown() string {
	if len(al) == 0 {
		return "*No attachments found*"
	}
	markdown := ""
	for _, attachment := range al {
		markdown += "- " + attachment.ToMarkdown() + "\n"
	}
	return markdown
}
