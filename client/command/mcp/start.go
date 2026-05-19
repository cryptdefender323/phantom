package mcp

import (
	"github.com/cryptdefender3232/phantom/client/command/settings"
	"github.com/cryptdefender3232/phantom/client/console"
	phantommcp "github.com/cryptdefender3232/phantom/client/mcp"
	"github.com/spf13/cobra"
)

// McpStartCmd starts the local MCP server.
func McpStartCmd(cmd *cobra.Command, con *console.PhantomClient, args []string) {
	rawTransport, _ := cmd.Flags().GetString("transport")
	transport, err := phantommcp.ParseTransport(rawTransport)
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}

	listen, _ := cmd.Flags().GetString("listen")
	name, _ := cmd.Flags().GetString("name")
	version, _ := cmd.Flags().GetString("version")

	cfg := phantommcp.Config{
		Transport:     transport,
		ListenAddress: listen,
		ServerName:    name,
		ServerVersion: version,
	}.WithDefaults()

	msg := `Do you know what prompt injection is and are you an adult?`
	if !settings.IsUserAnAdultWithPrompt(con, msg) {
		con.PrintErrorf("Failed to start MCP server, the user is not qualified to use feature\n")
		return
	}

	if err := phantommcp.Start(cfg, con.Rpc); err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}

	con.PrintInfof("Starting MCP server (%s) on %s\n", cfg.Transport, cfg.ListenAddress)
	endpoint, err := cfg.EndpointURL()
	if err == nil {
		con.PrintInfof("Endpoint: %s\n", endpoint)
	}
	status := phantommcp.GetStatus()
	if status.AuthHeader != "" {
		con.PrintInfof("Auth Header: %s\n", status.AuthHeader)
	}
	if status.AuthToken != "" {
		con.PrintInfof("Auth Token: %s\n", status.AuthToken)
	}
}
