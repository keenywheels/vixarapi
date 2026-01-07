package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres/user"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

const (
	methodDenormalized   = "denormalized"
	methodGlobalMedian   = "global_median"
	methodCategoryMedian = "category_median"
)

// SubscribeToTokenParams represents parameters for subscribing to token updates
type SubscribeToTokenParams struct {
	UserID    string
	Token     string
	Category  string
	Method    string
	Threshold float64
}

// SubscribeToToken subscribe user to token updates
func (s *Service) SubscribeToToken(ctx context.Context, params *SubscribeToTokenParams) (string, error) {
	op := "Service.SubscribeToToken"

	token, err := s.srch.GetLatestToken(ctx, &search.GetTokenParams{
		Token:    params.Token,
		Category: params.Category,
	})
	if err != nil {
		return "", service.ParseRepositoryError(op, err)
	}

	// TODO: подумать, как все будет готово: мб стоит вынести эту логику в sql квери?
	interest, err := parseInterest(token, params.Method)
	if err != nil {
		return "", fmt.Errorf("[%s] got unexpected method: %w", op, err)
	}

	subID, err := s.repo.AddTokenSub(ctx, &user.AddTokenSubParams{
		UserID:    params.UserID,
		Token:     params.Token,
		Category:  params.Category,
		Interest:  interest,
		Threshold: params.Threshold,
		Method:    params.Method,
		ScanDate:  token.ScrapeDate,
	})
	if err != nil {
		return "", service.ParseRepositoryError(op, err)
	}

	return subID, nil
}

// TokenSubInfo represents token subscription info
type TokenSubInfo struct {
	ID               string
	Token            string
	Category         string
	Method           string
	CurrentInterest  int64
	PreviousInterest int64
	ScanDate         time.Time
}

// GetSubscribedTokens get all subscribed tokens for user
func (s *Service) GetSubscribedTokens(ctx context.Context, userID string, limit, offset uint64) ([]*TokenSubInfo, error) {
	op := "Service.GetSubscribedTokens"

	subs, err := s.repo.GetTokenSubs(ctx, userID, limit, offset)
	if err != nil {
		return nil, service.ParseRepositoryError(op, err)
	}

	return convertTokenSubs(subs), nil
}

// UnsubscribeFromToken unsubscribe user from token updates
func (s *Service) UnsubscribeFromToken(ctx context.Context, id string) error {
	op := "Service.UnsubscribeFromToken"

	if err := s.repo.DeleteTokenSub(ctx, id); err != nil {
		return service.ParseRepositoryError(op, err)
	}

	return nil
}

// parseInterest return interest based on the chosen method
func parseInterest(token *models.Token, method string) (int64, error) {
	var interest int64

	switch method {
	case methodDenormalized:
		interest = token.Interest
	case methodGlobalMedian:
		interest = token.Interest / token.GlobalMedian
	case methodCategoryMedian:
		interest = token.Interest / token.CategoryMedian
	default:
		return 0, errors.New("got unexpected method")
	}

	return interest, nil
}

// convertTokenSubs converts token subs from the repository layer to the service layer format
func convertTokenSubs(subs []*models.UserTokenSub) []*TokenSubInfo {
	res := make([]*TokenSubInfo, 0, len(subs))
	for _, s := range subs {
		res = append(res, &TokenSubInfo{
			ID:               s.ID,
			Token:            s.Token,
			Category:         s.Category,
			Method:           s.Method,
			CurrentInterest:  s.CurrentInterest,
			PreviousInterest: s.PreviousInterest,
			ScanDate:         s.ScanDate,
		})
	}

	return res
}
