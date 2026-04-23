// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// MyCreatePullReviewCommentOptions is the payload for adding a comment to an
// existing pull request review. The Forgejo SDK does not expose this endpoint,
// so it is implemented here against the REST API directly.
type MyCreatePullReviewCommentOptions struct {
	Body        string `json:"body"`
	Path        string `json:"path"`
	NewPosition int64  `json:"new_position,omitempty"`
	OldPosition int64  `json:"old_position,omitempty"`
}

// MyCreatePullReviewComment adds a new inline comment to an existing pull
// request review, which is how Forgejo continues an inline conversation
// thread at a given file and line.
// POST /repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments
func (c *Client) MyCreatePullReviewComment(owner, repo string, index, reviewID int64, opt MyCreatePullReviewCommentOptions) (*forgejo.PullReviewComment, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/reviews/%d/comments", owner, repo, index, reviewID)

	var result forgejo.PullReviewComment
	if err := c.sendSimpleRequest("POST", endpoint, opt, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
