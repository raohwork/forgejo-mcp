// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const forgejo_version_to_test = "11.0.1+gitea-1.22.0"

// sendSimpleRequest Specification:
//
// Responsibility: Handle all pure JSON API requests, including Issue Dependencies,
// Wiki Pages, Forgejo Actions, and Issue Attachments non-upload operations
//
// Business Logic:
// 1. Create HTTP request (using provided method and endpoint)
// 2. If paramObj is not nil, serialize it as JSON for request body
// 3. Use Forgejo SDK's SignRequest method to add authentication headers
// 4. Send request and receive response
// 5. Deserialize JSON response to respObj
// 6. Handle HTTP error status codes
//
// Parameters:
// - method: HTTP method (GET, POST, PATCH, DELETE)
// - endpoint: API endpoint path (relative to base URL)
// - paramObj: request parameter object (JSON serialized), can be nil for GET/DELETE
// - respObj: response data receiver object (JSON deserialized)
//
// Returns: error if request fails or response cannot be parsed
//
// Design Philosophy:
// - Focus on our project needs, not pursuing genericity
// - Rely on Forgejo SDK for authentication, we only handle request/response serialization
// - Simplified error handling, mainly focus on network errors and JSON parsing errors
//
// Implementation Notes:
// - Need to combine c.base and endpoint to form complete URL
// - Content-Type should be set to application/json (when request body exists)
// - Accept should be set to application/json
// - HTTP 4xx/5xx status codes should return errors
func TestClient_sendSimpleRequest(t *testing.T) {
	// GET request test - no request body
	t.Run("GET_success", func(t *testing.T) {
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Expected GET method, got %s", r.Method)
			}
			if r.URL.Path != "/api/v1/repos/owner/repo/issues/1/dependencies" {
				t.Errorf("Expected specific path, got %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    1,
				"title": "Test Issue",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		var result map[string]interface{}
		err = client.sendSimpleRequest("GET", "/api/v1/repos/owner/repo/issues/1/dependencies", nil, &result)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result["id"] != float64(1) {
			t.Errorf("Expected id=1, got %v", result["id"])
		}
		if result["title"] != "Test Issue" {
			t.Errorf("Expected title='Test Issue', got %v", result["title"])
		}
	})

	// POST request test - with request body
	t.Run("POST_success", func(t *testing.T) {
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
			}

			var reqBody map[string]interface{}
			json.NewDecoder(r.Body).Decode(&reqBody)
			if reqBody["index"] != float64(2) {
				t.Errorf("Expected index=2, got %v", reqBody["index"])
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "dependency added",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		requestData := map[string]interface{}{"index": 2}
		var result map[string]interface{}
		err = client.sendSimpleRequest("POST", "/api/v1/repos/owner/repo/issues/1/dependencies", requestData, &result)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result["message"] != "dependency added" {
			t.Errorf("Expected message='dependency added', got %v", result["message"])
		}
	})

	// HTTP error handling test
	t.Run("HTTP_error", func(t *testing.T) {
		// Mock server returning 404
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Not found",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		var result map[string]interface{}
		err = client.sendSimpleRequest("GET", "/api/v1/repos/owner/repo/nonexistent", nil, &result)

		if err == nil {
			t.Error("Expected error for 404 response, got nil")
		}
	})

	// JSON parsing error test
	t.Run("JSON_parse_error", func(t *testing.T) {
		// Mock server returning invalid JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		var result map[string]interface{}
		err = client.sendSimpleRequest("GET", "/api/v1/repos/owner/repo/issues", nil, &result)

		if err == nil {
			t.Error("Expected JSON parsing error, got nil")
		}
	})
}

// sendUploadRequest Specification:
//
// Responsibility: Handle file upload requests, currently mainly for Issue Attachment creation
//
// Business Logic:
// 1. Create multipart/form-data format HTTP POST request
// 2. Add file to multipart writer (using filename)
// 3. Add additional form fields from extraFields
// 4. Use Forgejo SDK's SignRequest method to add authentication headers
// 5. Send request and receive response
// 6. Deserialize JSON response to respObj
//
// Parameters:
// - endpoint: API endpoint path (fixed to use POST method)
// - filename: upload file name
// - file: file content (io.Reader)
// - extraFields: additional form fields (e.g. name, updated_at)
// - respObj: response data receiver object
//
// Returns: error if upload fails or response cannot be parsed
//
// Design Philosophy:
// - Focus on Issue Attachment upload requirements
// - Use standard multipart/form-data format
// - Rely on Forgejo SDK for authentication
//
// Implementation Notes:
// - Content-Type will be automatically set by multipart.Writer
// - File field name should be "attachment" (according to Forgejo API)
// - Accept should be set to application/json
// - Need to correctly handle multipart boundary
func TestClient_sendUploadRequest(t *testing.T) {
	// Successful upload test
	t.Run("upload_success", func(t *testing.T) {
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}
			if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				t.Errorf("Expected multipart/form-data Content-Type, got %s", r.Header.Get("Content-Type"))
			}

			// Parse multipart form
			err := r.ParseMultipartForm(32 << 20) // 32MB
			if err != nil {
				t.Errorf("Failed to parse multipart form: %v", err)
			}

			// Check extra fields
			if r.FormValue("name") != "test.txt" {
				t.Errorf("Expected name='test.txt', got %s", r.FormValue("name"))
			}

			// Check file
			file, header, err := r.FormFile("attachment")
			if err != nil {
				t.Errorf("Failed to get file: %v", err)
			}
			defer file.Close()

			if header.Filename != "test.txt" {
				t.Errorf("Expected filename='test.txt', got %s", header.Filename)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":           "123",
				"name":         "test.txt",
				"download_url": "http://example.com/download/123",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		file := strings.NewReader("test file content")
		extraFields := map[string]string{"name": "test.txt"}
		var result map[string]interface{}

		err = client.sendUploadRequest("/api/v1/repos/owner/repo/issues/1/assets", "test.txt", file, extraFields, &result)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result["id"] != "123" {
			t.Errorf("Expected id='123', got %v", result["id"])
		}
		if result["name"] != "test.txt" {
			t.Errorf("Expected name='test.txt', got %v", result["name"])
		}
	})

	// Empty file upload test
	t.Run("empty_file", func(t *testing.T) {
		// Mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST method, got %s", r.Method)
			}

			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				t.Errorf("Failed to parse multipart form: %v", err)
			}

			file, header, err := r.FormFile("attachment")
			if err != nil {
				t.Errorf("Failed to get file: %v", err)
			}
			defer file.Close()

			if header.Filename != "empty.txt" {
				t.Errorf("Expected filename='empty.txt', got %s", header.Filename)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":   "124",
				"name": "empty.txt",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		file := strings.NewReader("")
		var result map[string]interface{}

		err = client.sendUploadRequest("/api/v1/repos/owner/repo/issues/1/assets", "empty.txt", file, nil, &result)

		if err != nil {
			t.Errorf("Expected no error for empty file, got %v", err)
		}
		if result["id"] != "124" {
			t.Errorf("Expected id='124', got %v", result["id"])
		}
	})

	// Upload error test
	t.Run("upload_error", func(t *testing.T) {
		// Mock server returning 500 error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Internal server error",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		file := strings.NewReader("test content")
		var result map[string]interface{}

		err = client.sendUploadRequest("/api/v1/repos/owner/repo/issues/1/assets", "test.txt", file, nil, &result)

		if err == nil {
			t.Error("Expected error for 500 response, got nil")
		}
	})
}

func TestCustomClient_Integral(t *testing.T) {
	server := os.Getenv("FORGEJO_TEST_SERVER")
	token := os.Getenv("FORGEJO_TEST_TOKEN")
	repo := os.Getenv("FORGEJO_TEST_REPO")
	if server == "" || token == "" || repo == "" {
		t.Skip("Skipping test, FORGEJO_TEST_SERVER, FORGEJO_TEST_TOKEN, and FORGEJO_TEST_REPO must be set")
	}

	t.Logf("Using server: %s, repo: %s", server, repo)
	cl, err := NewClient(server, token, "", nil)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var resp map[string]any
	err = cl.sendSimpleRequest("GET", "/api/v1/repos/"+repo+"/actions/tasks", nil, &resp)
	if err != nil {
		t.Fatalf("Failed to get response: %v", err)
	}

	if resp["total_count"] == nil {
		t.Error("Expected total_count in response, got nil")
	}
	if resp["workflow_runs"] == nil {
		t.Error("Expected tasks in response, got nil")
	}
}
