// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

import (
	"strings"
	"testing"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// assertContains checks if all required strings are present in output
func assertContains(t *testing.T, output string, required []string) {
	t.Helper()
	for _, req := range required {
		if !strings.Contains(output, req) {
			t.Errorf("Expected output to contain %q, but it didn't. Output: %s", req, output)
		}
	}
}

// testTime returns a consistent test time
func testTime() time.Time {
	return time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
}

// testUser creates a test user with typical data
func testUser() *forgejo.User {
	return &forgejo.User{
		ID:       123,
		UserName: "testuser",
		FullName: "Test User",
		Email:    "test@example.com",
	}
}

// testMilestone creates a test milestone with typical data
func testMilestone() *forgejo.Milestone {
	deadline := testTime()
	return &forgejo.Milestone{
		ID:           1,
		Title:        "v1.0.0",
		Description:  "Major release",
		State:        "open",
		OpenIssues:   5,
		ClosedIssues: 10,
		Deadline:     &deadline,
	}
}

// testLabel creates a test label with typical data
func testLabel() *forgejo.Label {
	return &forgejo.Label{
		ID:          1,
		Name:        "bug",
		Color:       "ff0000",
		Description: "Something isn't working",
	}
}
