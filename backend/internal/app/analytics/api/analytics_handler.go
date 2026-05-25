package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/analytics/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type AnalyticsHandler struct {
	svc       *service.AnalyticsService
	exportSvc *service.ExportService
}

func NewAnalyticsHandler(svc *service.AnalyticsService, exportSvc *service.ExportService) *AnalyticsHandler {
	return &AnalyticsHandler{
		svc:       svc,
		exportSvc: exportSvc,
	}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.OverviewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.GetOverview(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *AnalyticsHandler) GetTrend(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.TrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.GetTrend(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *AnalyticsHandler) GetRank(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.RankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.GetRank(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *AnalyticsHandler) GetCompare(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.CompareRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.GetCompare(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

func (h *AnalyticsHandler) GetDetail(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.DetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items, total, err := h.svc.GetDetail(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
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

	response.SuccessWithPage(c, items, total, page, pageSize)
}

func (h *AnalyticsHandler) CreateExport(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)
	scopeIDs := getScopeAccountIDs(c)

	var req service.CreateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task, err := h.exportSvc.CreateExport(c.Request.Context(), tenantID, userID, scopeIDs, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, task)
}

func (h *AnalyticsHandler) ListExports(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)

	tasks, err := h.exportSvc.ListExports(c.Request.Context(), tenantID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tasks)
}

func (h *AnalyticsHandler) DownloadExport(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid export id")
		return
	}

	task, err := h.exportSvc.GetExport(c.Request.Context(), tenantID, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	if task.Status != "done" {
		response.BadRequest(c, "export task not completed")
		return
	}

	if !strings.HasPrefix(task.FilePath, "/tmp/exports/") {
		response.BadRequest(c, "invalid file path")
		return
	}

	c.File(task.FilePath)
}

func getScopeAccountIDs(c *gin.Context) []uint64 {
	v, exists := c.Get("scope_account_ids")
	if !exists {
		return nil
	}
	ids, _ := v.([]uint64)
	return ids
}
