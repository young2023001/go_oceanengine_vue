package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/project/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}

	var req service.ListProjectsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	projects, total, err := h.svc.ListProjects(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	response.OKWithList(c, projects, total, page, pageSize)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	project, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithData(c, project)
}

func (h *ProjectHandler) UpdateProjectStatus(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), tenantID, id, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c)
}

func (h *ProjectHandler) ListPromotions(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}

	var req service.ListPromotionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	promotions, total, err := h.svc.ListPromotions(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	response.OKWithList(c, promotions, total, page, pageSize)
}
