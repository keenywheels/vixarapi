package search

import (
	"context"
	"fmt"
	"slices"

	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

const (
	IntervalDays  = "days"
	IntervalHours = "hours"
)

// UpdateSearchTable performs the update of the search table
func (r *Repository) UpdateSearchTable(ctx context.Context) error {
	var (
		op    = "Repository.updateSearchTable"
		query = fmt.Sprintf("REFRESH MATERIALIZED VIEW CONCURRENTLY %s;", r.tbls.search.Name)
	)

	if _, err := r.db.Pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("[%s] failed to refresh materialized view: %w", op, err)
	}

	return nil
}

// UpdateUserTokenSubs updates all token subs
func (r *Repository) UpdateUserTokenSubs(ctx context.Context, intervalType string, amount int) error {
	var (
		op        = "Repository.UpdateUserTokenSubs"
		log       = ctxutils.GetLogger(ctx)
		queryTmpl = `
			WITH
				new_token_interest AS (SELECT uts.id  AS user_token_sub_id,
											  uts.method,
											  ts.scrape_date,
											  CASE uts.method
												  WHEN 'global_median' THEN ts.interest / ts.global_median
												  WHEN 'category_median' THEN ts.interest / ts.category_median
												  ELSE ts.interest
												  END AS interest
									   FROM user_token_sub uts
												JOIN mv_token_search ts
													 ON uts.token = ts.token_name AND uts.category = ts.category AND
														uts.scan_date + INTERVAL '%s' = ts.scrape_date)
			UPDATE user_token_sub uts
			SET curr_interest = nts.interest,
				prv_interest  = curr_interest,
				scan_date = nts.scrape_date
			FROM new_token_interest nts
			WHERE uts.id = nts.user_token_sub_id;
		`
	)

	// validate interval
	validIntervals := []string{IntervalDays, IntervalHours}
	if !slices.Contains(validIntervals, intervalType) {
		return fmt.Errorf("[%s] invalid interval type: %s", op, intervalType)
	} else if amount <= 0 {
		return fmt.Errorf("[%s] invalid amount: %d", op, amount)
	}

	// build interval string
	interval := fmt.Sprintf("%d %s", amount, intervalType)

	tag, err := r.db.Pool.Exec(ctx, fmt.Sprintf(queryTmpl, interval))
	if err != nil {
		return fmt.Errorf("[%s] failed to update user token subs: %w", op, err)
	}

	log.Infof("[%s] successfully updated %d records, interval=%s", op, tag.RowsAffected(), interval)

	return nil
}

// IncreasedTokenSubInfo represents info about increased token subs
type IncreasedTokenSubInfo struct {
	UserID           string
	Email            string
	Username         string
	Token            string
	Category         string
	CurrentInterest  float64
	PreviousInterest float64
	Threshold        float64
}

// GetIncreasedTokenSubs returns all token subs that were increased
func (r *Repository) GetIncreasedTokenSubs(
	ctx context.Context,
	limit uint64,
	offset uint64,
) ([]*IncreasedTokenSubInfo, error) {
	var (
		op    = "Repository.GetIncreasedTokenSubs"
		query = `
			SELECT
				u.id AS user_id,
				u.email,
				u.username,
				uts.token,
				uts.category,
				uts.curr_interest,
				uts.prv_interest,
				uts.threshold
			FROM user_token_sub uts
			JOIN users u ON uts.user_id = u.id
			WHERE curr_interest / prv_interest > threshold
			ORDER BY u.id, uts.id
			LIMIT $1 OFFSET $2;
		`
	)

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}
	defer rows.Close()

	// get user token subs
	var subs []*IncreasedTokenSubInfo
	for rows.Next() {
		var sub IncreasedTokenSubInfo

		if err := rows.Scan(
			&sub.UserID,
			&sub.Email,
			&sub.Username,
			&sub.Token,
			&sub.Category,
			&sub.CurrentInterest,
			&sub.PreviousInterest,
			&sub.Threshold,
		); err != nil {
			return nil, commonRepo.ParsePostgresError(op, err)
		}

		subs = append(subs, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	if len(subs) == 0 {
		return nil, fmt.Errorf("[%s] no token subs were increased: %w", op, commonRepo.ErrNotFound)
	}

	return subs, nil
}
