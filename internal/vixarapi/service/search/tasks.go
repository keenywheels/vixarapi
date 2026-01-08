package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	repo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

const typeEmail = "email"

const defaultLimit = 1000

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
	// TODO: вынести в конфиг, если надо будет менять интервал
	if err := s.r.UpdateUserTokenSubs(ctx, repo.IntervalDays, 1); err != nil {
		// return error cuz if we fail to update the token subs, no need to trigger users notification
		return fmt.Errorf("[%s] failed to update user token subs: %w", op, err)
	}

	// put notification tasks into the queue
	parsed, err := s.putNotificationTasks(ctx)
	if err != nil {
		return fmt.Errorf("[%s] failed to put notification tasks: %w", op, err)
	}

	log.Infof("[%s] successfully put %d notification tasks", op, parsed)

	return nil
}

// putNotificationTasks puts notification tasks into the queue
func (s *Service) putNotificationTasks(ctx context.Context) (uint64, error) {
	var (
		op  = "Service.putNotificationTasks"
		log = ctxutils.GetLogger(ctx)
	)

	var (
		err    error
		offset uint64
		subs   []*repo.IncreasedTokenSubInfo
	)

	for {
		subs, err = s.r.GetIncreasedTokenSubs(ctx, defaultLimit, uint64(offset))
		if err != nil {
			// got error -> break
			break
		}

		// put task for every sub
		for _, sub := range subs {
			// TODO: надо доработать логику, чтобы была защита от повторных отправок уведомлений
			// TODO: (если при обходе ничего не обновили, но старые значения удовлетворяют условию отправки)
			if err := s.broker.SendNotification(models.Notification{
				Type:             typeEmail,
				UserID:           sub.UserID,
				Username:         sub.Username,
				Email:            sub.Email,
				Token:            sub.Token,
				Category:         sub.Category,
				Threshold:        sub.Threshold,
				PreviousInterest: sub.PreviousInterest,
				CurrentInterest:  sub.CurrentInterest,
			}); err != nil {
				log.Errorf("[%s] failed to send notification for user_id=%s, token=%s: %v",
					op, sub.UserID, sub.Token, err,
				)
			}
		}

		offset += uint64(len(subs)) // update offset

		// break if got less than the limit
		if len(subs) < defaultLimit {
			break
		}
	}

	// handle unexpected error
	if err != nil && !errors.Is(err, commonRepo.ErrNotFound) {
		return offset, fmt.Errorf("[%s] failed to get increased token subs: %w", op, err)
	}

	// just log if did not put any tasks
	if offset == 0 {
		log.Infof("[%s] no token subs were increased -> no notification tasks", op)
	}

	return offset, nil
}
