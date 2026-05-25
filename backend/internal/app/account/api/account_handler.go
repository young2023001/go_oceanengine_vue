package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/account/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Import(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var req service.ImportAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	count, err := h.svc.Import(c.Request.Context(), tenantID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, gin.H{"imported": count})
}

func (h *AccountHandler) List(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	accounts, err := h.svc.List(c.Request.Context(), tenantID, scopeIDs)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, accounts)
}

func (h *AccountHandler) GetByID(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	account, err := h.svc.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, account)
}
