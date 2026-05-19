package engagements

import (
	"fmt"
	"strings"
	"time"

	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// ListEngagementsCmd - List all engagements
func ListEngagementsCmd(cmd *cobra.Command, con *console.PhantomClient) {
	engs, err := con.Rpc.GetEngagements(context.Background(), &commonpb.Empty{})
	if err != nil {
		con.PrintErrorf("Failed to get engagements: %s\n", err)
		return
	}
	if len(engs.Engagements) == 0 {
		con.PrintInfof("No engagements found. Create one with: engagements create --name <name>\n")
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.AppendHeader(table.Row{"ID", "Name", "Status", "Sessions", "Beacons", "Findings", "Created By", "Start"})

	for _, e := range engs.Engagements {
		start := "-"
		if e.StartDate > 0 {
			start = time.Unix(e.StartDate, 0).Format("2006-01-02")
		}
		shortID := e.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		tw.AppendRow(table.Row{
			shortID,
			e.Name,
			statusBadge(e.Status),
			len(e.SessionIDs),
			len(e.BeaconIDs),
			len(e.Findings),
			e.CreatedBy,
			start,
		})
	}
	fmt.Println(tw.Render())
}

// CreateEngagementCmd - Create a new engagement
func CreateEngagementCmd(cmd *cobra.Command, con *console.PhantomClient) {
	name, _ := cmd.Flags().GetString("name")
	desc, _ := cmd.Flags().GetString("description")
	scope, _ := cmd.Flags().GetString("scope")
	startStr, _ := cmd.Flags().GetString("start")
	endStr, _ := cmd.Flags().GetString("end")
	tags, _ := cmd.Flags().GetString("tags")

	req := &clientpb.Engagement{
		Name:        name,
		Description: desc,
		Scope:       scope,
		Tags:        tags,
		Status:      "active",
	}

	if startStr != "" {
		t, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			con.PrintErrorf("Invalid --start date, use YYYY-MM-DD\n")
			return
		}
		req.StartDate = t.Unix()
	}
	if endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			con.PrintErrorf("Invalid --end date, use YYYY-MM-DD\n")
			return
		}
		req.EndDate = t.Unix()
	}

	eng, err := con.Rpc.CreateEngagement(context.Background(), req)
	if err != nil {
		con.PrintErrorf("Failed to create engagement: %s\n", err)
		return
	}
	con.PrintInfof("Created engagement: %s (ID: %s)\n", eng.Name, eng.ID[:8])
}

// EngagementInfoCmd - Show engagement details
func EngagementInfoCmd(cmd *cobra.Command, con *console.PhantomClient, id string) {
	eng, err := con.Rpc.GetEngagement(context.Background(), &clientpb.EngagementReq{ID: id})
	if err != nil {
		con.PrintErrorf("Engagement not found: %s\n", err)
		return
	}

	fmt.Printf("\n  Name:        %s\n", eng.Name)
	fmt.Printf("  ID:          %s\n", eng.ID)
	fmt.Printf("  Status:      %s\n", statusBadge(eng.Status))
	fmt.Printf("  Created by:  %s\n", eng.CreatedBy)
	fmt.Printf("  Created at:  %s\n", time.Unix(eng.CreatedAt, 0).Format("2006-01-02 15:04:05"))
	if eng.Description != "" {
		fmt.Printf("  Description: %s\n", eng.Description)
	}
	if eng.Scope != "" {
		fmt.Printf("  Scope:       %s\n", eng.Scope)
	}
	if eng.Tags != "" {
		fmt.Printf("  Tags:        %s\n", eng.Tags)
	}
	if eng.StartDate > 0 {
		fmt.Printf("  Start:       %s\n", time.Unix(eng.StartDate, 0).Format("2006-01-02"))
	}
	if eng.EndDate > 0 {
		fmt.Printf("  End:         %s\n", time.Unix(eng.EndDate, 0).Format("2006-01-02"))
	}

	fmt.Printf("\n  Sessions (%d): %s\n", len(eng.SessionIDs), strings.Join(eng.SessionIDs, ", "))
	fmt.Printf("  Beacons  (%d): %s\n", len(eng.BeaconIDs), strings.Join(eng.BeaconIDs, ", "))

	if len(eng.Findings) > 0 {
		fmt.Printf("\n  Findings (%d):\n", len(eng.Findings))
		for _, f := range eng.Findings {
			fmt.Printf("    [%s] %s — %s\n", severityBadge(f.Severity), f.Title, f.Host)
		}
	}
	fmt.Println()
}

// UpdateEngagementCmd - Update an engagement
func UpdateEngagementCmd(cmd *cobra.Command, con *console.PhantomClient, id string) {
	eng, err := con.Rpc.GetEngagement(context.Background(), &clientpb.EngagementReq{ID: id})
	if err != nil {
		con.PrintErrorf("Engagement not found: %s\n", err)
		return
	}

	if name, _ := cmd.Flags().GetString("name"); name != "" {
		eng.Name = name
	}
	if desc, _ := cmd.Flags().GetString("description"); desc != "" {
		eng.Description = desc
	}
	if scope, _ := cmd.Flags().GetString("scope"); scope != "" {
		eng.Scope = scope
	}
	if status, _ := cmd.Flags().GetString("status"); status != "" {
		eng.Status = status
	}
	if tags, _ := cmd.Flags().GetString("tags"); tags != "" {
		eng.Tags = tags
	}

	updated, err := con.Rpc.UpdateEngagement(context.Background(), eng)
	if err != nil {
		con.PrintErrorf("Failed to update engagement: %s\n", err)
		return
	}
	con.PrintInfof("Updated engagement: %s\n", updated.Name)
}

