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

// Release represents a release response with embedded SDK release
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases (list)
// - POST /repos/{owner}/{repo}/releases (create)
// - PATCH /repos/{owner}/{repo}/releases/{id} (edit)
type Release struct {
	*forgejo.Release
}

// ToMarkdown renders release with tag, name, draft/prerelease status and description
// Example: **v1.0.0** - Major Release `PRERELEASE` (2024-01-15)
// This release includes new authentication system and bug fixes...
func (r *Release) ToMarkdown() string {
	if r.Release == nil {
		return "*Invalid release*"
	}
	markdown := "**" + r.TagName + "**"
	if r.Title != "" {
		markdown += " - " + r.Title
	}
	if r.IsDraft {
		markdown += " `DRAFT`"
	}
	if r.IsPrerelease {
		markdown += " `PRERELEASE`"
	}
	if !r.CreatedAt.IsZero() {
		markdown += " (" + r.CreatedAt.Format("2006-01-02") + ")"
	}
	if r.Note != "" {
		markdown += "\n" + r.Note
	}
	return markdown
}

// ReleaseList represents a list of releases response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/releases
type ReleaseList []*Release

// ToMarkdown renders releases as a numbered list with details
// Example:
// 1. **v1.0.0** - Major Release `PRERELEASE` (2024-01-15)
// This release includes new authentication system...
// 2. **v0.9.0** - Beta Release (2024-01-01)
// Initial beta version with core features
func (rl ReleaseList) ToMarkdown() string {
	if len(rl) == 0 {
		return "*No releases found*"
	}
	markdown := ""
	for i, release := range rl {
		markdown += fmt.Sprintf("%d. %s\n", i+1, release.ToMarkdown())
	}
	return markdown
}
