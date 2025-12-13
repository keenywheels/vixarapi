package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

// UserInfo represents user session information
type UserInfo struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	TgUser   *string `json:"tguser"`
	VKID     *int64  `json:"vkid"`
}

// SaveUserSession saves a user session in Redis
func (r *Repository) SaveUserSession(ctx context.Context, key string, userInfo *UserInfo) error {
	if userInfo == nil {
		return fmt.Errorf("failed to save session: %w", ErrNilUserInfo)
	}

	data, err := json.Marshal(*userInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %w", err)
	}

	if err := r.redis.Set(ctx, key, data, r.ttl); err != nil {
		return fmt.Errorf("failed to save session in redis: %w", err)
	}

	return nil
}

// GetUserSession retrieves a user session from Redis
func (r *Repository) GetUserSession(ctx context.Context, key string) (*UserInfo, error) {
	data, ok, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("failed to get session: %w", ErrNotFound)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// DeleteUserSession deletes a user session from Redis
func (r *Repository) DeleteUserSession(ctx context.Context, userID string) error {
	if err := r.redis.Del(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}
	return nil
}

// VkTokens represents VK OAuth tokens
type VkTokens struct {
	VKID         int64  `json:"vkid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
	ExpiresIn    int64  `json:"expires_in"`
}

// SaveVkTokens saves VK tokens in Redis
func (r *Repository) SaveVkTokens(ctx context.Context, key string, tokens *VkTokens) error {
	if tokens == nil {
		return fmt.Errorf("failed to save tokens: %w", ErrNilTokens)
	}

	data, err := json.Marshal(*tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal vk tokens: %w", err)
	}

	if err := r.redis.Set(ctx, key, data, r.ttl); err != nil {
		return fmt.Errorf("failed to save vk tokens in redis: %w", err)
	}

	return nil
}

// GetVkTokens retrieves VK tokens from Redis
func (r *Repository) GetVkTokens(ctx context.Context, key string) (*VkTokens, error) {
	data, ok, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get vk tokens from redis: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("failed to get vk tokens: %w", ErrNotFound)
	}

	var tokens VkTokens
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vk tokens: %w", err)
	}

	return &tokens, nil
}

// DeleteVkTokens deletes VK tokens from Redis
func (r *Repository) DeleteVkTokens(ctx context.Context, key string) error {
	if err := r.redis.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to delete vk tokens from redis: %w", err)
	}
	return nil
}
