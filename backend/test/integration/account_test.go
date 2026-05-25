package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
)

func TestAccountImport(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	// 创建租户
	tenant := &tenantModel.Tenant{
		Name:        "账户测试租户",
		OAuthAppID:  "app_acct_001",
		OAuthSecret: "secret_acct_001",
		TokenStatus: tenantModel.TokenStatusActive,
		Status:      1,
	}
	require.NoError(t, ts.DB.Create(tenant).Error)

	token, err := ts.GenerateTestTokenWithTenant(1, "admin", tenant.ID)
	require.NoError(t, err)

	body := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"local_account_id": 100001,
				"name":             "测试账户A",
				"store_name":       "门店A",
				"city":             "北京",
			},
			{
				"local_account_id": 100002,
				"name":             "测试账户B",
				"store_name":       "门店B",
				"city":             "上海",
			},
		},
	}

	w := ts.MakeRequest(http.MethodPost, "/api/v1/accounts/import", body, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestAccountList(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	// 创建租户
	tenant := &tenantModel.Tenant{
		Name:        "账户列表租户",
		OAuthAppID:  "app_acct_list",
		OAuthSecret: "secret_acct_list",
		TokenStatus: tenantModel.TokenStatusActive,
		Status:      1,
	}
	require.NoError(t, ts.DB.Create(tenant).Error)

	token, err := ts.GenerateTestTokenWithTenant(1, "admin", tenant.ID)
	require.NoError(t, err)

	// 先导入账户
	importBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"local_account_id": 200001,
				"name":             "列表测试账户",
			},
		},
	}
	w := ts.MakeRequest(http.MethodPost, "/api/v1/accounts/import", importBody, token)
	require.Equal(t, http.StatusOK, w.Code)

	// 获取列表
	w = ts.MakeRequest(http.MethodGet, "/api/v1/accounts", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
}
