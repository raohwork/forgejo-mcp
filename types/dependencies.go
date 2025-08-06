// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package types

// MyIssueMeta represents basic issue information for dependency operations.
// This type is not available in the Forgejo SDK.
type MyIssueMeta struct {
	Index int64  `json:"index"`
	Owner string `json:"owner,omitempty"`
	Name  string `json:"repo,omitempty"`
}
