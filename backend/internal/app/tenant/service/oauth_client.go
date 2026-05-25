package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OAuthClient interface {
	GetAuthURL(appID, redirectURI, state string) string
	ExchangeToken(ctx context.Context, appID, secret, authCode string) (accessToken, refreshToken string, expiresIn int64, err error)
	RefreshToken(ctx context.Context, appID, secret, refreshToken string) (newAccessToken, newRefreshToken string, expiresIn int64, err error)
}

type oceanOAuthClient struct {
	httpClient *http.Client
}

func NewOAuthClient() OAuthClient {
	return &oceanOAuthClient{httpClient: &http.Client{Timeout: 30 * time.Second}}
}

func (c *oceanOAuthClient) GetAuthURL(appID, redirectURI, state string) string {
	return fmt.Sprintf(
		"https://open.oceanengine.com/audit/oauth.html?app_id=%s&redirect_uri=%s&state=%s&scope=",
		appID, url.QueryEscape(redirectURI), state,
	)
}

type oauthTokenResponse struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	} `json:"data"`
	Message string `json:"message"`
}

func (c *oceanOAuthClient) ExchangeToken(ctx context.Context, appID, secret, authCode string) (string, string, int64, error) {
	reqURL := "https://open.oceanengine.com/open_api/oauth2/access_token/"
	payload := map[string]string{
		"app_id":     appID,
		"secret":     secret,
		"grant_type": "auth_code",
		"auth_code":  authCode,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var result oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, err
	}
	if result.Data.AccessToken == "" {
		return "", "", 0, fmt.Errorf("oauth exchange failed: %s", result.Message)
	}
	return result.Data.AccessToken, result.Data.RefreshToken, result.Data.ExpiresIn, nil
}

func (c *oceanOAuthClient) RefreshToken(ctx context.Context, appID, secret, refreshToken string) (string, string, int64, error) {
	reqURL := "https://open.oceanengine.com/open_api/oauth2/refresh_token/"
	payload := map[string]string{
		"app_id":        appID,
		"secret":        secret,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var result oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, err
	}
	if result.Data.AccessToken == "" {
		return "", "", 0, fmt.Errorf("oauth refresh failed: %s", result.Message)
	}
	return result.Data.AccessToken, result.Data.RefreshToken, result.Data.ExpiresIn, nil
}
