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

// MyEditAttachmentOptions extends the SDK version with missing fields.
type MyEditAttachmentOptions struct {
	Name        string `json:"name,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

// MyListIssueAttachments lists all attachments of an issue.
// GET /repos/{owner}/{repo}/issues/{index}/assets
func (c *Client) MyListIssueAttachments(owner, repo string, index int64) ([]*forgejo.Attachment, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/assets", owner, repo, index)

	var result []*forgejo.Attachment
	err := c.sendSimpleRequest("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MyDeleteIssueAttachment deletes an attachment from an issue.
// DELETE /repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}
func (c *Client) MyDeleteIssueAttachment(owner, repo string, index, attachmentID int64) error {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/assets/%d", owner, repo, index, attachmentID)

	// DELETE returns 204 No Content on success, so we can use nil as response target
	var result interface{}
	err := c.sendSimpleRequest("DELETE", endpoint, nil, &result)
	if err != nil {
		return err
	}

	return nil
}

// MyEditIssueAttachment edits an attachment of an issue.
// PATCH /repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}
func (c *Client) MyEditIssueAttachment(owner, repo string, index, attachmentID int64, options MyEditAttachmentOptions) (*forgejo.Attachment, error) {
	endpoint := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/assets/%d", owner, repo, index, attachmentID)

	var result forgejo.Attachment
	err := c.sendSimpleRequest("PATCH", endpoint, options, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
