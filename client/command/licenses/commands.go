package licenses

import (
	"github.com/spf13/cobra"

	"github.com/cryptdefender3232/phantom/client/command/help"
	"github.com/cryptdefender3232/phantom/client/console"
	consts "github.com/cryptdefender3232/phantom/client/constants"
	"github.com/cryptdefender3232/phantom/client/licenses"
)

// Commands returns the `licences` command.
func Commands(con *console.PhantomClient) []*cobra.Command {
	licensesCmd := &cobra.Command{
		Use:   consts.LicensesStr,
		Short: "Open source licenses",
		Long:  help.GetHelpFor([]string{consts.LicensesStr}),
		Run: func(cmd *cobra.Command, args []string) {
			con.Println(licenses.All)
		},
		GroupID: consts.GenericHelpGroup,
	}

	return []*cobra.Command{licensesCmd}
}
