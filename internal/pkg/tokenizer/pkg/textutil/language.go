package textutil

import "unicode"

// Language represents a language type
type Language string

// Supported languages
const (
	English Language = "english"
	Russian Language = "russian"
)

var DefaultLanguage = Russian

// DetectLanguage detects the language of the given token based on its characters
func DetectLanguage(token string) Language {
	for _, t := range token {
		switch {
		case unicode.In(t, unicode.Cyrillic):
			return Russian
		case unicode.In(t, unicode.Latin):
			return English
		}
	}

	return DefaultLanguage
}
