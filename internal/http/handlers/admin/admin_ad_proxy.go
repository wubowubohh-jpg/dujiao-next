package admin

import (
	"net/http"

	"github.com/dujiao-next/internal/http/response"
	"github.com/gin-gonic/gin"
)

// GetAdRender 代理广告位渲染请求到 ad-system
func (h *Handler) GetAdRender(c *gin.Context) {
	slotCode := c.Param("slotCode")
	if slotCode == "" {
		response.Error(c, http.StatusBadRequest, "slot_code is required")
		return
	}

	params := make(map[string]string)
	for _, key := range []string{"tenant", "client", "locale"} {
		if v := c.Query(key); v != "" {
			params[key] = v
		}
	}

	data, err := h.AdProxyService.RenderSlot(c.Request.Context(), slotCode, params)
	if err != nil {
		// 广告请求失败时静默返回空数据，不影响主业务
		response.Success(c, nil)
		return
	}

	response.Success(c, data)
}

// PostAdImpression 保留管理端曝光接口兼容前端，但不再向外部广告服务上报。
func (h *Handler) PostAdImpression(c *gin.Context) {
	response.Success(c, nil)
}
