package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

type SearchTokenParams struct {
	Token string
	Start time.Time
	End   time.Time
}

// SearchTokenInfo searches for token information in the repository
func (r *Repository) SearchTokenInfo(
	ctx context.Context,
	params *SearchTokenParams,
) ([]models.TokenInfo, error) {
	op := "Repository.SearchTokenInfo"

	// prepare query
	query, args, err := r.db.Builder.
		Select(
			r.tbl.Fields.TokenName,
			r.tbl.Fields.ScrapeDate,
			r.tbl.Fields.Interest,
			fmt.Sprintf("1.0 * %s / %s", r.tbl.Fields.Interest, r.tbl.Fields.MedianInterest),
			r.tbl.Fields.Sentiment,
		).
		From(r.tbl.Name).
		Where(
			sq.And{
				sq.Expr(fmt.Sprintf("%s %% lower(?)", r.tbl.Fields.TokenName), params.Token),
				sq.GtOrEq{r.tbl.Fields.ScrapeDate: params.Start},
				sq.LtOrEq{r.tbl.Fields.ScrapeDate: params.End},
			},
		).
		OrderBy(
			fmt.Sprintf("similarity(%s, lower($1)) DESC", r.tbl.Fields.TokenName),
			fmt.Sprintf("%s ASC", r.tbl.Fields.ScrapeDate),
		).
		Limit(searchLimit).
		ToSql()
	if err != nil {
		return nil, parsePostgresError(op, err)
	}

	ctxutils.GetLogger(ctx).Debugf("[%s] search token info query: %s, args: %v", op, query, args)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, parsePostgresError(op, err)
	}
	defer rows.Close()

	var (
		currTokenName = ""
		prvTokenName  = ""
		res           = make([]models.TokenInfo, 0)
		records       = make([]models.TokenRecord, 0)
	)

	for rows.Next() {
		var record models.TokenRecord

		if err := rows.Scan(
			&currTokenName,
			&record.ScrapeDate,
			&record.Interest,
			&record.NormalizedInterest,
			&record.Sentiment,
		); err != nil {
			return nil, parsePostgresError(op, err)
		}

		// is the same token as previous row -> add record to current token info
		if currTokenName == prvTokenName {
			records = append(records, record)
			continue
		}

		// new token encountered -> save previous token info
		if len(records) > 0 {
			res = append(res, models.TokenInfo{
				TokenName: prvTokenName,
				Records:   records,
			})
		}

		records = make([]models.TokenRecord, 1)
		records[0] = record

		prvTokenName = currTokenName
	}

	// save the last token info
	if len(records) > 0 {
		res = append(res, models.TokenInfo{
			TokenName: prvTokenName,
			Records:   records,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, parsePostgresError(op, err)
	}

	// check if any rows were returned
	if len(res) == 0 {
		return nil, fmt.Errorf("[%s] failed to find token info: %w", op, ErrNotFound)
	}

	return res, nil
}
