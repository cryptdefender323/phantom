package monitor

import (
	"github.com/spf13/cobra"

	"github.com/cryptdefender3232/phantom/client/console"
	consts "github.com/cryptdefender3232/phantom/client/constants"
)

// Commands returns the “ command and its subcommands.
func Commands(con *console.PhantomClient) []*cobra.Command {
	monitorCmd := &cobra.Command{
		Use:     consts.MonitorStr,
		Short:   "Monitor threat intel platforms for Phantom implants",
		GroupID: consts.PhantomHelpGroup,
	}
	monitorCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the monitoring loops",
		Run: func(cmd *cobra.Command, args []string) {
			MonitorStartCmd(cmd, con, args)
		},
	})
	monitorCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop the monitoring loops",
		Run: func(cmd *cobra.Command, args []string) {
			MonitorStopCmd(cmd, con, args)
		},
	})

	return []*cobra.Command{monitorCmd}
}
