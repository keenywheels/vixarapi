package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/keenywheels/backend/pkg/ctxutils"
)

func (c *Client) GetUserInfo(
	ctx context.Context,
	accessToken string,
) (*UserInfoResponse, error) {
	// https://id.vk.com/about/business/go/docs/ru/vkid/latest/vk-id/connection/api-description#Poluchenie-nemaskirovannyh-dannyh
	vals := url.Values{
		"client_id":    []string{c.cfg.Auth.ClientID},
		"access_token": []string{accessToken},
	}

	ctxutils.GetLogger(ctx).Debugf("[VkClient.GetUserInfo] got request to with params=%v", vals.Encode())

	respRaw, err := c.makeRequest(
		ctx,
		http.MethodPost,
		c.endpoints.getUserInfo,
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		[]byte(vals.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", c.endpoints.getUserInfo, err)
	}

	var resp UserInfoResponse
	if err := json.Unmarshal(respRaw, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// check for error
	if resp.Error != "" {
		return nil, fmt.Errorf("error from vk %s: description=%s", resp.Error, resp.ErrorDescription)
	}

	return &resp, nil
}
