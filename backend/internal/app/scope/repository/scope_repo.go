package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"oceanengine-backend/internal/app/scope/model"
)

type ScopeRepository struct {
	db *gorm.DB
}

func NewScopeRepository(db *gorm.DB) *ScopeRepository {
	return &ScopeRepository{db: db}
}

func (r *ScopeRepository) SetUserScope(ctx context.Context, userID uint64, accountIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserAccountScope{}).Error; err != nil {
			return err
		}
		if len(accountIDs) == 0 {
			return nil
		}
		scopes := make([]model.UserAccountScope, 0, len(accountIDs))
		for _, aid := range accountIDs {
			scopes = append(scopes, model.UserAccountScope{UserID: userID, AccountID: aid})
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(scopes, 100).Error
	})
}

func (r *ScopeRepository) GetUserScope(ctx context.Context, userID uint64) ([]uint64, error) {
	var accountIDs []uint64
	err := r.db.WithContext(ctx).Model(&model.UserAccountScope{}).Where("user_id = ?", userID).Pluck("account_id", &accountIDs).Error
	return accountIDs, err
}
