package search

import (
	"context"
	"fmt"

	"github.com/keenywheels/backend/pkg/ctxutils"
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
func (r *Repository) UpdateUserTokenSubs(ctx context.Context) error {
	var (
		op    = "Repository.UpdateUserTokenSubs"
		log   = ctxutils.GetLogger(ctx)
		query = `
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
														uts.scan_date::date + 1 = ts.scrape_date)
			UPDATE user_token_sub uts
			SET curr_interest = nts.interest,
				prv_interest  = curr_interest,
				scan_date = nts.scrape_date
			FROM new_token_interest nts
			WHERE uts.id = nts.user_token_sub_id;
		`
	)

	tag, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("[%s] failed to update user token subs: %w", op, err)
	}

	log.Infof("[%s] successfully updated %d records", op, tag.RowsAffected())

	return nil
}
