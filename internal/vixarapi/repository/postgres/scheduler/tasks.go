package scheduler

import (
	"context"
	"fmt"
)

// updateSearchTable performs the update of the search table
func (r *Repository) updateSearchTable(ctx context.Context) error {
	var (
		op    = "Repository.updateSearchTable"
		query = fmt.Sprintf("REFRESH MATERIALIZED VIEW CONCURRENTLY %s;", r.tbls.search.Name)
	)

	if _, err := r.db.Pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("[%s] failed to refresh materialized view: %w", op, err)
	}

	return nil
}
