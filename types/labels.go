// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Label represents a label response with embedded SDK label
// Used by endpoints:
// - GET /repos/{owner}/{repo}/labels (list)
// - POST /repos/{owner}/{repo}/labels (create)
// - PATCH /repos/{owner}/{repo}/labels/{id} (edit)
// - POST /repos/{owner}/{repo}/issues/{index}/labels (add to issue)
// - PUT /repos/{owner}/{repo}/issues/{index}/labels (replace issue labels)
// - DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id} (remove from issue)
type Label struct {
	*forgejo.Label
}

// ToMarkdown renders a label as a colored badge with name and description
// Example: **bug** `#ff0000` - Something isn't working
func (l *Label) ToMarkdown() string {
	if l.Label == nil {
		return "*Invalid label*"
	}
	markdown := "**" + l.Name + "**"
	if l.Color != "" {
		markdown += " `#" + l.Color + "`"
	}
	if l.Description != "" {
		markdown += " - " + l.Description
	}
	return markdown
}

// LabelList represents a list of labels response
// Used by endpoints:
// - GET /repos/{owner}/{repo}/labels
// - POST /repos/{owner}/{repo}/issues/{index}/labels
// - PUT /repos/{owner}/{repo}/issues/{index}/labels
type LabelList []*Label

// ToMarkdown renders labels as a bullet list of colored badges
// Example:
// - **bug** `#ff0000` - Something isn't working
// - **enhancement** `#a2eeef` - New feature or request
func (ll LabelList) ToMarkdown() string {
	if len(ll) == 0 {
		return "*No labels found*"
	}
	markdown := ""
	for _, label := range ll {
		markdown += "- " + label.ToMarkdown() + "\n"
	}
	return markdown
}
