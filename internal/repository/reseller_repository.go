package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/dujiao-next/internal/models"
	"gorm.io/gorm"
)

// ResellerRepository 分销商数据访问接口。
type ResellerRepository interface {
	Transaction(fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) ResellerRepository
	CreateProfile(profile *models.ResellerProfile) error
	GetProfileByID(id uint) (*models.ResellerProfile, error)
	GetProfileByUserID(userID uint) (*models.ResellerProfile, error)
	UpsertDomain(domain models.ResellerDomain) (*models.ResellerDomain, error)
	FindDomainByHost(host string) (*models.ResellerDomain, error)
	FindActiveVerifiedDomain(host string) (*models.ResellerDomain, error)
	UpsertSiteConfig(config models.ResellerSiteConfig) (*models.ResellerSiteConfig, error)
}

// GormResellerRepository GORM 分销商仓储。
type GormResellerRepository struct {
	BaseRepository
}

// NewResellerRepository 创建分销商仓储。
func NewResellerRepository(db *gorm.DB) *GormResellerRepository {
	return &GormResellerRepository{BaseRepository: BaseRepository{db: db}}
}

// WithTx 绑定事务。
func (r *GormResellerRepository) WithTx(tx *gorm.DB) ResellerRepository {
	if tx == nil {
		return r
	}
	return &GormResellerRepository{BaseRepository: BaseRepository{db: tx}}
}

// CreateProfile 创建分销商资料。
func (r *GormResellerRepository) CreateProfile(profile *models.ResellerProfile) error {
	if profile == nil {
		return errors.New("reseller profile is nil")
	}
	return r.db.Create(profile).Error
}

// GetProfileByID 按 ID 获取分销商资料。
func (r *GormResellerRepository) GetProfileByID(id uint) (*models.ResellerProfile, error) {
	if id == 0 {
		return nil, nil
	}
	var profile models.ResellerProfile
	if err := r.db.Preload("User").First(&profile, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// GetProfileByUserID 按用户 ID 获取分销商资料。
func (r *GormResellerRepository) GetProfileByUserID(userID uint) (*models.ResellerProfile, error) {
	if userID == 0 {
		return nil, nil
	}
	var profile models.ResellerProfile
	if err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// UpsertDomain 创建域名，或恢复同域名的软删除记录。
func (r *GormResellerRepository) UpsertDomain(input models.ResellerDomain) (*models.ResellerDomain, error) {
	input.Domain = normalizeDomainForRepository(input.Domain)
	if input.ResellerID == 0 || input.Domain == "" {
		return nil, errors.New("invalid reseller domain")
	}
	now := time.Now()
	var existing models.ResellerDomain
	err := r.db.Unscoped().Where("domain = ?", input.Domain).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		input.CreatedAt = now
		input.UpdatedAt = now
		if err := r.db.Create(&input).Error; err != nil {
			return nil, err
		}
		return &input, nil
	}
	if !existing.DeletedAt.Valid {
		return nil, errors.New("reseller domain already exists")
	}
	existing.ResellerID = input.ResellerID
	existing.Type = input.Type
	existing.VerificationToken = input.VerificationToken
	existing.VerificationStatus = input.VerificationStatus
	existing.Status = input.Status
	existing.IsPrimary = input.IsPrimary
	existing.VerifiedAt = input.VerifiedAt
	existing.DeletedAt = gorm.DeletedAt{}
	existing.UpdatedAt = now
	if err := r.db.Unscoped().Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// FindDomainByHost 按域名获取未删除域名记录。
func (r *GormResellerRepository) FindDomainByHost(host string) (*models.ResellerDomain, error) {
	domain := normalizeDomainForRepository(host)
	if domain == "" {
		return nil, nil
	}
	var row models.ResellerDomain
	err := r.db.Preload("Profile").Where("domain = ?", domain).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// FindActiveVerifiedDomain 按域名获取已验证且启用的分销域名。
func (r *GormResellerRepository) FindActiveVerifiedDomain(host string) (*models.ResellerDomain, error) {
	domain := normalizeDomainForRepository(host)
	if domain == "" {
		return nil, nil
	}
	var row models.ResellerDomain
	err := r.db.Preload("Profile").
		Where("domain = ? AND status = ? AND verification_status = ?", domain, models.ResellerDomainStatusActive, models.ResellerDomainVerificationVerified).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// UpsertSiteConfig 创建或恢复分销站点配置。
func (r *GormResellerRepository) UpsertSiteConfig(input models.ResellerSiteConfig) (*models.ResellerSiteConfig, error) {
	if input.ResellerID == 0 {
		return nil, errors.New("invalid reseller site config")
	}
	now := time.Now()
	var existing models.ResellerSiteConfig
	err := r.db.Unscoped().Where("reseller_id = ?", input.ResellerID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		input.CreatedAt = now
		input.UpdatedAt = now
		if err := r.db.Create(&input).Error; err != nil {
			return nil, err
		}
		return &input, nil
	}
	existing.SiteName = input.SiteName
	existing.Logo = input.Logo
	existing.Favicon = input.Favicon
	existing.AnnouncementJSON = input.AnnouncementJSON
	existing.SupportJSON = input.SupportJSON
	existing.SEOJSON = input.SEOJSON
	existing.FooterLinksJSON = input.FooterLinksJSON
	existing.NavConfigJSON = input.NavConfigJSON
	existing.ThemeJSON = input.ThemeJSON
	existing.DeletedAt = gorm.DeletedAt{}
	existing.UpdatedAt = now
	if err := r.db.Unscoped().Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func normalizeDomainForRepository(raw string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(raw)), ".")
}