// DeleteEngagementCmd - Delete an engagement
func DeleteEngagementCmd(cmd *cobra.Command, con *console.PhantomClient, id string) {
	_, err := con.Rpc.DeleteEngagement(context.Background(), &clientpb.EngagementReq{ID: id})
	if err != nil {
		con.PrintErrorf("Failed to delete engagement: %s\n", err)
		return
	}
	con.PrintInfof("Deleted engagement %s\n", id[:8])
}

// AssignSessionCmd - Assign a session to an engagement
func AssignSessionCmd(cmd *cobra.Command, con *console.PhantomClient, engID, sessionID string) {
	_, err := con.Rpc.AssignSessionToEngagement(context.Background(), &clientpb.AssignTargetReq{
		EngagementID: engID,
		TargetID:     sessionID,
	})
	if err != nil {
		con.PrintErrorf("Failed to assign session: %s\n", err)
		return
	}
	con.PrintInfof("Session %s assigned to engagement %s\n", sessionID[:8], engID[:8])
}

// AssignBeaconCmd - Assign a beacon to an engagement
func AssignBeaconCmd(cmd *cobra.Command, con *console.PhantomClient, engID, beaconID string) {
	_, err := con.Rpc.AssignBeaconToEngagement(context.Background(), &clientpb.AssignTargetReq{
		EngagementID: engID,
		TargetID:     beaconID,
	})
	if err != nil {
		con.PrintErrorf("Failed to assign beacon: %s\n", err)
		return
	}
	con.PrintInfof("Beacon %s assigned to engagement %s\n", beaconID[:8], engID[:8])
}

// RemoveSessionCmd - Remove a session from an engagement
func RemoveSessionCmd(cmd *cobra.Command, con *console.PhantomClient, engID, sessionID string) {
	_, err := con.Rpc.RemoveSessionFromEngagement(context.Background(), &clientpb.AssignTargetReq{
		EngagementID: engID,
		TargetID:     sessionID,
	})
	if err != nil {
		con.PrintErrorf("Failed to remove session: %s\n", err)
		return
	}
	con.PrintInfof("Session removed from engagement\n")
}

// RemoveBeaconCmd - Remove a beacon from an engagement
func RemoveBeaconCmd(cmd *cobra.Command, con *console.PhantomClient, engID, beaconID string) {
	_, err := con.Rpc.RemoveBeaconFromEngagement(context.Background(), &clientpb.AssignTargetReq{
		EngagementID: engID,
		TargetID:     beaconID,
	})
	if err != nil {
		con.PrintErrorf("Failed to remove beacon: %s\n", err)
		return
	}
	con.PrintInfof("Beacon removed from engagement\n")
}

// AddFindingCmd - Add a finding to an engagement
func AddFindingCmd(cmd *cobra.Command, con *console.PhantomClient, engID string) {
	title, _ := cmd.Flags().GetString("title")
	severity, _ := cmd.Flags().GetString("severity")
	desc, _ := cmd.Flags().GetString("description")
	host, _ := cmd.Flags().GetString("host")
	evidence, _ := cmd.Flags().GetString("evidence")
	remediation, _ := cmd.Flags().GetString("remediation")

	f, err := con.Rpc.AddFinding(context.Background(), &clientpb.Finding{
		EngagementID: engID,
		Title:        title,
		Severity:     severity,
		Description:  desc,
		Host:         host,
		Evidence:     evidence,
		Remediation:  remediation,
	})
	if err != nil {
		con.PrintErrorf("Failed to add finding: %s\n", err)
		return
	}
	con.PrintInfof("Added finding [%s] %s (ID: %s)\n", severityBadge(f.Severity), f.Title, f.ID[:8])
}

// ListFindingsCmd - List findings for an engagement
func ListFindingsCmd(cmd *cobra.Command, con *console.PhantomClient, engID string) {
	eng, err := con.Rpc.GetEngagement(context.Background(), &clientpb.EngagementReq{ID: engID})
	if err != nil {
		con.PrintErrorf("Engagement not found: %s\n", err)
		return
	}
	if len(eng.Findings) == 0 {
		con.PrintInfof("No findings for this engagement\n")
		return
	}

	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.AppendHeader(table.Row{"ID", "Severity", "Title", "Host", "Created By", "Date"})

	for _, f := range eng.Findings {
		shortID := f.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		tw.AppendRow(table.Row{
			shortID,
			severityBadge(f.Severity),
			f.Title,
			f.Host,
			f.CreatedBy,
			time.Unix(f.CreatedAt, 0).Format("2006-01-02"),
		})
	}
	fmt.Println(tw.Render())
}

func statusBadge(status string) string {
	switch status {
	case "active":
		return "● active"
	case "completed":
		return "✓ completed"
	case "paused":
		return "⏸ paused"
	default:
		return status
	}
}

func severityBadge(severity string) string {
	switch strings.ToLower(severity) {
	case "critical":
		return "CRITICAL"
	case "high":
		return "HIGH"
	case "medium":
		return "MEDIUM"
	case "low":
		return "LOW"
	case "info":
		return "INFO"
	default:
		return severity
	}
}
