package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	accountModel "oceanengine-backend/internal/app/account/model"
	"oceanengine-backend/internal/app/batch/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type BatchHandler struct {
	svc       *service.BatchService
	accountDB *gorm.DB
}

func NewBatchHandler(svc *service.BatchService, accountDB *gorm.DB) *BatchHandler {
	return &BatchHandler{svc: svc, accountDB: accountDB}
}

func (h *BatchHandler) CreateTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)

	var req service.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if len(req.AccountIDs) == 0 {
		response.BadRequest(c, "账户列表不能为空")
		return
	}

	// 查询 local_account 表构建 ID 映射
	var accounts []accountModel.LocalAccount
	if err := h.accountDB.Where("tenant_id = ? AND local_account_id IN ?", tenantID, req.AccountIDs).
		Find(&accounts).Error; err != nil {
		response.InternalError(c, "查询账户失败")
		return
	}

	localAccountMap := make(map[uint64]uint64, len(accounts))
	for _, acc := range accounts {
		localAccountMap[acc.LocalAccountID] = acc.ID
	}

	task, err := h.svc.CreateTask(c.Request.Context(), tenantID, userID, localAccountMap, req)
	if err != nil {
		response.InternalError(c, "创建任务失败: "+err.Error())
		return
	}

	response.Success(c, task)
}

func (h *BatchHandler) GetTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	task, err := h.svc.GetTask(c.Request.Context(), tenantID, taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	response.Success(c, task)
}

func (h *BatchHandler) ListTasks(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tasks, total, err := h.svc.ListTasks(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		response.InternalError(c, "查询任务列表失败")
		return
	}

	response.SuccessWithPage(c, tasks, total, page, pageSize)
}

func (h *BatchHandler) CancelTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	if err := h.svc.CancelTask(c.Request.Context(), tenantID, taskID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c)
}

func (h *BatchHandler) RetryTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	if err := h.svc.RetryTask(c.Request.Context(), tenantID, taskID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c)
}
