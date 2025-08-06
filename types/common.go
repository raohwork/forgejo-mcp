// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

// EmptyResponse represents an empty response for endpoints that don't return data
// Used by endpoints that only return status codes:
// - DELETE /repos/{owner}/{repo}/labels/{id}
// - DELETE /repos/{owner}/{repo}/milestones/{id}
// - DELETE /repos/{owner}/{repo}/releases/{id}
// - DELETE /repos/{owner}/{repo}/releases/assets/{id}
// - DELETE /repos/{owner}/{repo}/issues/{index}/labels/{id}
// - DELETE /repos/{owner}/{repo}/issues/{index}/attachments/{attachment_id}
// - DELETE /repos/{owner}/{repo}/wiki/page/{pageName}
type EmptyResponse struct{}

// ToMarkdown renders simple success message for empty responses
// Example: *Operation completed successfully*
func (er EmptyResponse) ToMarkdown() string {
	return "*Operation completed successfully*"
}
