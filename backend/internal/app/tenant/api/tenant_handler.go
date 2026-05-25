package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

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

const oauthStateSecret = "oauth-state-secret"

func signState(tenantID uint64, secret string) string {
	msg := fmt.Sprintf("%d", tenantID)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))[:16]
	return fmt.Sprintf("%d:%s", tenantID, sig)
}

func verifyState(state, secret string) (uint64, bool) {
	parts := strings.SplitN(state, ":", 2)
	if len(parts) != 2 {
		return 0, false
	}
	tenantID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, false
	}
	expected := signState(tenantID, secret)
	return tenantID, state == expected
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
	scheme := "https"
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	redirectURI := fmt.Sprintf("%s://%s/api/v1/tenants/oauth/callback", scheme, c.Request.Host)
	state := signState(tenant.ID, oauthStateSecret)
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

	tenantID, valid := verifyState(state, oauthStateSecret)
	if !valid {
		response.BadRequest(c, "invalid state signature")
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
