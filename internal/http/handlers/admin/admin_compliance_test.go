package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dujiao-next/internal/provider"
	"github.com/dujiao-next/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupComplianceHandler() *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := &Handler{Container: &provider.Container{ComplianceService: service.NewComplianceService(nil)}}

	r := gin.New()
	r.GET("/compliance/status", h.GetComplianceStatus)
	r.POST("/compliance/acknowledge", h.AcknowledgeCompliance)
	return r
}

func TestGetComplianceStatusDisabled(t *testing.T) {
	r := setupComplianceHandler()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/compliance/status", nil))

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"acknowledged":true`)
	assert.Contains(t, w.Body.String(), `"version":"disabled"`)
}

func TestAcknowledgeComplianceDisabled(t *testing.T) {
	r := setupComplianceHandler()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/compliance/acknowledge", nil))

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"already_acknowledged":true`)
}
