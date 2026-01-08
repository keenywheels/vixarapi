package search

import (
	"context"
	"fmt"

	"github.com/keenywheels/backend/pkg/ctxutils"
)

// updateSearchTask updates the search table
func (s *Service) updateSearchTask(ctx context.Context) error {
	var (
		op  = "Service.updateSearchTask"
		log = ctxutils.GetLogger(ctx)
	)

	// update search table
	log.Infof("[%s] updating search table", op)

	if err := s.r.UpdateSearchTable(ctx); err != nil {
		// return error cuz if we fail to update the search table, no need to proceed further
		return fmt.Errorf("[%s] failed to update search table: %w", op, err)
	}

	// update token subs values
	if err := s.r.UpdateUserTokenSubs(ctx); err != nil {
		// return error cuz if we fail to update the token subs, no need to trigger users notification
		return fmt.Errorf("[%s] failed to update user token subs: %w", op, err)
	}

	// notify users
	// TODO: тут надо будет написать функцию, которая будет делать запрос к подпискам
	// TODO: и будет выбирать те, у который превысилось пороговое значение,
	// TODO: для таких подписок будем класть задачу в очередь

	return nil
}
