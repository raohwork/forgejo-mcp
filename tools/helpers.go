// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package tools

// Helper functions for creating pointers to basic types

// boolPtr creates a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}

// intPtr creates a pointer to an int value
func intPtr(i int) *int {
	return &i
}

// float64Ptr creates a pointer to a float64 value
func float64Ptr(f float64) *float64 {
	return &f
}
