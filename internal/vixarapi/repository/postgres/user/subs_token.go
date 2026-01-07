package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
)

// AddTokenSubParams represents parameters for adding a token subscription
type AddTokenSubParams struct {
	UserID    string
	Token     string
	Category  string
	Interest  int64
	Threshold float64
	Method    string
	ScanDate  time.Time
}

// AddTokenSub adds a new token subscription in database
func (r *Repository) AddTokenSub(ctx context.Context, params *AddTokenSubParams) (string, error) {
	op := "Repository.AddTokenSub"

	query, args, err := r.db.Builder.
		Insert(r.tbls.userTokenSub.Name).
		Columns(
			r.tbls.userTokenSub.Fields.UserID,
			r.tbls.userTokenSub.Fields.Token,
			r.tbls.userTokenSub.Fields.Category,
			r.tbls.userTokenSub.Fields.CurrentInterest,
			r.tbls.userTokenSub.Fields.Threshold,
			r.tbls.userTokenSub.Fields.Method,
			r.tbls.userTokenSub.Fields.ScanDate,
		).
		Values(
			params.UserID,
			params.Token,
			params.Category,
			params.Interest,
			params.Threshold,
			params.Method,
			params.ScanDate,
		).
		Suffix(fmt.Sprintf("RETURNING %s", r.tbls.userTokenSub.Fields.ID)).
		ToSql()
	if err != nil {
		return "", commonRepo.ParsePostgresError(op, err)
	}

	var id string

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return "", commonRepo.ParsePostgresError(op, err)
	}

	return id, nil
}

// GetTokenSubs return all user's token subs
func (r *Repository) GetTokenSubs(
	ctx context.Context,
	userID string,
	limit uint64,
	offset uint64,
) ([]*models.UserTokenSub, error) {
	op := "Repository.GetTokenSubs"

	query, args, err := r.db.Builder.
		Select(
			r.tbls.userTokenSub.Fields.ID,
			r.tbls.userTokenSub.Fields.UserID,
			r.tbls.userTokenSub.Fields.Token,
			r.tbls.userTokenSub.Fields.Category,
			r.tbls.userTokenSub.Fields.CurrentInterest,
			r.tbls.userTokenSub.Fields.PreviousInterest,
			r.tbls.userTokenSub.Fields.Threshold,
			r.tbls.userTokenSub.Fields.Method,
			r.tbls.userTokenSub.Fields.ScanDate,
			r.tbls.userTokenSub.Fields.CreatedAt,
		).
		Where(sq.Eq{r.tbls.userTokenSub.Fields.UserID: userID}).
		OrderBy(fmt.Sprintf("%s DESC", r.tbls.userTokenSub.Fields.CreatedAt)).
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}
	defer rows.Close()

	var subs []*models.UserTokenSub

	for rows.Next() {
		var sub models.UserTokenSub

		if err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.Token,
			&sub.Category,
			&sub.CurrentInterest,
			&sub.PreviousInterest,
			&sub.Threshold,
			&sub.Method,
			&sub.ScanDate,
			&sub.CreatedAt,
		); err != nil {
			return nil, commonRepo.ParsePostgresError(op, err)
		}

		subs = append(subs, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	if len(subs) == 0 {
		return nil, fmt.Errorf("[%s] user does not have any token subs: %w", op, commonRepo.ErrNotFound)
	}

	return subs, nil
}

// DeleteTokenSub deletes a token subscription from database
func (r *Repository) DeleteTokenSub(ctx context.Context, id string) error {
	op := "Repository.DeleteTokenSub"

	query, args, err := r.db.Builder.
		Delete(r.tbls.userTokenSub.Name).
		Where(sq.Eq{r.tbls.userTokenSub.Fields.ID: id}).
		ToSql()
	if err != nil {
		return commonRepo.ParsePostgresError(op, err)
	}

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return commonRepo.ParsePostgresError(op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("[%s] token sub does not exist: %w", op, commonRepo.ErrNotFound)
	}

	return nil
}
