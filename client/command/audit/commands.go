package audit

import (
	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/spf13/cobra"
)

// Commands returns audit log commands
func Commands(con *console.PhantomClient) []*cobra.Command {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Operator audit log",
		Long:  "View and filter the operator audit log — every action taken by every operator",
		Run: func(cmd *cobra.Command, args []string) {
			AuditLogsCmd(cmd, con)
		},
	}

	auditCmd.Flags().StringP("operator", "o", "", "filter by operator name")
	auditCmd.Flags().StringP("action", "a", "", "filter by action (e.g. session, implant, engagement)")
	auditCmd.Flags().StringP("type", "t", "", "filter by target type (session, beacon, implant, engagement, finding)")
	auditCmd.Flags().IntP("limit", "n", 50, "max number of entries to show (0 = all)")
	auditCmd.Flags().StringP("since", "s", "", "show entries since date (YYYY-MM-DD)")
	auditCmd.Flags().StringP("until", "u", "", "show entries until date (YYYY-MM-DD)")

	return []*cobra.Command{auditCmd}
}
