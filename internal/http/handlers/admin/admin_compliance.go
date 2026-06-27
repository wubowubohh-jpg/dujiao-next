package admin

import (
	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"

	"github.com/gin-gonic/gin"
)

// GetComplianceStatus GET /admin/compliance/status
func (h *Handler) GetComplianceStatus(c *gin.Context) {
	status, err := h.ComplianceService.Status()
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.internal", err)
		return
	}
	response.Success(c, status)
}

// ComplianceAcknowledgeRequest 请求体
type ComplianceAcknowledgeRequest struct {
	Segment1 string `json:"segment1" binding:"required"`
	Segment2 string `json:"segment2" binding:"required"`
	Segment3 string `json:"segment3" binding:"required"`
}

// AcknowledgeCompliance POST /admin/compliance/acknowledge
func (h *Handler) AcknowledgeCompliance(c *gin.Context) {
	response.Success(c, gin.H{"already_acknowledged": true})
}
