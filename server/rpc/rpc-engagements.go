package rpc

import (
	"context"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/server/db"
	"github.com/cryptdefender3232/phantom/server/log"
)

var engagementRpcLog = log.NamedLogger("rpc", "engagements")

// CreateEngagement - Create a new engagement
func (rpc *Server) CreateEngagement(ctx context.Context, req *clientpb.Engagement) (*clientpb.Engagement, error) {
	operator := rpc.getClientCommonName(ctx)
	req.CreatedBy = operator

	eng, err := db.CreateEngagement(req)
	if err != nil {
		engagementRpcLog.Errorf("Failed to create engagement: %s", err)
		return nil, ErrDatabaseFailure
	}

	db.AuditLog(operator, "engagement.create", eng.ID.String(), "engagement",
		"Created engagement: "+eng.Name, true, "")

	return eng.ToProtobuf(), nil
}

// GetEngagements - List all engagements
func (rpc *Server) GetEngagements(ctx context.Context, _ *commonpb.Empty) (*clientpb.Engagements, error) {
	engs, err := db.GetEngagements()
	if err != nil {
		engagementRpcLog.Errorf("Failed to list engagements: %s", err)
		return nil, ErrDatabaseFailure
	}
	return &clientpb.Engagements{Engagements: engs}, nil
}

// GetEngagement - Get a single engagement by ID
func (rpc *Server) GetEngagement(ctx context.Context, req *clientpb.EngagementReq) (*clientpb.Engagement, error) {
	eng, err := db.GetEngagementByID(req.ID)
	if err != nil {
		return nil, rpcError(err)
	}
	return eng.ToProtobuf(), nil
}

// UpdateEngagement - Update an existing engagement
func (rpc *Server) UpdateEngagement(ctx context.Context, req *clientpb.Engagement) (*clientpb.Engagement, error) {
	operator := rpc.getClientCommonName(ctx)
	eng, err := db.UpdateEngagement(req)
	if err != nil {
		engagementRpcLog.Errorf("Failed to update engagement: %s", err)
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.update", eng.ID.String(), "engagement",
		"Updated engagement: "+eng.Name, true, "")
	return eng.ToProtobuf(), nil
}

// DeleteEngagement - Delete an engagement
func (rpc *Server) DeleteEngagement(ctx context.Context, req *clientpb.EngagementReq) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.DeleteEngagement(req.ID); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.delete", req.ID, "engagement",
		"Deleted engagement", true, "")
	return &commonpb.Empty{}, nil
}

// AssignSessionToEngagement - Link a session to an engagement
func (rpc *Server) AssignSessionToEngagement(ctx context.Context, req *clientpb.AssignTargetReq) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.AssignSessionToEngagement(req.EngagementID, req.TargetID, operator); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.assign_session", req.TargetID, "session",
		"Assigned session to engagement "+req.EngagementID, true, "")
	return &commonpb.Empty{}, nil
}

// AssignBeaconToEngagement - Link a beacon to an engagement
func (rpc *Server) AssignBeaconToEngagement(ctx context.Context, req *clientpb.AssignTargetReq) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.AssignBeaconToEngagement(req.EngagementID, req.TargetID, operator); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.assign_beacon", req.TargetID, "beacon",
		"Assigned beacon to engagement "+req.EngagementID, true, "")
	return &commonpb.Empty{}, nil
}

// RemoveSessionFromEngagement - Unlink a session from an engagement
func (rpc *Server) RemoveSessionFromEngagement(ctx context.Context, req *clientpb.AssignTargetReq) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.RemoveSessionFromEngagement(req.EngagementID, req.TargetID); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.remove_session", req.TargetID, "session",
		"Removed session from engagement "+req.EngagementID, true, "")
	return &commonpb.Empty{}, nil
}

// RemoveBeaconFromEngagement - Unlink a beacon from an engagement
func (rpc *Server) RemoveBeaconFromEngagement(ctx context.Context, req *clientpb.AssignTargetReq) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.RemoveBeaconFromEngagement(req.EngagementID, req.TargetID); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "engagement.remove_beacon", req.TargetID, "beacon",
		"Removed beacon from engagement "+req.EngagementID, true, "")
	return &commonpb.Empty{}, nil
}

// AddFinding - Add a finding to an engagement
func (rpc *Server) AddFinding(ctx context.Context, req *clientpb.Finding) (*clientpb.Finding, error) {
	operator := rpc.getClientCommonName(ctx)
	req.CreatedBy = operator
	f, err := db.AddFinding(req)
	if err != nil {
		engagementRpcLog.Errorf("Failed to add finding: %s", err)
		return nil, ErrDatabaseFailure
	}
	db.AuditLog(operator, "finding.add", f.ID.String(), "finding",
		"Added finding: "+f.Title+" ["+f.Severity+"]", true, "")
	return f.ToProtobuf(), nil
}

// UpdateFinding - Update a finding
func (rpc *Server) UpdateFinding(ctx context.Context, req *clientpb.Finding) (*clientpb.Finding, error) {
	operator := rpc.getClientCommonName(ctx)
	f, err := db.UpdateFinding(req)
	if err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "finding.update", f.ID.String(), "finding",
		"Updated finding: "+f.Title, true, "")
	return f.ToProtobuf(), nil
}

// DeleteFinding - Delete a finding
func (rpc *Server) DeleteFinding(ctx context.Context, req *clientpb.Finding) (*commonpb.Empty, error) {
	operator := rpc.getClientCommonName(ctx)
	if err := db.DeleteFinding(req.ID); err != nil {
		return nil, rpcError(err)
	}
	db.AuditLog(operator, "finding.delete", req.ID, "finding",
		"Deleted finding", true, "")
	return &commonpb.Empty{}, nil
}
