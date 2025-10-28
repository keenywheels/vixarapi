package models

import (
	"strings"
	"time"
)

const (
	DefaultTokenSource   = "undefined"
	DefaultContextWindow = 5
)

// Token represents a token with its context and metadata
type Token struct {
	Target    string
	Context   []Token
	Source    string
	Metadata  map[string]any
	Timestamp time.Time
}

// GetTokens tokenizes the input text and returns a slice of Tokens with context
func GetTokens(text string, tokenConfig *TokenConfig) []Token {
	now := time.Now()

	words := strings.Fields(text)
	tokens := make([]Token, len(words))

	for i, word := range words {
		tokens[i] = Token{
			Target:    word,
			Source:    tokenConfig.TokenSource,
			Metadata:  make(map[string]any),
			Timestamp: now,
		}
	}

	collectContext(tokens, tokenConfig)

	return tokens
}

// Filter marks the token as filtered in its metadata
func (t *Token) Filter() {
	if t.Metadata == nil {
		t.Metadata = make(map[string]any)
	}
	t.Metadata["filtered"] = true
}

// collectContext populates the Context field for each token based on the context window
func collectContext(tokens []Token, tokenConfig *TokenConfig) {
	n := len(tokens)

	for i := range tokens {
		start := max(0, i-tokenConfig.ContextWindow)
		end := min(n, i+tokenConfig.ContextWindow+1)

		tokens[i].Context = make([]Token, end-start)

		for j := start; j < end; j++ {
			tokens[i].Context[j-start] = tokens[j]
		}
	}
}

// TokenConfig holds configuration for tokenization
type TokenConfig struct {
	TokenSource   string
	ContextWindow int
}

// NewTokenConfig creates a new TokenConfig with defaults if necessary
func NewTokenConfig(source string, window int) *TokenConfig {
	return &TokenConfig{
		TokenSource:   getTokenSource(source),
		ContextWindow: getContextWindow(window),
	}
}

// getTokenSource returns the token source or the default if empty
func getTokenSource(source string) string {
	if source == "" {
		return DefaultTokenSource
	}
	return source
}

// getContextWindow returns the context window or the default if non-positive
func getContextWindow(window int) int {
	if window <= 0 {
		return DefaultContextWindow
	}
	return window
}
