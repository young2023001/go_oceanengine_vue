package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/tenant/service"
	"oceanengine-backend/pkg/response"
)

type TenantHandler struct {
	svc   *service.TenantService
	oauth service.OAuthClient
}

func NewTenantHandler(svc *service.TenantService, oauth service.OAuthClient) *TenantHandler {
	return &TenantHandler{svc: svc, oauth: oauth}
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req service.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tenant, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, tenant)
}

func (h *TenantHandler) List(c *gin.Context) {
	tenants, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, tenants)
}

func (h *TenantHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, tenant)
}

func (h *TenantHandler) GetOAuthURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	redirectURI := fmt.Sprintf("%s/api/v1/tenants/oauth/callback", c.Request.Host)
	state := fmt.Sprintf("%d", tenant.ID)
	authURL := h.oauth.GetAuthURL(tenant.OAuthAppID, redirectURI, state)
	response.OKWithData(c, gin.H{"auth_url": authURL})
}

func (h *TenantHandler) OAuthCallback(c *gin.Context) {
	authCode := c.Query("auth_code")
	state := c.Query("state")
	if authCode == "" || state == "" {
		response.BadRequest(c, "missing auth_code or state")
		return
	}

	tenantID, err := strconv.ParseUint(state, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid state")
		return
	}

	tenant, err := h.svc.GetByID(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	accessToken, refreshToken, expiresIn, err := h.oauth.ExchangeToken(
		c.Request.Context(), tenant.OAuthAppID, tenant.OAuthSecret, authCode,
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if err := h.svc.SaveOAuthToken(c.Request.Context(), tenantID, accessToken, refreshToken, expiresIn); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OKWithData(c, gin.H{"message": "授权成功"})
}
