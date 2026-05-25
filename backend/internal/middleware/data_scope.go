package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	scopeModel "oceanengine-backend/internal/app/scope/model"
	"oceanengine-backend/pkg/ginutil"
)

func DataScope(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleKey := ginutil.GetRoleKey(c)
		if roleKey == "admin" || roleKey == "super_admin" {
			c.Next()
			return
		}

		userID := ginutil.GetUserID(c)
		if userID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if db == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "no account scope assigned"})
			return
		}

		var accountIDs []uint64
		db.Model(&scopeModel.UserAccountScope{}).
			Where("user_id = ?", userID).
			Pluck("account_id", &accountIDs)

		if len(accountIDs) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "no account scope assigned"})
			return
		}

		c.Set("scope_account_ids", accountIDs)
		c.Next()
	}
}
