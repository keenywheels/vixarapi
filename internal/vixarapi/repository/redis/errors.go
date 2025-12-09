package redis

import "errors"

var (
	ErrNotFound    = errors.New("data not found")
	ErrNilUserInfo = errors.New("user info is nil")
	ErrNilTokens   = errors.New("vk tokens is nil")
)
