package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDataScope_AdminBypassesScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Set("tenant_id", uint64(1))
		c.Set("role_key", "admin")
		c.Next()
	})
	r.Use(DataScope(nil))
	r.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("scope_account_ids")
		assert.False(t, exists, "admin should not have scope_account_ids set")
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDataScope_NonAdminWithoutScope_Returns403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(99))
		c.Set("tenant_id", uint64(1))
		c.Set("role_key", "operator")
		c.Next()
	})
	r.Use(DataScope(nil))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
