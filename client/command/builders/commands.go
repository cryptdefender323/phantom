package builders

import (
	"github.com/cryptdefender3232/phantom/client/command/flags"
	"github.com/cryptdefender3232/phantom/client/command/help"
	"github.com/cryptdefender3232/phantom/client/console"
	consts "github.com/cryptdefender3232/phantom/client/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Commands returns the “ command and its subcommands.
func Commands(con *console.PhantomClient) []*cobra.Command {
	buildersCmd := &cobra.Command{
		Use:   consts.BuildersStr,
		Short: "List external builders",
		Long:  help.GetHelpFor([]string{consts.BuildersStr}),
		Run: func(cmd *cobra.Command, args []string) {
			BuildersCmd(cmd, con, args)
		},
		GroupID: consts.PayloadsHelpGroup,
	}
	flags.Bind("builders", false, buildersCmd, func(f *pflag.FlagSet) {
		f.Int64P("timeout", "t", flags.DefaultTimeout, "grpc timeout in seconds")
	})

	return []*cobra.Command{buildersCmd}
}
