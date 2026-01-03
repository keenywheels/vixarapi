package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/keenywheels/backend/internal/processor/models"
)

const (
	maxBatchSize = 1000
)

// InsertTokens inserts multiple token data records into the database
func (r *Repository) InsertTokens(ctx context.Context, tokens []models.TokenData) error {
	var errs []error

	batch := make([]*models.TokenData, 0, maxBatchSize)
	for _, token := range tokens {
		batch = append(batch, &token)

		if len(batch) >= maxBatchSize {
			if err := r.insertBatch(ctx, batch); err != nil {
				errs = append(errs, err) // collect errors but continue processing
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := r.insertBatch(ctx, batch); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// insertBatch inserts a batch of token data records into the database
func (r *Repository) insertBatch(ctx context.Context, tokens []*models.TokenData) error {
	op := "Repository.insertBatch"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[%s] failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, token := range tokens {
		// create query
		query, args, err := r.db.Builder.Insert(r.tbl.Name).
			Columns(
				r.tbl.Fields.TokenName,
				r.tbl.Fields.Interest,
				r.tbl.Fields.Category,
				r.tbl.Fields.SiteName,
				r.tbl.Fields.Date,
				r.tbl.Fields.Sentiment,
			).
			Values(
				token.TokenName,
				token.Interest,
				token.Category,
				token.SiteName,
				token.Date,
				token.Sentiment,
			).
			ToSql()
		if err != nil {
			return fmt.Errorf("[%s] failed to build insert query: %w", op, err)
		}

		batch.Queue(query, args...)
	}

	res := tx.SendBatch(ctx, batch)

	if err := res.Close(); err != nil {
		return fmt.Errorf("[%s] failed to execute batch insert: %w", op, err)
	}

	return tx.Commit(ctx)
}
