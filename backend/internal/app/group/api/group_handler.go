package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/group/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type GroupHandler struct {
	svc *service.GroupService
}

func NewGroupHandler(svc *service.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) Create(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var req service.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.svc.Create(c.Request.Context(), tenantID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, group)
}

func (h *GroupHandler) List(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	groupType := c.Query("type")
	groups, err := h.svc.List(c.Request.Context(), tenantID, groupType)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

func (h *GroupHandler) Update(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req service.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(c.Request.Context(), tenantID, id, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}

func (h *GroupHandler) Delete(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}

func (h *GroupHandler) AddMembers(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req service.MembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.AddMembers(c.Request.Context(), tenantID, id, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}

func (h *GroupHandler) RemoveMembers(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req service.MembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.RemoveMembers(c.Request.Context(), tenantID, id, &req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}
