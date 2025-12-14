package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
)

// SaveSearchQuery saves a search query for a user
func (r *Repository) SaveSearchQuery(ctx context.Context, userID string, query string) (*models.UserQuery, error) {
	op := "Repository.SaveSearchQuery"

	query, args, err := r.db.Builder.
		Insert(r.tbls.userQuery.Name).
		Columns(
			r.tbls.userQuery.Fields.UserID,
			r.tbls.userQuery.Fields.Query,
		).
		Values(
			userID,
			query,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s",
			r.tbls.userQuery.Fields.ID,
			r.tbls.userQuery.Fields.UserID,
			r.tbls.userQuery.Fields.Query,
			r.tbls.userQuery.Fields.CreatedAt,
		)).
		ToSql()
	if err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	var userQuery models.UserQuery

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&userQuery.ID,
		&userQuery.UserID,
		&userQuery.Query,
		&userQuery.CreatedAt,
	); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	return &userQuery, nil
}

// DeleteSearchQuery deletes a search query for a user
func (r *Repository) DeleteSearchQuery(ctx context.Context, id string) error {
	op := "Repository.DeleteSearchQuery"

	query, args, err := r.db.Builder.
		Delete(r.tbls.userQuery.Name).
		Where(sq.Eq{r.tbls.userQuery.Fields.ID: id}).
		ToSql()
	if err != nil {
		return commonRepo.ParsePostgresError(op, err)
	}

	tag, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return commonRepo.ParsePostgresError(op, err)
	}

	// nothing was deleted -> return not found
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("[%s] search query not found: %w", op, commonRepo.ErrNotFound)
	}

	return nil
}

// GetSearchQueries retrieves all search queries for a user
func (r *Repository) GetSearchQueries(
	ctx context.Context,
	userID string,
	limit uint64,
	offset uint64,
) ([]*models.UserQuery, error) {
	op := "Repository.GetSearchQueries"

	query, args, err := r.db.Builder.
		Select(
			r.tbls.userQuery.Fields.ID,
			r.tbls.userQuery.Fields.UserID,
			r.tbls.userQuery.Fields.Query,
			r.tbls.userQuery.Fields.CreatedAt,
		).
		From(r.tbls.userQuery.Name).
		Where(sq.Eq{r.tbls.userQuery.Fields.UserID: userID}).
		OrderBy(fmt.Sprintf("%s DESC", r.tbls.userQuery.Fields.CreatedAt)).
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

	var queries []*models.UserQuery

	for rows.Next() {
		var userQuery models.UserQuery

		if err := rows.Scan(
			&userQuery.ID,
			&userQuery.UserID,
			&userQuery.Query,
			&userQuery.CreatedAt,
		); err != nil {
			return nil, commonRepo.ParsePostgresError(op, err)
		}

		queries = append(queries, &userQuery)
	}

	if err := rows.Err(); err != nil {
		return nil, commonRepo.ParsePostgresError(op, err)
	}

	if len(queries) == 0 {
		return nil, fmt.Errorf("[%s] user does not have any queries: %w", op, commonRepo.ErrNotFound)
	}

	return queries, nil
}
