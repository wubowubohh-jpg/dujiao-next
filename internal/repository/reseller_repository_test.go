package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/dujiao-next/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func openResellerRepoTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:reseller_repo_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.ResellerProfile{}, &models.ResellerDomain{}, &models.ResellerSiteConfig{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_reseller_domains_active_domain ON reseller_domains(domain) WHERE deleted_at IS NULL").Error; err != nil {
		t.Fatalf("create domain index failed: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_reseller_site_configs_active_reseller ON reseller_site_configs(reseller_id) WHERE deleted_at IS NULL").Error; err != nil {
		t.Fatalf("create site config index failed: %v", err)
	}
	return db
}

func seedResellerProfile(t *testing.T, db *gorm.DB, email string) models.ResellerProfile {
	t.Helper()
	user := models.User{Email: email, PasswordHash: "hash"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}
	profile := models.ResellerProfile{UserID: user.ID, Status: models.ResellerProfileStatusActive, SettlementStatus: models.ResellerSettlementStatusNormal}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("create profile failed: %v", err)
	}
	return profile
}

func TestResellerRepositoryUpsertDomainRestoresSoftDeleted(t *testing.T) {
	db := openResellerRepoTestDB(t)
	profile := seedResellerProfile(t, db, "owner@example.com")
	repo := NewResellerRepository(db)
	first, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             "shop.example.test",
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusActive,
		VerificationStatus: models.ResellerDomainVerificationVerified,
	})
	if err != nil {
		t.Fatalf("create domain failed: %v", err)
	}
	if err := db.Delete(first).Error; err != nil {
		t.Fatalf("soft delete domain failed: %v", err)
	}
	second, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             "shop.example.test",
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusDisabled,
		VerificationStatus: models.ResellerDomainVerificationPending,
	})
	if err != nil {
		t.Fatalf("restore domain failed: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected restore same row id=%d got id=%d", first.ID, second.ID)
	}
	if second.DeletedAt.Valid {
		t.Fatal("expected restored domain deleted_at cleared")
	}
	if second.Status != models.ResellerDomainStatusDisabled {
		t.Fatalf("expected restored status disabled, got %s", second.Status)
	}
}

func TestResellerRepositoryFindActiveVerifiedDomain(t *testing.T) {
	db := openResellerRepoTestDB(t)
	profile := seedResellerProfile(t, db, "owner2@example.com")
	repo := NewResellerRepository(db)
	if _, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             "inactive.example.test",
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusDisabled,
		VerificationStatus: models.ResellerDomainVerificationVerified,
	}); err != nil {
		t.Fatalf("create disabled domain failed: %v", err)
	}
	if _, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             "active.example.test",
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusActive,
		VerificationStatus: models.ResellerDomainVerificationVerified,
		IsPrimary:          true,
	}); err != nil {
		t.Fatalf("create active domain failed: %v", err)
	}
	disabled, err := repo.FindActiveVerifiedDomain("inactive.example.test")
	if err != nil {
		t.Fatalf("lookup disabled failed: %v", err)
	}
	if disabled != nil {
		t.Fatalf("disabled domain should not resolve: %+v", disabled)
	}
	active, err := repo.FindActiveVerifiedDomain("active.example.test")
	if err != nil {
		t.Fatalf("lookup active failed: %v", err)
	}
	if active == nil || active.Profile == nil {
		t.Fatalf("expected active domain with profile, got %+v", active)
	}
	if active.Profile.UserID != profile.UserID {
		t.Fatalf("profile user mismatch want %d got %d", profile.UserID, active.Profile.UserID)
	}
}
