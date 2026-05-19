package models

import (
	"time"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Engagement - Represents a red team engagement / operation
type Engagement struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;"`
	CreatedAt   time.Time `gorm:"->;<-:create;"`
	UpdatedAt   time.Time
	Name        string `gorm:"uniqueIndex"`
	Description string
	Scope       string // IP ranges, domains, etc.
	StartDate   time.Time
	EndDate     time.Time
	Status      string // "active", "completed", "paused"
	CreatedBy   string
	Tags        string // comma-separated tags

	Sessions []EngagementSession
	Beacons  []EngagementBeacon
	Findings []Finding
}

func (e *Engagement) BeforeCreate(tx *gorm.DB) error {
	var err error
	e.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
	if e.Status == "" {
		e.Status = "active"
	}
	return nil
}

func (e *Engagement) ToProtobuf() *clientpb.Engagement {
	pb := &clientpb.Engagement{
		ID:          e.ID.String(),
		CreatedAt:   e.CreatedAt.Unix(),
		UpdatedAt:   e.UpdatedAt.Unix(),
		Name:        e.Name,
		Description: e.Description,
		Scope:       e.Scope,
		StartDate:   e.StartDate.Unix(),
		EndDate:     e.EndDate.Unix(),
		Status:      e.Status,
		CreatedBy:   e.CreatedBy,
		Tags:        e.Tags,
	}
	for _, s := range e.Sessions {
		pb.SessionIDs = append(pb.SessionIDs, s.SessionID)
	}
	for _, b := range e.Beacons {
		pb.BeaconIDs = append(pb.BeaconIDs, b.BeaconID)
	}
	for _, f := range e.Findings {
		pb.Findings = append(pb.Findings, f.ToProtobuf())
	}
	return pb
}

// EngagementSession - Join table linking sessions to engagements
type EngagementSession struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;"`
	CreatedAt    time.Time `gorm:"->;<-:create;"`
	EngagementID uuid.UUID `gorm:"type:uuid;index"`
	SessionID    string    `gorm:"index"`
	AddedBy      string
}

func (e *EngagementSession) BeforeCreate(tx *gorm.DB) error {
	var err error
	e.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	e.CreatedAt = time.Now()
	return nil
}

// EngagementBeacon - Join table linking beacons to engagements
type EngagementBeacon struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;"`
	CreatedAt    time.Time `gorm:"->;<-:create;"`
	EngagementID uuid.UUID `gorm:"type:uuid;index"`
	BeaconID     string    `gorm:"index"`
	AddedBy      string
}

func (e *EngagementBeacon) BeforeCreate(tx *gorm.DB) error {
	var err error
	e.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	e.CreatedAt = time.Now()
	return nil
}

// Finding - A finding/vulnerability discovered during an engagement
type Finding struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;"`
	CreatedAt    time.Time `gorm:"->;<-:create;"`
	UpdatedAt    time.Time
	EngagementID uuid.UUID `gorm:"type:uuid;index"`
	Title        string
	Description  string
	Severity     string // "critical", "high", "medium", "low", "info"
	Host         string // affected host
	Evidence     string // commands run, screenshots, etc.
	Remediation  string
	CreatedBy    string
}

func (f *Finding) BeforeCreate(tx *gorm.DB) error {
	var err error
	f.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	return nil
}

func (f *Finding) ToProtobuf() *clientpb.Finding {
	return &clientpb.Finding{
		ID:           f.ID.String(),
		CreatedAt:    f.CreatedAt.Unix(),
		UpdatedAt:    f.UpdatedAt.Unix(),
		EngagementID: f.EngagementID.String(),
		Title:        f.Title,
		Description:  f.Description,
		Severity:     f.Severity,
		Host:         f.Host,
		Evidence:     f.Evidence,
		Remediation:  f.Remediation,
		CreatedBy:    f.CreatedBy,
	}
}
