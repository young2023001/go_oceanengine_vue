package ginutil

import "github.com/gin-gonic/gin"

func GetTenantID(c *gin.Context) uint64 {
	v, _ := c.Get("tenant_id")
	id, _ := v.(uint64)
	return id
}

func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get("user_id")
	id, _ := v.(uint64)
	return id
}

func GetRoleKey(c *gin.Context) string {
	v, _ := c.Get("role_key")
	s, _ := v.(string)
	return s
}
