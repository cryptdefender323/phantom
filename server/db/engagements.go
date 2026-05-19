package db

import (
	"time"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/server/db/models"
	"github.com/gofrs/uuid"
)

// CreateEngagement - Create a new engagement
func CreateEngagement(pb *clientpb.Engagement) (*models.Engagement, error) {
	eng := &models.Engagement{
		Name:        pb.Name,
		Description: pb.Description,
		Scope:       pb.Scope,
		Status:      pb.Status,
		CreatedBy:   pb.CreatedBy,
		Tags:        pb.Tags,
	}
	if pb.StartDate > 0 {
		eng.StartDate = time.Unix(pb.StartDate, 0)
	}
	if pb.EndDate > 0 {
		eng.EndDate = time.Unix(pb.EndDate, 0)
	}
	if eng.Status == "" {
		eng.Status = "active"
	}
	return eng, Session().Create(eng).Error
}

// GetEngagements - List all engagements
func GetEngagements() ([]*clientpb.Engagement, error) {
	var engs []*models.Engagement
	err := Session().
		Preload("Sessions").
		Preload("Beacons").
		Preload("Findings").
		Order("created_at desc").
		Find(&engs).Error
	if err != nil {
		return nil, err
	}
	var result []*clientpb.Engagement
	for _, e := range engs {
		result = append(result, e.ToProtobuf())
	}
	return result, nil
}

// GetEngagementByID - Get a single engagement with all relations
func GetEngagementByID(id string) (*models.Engagement, error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	var eng models.Engagement
	err = Session().
		Preload("Sessions").
		Preload("Beacons").
		Preload("Findings").
		Where("id = ?", uid).
		First(&eng).Error
	if err != nil {
		return nil, err
	}
	return &eng, nil
}

// UpdateEngagement - Update engagement fields
func UpdateEngagement(pb *clientpb.Engagement) (*models.Engagement, error) {
	eng, err := GetEngagementByID(pb.ID)
	if err != nil {
		return nil, err
	}
	eng.Name = pb.Name
	eng.Description = pb.Description
	eng.Scope = pb.Scope
	eng.Status = pb.Status
	eng.Tags = pb.Tags
	eng.UpdatedAt = time.Now()
	if pb.StartDate > 0 {
		eng.StartDate = time.Unix(pb.StartDate, 0)
	}
	if pb.EndDate > 0 {
		eng.EndDate = time.Unix(pb.EndDate, 0)
	}
	return eng, Session().Save(eng).Error
}

// DeleteEngagement - Delete an engagement and all related records
func DeleteEngagement(id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return ErrRecordNotFound
	}
	Session().Where("engagement_id = ?", uid).Delete(&models.EngagementSession{})
	Session().Where("engagement_id = ?", uid).Delete(&models.EngagementBeacon{})
	Session().Where("engagement_id = ?", uid).Delete(&models.Finding{})
	return Session().Where("id = ?", uid).Delete(&models.Engagement{}).Error
}

// AssignSessionToEngagement - Link a session to an engagement
func AssignSessionToEngagement(engagementID, sessionID, addedBy string) error {
	uid, err := uuid.FromString(engagementID)
	if err != nil {
		return ErrRecordNotFound
	}
	link := &models.EngagementSession{
		EngagementID: uid,
		SessionID:    sessionID,
		AddedBy:      addedBy,
	}
	return Session().Create(link).Error
}

// AssignBeaconToEngagement - Link a beacon to an engagement
func AssignBeaconToEngagement(engagementID, beaconID, addedBy string) error {
	uid, err := uuid.FromString(engagementID)
	if err != nil {
		return ErrRecordNotFound
	}
	link := &models.EngagementBeacon{
		EngagementID: uid,
		BeaconID:     beaconID,
		AddedBy:      addedBy,
	}
	return Session().Create(link).Error
}

// RemoveSessionFromEngagement - Unlink a session from an engagement
func RemoveSessionFromEngagement(engagementID, sessionID string) error {
	uid, err := uuid.FromString(engagementID)
	if err != nil {
		return ErrRecordNotFound
	}
	return Session().
		Where("engagement_id = ? AND session_id = ?", uid, sessionID).
		Delete(&models.EngagementSession{}).Error
}

// RemoveBeaconFromEngagement - Unlink a beacon from an engagement
func RemoveBeaconFromEngagement(engagementID, beaconID string) error {
	uid, err := uuid.FromString(engagementID)
	if err != nil {
		return ErrRecordNotFound
	}
	return Session().
		Where("engagement_id = ? AND beacon_id = ?", uid, beaconID).
		Delete(&models.EngagementBeacon{}).Error
}

// AddFinding - Add a finding to an engagement
func AddFinding(pb *clientpb.Finding) (*models.Finding, error) {
	uid, err := uuid.FromString(pb.EngagementID)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	f := &models.Finding{
		EngagementID: uid,
		Title:        pb.Title,
		Description:  pb.Description,
		Severity:     pb.Severity,
		Host:         pb.Host,
		Evidence:     pb.Evidence,
		Remediation:  pb.Remediation,
		CreatedBy:    pb.CreatedBy,
	}
	return f, Session().Create(f).Error
}

// UpdateFinding - Update a finding
func UpdateFinding(pb *clientpb.Finding) (*models.Finding, error) {
	uid, err := uuid.FromString(pb.ID)
	if err != nil {
		return nil, ErrRecordNotFound
	}
	var f models.Finding
	if err := Session().Where("id = ?", uid).First(&f).Error; err != nil {
		return nil, err
	}
	f.Title = pb.Title
	f.Description = pb.Description
	f.Severity = pb.Severity
	f.Host = pb.Host
	f.Evidence = pb.Evidence
	f.Remediation = pb.Remediation
	f.UpdatedAt = time.Now()
	return &f, Session().Save(&f).Error
}

// DeleteFinding - Delete a finding
func DeleteFinding(id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return ErrRecordNotFound
	}
	return Session().Where("id = ?", uid).Delete(&models.Finding{}).Error
}
