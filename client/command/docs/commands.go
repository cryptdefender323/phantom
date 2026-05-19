package docs

import (
	"github.com/cryptdefender3232/phantom/client/command/help"
	"github.com/cryptdefender3232/phantom/client/console"
	consts "github.com/cryptdefender3232/phantom/client/constants"
	"github.com/spf13/cobra"
)

// Commands returns the docs command.
func Commands(con *console.PhantomClient) []*cobra.Command {
	return []*cobra.Command{newDocsCommand(consts.PhantomCoreHelpGroup, con)}
}

// ServerCommands returns the docs command for the top-level client REPL.
func ServerCommands(con *console.PhantomClient) []*cobra.Command {
	return []*cobra.Command{newDocsCommand(consts.GenericHelpGroup, con)}
}

func newDocsCommand(groupID string, con *console.PhantomClient) *cobra.Command {
	return &cobra.Command{
		Use:     consts.DocsStr,
		Short:   "Browse the embedded Phantom docs in a TUI",
		Long:    help.GetHelpFor([]string{consts.DocsStr}),
		Args:    cobra.NoArgs,
		GroupID: groupID,
		Run: func(cmd *cobra.Command, args []string) {
			DocsCmd(cmd, con, args)
		},
	}
}
