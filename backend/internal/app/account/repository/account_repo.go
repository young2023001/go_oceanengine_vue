package repository

import (
	"context"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/account/model"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) BatchCreate(ctx context.Context, accounts []model.LocalAccount) error {
	return r.db.WithContext(ctx).CreateInBatches(accounts, 100).Error
}

func (r *AccountRepository) ListByTenant(ctx context.Context, tenantID uint64, accountIDs []uint64) ([]model.LocalAccount, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if len(accountIDs) > 0 {
		query = query.Where("id IN ?", accountIDs)
	}
	var accounts []model.LocalAccount
	err := query.Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalAccount, error) {
	var account model.LocalAccount
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) CreateStore(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *AccountRepository) ListStores(ctx context.Context, tenantID uint64, accountID uint64) ([]model.Store, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}
	var stores []model.Store
	err := query.Find(&stores).Error
	return stores, err
}
