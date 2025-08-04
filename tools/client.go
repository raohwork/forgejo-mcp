// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"errors"
	"io"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

const UserAgent = "Forgejo-MCP/0.0.1"

type Client struct {
	*forgejo.Client
	cl   *http.Client
	base string
}

func NewClient(base, token string, cl *http.Client) (*Client, error) {
	if cl == nil {
		cl = http.DefaultClient
	}
	sdk, err := forgejo.NewClient(
		base,
		forgejo.SetHTTPClient(cl),
		forgejo.SetToken(token),
		forgejo.SetUserAgent(UserAgent),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: sdk,
		cl:     cl,
		base:   base,
	}, nil
}

// sendSimpleRequest handles pure JSON API requests
// method: HTTP method (GET, POST, PATCH, DELETE)
// endpoint: API endpoint path (relative to base URL)
// paramObj: request parameter object (JSON serialized), can be nil for GET/DELETE
// respObj: response data receiver object (JSON deserialized)
func (c *Client) sendSimpleRequest(method, endpoint string, paramObj, respObj any) error {
	return errors.New("not implemented")
}

// sendUploadRequest handles file upload requests (multipart/form-data)
// endpoint: API endpoint path (fixed to use POST)
// filename: upload file name
// file: file content
// extraFields: additional form fields
// respObj: response data receiver object (JSON deserialized)
func (c *Client) sendUploadRequest(endpoint, filename string, file io.Reader, extraFields map[string]string, respObj any) error {
	return errors.New("not implemented")
}
