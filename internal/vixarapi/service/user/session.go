package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/redis"
)

// UserSessionInfo represents user information stored in session
type UserSessionInfo struct {
	ID       string
	Username string
	Email    string
}

// ValidateSession checks if the session is valid
func (s *Service) ValidateSession(ctx context.Context, session string) (bool, *UserSessionInfo, error) {
	userInfo, err := s.redis.GetUserSession(ctx, session)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return false, nil, nil
		}

		return false, nil, fmt.Errorf("failed to get user session: %w", err)
	}

	return true, &UserSessionInfo{
		ID:       userInfo.ID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}, nil
}

// saveSession saves the user session in Redis
func (s *Service) saveSession(
	ctx context.Context,
	user *models.User,
	session *string,
) (string, error) {
	var (
		tguser *string
		vkid   *int64
	)

	// handle optional fields
	if user.TgUser.Valid {
		tguser = &user.TgUser.String
	}

	if user.VKID.Valid {
		vkid = &user.VKID.Int64
	}

	// create and save session
	sessionID := createSession([]byte(s.cfg.SessionSecret))
	if session != nil {
		// if session is provided, use it instead of creating a new one
		sessionID = *session
	}

	err := s.redis.SaveUserSession(ctx, sessionID, &redis.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		TgUser:   tguser,
		VKID:     vkid,
	})
	if err != nil {
		return "", fmt.Errorf("failed to save user session: %w", err)
	}

	return sessionID, nil
}

// saveVkTokens saves VK OAuth tokens in Redis
func (s *Service) saveVkTokens(ctx context.Context, tokens *vkTokens) error {
	// save vk tokens by vkid for long term access
	if err := s.redis.SaveVkTokens(ctx, fmt.Sprintf("%d", tokens.VKID), &redis.VkTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		DeviceID:     tokens.DeviceID,
		ExpiresIn:    tokens.ExpiresIn,
	}); err != nil {
		return fmt.Errorf("failed to save vk tokens: %w", err)
	}

	return nil
}

// createSession creates a new session identifier
func createSession(secret []byte) string {
	vals := url.Values{
		"random_uuid": []string{uuid.New().String()},
		"ts":          []string{fmt.Sprintf("%d", time.Now().Unix())},
	}

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(vals.Encode()))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
