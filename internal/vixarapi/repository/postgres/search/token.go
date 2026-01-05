package search

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

// tokenHeader helper struct which represents a token search result header
type tokenHeader struct {
	name     string
	category string
}

// isEqTokenHeader checks if two token headers are equal
func isEqTokenHeader(a, b *tokenHeader) bool {
	return a.name == b.name && a.category == b.category
}

// SearchTokenParams parameters for token search query
type SearchTokenParams struct {
	Token    string
	Category *string
	Start    time.Time
	End      time.Time
}

// SearchTokenInfo searches for token information in the repository
func (r *Repository) SearchTokenInfo(
	ctx context.Context,
	params *SearchTokenParams,
) ([]models.TokenInfo, error) {
	op := "Repository.SearchTokenInfo"

	// create filter for where statement in the query
	filter := []sq.Sqlizer{
		sq.Expr(fmt.Sprintf("%s %% lower(?)", r.tbls.search.Fields.TokenName), params.Token),
		sq.GtOrEq{r.tbls.search.Fields.ScrapeDate: params.Start},
		sq.LtOrEq{r.tbls.search.Fields.ScrapeDate: params.End},
	}

	if params.Category != nil {
		filter = append(filter, sq.Eq{r.tbls.search.Fields.Category: *params.Category})
	}

	// prepare query
	query, args, err := r.db.Builder.
		Select(
			r.tbls.search.Fields.TokenName,
			r.tbls.search.Fields.Category,
			r.tbls.search.Fields.ScrapeDate,
			r.tbls.search.Fields.Interest,
			fmt.Sprintf("1.0 * %s / %s", r.tbls.search.Fields.Interest, r.tbls.search.Fields.GlobalMedian),
			fmt.Sprintf("1.0 * %s / %s", r.tbls.search.Fields.Interest, r.tbls.search.Fields.CategoryMedian),
			r.tbls.search.Fields.Sentiment,
		).
		From(r.tbls.search.Name).
		Where(sq.And(filter)).
		OrderBy(
			fmt.Sprintf("similarity(%s, lower($1)) DESC", r.tbls.search.Fields.TokenName),
			fmt.Sprintf("%s", r.tbls.search.Fields.Category),
			fmt.Sprintf("%s ASC", r.tbls.search.Fields.ScrapeDate),
		).
		Limit(searchLimit).
		ToSql()
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	ctxutils.GetLogger(ctx).Debugf("[%s] search token info query: %s, args: %v", op, query, args)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}
	defer rows.Close()

	var (
		cth     = tokenHeader{}
		pth     = tokenHeader{}
		res     = make([]models.TokenInfo, 0)
		records = make([]models.TokenRecord, 0)
	)

	for rows.Next() {
		var record models.TokenRecord

		if err := rows.Scan(
			&cth.name,
			&cth.category,
			&record.ScrapeDate,
			&record.Interest,
			&record.GlobalInterest,
			&record.CategoryInterest,
			&record.Sentiment,
		); err != nil {
			return nil, commonRepo.ParsePostgresError(op, err)
		}

		// is the same token as previous row -> add record to current token info
		if isEqTokenHeader(&cth, &pth) {
			records = append(records, record)
			continue
		}

		// new token encountered -> save previous token info
		if len(records) > 0 {
			res = append(res, models.TokenInfo{
				TokenName: pth.name,
				Category:  pth.category,
				Records:   records,
			})
		}

		records = make([]models.TokenRecord, 1)
		records[0] = record

		pth = cth
	}

	// save the last token info
	if len(records) > 0 {
		res = append(res, models.TokenInfo{
			TokenName: pth.name,
			Category:  pth.category,
			Records:   records,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	// check if any rows were returned
	if len(res) == 0 {
		return nil, fmt.Errorf("[%s] failed to find token info: %w", op, commonRepo.ErrNotFound)
	}

	return res, nil
}
