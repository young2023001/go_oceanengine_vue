package ocean

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

type CreateProjectRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	Name           string `json:"name"`
	MarketingGoal  string `json:"marketing_goal"`
	BidType        string `json:"bid_type"`
	Bid            int    `json:"bid"`
	BudgetMode     string `json:"budget_mode"`
	Budget         int    `json:"budget"`
	DeliveryScene  string `json:"delivery_scene"`
	ExternalAction string `json:"external_action"`
}

type CreateProjectResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ProjectID uint64 `json:"project_id"`
	} `json:"data"`
}

func (c *Client) CreateProject(ctx context.Context, accessToken string, req CreateProjectRequest) (*CreateProjectResponse, error) {
	return c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/create/", req)
}

type UpdateBudgetRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	Budget         int    `json:"budget"`
}

func (c *Client) UpdateProjectBudget(ctx context.Context, accessToken string, req UpdateBudgetRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/", req)
	return err
}

type UpdateBidRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	Bid            int    `json:"bid"`
}

func (c *Client) UpdateProjectBid(ctx context.Context, accessToken string, req UpdateBidRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/", req)
	return err
}

type UpdateStatusRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	OptStatus      string `json:"opt_status"`
}

func (c *Client) UpdateProjectStatus(ctx context.Context, accessToken string, req UpdateStatusRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/status/", req)
	return err
}

type ListProjectsRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
}

type ProjectInfo struct {
	ProjectID    uint64 `json:"project_id"`
	Name         string `json:"name"`
	StatusFirst  string `json:"first_status"`
	StatusSecond string `json:"second_status"`
}

type ListProjectsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List     []ProjectInfo `json:"list"`
		PageInfo struct {
			Total int `json:"total"`
		} `json:"page_info"`
	} `json:"data"`
}

func (c *Client) ListProjects(ctx context.Context, accessToken string, req ListProjectsRequest) (*ListProjectsResponse, error) {
	url := fmt.Sprintf("https://open.oceanengine.com/open_api/v1.0/local/project/list/?advertiser_id=%d&page=%d&page_size=%d",
		req.LocalAccountID, req.Page, req.PageSize)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Access-Token", accessToken)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ListProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) doPost(ctx context.Context, accessToken, url string, body interface{}) (*CreateProjectResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result CreateProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("ocean API error: %s", result.Message)
	}
	return &result, nil
}
