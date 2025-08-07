// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/raohwork/forgejo-mcp/types"
)

var UserAgent = "Forgejo-MCP/" + types.VERSION

// Client wraps the Forgejo SDK client with additional functionality for
// unsupported API endpoints. It provides methods for JSON requests and
// multipart file uploads with manual authentication.
type Client struct {
	*forgejo.Client
	cl    *http.Client
	base  string
	token string
}

// NewClient creates a new Client instance with extended functionality beyond
// the standard Forgejo SDK. This client supports custom API endpoints that
// are not available in the SDK, such as issue dependencies, wiki pages,
// and Forgejo Actions.
//
// Parameters:
//   - base: Forgejo server base URL (e.g., "https://git.example.com")
//   - token: API access token for authentication
//   - version: Forgejo version string to skip version check, empty to auto-detect
//   - cl: HTTP client to use, nil for http.DefaultClient
//
// The client uses manual token authentication for custom endpoints while
// preserving full SDK functionality for supported operations.
func NewClient(base, token, version string, cl *http.Client) (*Client, error) {
	if cl == nil {
		cl = http.DefaultClient
	}
	opts := make([]forgejo.ClientOption, 0, 4)
	opts = append(
		opts,
		forgejo.SetHTTPClient(cl),
		forgejo.SetToken(token),
		forgejo.SetUserAgent(UserAgent),
	)
	if version != "" {
		opts = append(opts, forgejo.SetForgejoVersion(version))
	}
	sdk, err := forgejo.NewClient(base, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: sdk,
		cl:     cl,
		base:   base,
		token:  token,
	}, nil
}

// sendSimpleRequest handles pure JSON API requests
// method: HTTP method (GET, POST, PATCH, DELETE)
// endpoint: API endpoint path (relative to base URL)
// paramObj: request parameter object (JSON serialized), can be nil for GET/DELETE
// respObj: response data receiver object (JSON deserialized)
func (c *Client) sendSimpleRequest(method, endpoint string, paramObj, respObj any) error {
	// Build complete URL
	u, err := url.Parse(c.base + endpoint)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare request body
	var body io.Reader
	if paramObj != nil {
		jsonData, err := json.Marshal(paramObj)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	if paramObj != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set authentication header manually
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	// Send request
	resp, err := c.cl.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse JSON response
	if err := json.NewDecoder(resp.Body).Decode(respObj); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// sendUploadRequest handles file upload requests (multipart/form-data)
// endpoint: API endpoint path (fixed to use POST)
// filename: upload file name
// file: file content
// extraFields: additional form fields
// respObj: response data receiver object (JSON deserialized)
func (c *Client) sendUploadRequest(endpoint, filename string, file io.Reader, extraFields map[string]string, respObj any) error {
	// Build complete URL
	u, err := url.Parse(c.base + endpoint)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add extra fields
	for key, value := range extraFields {
		if err := writer.WriteField(key, value); err != nil {
			return fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// Add file
	part, err := writer.CreateFormFile("attachment", filename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", u.String(), &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set authentication header manually
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	// Send request
	resp, err := c.cl.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse JSON response
	if err := json.NewDecoder(resp.Body).Decode(respObj); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
