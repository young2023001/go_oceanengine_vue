package service

import (
	"context"
	"fmt"

	"oceanengine-backend/internal/app/account/model"
	"oceanengine-backend/internal/app/account/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

type ImportAccountItem struct {
	LocalAccountID uint64 `json:"local_account_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	StoreName      string `json:"store_name"`
	City           string `json:"city"`
	District       string `json:"district"`
	Address        string `json:"address"`
}

type ImportAccountsRequest struct {
	Items []ImportAccountItem `json:"items" binding:"required,min=1"`
}

func (s *AccountService) Import(ctx context.Context, tenantID uint64, req ImportAccountsRequest) (int, error) {
	accounts := make([]model.LocalAccount, 0, len(req.Items))
	for _, item := range req.Items {
		accounts = append(accounts, model.LocalAccount{
			TenantID:       tenantID,
			LocalAccountID: item.LocalAccountID,
			Name:           item.Name,
			Status:         "active",
		})
	}
	if err := s.repo.BatchCreate(ctx, accounts); err != nil {
		return 0, fmt.Errorf("batch create accounts: %w", err)
	}
	return len(accounts), nil
}

func (s *AccountService) List(ctx context.Context, tenantID uint64, scopeAccountIDs []uint64) ([]model.LocalAccount, error) {
	return s.repo.ListByTenant(ctx, tenantID, scopeAccountIDs)
}

func (s *AccountService) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalAccount, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}
