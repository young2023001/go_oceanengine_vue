package service

import (
	"context"

	"oceanengine-backend/internal/app/scope/repository"
)

type ScopeService struct {
	repo *repository.ScopeRepository
}

func NewScopeService(repo *repository.ScopeRepository) *ScopeService {
	return &ScopeService{repo: repo}
}

type SetScopeRequest struct {
	AccountIDs []uint64 `json:"account_ids" binding:"required"`
}

func (s *ScopeService) SetUserScope(ctx context.Context, userID uint64, req SetScopeRequest) error {
	return s.repo.SetUserScope(ctx, userID, req.AccountIDs)
}

func (s *ScopeService) GetUserScope(ctx context.Context, userID uint64) ([]uint64, error) {
	return s.repo.GetUserScope(ctx, userID)
}
