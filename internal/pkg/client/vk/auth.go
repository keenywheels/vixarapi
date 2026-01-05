package vk

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidStateInResponse = fmt.Errorf("invalid state in response")
)

// ExchangeCodeToTokensParams parameters for ExchangeCodeToTokens
type ExchangeCodeToTokensParams struct {
	Code         string
	State        string
	CodeVerifier string
	DeviceID     string
	RedirectURI  string
}

// ExchangeCodeToTokens exchange authorization code to tokens
func (c *Client) ExchangeCodeToTokens(
	ctx context.Context,
	params *ExchangeCodeToTokensParams,
) (*ExchangeCodeToTokensResponse, error) {
	// https://id.vk.com/about/business/go/docs/ru/vkid/latest/vk-id/connection/api-description#Poluchenie-cherez-kod-podtverzhdeniya
	vals := url.Values{
		"grant_type":    []string{"authorization_code"},
		"code_verifier": []string{params.CodeVerifier},
		"redirect_uri":  []string{params.RedirectURI},
		"code":          []string{params.Code},
		"client_id":     []string{c.cfg.Auth.ClientID},
		"device_id":     []string{params.DeviceID},
		"state":         []string{params.State},
	}

	respRaw, err := c.makeRequest(
		ctx,
		http.MethodPost,
		c.endpoints.getTokenByGrantType,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		[]byte(vals.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", c.endpoints.getTokenByGrantType, err)
	}

	var resp ExchangeCodeToTokensResponse
	if err := json.Unmarshal(respRaw, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// check for error
	if resp.Error != "" {
		return nil, fmt.Errorf("error from vk %s: description=%s", resp.Error, resp.ErrorDescription)
	}

	// validate state
	if resp.State != params.State {
		return nil, fmt.Errorf("unexpected response for state: %w", ErrInvalidStateInResponse)
	}

	return &resp, nil
}

// RefreshTokensParams parameters for RefreshTokens
type RefreshTokensParams struct {
	RefreshToken string
	DeviceID     string
}

// RefreshTokens refresh tokens using refresh token
func (c *Client) RefreshTokens(
	ctx context.Context,
	params *RefreshTokensParams,
) (*RefreshTokensResponse, error) {
	state := generateState([]string{params.RefreshToken, params.DeviceID}, c.cfg.Secret)

	// https://id.vk.com/about/business/go/docs/ru/vkid/latest/vk-id/connection/api-description#Poluchenie-cherez-Refresh-token
	vals := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{params.RefreshToken},
		"client_id":     []string{c.cfg.Auth.ClientID},
		"device_id":     []string{params.DeviceID},
		"state":         []string{state},
	}

	respRaw, err := c.makeRequest(
		ctx,
		http.MethodPost,
		c.endpoints.getTokenByGrantType,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		[]byte(vals.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", c.endpoints.getTokenByGrantType, err)
	}

	var resp RefreshTokensResponse
	if err := json.Unmarshal(respRaw, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// check for error
	if resp.Error != "" {
		return nil, fmt.Errorf("error from vk %s: description=%s", resp.Error, resp.ErrorDescription)
	}

	// validate state
	if resp.State != state {
		return nil, fmt.Errorf("unexpected response for state: %w", ErrInvalidStateInResponse)
	}

	return &resp, nil
}

// generateState generate random state string
func generateState(vals []string, secret string) string {
	vals = append(vals, fmt.Sprintf("%d", time.Now().Unix()))
	vals = append(vals, secret)

	h := sha256.New()
	h.Write([]byte(strings.Join(vals, "|")))

	return string(h.Sum(nil))
}
