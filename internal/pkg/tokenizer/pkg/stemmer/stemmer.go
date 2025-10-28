package stemmer

import (
	"fmt"
	"strings"

	"github.com/keenywheels/backend/internal/pkg/tokenizer/pkg/textutil"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/stages"

	"github.com/kljensen/snowball"
)

var (
	DefaultStemmer = New(textutil.DefaultLanguage)

	_ = stages.Stemmer(&Stemmer{})
)

// Stemmer provides stemming functionality for tokens.
type Stemmer struct {
	defaultLanguage textutil.Language
}

// New creates a new Stemmer with the specified default language.
func New(defaultLanguage textutil.Language) *Stemmer {
	return &Stemmer{
		defaultLanguage: getLanguage(defaultLanguage),
	}
}

// Stem executes the stemming process on the provided token.
func (s *Stemmer) Stem(token string) (string, error) {
	token = strings.ToLower(token)
	language := textutil.DetectLanguage(token)

	stemmedToken, err := snowball.Stem(token, string(language), true)
	if err != nil {
		return "", fmt.Errorf("failed to stem token: %w", err)
	}

	return stemmedToken, nil
}

// getLanguage returns the provided language or the default language if none is provided.
func getLanguage(language textutil.Language) textutil.Language {
	if language == "" {
		return textutil.DefaultLanguage
	}
	return language
}
