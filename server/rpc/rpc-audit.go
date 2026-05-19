package rpc

import (
	"context"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/server/db"
	"github.com/cryptdefender3232/phantom/server/log"
)

var auditRpcLog = log.NamedLogger("rpc", "audit")

// GetAuditLogs - Return audit log entries with optional filtering
func (rpc *Server) GetAuditLogs(ctx context.Context, filter *clientpb.AuditLogFilter) (*clientpb.AuditLogs, error) {
	logs, err := db.GetAuditLogs(filter)
	if err != nil {
		auditRpcLog.Errorf("Failed to query audit logs: %s", err)
		return nil, ErrDatabaseFailure
	}
	return &clientpb.AuditLogs{Logs: logs}, nil
}
