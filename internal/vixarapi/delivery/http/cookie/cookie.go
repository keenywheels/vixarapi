package cookie

import (
	"net/http"
	"time"
)

// CookieManager manages cookie-related operations
type CookieManager struct {
	cfg *Config
}

// New creates a new CookieManager instance
func New(cfg *Config) *CookieManager {
	fixConfig(cfg)

	return &CookieManager{
		cfg: cfg,
	}
}

// SessionCookie creates a session cookie with the given key and value
func (m *CookieManager) SessionCookie(session string) *http.Cookie {
	return &http.Cookie{
		Name:     m.cfg.Session.Name,
		Value:    session,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().Add(m.cfg.Session.Expiration),
	}
}
