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
		From(r.tbls.userTokenSub.Name).
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

// UpdateTokenSubParams represents parameters for updating a token subscription in database
type UpdateTokenSubParams struct {
	ID        string
	Threshold float64
	Method    string
}

type UpdateTokenSubResult struct {
	CurrInterest float64
	PrvInterest  float64
}

// UpdateTokenSub updates a token subscription in database
func (r *Repository) UpdateTokenSub(ctx context.Context, params *UpdateTokenSubParams) (*UpdateTokenSubResult, error) {
	// TODO: parameterize interval in query if change scrape interval
	var (
		op    = "Repository.UpdateTokenSub"
		query = `
			WITH
				curr_token_info AS (SELECT *
									FROM user_token_sub uts
											 JOIN mv_token_search ts
												  ON uts.token = ts.token_name AND uts.category = ts.category
									WHERE uts.id = $1
									  AND (ts.scrape_date = uts.scan_date)),
				prv_token_info AS (SELECT *
								   FROM user_token_sub uts
											JOIN mv_token_search ts
												 ON uts.token = ts.token_name AND uts.category = ts.category
								   WHERE uts.id = $1
									 AND (ts.scrape_date = uts.scan_date - INTERVAL '1 days')),
				new_data AS (SELECT $2::numeric            AS threshold,
									$3::text               AS method,
									(SELECT CASE $3
												WHEN 'global_median' THEN interest / global_median
												WHEN 'category_median' THEN interest / category_median
												ELSE interest
												END
									 FROM curr_token_info) AS curr_interest,
									(SELECT CASE $3
												WHEN 'global_median' THEN interest / global_median
												WHEN 'category_median' THEN interest / category_median
												ELSE interest
												END
									 FROM prv_token_info)  AS prv_interest)
			UPDATE user_token_sub uts
			SET method        = nd.method,
				threshold     = nd.threshold,
				curr_interest = nd.curr_interest,
				prv_interest  = nd.prv_interest
			FROM new_data nd
			WHERE uts.id = $1
			RETURNING nd.curr_interest, nd.prv_interest;
		`
		args = []any{params.ID, params.Threshold, params.Method}
	)

	// update user token sub
	var res UpdateTokenSubResult
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&res.CurrInterest,
		&res.PrvInterest,
	); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	return &res, nil
}
