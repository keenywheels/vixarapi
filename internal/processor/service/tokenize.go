package service

import (
	"context"

	"github.com/keenywheels/backend/pkg/ctxutils"
)

// TokenizeMessage processes and tokenizes the given message
func (s *Service) TokenizeMessage(ctx context.Context, message string) error {
	log := ctxutils.GetLogger(ctx)

	log.Infof("tokenizing message: %s", message)

	// TODO: добавить логику: прогон через токенайзер + запись в БД

	return nil
}
