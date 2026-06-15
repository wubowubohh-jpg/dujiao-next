package repository

import (
	"os"
	"testing"
	"time"

	"github.com/dujiao-next/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestResellerRepositoryPostgresActiveUniqueIndex(t *testing.T) {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("skip postgres integration test: TEST_POSTGRES_DSN is empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open postgres failed: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.ResellerProfile{}, &models.ResellerDomain{}, &models.ResellerSiteConfig{}); err != nil {
		t.Fatalf("migrate postgres failed: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_reseller_domains_active_domain ON reseller_domains(domain) WHERE deleted_at IS NULL").Error; err != nil {
		t.Fatalf("create postgres active domain index failed: %v", err)
	}
	suffix := time.Now().Format("20060102150405")
	profile := seedResellerProfile(t, db, "pg-reseller-"+suffix+"@example.com")
	repo := NewResellerRepository(db)
	first, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             "pg-" + suffix + ".example.test",
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusActive,
		VerificationStatus: models.ResellerDomainVerificationVerified,
	})
	if err != nil {
		t.Fatalf("create postgres domain failed: %v", err)
	}
	if err := db.Delete(first).Error; err != nil {
		t.Fatalf("soft delete postgres domain failed: %v", err)
	}
	second, err := repo.UpsertDomain(models.ResellerDomain{
		ResellerID:         profile.ID,
		Domain:             first.Domain,
		Type:               models.ResellerDomainTypeCustom,
		Status:             models.ResellerDomainStatusActive,
		VerificationStatus: models.ResellerDomainVerificationVerified,
	})
	if err != nil {
		t.Fatalf("restore postgres domain failed: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected restore same row id=%d got id=%d", first.ID, second.ID)
	}
}
