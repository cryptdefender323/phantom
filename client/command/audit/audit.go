package audit

import (
	"fmt"
	"time"

	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// AuditLogsCmd - Display audit log entries
func AuditLogsCmd(cmd *cobra.Command, con *console.PhantomClient) {
	operatorFilter, _ := cmd.Flags().GetString("operator")
	actionFilter, _ := cmd.Flags().GetString("action")
	typeFilter, _ := cmd.Flags().GetString("type")
	limit, _ := cmd.Flags().GetInt("limit")
	sinceStr, _ := cmd.Flags().GetString("since")
	untilStr, _ := cmd.Flags().GetString("until")

	filter := &clientpb.AuditLogFilter{
		OperatorName: operatorFilter,
		Action:       actionFilter,
		TargetType:   typeFilter,
		Limit:        int32(limit),
	}

	if sinceStr != "" {
		t, err := time.Parse("2006-01-02", sinceStr)
		if err != nil {
			con.PrintErrorf("Invalid --since date format, use YYYY-MM-DD\n")
			return
		}
		filter.Since = t.Unix()
	}
	if untilStr != "" {
		t, err := time.Parse("2006-01-02", untilStr)
		if err != nil {
			con.PrintErrorf("Invalid --until date format, use YYYY-MM-DD\n")
			return
		}
		filter.Until = t.Add(24 * time.Hour).Unix()
	}

	logs, err := con.Rpc.GetAuditLogs(context.Background(), filter)
	if err != nil {
		con.PrintErrorf("Failed to get audit logs: %s\n", err)
		return
	}

	if len(logs.Logs) == 0 {
		con.PrintInfof("No audit log entries found\n")
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.AppendHeader(table.Row{"Time", "Operator", "Action", "Target", "Type", "Status"})

	for _, entry := range logs.Logs {
		ts := time.Unix(entry.CreatedAt, 0).Format("2006-01-02 15:04:05")
		status := "✓"
		if !entry.Success {
			status = "✗"
		}
		target := entry.Target
		if len(target) > 20 {
			target = target[:20] + "..."
		}
		tw.AppendRow(table.Row{
			ts,
			entry.OperatorName,
			entry.Action,
			target,
			entry.TargetType,
			status,
		})
	}

	fmt.Println(tw.Render())
	fmt.Printf("\n%d entries\n", len(logs.Logs))
}
