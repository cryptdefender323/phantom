package engagements

import (
	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/spf13/cobra"
)

// Commands returns engagement management commands
func Commands(con *console.PhantomClient) []*cobra.Command {
	engCmd := &cobra.Command{
		Use:   "engagements",
		Short: "Manage red team engagements",
		Long:  "Create and manage engagements, assign sessions/beacons, track findings",
		Run: func(cmd *cobra.Command, args []string) {
			ListEngagementsCmd(cmd, con)
		},
	}

	// list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all engagements",
		Run: func(cmd *cobra.Command, args []string) {
			ListEngagementsCmd(cmd, con)
		},
	}

	// create
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new engagement",
		Run: func(cmd *cobra.Command, args []string) {
			CreateEngagementCmd(cmd, con)
		},
	}
	createCmd.Flags().StringP("name", "n", "", "engagement name (required)")
	createCmd.Flags().StringP("description", "d", "", "description")
	createCmd.Flags().StringP("scope", "s", "", "scope (IP ranges, domains, etc.)")
	createCmd.Flags().StringP("start", "", "", "start date (YYYY-MM-DD)")
	createCmd.Flags().StringP("end", "", "", "end date (YYYY-MM-DD)")
	createCmd.Flags().StringP("tags", "t", "", "comma-separated tags")
	createCmd.MarkFlagRequired("name")

	// info
	infoCmd := &cobra.Command{
		Use:   "info <engagement-id>",
		Short: "Show engagement details",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			EngagementInfoCmd(cmd, con, args[0])
		},
	}

	// update
	updateCmd := &cobra.Command{
		Use:   "update <engagement-id>",
		Short: "Update an engagement",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			UpdateEngagementCmd(cmd, con, args[0])
		},
	}
	updateCmd.Flags().StringP("name", "n", "", "new name")
	updateCmd.Flags().StringP("description", "d", "", "new description")
	updateCmd.Flags().StringP("scope", "s", "", "new scope")
	updateCmd.Flags().StringP("status", "", "", "status: active, paused, completed")
	updateCmd.Flags().StringP("tags", "t", "", "new tags")

	// delete
	deleteCmd := &cobra.Command{
		Use:   "delete <engagement-id>",
		Short: "Delete an engagement",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			DeleteEngagementCmd(cmd, con, args[0])
		},
	}

	// assign-session
	assignSessionCmd := &cobra.Command{
		Use:   "assign-session <engagement-id> <session-id>",
		Short: "Assign a session to an engagement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			AssignSessionCmd(cmd, con, args[0], args[1])
		},
	}

	// assign-beacon
	assignBeaconCmd := &cobra.Command{
		Use:   "assign-beacon <engagement-id> <beacon-id>",
		Short: "Assign a beacon to an engagement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			AssignBeaconCmd(cmd, con, args[0], args[1])
		},
	}

	// remove-session
	removeSessionCmd := &cobra.Command{
		Use:   "remove-session <engagement-id> <session-id>",
		Short: "Remove a session from an engagement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			RemoveSessionCmd(cmd, con, args[0], args[1])
		},
	}

	// remove-beacon
	removeBeaconCmd := &cobra.Command{
		Use:   "remove-beacon <engagement-id> <beacon-id>",
		Short: "Remove a beacon from an engagement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			RemoveBeaconCmd(cmd, con, args[0], args[1])
		},
	}

	// add-finding
	addFindingCmd := &cobra.Command{
		Use:   "add-finding <engagement-id>",
		Short: "Add a finding to an engagement",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			AddFindingCmd(cmd, con, args[0])
		},
	}
	addFindingCmd.Flags().StringP("title", "t", "", "finding title (required)")
	addFindingCmd.Flags().StringP("severity", "s", "medium", "severity: critical, high, medium, low, info")
	addFindingCmd.Flags().StringP("description", "d", "", "finding description")
	addFindingCmd.Flags().StringP("host", "", "", "affected host")
	addFindingCmd.Flags().StringP("evidence", "e", "", "evidence (commands, output, etc.)")
	addFindingCmd.Flags().StringP("remediation", "r", "", "remediation recommendation")
	addFindingCmd.MarkFlagRequired("title")

	// findings
	findingsCmd := &cobra.Command{
		Use:   "findings <engagement-id>",
		Short: "List findings for an engagement",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ListFindingsCmd(cmd, con, args[0])
		},
	}

	engCmd.AddCommand(listCmd, createCmd, infoCmd, updateCmd, deleteCmd,
		assignSessionCmd, assignBeaconCmd, removeSessionCmd, removeBeaconCmd,
		addFindingCmd, findingsCmd)

	return []*cobra.Command{engCmd}
}
