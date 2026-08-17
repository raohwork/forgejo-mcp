// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package action

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatJobLogs(t *testing.T) {
	tests := []struct {
		name       string
		logText    string
		tail       int
		required   []string
		unexpected []string
	}{
		{
			name:     "short log is returned in full",
			logText:  "line one\nline two\nline three\n",
			tail:     500,
			required: []string{"Logs of job 42 (3 lines)", "line one", "line two", "line three"},
		},
		{
			name:       "long log is tailed with truncation notice",
			logText:    "first\nsecond\nthird\nfourth\nfifth",
			tail:       2,
			required:   []string{"last 2 of 5 lines", "fourth", "fifth"},
			unexpected: []string{"first", "second", "third"},
		},
		{
			name:     "empty log",
			logText:  "",
			tail:     500,
			required: []string{"Logs of job 42 are empty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := formatJobLogs(tt.logText, 42, tt.tail)
			for _, req := range tt.required {
				if !strings.Contains(output, req) {
					t.Errorf("Expected output to contain %q, but it didn't. Output: %s", req, output)
				}
			}
			for _, un := range tt.unexpected {
				if strings.Contains(output, un) {
					t.Errorf("Expected output NOT to contain %q, but it did. Output: %s", un, output)
				}
			}
		})
	}
}

func TestFormatJobLogsExactTail(t *testing.T) {
	// A log with exactly `tail` lines must not report truncation.
	logText := "a\nb\nc"
	output := formatJobLogs(logText, 7, 3)
	if strings.Contains(output, "last") {
		t.Errorf("Expected no truncation notice for exact-tail log. Output: %s", output)
	}
	if !strings.Contains(output, fmt.Sprintf("Logs of job %d (3 lines)", 7)) {
		t.Errorf("Expected full-log header. Output: %s", output)
	}
}
