/*
Copyright © 2025 Ronmi Ren <ronmi.ren@gmail.com
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/raohwork/forgejo-mcp/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run MCP server in SSE mode",
	Long: `Run the Forgejo MCP server with an HTTP interface.

This command starts a web server that allows MCP clients to connect over HTTP.
It supports both Server-Sent Events (SSE) on the /sse endpoint and a standard
request/response model on the / endpoint.

The server can operate in two modes:
  - Single-user mode: When a --token is provided at startup, all operations
    are performed using this single token. This is suitable for personal use
    or dedicated services where one identity is sufficient.
  - Multi-user mode: If no --token is provided, the server requires clients
    to authenticate by providing their own token in the 'Authorization'
    header of each request. This allows the server to act as a gateway for
    multiple users.

This HTTP mode is ideal for:
  - Web-based clients and services.
  - Remote access to the Forgejo instance through MCP.
  - Environments where a central MCP gateway is needed.

Example:
  forgejo-mcp http --address :8080 --server https://git.example.com --token your_token`,
	Run: func(cmd *cobra.Command, args []string) {
		base := viper.GetString("server")
		token := viper.GetString("token")
		addr := viper.GetString("address")
		if addr == "" {
			addr = ":8080"
		}

		if base == "" {
			cmd.Help()
			os.Exit(1)
		}
		singleMode := token != ""

		var cl *tools.Client
		if singleMode {
			c, err := tools.NewClient(base, token, "", nil)
			if err != nil {
				fmt.Printf("Error creating SDK client: %v\n", err)
				os.Exit(1)
			}
			cl = c
		} else {
			cl, _ = tools.NewClient(base, "", "9", nil)
		}

		getServer := func(q *http.Request) *mcp.Server {
			if singleMode {
				return createServer(cl)
			}

			mycl := cl
			myToken := q.Header.Get("Authorization")
			if myToken != "" {
				if strings.HasPrefix(myToken, "Bearer ") {
					myToken = myToken[7:]
				}
				c, err := tools.NewClient(base, myToken, "", nil)
				if err == nil {
					mycl = c
				}
			}

			return createServer(mycl)
		}

		mux := http.NewServeMux()
		mux.Handle("/sse", mcp.NewSSEHandler(getServer))
		mux.Handle("/", mcp.NewStreamableHTTPHandler(getServer, nil))

		mode := "single"
		if !singleMode {
			mode = "multiuser"
		}
		fmt.Printf("Starting %s mode MCP server on %s\n", mode, addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			fmt.Printf("Server exited with error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	f := httpCmd.Flags()
	f.String("address", ":8080", "Address to listen on for incoming connections")
	viper.BindPFlags(f)
}
