package command

/*
	Phantom Implant Framework
	Copyright (C) 2023  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"os"

	"github.com/cryptdefender3232/phantom/client/command/ai"
	"github.com/cryptdefender3232/phantom/client/command/aka"
	"github.com/cryptdefender3232/phantom/client/command/alias"
	"github.com/cryptdefender3232/phantom/client/command/armory"
	"github.com/cryptdefender3232/phantom/client/command/beacons"
	"github.com/cryptdefender3232/phantom/client/command/builders"
	"github.com/cryptdefender3232/phantom/client/command/c2profiles"
	"github.com/cryptdefender3232/phantom/client/command/certificates"
	"github.com/cryptdefender3232/phantom/client/command/clean"
	"github.com/cryptdefender3232/phantom/client/command/crack"
	"github.com/cryptdefender3232/phantom/client/command/creds"
	docscmd "github.com/cryptdefender3232/phantom/client/command/docs"
	"github.com/cryptdefender3232/phantom/client/command/exit"
	"github.com/cryptdefender3232/phantom/client/command/extensions"
	"github.com/cryptdefender3232/phantom/client/command/generate"
	"github.com/cryptdefender3232/phantom/client/command/hosts"
	"github.com/cryptdefender3232/phantom/client/command/info"
	"github.com/cryptdefender3232/phantom/client/command/jobs"
	"github.com/cryptdefender3232/phantom/client/command/licenses"
	"github.com/cryptdefender3232/phantom/client/command/loot"
	"github.com/cryptdefender3232/phantom/client/command/mcp"
	"github.com/cryptdefender3232/phantom/client/command/monitor"
	"github.com/cryptdefender3232/phantom/client/command/operators"
	"github.com/cryptdefender3232/phantom/client/command/reaction"
	"github.com/cryptdefender3232/phantom/client/command/serverctx"
	"github.com/cryptdefender3232/phantom/client/command/audit"
	"github.com/cryptdefender3232/phantom/client/command/engagements"
	"github.com/cryptdefender3232/phantom/client/command/sessions"
	"github.com/cryptdefender3232/phantom/client/command/settings"
	shellcodeencoders "github.com/cryptdefender3232/phantom/client/command/shellcode-encoders"
	sgn "github.com/cryptdefender3232/phantom/client/command/shikata-ga-nai"
	"github.com/cryptdefender3232/phantom/client/command/socks"
	"github.com/cryptdefender3232/phantom/client/command/taskmany"
	"github.com/cryptdefender3232/phantom/client/command/update"
	"github.com/cryptdefender3232/phantom/client/command/use"
	"github.com/cryptdefender3232/phantom/client/command/websites"
	"github.com/cryptdefender3232/phantom/client/command/wireguard"
	client "github.com/cryptdefender3232/phantom/client/console"
	consts "github.com/cryptdefender3232/phantom/client/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// ServerCommands returns all commands bound to the server menu, optionally
// accepting a function returning a list of additional (admin) commands.
func ServerCommands(con *client.PhantomClient, serverCmds func() []*cobra.Command) console.Commands {
	serverCommands := func() *cobra.Command {
		server := &cobra.Command{
			Short: "Server commands",
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		}
		if !con.IsCLI {
			server.SilenceErrors = true
			server.SilenceUsage = true
		}

		// Utility function to be used for binding new commands to
		// the phantom menu: call the function with the name of the
		// group under which this/these commands should be added,
		// and the group will be automatically created if needed.
		bind := makeBind(server, con)

		if serverCmds != nil {
			server.AddGroup(&cobra.Group{ID: consts.MultiplayerHelpGroup, Title: consts.MultiplayerHelpGroup})
			server.AddCommand(serverCmds()...)
		}

		// [ Bind commands ] --------------------------------------------------------

		// Below are bounds all commands of the server menu, gathered by the group
		// under which they should be printed in help messages and/or completions.
		// You can either add a new bindCommands() call with a new group (which will
		// be automatically added to the command tree), or add your commands in one of
		// the present calls.

		// Core
		bind(consts.GenericHelpGroup,
			exit.Command,
			ai.ServerCommands,
			docscmd.ServerCommands,
			serverctx.Commands,
			licenses.Commands,
			settings.Commands,
			alias.Commands,
			extensions.Commands,
			armory.Commands,
			update.Commands,
			operators.Commands,
			creds.Commands,
			crack.Commands,
			certificates.Commands,
			clean.Command,
			aka.ServerCommands,
			audit.Commands,
			engagements.Commands,
		)

		// C2 Network
		bind(consts.NetworkHelpGroup,
			jobs.Commands,
			mcp.Commands,
			websites.Commands,
			wireguard.Commands,
			c2profiles.Commands,
			socks.RootCommands,
		)

		// Payloads
		bind(consts.PayloadsHelpGroup,
			sgn.Commands,
			shellcodeencoders.Commands,
			generate.Commands,
			builders.Commands,
		)

		// Phantoms
		bind(consts.PhantomHelpGroup,
			use.Commands,
			info.Commands,
			sessions.Commands,
			beacons.Commands,
			monitor.Commands,
			loot.Commands,
			hosts.Commands,
			reaction.Commands,
			taskmany.Command,
		)

		// [ Post-command declaration setup ]-----------------------------------------

		// Load Extensions
		// Similar to the PhantomCommand loading, without adding the commands to the
		// Server command tree. This is done to ensure that the extensions are loaded
		// before the server is started, so that the extensions are registered.
		extensionManifests := extensions.GetAllExtensionManifests()
		for _, manifest := range extensionManifests {
			_, err := extensions.LoadExtensionManifest(manifest)
			// Absorb error in case there's no extensions manifest
			if err != nil {
				//con doesn't appear to be initialised here?
				//con.PrintErrorf("Failed to load extension: %s", err)
				fmt.Printf("Failed to load extension: %s\n", err)
				continue
			}

			//for _, ext := range mext.ExtCommand {
			//	extensions.ExtensionRegisterCommand(ext, phantom, con)
			//}
		}

		// Everything below this line should preferably not be any command binding
		// (although you can do so without fear). If there are any final modifications
		// to make to the server menu command tree, it time to do them here.

		// Only load reactions when the console is going to be started.
		if !con.IsCLI {
			n, err := reaction.LoadReactions()
			if err != nil && !os.IsNotExist(err) {
				con.PrintErrorf("Failed to load reactions: %s\n", err)
			} else if n > 0 {
				con.PrintInfof("Loaded %d reaction(s) from disk\n", n)
			}
		}

		server.InitDefaultHelpCmd()
		server.SetHelpCommandGroupID(consts.GenericHelpGroup)

		return server
	}

	return serverCommands
}
