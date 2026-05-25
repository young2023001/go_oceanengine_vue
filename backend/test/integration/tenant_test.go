package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantCreate(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	token, err := ts.GenerateTestToken(1, "admin")
	require.NoError(t, err)

	body := map[string]interface{}{
		"name":         "测试租户",
		"oauth_app_id": "app_123456",
		"oauth_secret": "secret_abcdef",
	}

	w := ts.MakeRequest(http.MethodPost, "/api/v1/tenants", body, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestTenantList(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	token, err := ts.GenerateTestToken(1, "admin")
	require.NoError(t, err)

	// 先创建一个租户
	body := map[string]interface{}{
		"name":         "列表测试租户",
		"oauth_app_id": "app_list_001",
		"oauth_secret": "secret_list_001",
	}
	w := ts.MakeRequest(http.MethodPost, "/api/v1/tenants", body, token)
	require.Equal(t, http.StatusOK, w.Code)

	// 获取列表
	w = ts.MakeRequest(http.MethodGet, "/api/v1/tenants", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
}
