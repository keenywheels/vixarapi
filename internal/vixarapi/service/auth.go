package service

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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/keenywheels/backend/internal/pkg/client/vk"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/redis"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

// vkTokens is a helper struct for VK OAuth tokens
type vkTokens struct {
	VKID         int64
	AccessToken  string
	RefreshToken string
	DeviceID     string
	ExpiresIn    int64
}

// VkAuthCallbackParams represents the parameters received in VK OAuth callback
type VkAuthCallbackParams struct {
	Code         string
	State        string
	CodeVerifier string
	DeviceID     string
	RedirectURI  string
}

// VkAuthCallbackResult represents the result of VK OAuth callback processing
type VkAuthCallbackResult struct {
	UserExists bool
	Session    string
	Username   string
	Email      string
	VKID       int64
}

// HandleVkAuthCallback processes the VK OAuth callback
func (s *Service) HandleVkAuthCallback(
	ctx context.Context,
	params *VkAuthCallbackParams,
) (*VkAuthCallbackResult, error) {
	op := "Service.HandleVkAuthCallback"

	// get vk tokens
	tokens, err := s.exchangeCode(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to exchange code: %w", op, err)
	}

	// try to get user from db
	user, err := s.repo.GetUserByVKID(ctx, tokens.VKID)
	if err == nil {
		// if user exists, saves session and vk tokens
		sessionID, err := s.saveSession(ctx, user, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to save session: %w", op, err)
		}

		if err := s.saveVkTokens(ctx, tokens); err != nil {
			return nil, fmt.Errorf("%s: failed to save vk tokens: %w", op, err)
		}

		return &VkAuthCallbackResult{
			UserExists: true,
			Session:    sessionID,
			Username:   user.Username,
			Email:      user.Email,
			VKID:       user.VKID.Int64,
		}, nil
	}

	// get user info from vk
	userInfo, err := s.vk.GetUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// saves vk tokens for new user
	if err := s.saveVkTokens(ctx, tokens); err != nil {
		return nil, fmt.Errorf("%s: failed to save vk tokens: %w", op, err)
	}

	// save session for new user
	username := fmt.Sprintf("user_%s", userInfo.User.UserID)

	sessionID, err := s.saveSession(ctx, &models.User{
		Username: username,
		Email:    userInfo.User.Email,
		VKID:     pgtype.Int8{Int64: tokens.VKID, Valid: true},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to save session for new user: %w", op, err)
	}

	return &VkAuthCallbackResult{
		UserExists: false,
		Session:    sessionID,
		Username:   username,
		Email:      userInfo.User.Email,
		VKID:       tokens.VKID,
	}, nil
}

// RegisterVkUserParams represents the parameters required to register a VK user
type RegisterVkUserParams struct {
	SessionID string
	Email     string
	Username  string
	VKID      int64
}

// RegisterVkUser registers a new user using VK OAuth information
func (s *Service) RegisterVkUser(
	ctx context.Context,
	params *RegisterVkUserParams,
) error {
	op := "Service.RegisterVkUser"

	user := &models.User{
		Username: params.Username,
		Email:    params.Email,
		VKID:     pgtype.Int8{Int64: params.VKID, Valid: true},
	}

	user, err := s.repo.RegisterVKUser(ctx, user)
	if err != nil {
		return parseRepositoryError(op, err)
	}

	// update session with new user info
	_, err = s.saveSession(ctx, user, &params.SessionID)
	if err != nil {
		return fmt.Errorf("%s: failed to save session: %w", op, err)
	}

	return nil
}

// LogoutUser logs out the user by clearing the session
func (s *Service) LogoutUser(ctx context.Context, sessionID string) error {
	var (
		op  = "Service.LogoutUser"
		log = ctxutils.GetLogger(ctx)
	)

	// get user info from session
	user, err := s.redis.GetUserSession(ctx, sessionID)
	if err == nil && user.VKID != nil {
		// if retrieved user session, than get vkid and delete vk tokens
		if err := s.redis.DeleteVkTokens(ctx, fmt.Sprintf("%d", *user.VKID)); err != nil {
			log.Errorf("[%s] failed to delete vk tokens for vkid: %v", op, err)
		}
	} else {
		// if not found, log the error but continue
		log.Errorf("[%s] failed to get user session: %v", op, err)
	}

	// delete user session from redis
	if err := s.redis.DeleteUserSession(ctx, sessionID); err != nil {
		return fmt.Errorf("%s: failed to delete user session: %w", op, err)
	}

	return nil
}

// ValidateSession checks if the session is valid
func (s *Service) ValidateSession(ctx context.Context, session string) (bool, error) {
	_, err := s.redis.GetUserSession(ctx, session)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("failed to get user session: %w", err)
	}

	return true, nil
}

// exchangeCode exchanges the authorization code for VK OAuth tokens
func (s *Service) exchangeCode(ctx context.Context, params *VkAuthCallbackParams) (*vkTokens, error) {
	var tokens vkTokens

	// make request using vk client to exchange code to tokens
	tokensResp, err := s.vk.ExchangeCodeToTokens(ctx, &vk.ExchangeCodeToTokensParams{
		Code:         params.Code,
		State:        params.State,
		CodeVerifier: params.CodeVerifier,
		DeviceID:     params.DeviceID,
		RedirectURI:  params.RedirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code to tokens: %w", err)
	}

	tokens.VKID = tokensResp.UserID
	tokens.AccessToken = tokensResp.AccessToken
	tokens.RefreshToken = tokensResp.RefreshToken
	tokens.ExpiresIn = tokensResp.ExpiresIn
	tokens.DeviceID = params.DeviceID

	return &tokens, nil
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
