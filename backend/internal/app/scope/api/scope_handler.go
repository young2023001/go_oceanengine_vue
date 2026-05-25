package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/scope/service"
	"oceanengine-backend/pkg/response"
)

type ScopeHandler struct {
	svc *service.ScopeService
}

func NewScopeHandler(svc *service.ScopeService) *ScopeHandler {
	return &ScopeHandler{svc: svc}
}

func (h *ScopeHandler) SetScope(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	var req service.SetScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SetUserScope(c.Request.Context(), userID, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}

func (h *ScopeHandler) GetScope(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	accountIDs, err := h.svc.GetUserScope(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, gin.H{"account_ids": accountIDs})
}
