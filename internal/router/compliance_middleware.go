package router

import (
	"github.com/dujiao-next/internal/service"

	"github.com/gin-gonic/gin"
)

// PaymentComplianceRequired no longer blocks payment/finance routes.
func PaymentComplianceRequired(cs *service.ComplianceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = cs
		c.Next()
	}
}
