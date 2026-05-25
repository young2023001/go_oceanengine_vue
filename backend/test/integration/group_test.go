package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	accountModel "oceanengine-backend/internal/app/account/model"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
)

func TestGroupCreate(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	// 创建租户
	tenant := &tenantModel.Tenant{
		Name:        "分组测试租户",
		OAuthAppID:  "app_grp_001",
		OAuthSecret: "secret_grp_001",
		TokenStatus: tenantModel.TokenStatusActive,
		Status:      1,
	}
	require.NoError(t, ts.DB.Create(tenant).Error)

	token, err := ts.GenerateTestTokenWithTenant(1, "admin", tenant.ID)
	require.NoError(t, err)

	body := map[string]interface{}{
		"name": "华北区域",
		"type": "region",
		"sort": 1,
	}

	w := ts.MakeRequest(http.MethodPost, "/api/v1/groups", body, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestGroupAddMembers(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	// 创建租户
	tenant := &tenantModel.Tenant{
		Name:        "分组成员租户",
		OAuthAppID:  "app_grp_mem",
		OAuthSecret: "secret_grp_mem",
		TokenStatus: tenantModel.TokenStatusActive,
		Status:      1,
	}
	require.NoError(t, ts.DB.Create(tenant).Error)

	token, err := ts.GenerateTestTokenWithTenant(1, "admin", tenant.ID)
	require.NoError(t, err)

	// 创建分组
	groupBody := map[string]interface{}{
		"name": "加盟商A",
		"type": "franchisee",
	}
	w := ts.MakeRequest(http.MethodPost, "/api/v1/groups", groupBody, token)
	require.Equal(t, http.StatusOK, w.Code)

	var createResp struct {
		Code int `json:"code"`
		Data struct {
			ID uint64 `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, ParseResponse(w, &createResp))
	groupID := createResp.Data.ID
	require.NotZero(t, groupID)

	// 创建测试账户
	accounts := []accountModel.LocalAccount{
		{TenantID: tenant.ID, LocalAccountID: 300001, Name: "账户1", Status: "active"},
		{TenantID: tenant.ID, LocalAccountID: 300002, Name: "账户2", Status: "active"},
	}
	require.NoError(t, ts.DB.Create(&accounts).Error)

	// 添加成员
	memberBody := map[string]interface{}{
		"account_ids": []uint64{accounts[0].ID, accounts[1].ID},
	}
	w = ts.MakeRequest(http.MethodPost, fmt.Sprintf("/api/v1/groups/%d/members", groupID), memberBody, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
}

func TestGroupList(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Cleanup()
	ts.SeedTestData(t)

	// 创建租户
	tenant := &tenantModel.Tenant{
		Name:        "分组列表租户",
		OAuthAppID:  "app_grp_list",
		OAuthSecret: "secret_grp_list",
		TokenStatus: tenantModel.TokenStatusActive,
		Status:      1,
	}
	require.NoError(t, ts.DB.Create(tenant).Error)

	token, err := ts.GenerateTestTokenWithTenant(1, "admin", tenant.ID)
	require.NoError(t, err)

	// 创建分组
	groupBody := map[string]interface{}{
		"name": "自定义分组",
		"type": "custom",
	}
	w := ts.MakeRequest(http.MethodPost, "/api/v1/groups", groupBody, token)
	require.Equal(t, http.StatusOK, w.Code)

	// 获取列表
	w = ts.MakeRequest(http.MethodGet, "/api/v1/groups", nil, token)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err = ParseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
}
