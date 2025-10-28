package stages

import (
	"strings"
	"unicode"

	"github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/models"
	"golang.org/x/text/unicode/norm"
)

// NewNormalizerStage creates a new normalization stage that normalizes token strings.
func NewNormalizerStage() *tokenizer.Stage {
	stage := &tokenizer.Stage{}

	stage.CallbackFunc = func(token *models.Token) error {
		token.Target = normalizeString(token.Target)
		return nil
	}

	return stage
}

// normalizeString normalizes the input string by applying Unicode NFC normalization,
func normalizeString(str string) string {
	str = norm.NFC.String(str)

	builder := strings.Builder{}
	builder.Grow(len(str))

	for _, r := range str {
		switch {
		case unicode.IsLetter(r):
			builder.WriteRune(unicode.ToLower(r))
		}
	}

	return builder.String()
}
