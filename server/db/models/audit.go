package models

import (
	"time"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// AuditLog - Immutable record of every operator action
type AuditLog struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;"`
	CreatedAt    time.Time `gorm:"->;<-:create;index"`
	OperatorName string    `gorm:"index"`
	Action       string    `gorm:"index"` // e.g. "session.execute", "implant.generate"
	Target       string    // session ID, beacon ID, implant name, etc.
	TargetType   string    // "session", "beacon", "implant", "server"
	Details      string    // JSON or human-readable detail
	Success      bool
	RemoteAddr   string
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	var err error
	a.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	a.CreatedAt = time.Now()
	return nil
}

func (a *AuditLog) ToProtobuf() *clientpb.AuditLog {
	return &clientpb.AuditLog{
		ID:           a.ID.String(),
		CreatedAt:    a.CreatedAt.Unix(),
		OperatorName: a.OperatorName,
		Action:       a.Action,
		Target:       a.Target,
		TargetType:   a.TargetType,
		Details:      a.Details,
		Success:      a.Success,
		RemoteAddr:   a.RemoteAddr,
	}
}
