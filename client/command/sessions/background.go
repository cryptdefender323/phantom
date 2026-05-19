package sessions

import (
	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/spf13/cobra"
)

// BackgroundCmd - Background the active session.
func BackgroundCmd(cmd *cobra.Command, con *console.PhantomClient, args []string) {
	con.ActiveTarget.Background()
	con.PrintInfof("Background ...\n")
}
