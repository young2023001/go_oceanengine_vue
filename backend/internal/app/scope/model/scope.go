package model

type UserAccountScope struct {
	UserID    uint64 `gorm:"primaryKey" json:"user_id"`
	AccountID uint64 `gorm:"primaryKey;index" json:"account_id"`
}

func (UserAccountScope) TableName() string {
	return "user_account_scope"
}
