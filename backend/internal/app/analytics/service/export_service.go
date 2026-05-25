package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"oceanengine-backend/internal/app/analytics/model"
	"oceanengine-backend/internal/app/analytics/repository"
)

const maxExportRows = 10000

type CreateExportRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type ExportService struct {
	exportRepo *repository.ExportRepository
	reportRepo *repository.ReportRepository
}

func NewExportService(exportRepo *repository.ExportRepository, reportRepo *repository.ReportRepository) *ExportService {
	return &ExportService{
		exportRepo: exportRepo,
		reportRepo: reportRepo,
	}
}

func (s *ExportService) CreateExport(ctx context.Context, tenantID, userID uint64, scopeIDs []uint64, req CreateExportRequest) (*model.ExportTask, error) {
	expireAt := time.Now().Add(24 * time.Hour)
	task := &model.ExportTask{
		TenantID:  tenantID,
		Title:     req.Title,
		Status:    "pending",
		ExpireAt:  &expireAt,
		CreatedBy: userID,
	}

	if err := s.exportRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.exportRepo.UpdateStatus(context.Background(), task.ID, "failed", "", 0)
			}
		}()
		s.executeExport(task.ID, tenantID, scopeIDs, req)
	}()

	return task, nil
}

func (s *ExportService) executeExport(taskID, tenantID uint64, scopeIDs []uint64, req CreateExportRequest) {
	ctx := context.Background()

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	records, total, err := s.reportRepo.GetDetail(ctx, tenantID, scopeIDs, startDate, endDate, 1, maxExportRows)
	if err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	if total > maxExportRows {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	filePath := fmt.Sprintf("/tmp/exports/%d_%d.csv", tenantID, taskID)

	if err := os.MkdirAll("/tmp/exports", 0o755); err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}
	defer file.Close()

	header := "stat_date,account_id,cost,show_cnt,click_cnt,convert_cnt,ctr,cpc,cpm\n"
	if _, err := file.WriteString(header); err != nil {
		_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	for _, r := range records {
		line := fmt.Sprintf("%s,%d,%.2f,%d,%d,%d,%.4f,%.2f,%.2f\n",
			r.StatDate.Format(dateLayout),
			r.AccountID,
			r.Cost,
			r.ShowCnt,
			r.ClickCnt,
			r.ConvertCnt,
			r.Ctr,
			r.Cpc,
			r.Cpm,
		)
		if _, err := file.WriteString(line); err != nil {
			_ = s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
			return
		}
	}

	_ = s.exportRepo.UpdateStatus(ctx, taskID, "done", filePath, len(records))
}

func (s *ExportService) ListExports(ctx context.Context, tenantID uint64) ([]model.ExportTask, error) {
	return s.exportRepo.List(ctx, tenantID)
}

func (s *ExportService) GetExport(ctx context.Context, tenantID, id uint64) (*model.ExportTask, error) {
	return s.exportRepo.GetByID(ctx, tenantID, id)
}
