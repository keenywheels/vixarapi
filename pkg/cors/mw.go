package cors

import (
	"fmt"
	"net/http"
	"strings"
)

// varyHeaders headers used to set up vary header
var varyHeaders = []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}

// WithCORS mw which adds cors headers
func WithCORS(cfg *Config, next http.Handler) http.Handler {
	allowAllOrigin := false
	for _, origin := range cfg.AllowOrigins {
		if origin == "*" {
			allowAllOrigin = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if len(origin) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		host := r.Host
		if origin == "http://"+host || origin == "https://"+host {
			next.ServeHTTP(w, r)
			return
		}

		isValid := validateOrigin(origin, cfg.AllowOrigins)
		if allowAllOrigin || isValid {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
			w.Header().Set("Vary", strings.Join(varyHeaders, ", "))

			if !allowAllOrigin && cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		if r.Method == http.MethodOptions {
			if allowAllOrigin || isValid {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateOrigin validates origin
func validateOrigin(origin string, validOrigins []string) bool {
	for _, validOrigin := range validOrigins {
		if validOrigin == "*" || validOrigin == origin {
			return true
		}
	}

	return false
}
