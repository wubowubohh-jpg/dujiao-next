package service

import "github.com/dujiao-next/internal/repository"

const complianceVersion = "disabled"

// AcknowledgeRequest is kept for API compatibility with older admin frontends.
type AcknowledgeRequest struct {
	Segment1  string
	Segment2  string
	Segment3  string
	AdminID   uint
	Username  string
	ClientIP  string
	UserAgent string
}

// ComplianceStatus reports the compliance gate as disabled.
type ComplianceStatus struct {
	Acknowledged           bool   `json:"acknowledged"`
	AcknowledgedAt         string `json:"acknowledged_at,omitempty"`
	AcknowledgedByAdminID  uint   `json:"acknowledged_by_admin_id,omitempty"`
	AcknowledgedByUsername string `json:"acknowledged_by_username,omitempty"`
	Version                string `json:"version,omitempty"`
}

// ComplianceService is retained as a no-op compatibility service.
type ComplianceService struct {
	settingRepo repository.SettingRepository
}

func NewComplianceService(repo repository.SettingRepository) *ComplianceService {
	return &ComplianceService{settingRepo: repo}
}

func (s *ComplianceService) IsAcknowledged() bool {
	return true
}

func (s *ComplianceService) Status() (*ComplianceStatus, error) {
	return &ComplianceStatus{
		Acknowledged: true,
		Version:      complianceVersion,
	}, nil
}

func (s *ComplianceService) Acknowledge(req AcknowledgeRequest) error {
	_ = req
	return nil
}
