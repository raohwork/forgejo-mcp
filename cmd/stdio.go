// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stdioCmd represents the stdio command
var stdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Run MCP server in stdio mode",
	Long: `Run the Forgejo MCP server using stdio transport for communication
with MCP clients through JSON-RPC over standard input/output.

This transport mode is ideal for:
  - Local integrations and development
  - Direct process communication
  - Applications that can spawn child processes

Example:
  forgejo-mcp stdio --server https://git.example.com --token your_token`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stdio called")
	},
}

func init() {
	rootCmd.AddCommand(stdioCmd)
}
