package db

import (
	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/server/db/models"
)

// AuditLog - Write an audit log entry
func AuditLog(operatorName, action, target, targetType, details string, success bool, remoteAddr string) error {
	entry := &models.AuditLog{
		OperatorName: operatorName,
		Action:       action,
		Target:       target,
		TargetType:   targetType,
		Details:      details,
		Success:      success,
		RemoteAddr:   remoteAddr,
	}
	return Session().Create(entry).Error
}

// GetAuditLogs - Query audit logs with optional filters
func GetAuditLogs(filter *clientpb.AuditLogFilter) ([]*clientpb.AuditLog, error) {
	query := Session().Model(&models.AuditLog{}).Order("created_at desc")

	if filter != nil {
		if filter.OperatorName != "" {
			query = query.Where("operator_name = ?", filter.OperatorName)
		}
		if filter.Action != "" {
			query = query.Where("action LIKE ?", "%"+filter.Action+"%")
		}
		if filter.TargetType != "" {
			query = query.Where("target_type = ?", filter.TargetType)
		}
		if filter.Since > 0 {
			query = query.Where("created_at >= datetime(?, 'unixepoch')", filter.Since)
		}
		if filter.Until > 0 {
			query = query.Where("created_at <= datetime(?, 'unixepoch')", filter.Until)
		}
		if filter.Limit > 0 {
			query = query.Limit(int(filter.Limit))
		} else {
			query = query.Limit(500) // default cap
		}
	}

	var entries []*models.AuditLog
	if err := query.Find(&entries).Error; err != nil {
		return nil, err
	}

	var result []*clientpb.AuditLog
	for _, e := range entries {
		result = append(result, e.ToProtobuf())
	}
	return result, nil
}
