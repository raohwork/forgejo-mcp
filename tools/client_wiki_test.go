// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raohwork/forgejo-mcp/types"
)

func TestMyListWikiPages_Pagination(t *testing.T) {
	t.Run("paginates_through_all_pages", func(t *testing.T) {
		// Mock server that returns 50 items on page 1, 30 on page 2
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			page := r.URL.Query().Get("page")

			var pages []*types.MyWikiPageMetaData
			if page == "1" {
				// Return 50 items (full page)
				for i := 0; i < 50; i++ {
					pages = append(pages, &types.MyWikiPageMetaData{
						Title:  "page-" + string(rune('A'+i%26)) + string(rune('0'+i/26)),
						SubURL: "page-" + string(rune('A'+i%26)) + string(rune('0'+i/26)),
					})
				}
			} else if page == "2" {
				// Return 30 items (partial page = last page)
				for i := 0; i < 30; i++ {
					pages = append(pages, &types.MyWikiPageMetaData{
						Title:  "page-extra-" + string(rune('A'+i%26)),
						SubURL: "page-extra-" + string(rune('A'+i%26)),
					})
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(pages)
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		pages, err := client.MyListWikiPages("owner", "repo")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(pages) != 80 {
			t.Errorf("Expected 80 pages (50+30), got %d", len(pages))
		}
		if requestCount != 2 {
			t.Errorf("Expected 2 API requests (2 pages), got %d", requestCount)
		}
	})

	t.Run("single_page_no_extra_request", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			// Return 10 items (less than limit=50, so no second page)
			var pages []*types.MyWikiPageMetaData
			for i := 0; i < 10; i++ {
				pages = append(pages, &types.MyWikiPageMetaData{
					Title: "page-" + string(rune('A'+i)),
				})
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(pages)
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		pages, err := client.MyListWikiPages("owner", "repo")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(pages) != 10 {
			t.Errorf("Expected 10 pages, got %d", len(pages))
		}
		if requestCount != 1 {
			t.Errorf("Expected 1 API request (single page), got %d", requestCount)
		}
	})

	t.Run("empty_wiki_returns_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "no such file or directory",
			})
		}))
		defer server.Close()

		client, err := NewClient(server.URL, "test-token", forgejo_version_to_test, server.Client())
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		_, err = client.MyListWikiPages("owner", "repo")
		if err == nil {
			t.Error("Expected error for 404 response, got nil")
		}
	})
}
