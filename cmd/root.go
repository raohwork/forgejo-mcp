// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com>

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "forgejo-mcp",
	Short: "Forgejo MCP Server for Model Context Protocol clients",
	Long: `Forgejo MCP Server provides Model Context Protocol integration
for managing Gitea/Forgejo repositories through MCP-compatible clients.

Supported operations:
  - Issues (create, edit, comment, close)
  - Labels (list, create, edit, delete)
  - Milestones (list, create, edit, delete)
  - Releases (list, create, edit, delete, manage assets)
  - Pull requests (list, view)
  - Repository search and listing
  - Wiki pages (create, edit, delete)
  - Forgejo Actions tasks (view)

Available transport modes:
  - stdio: Standard input/output (best for local integration)
  - sse: Server-Sent Events over HTTP (best for web apps)
  - http: HTTP POST requests (best for simple integrations)

Configure your Forgejo instance:
  forgejo-mcp [mode] --server https://git.example.com --token your_token`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	f := rootCmd.PersistentFlags()
	f.String("server", "", "Forgejo server URL (env: FORGEJOMCP_SERVER)")
	f.String("token", "", "Forgejo access token (env: FORGEJOMCP_TOKEN)")
	viper.BindPFlags(f)

	viper.SetEnvPrefix("FORGEJOMCP")
	viper.AutomaticEnv()
}
