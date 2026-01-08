package search

import (
	"context"
	"fmt"
	"slices"

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
