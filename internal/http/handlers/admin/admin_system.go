package admin

import (
	"context"
	"time"

	"github.com/dujiao-next/internal/http/handlers/shared"
	"github.com/dujiao-next/internal/http/response"
	"github.com/dujiao-next/internal/version"

	"github.com/gin-gonic/gin"
)

// CheckSystemUpdate returns local version status without external update checks.
// GET /api/v1/admin/system/version/check
func (h *Handler) CheckSystemUpdate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 12*time.Second)
	defer cancel()

	result, err := version.CheckLatestRelease(ctx)
	if err != nil {
		shared.RespondError(c, response.CodeInternal, "error.update_check_failed", err)
		return
	}

	response.Success(c, result)
}
