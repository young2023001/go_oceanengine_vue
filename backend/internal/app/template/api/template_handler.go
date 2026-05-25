package api

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"oceanengine-backend/internal/app/template/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	svc *service.TemplateService
}

// NewTemplateHandler 创建模板处理器实例
func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

// CreateProject 创建项目模板
func (h *TemplateHandler) CreateProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)

	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.CreateProject(c.Request.Context(), tenantID, userID, &req); err != nil {
		response.InternalError(c, "创建项目模板失败")
		return
	}

	response.OK(c)
}

// ListProjects 获取项目模板列表
func (h *TemplateHandler) ListProjects(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)

	list, err := h.svc.ListProjects(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, "获取项目模板列表失败")
		return
	}

	response.OKWithData(c, list)
}

// GetProject 获取单个项目模板
func (h *TemplateHandler) GetProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	t, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "项目模板不存在")
		return
	}

	response.OKWithData(c, t)
}

// UpdateProject 更新项目模板
func (h *TemplateHandler) UpdateProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	existing, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "项目模板不存在")
		return
	}

	existing.Name = req.Name
	if req.Config != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			response.InternalError(c, "配置序列化失败")
			return
		}
		existing.Config = string(configJSON)
	}

	if err := h.svc.UpdateProject(c.Request.Context(), existing); err != nil {
		response.InternalError(c, "更新项目模板失败")
		return
	}

	response.OK(c)
}

// DeleteProject 删除项目模板
func (h *TemplateHandler) DeleteProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.DeleteProject(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, "删除项目模板失败")
		return
	}

	response.OK(c)
}

// CreatePromotion 创建推广模板
func (h *TemplateHandler) CreatePromotion(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)

	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.CreatePromotion(c.Request.Context(), tenantID, userID, &req); err != nil {
		response.InternalError(c, "创建推广模板失败")
		return
	}

	response.OK(c)
}

// ListPromotions 获取推广模板列表
func (h *TemplateHandler) ListPromotions(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)

	list, err := h.svc.ListPromotions(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, "获取推广模板列表失败")
		return
	}

	response.OKWithData(c, list)
}

// DeletePromotion 删除推广模板
func (h *TemplateHandler) DeletePromotion(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.DeletePromotion(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, "删除推广模板失败")
		return
	}

	response.OK(c)
}
