package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/keenywheels/backend/internal/vixarapi/models"
)

// GetUserByVKID return user by his vkid
func (r *Repository) GetUserByVKID(ctx context.Context, vkid int64) (*models.User, error) {
	op := "Repository.GetUserByVKID"

	query, args, err := r.db.Builder.
		Select(
			r.tbls.user.Fields.ID,
			r.tbls.user.Fields.Username,
			r.tbls.user.Fields.Email,
			r.tbls.user.Fields.TgUser,
			r.tbls.user.Fields.VKID,
			r.tbls.user.Fields.CreatedAt,
		).
		From(r.tbls.user.Name).
		Where(sq.Eq{r.tbls.user.Fields.VKID: vkid}).
		ToSql()
	if err != nil {
		return nil, parsePostgresError(op, err)
	}

	var user models.User

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.TgUser,
		&user.VKID,
		&user.CreatedAt,
	); err != nil {
		return nil, parsePostgresError(op, err)
	}

	return &user, nil
}

// RegisterVKUser registers a new vk user
func (r *Repository) RegisterVKUser(ctx context.Context, user *models.User) (*models.User, error) {
	op := "Repository.RegisterVKUser"

	query, args, err := r.db.Builder.
		Insert(r.tbls.user.Name).
		Columns(
			r.tbls.user.Fields.Username,
			r.tbls.user.Fields.Email,
			r.tbls.user.Fields.TgUser,
			r.tbls.user.Fields.VKID,
		).
		Values(
			user.Username,
			user.Email,
			user.TgUser,
			user.VKID,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s, %s, %s, %s",
			r.tbls.user.Fields.ID,
			r.tbls.user.Fields.Username,
			r.tbls.user.Fields.Email,
			r.tbls.user.Fields.TgUser,
			r.tbls.user.Fields.VKID,
			r.tbls.user.Fields.CreatedAt,
		)).
		ToSql()
	if err != nil {
		return nil, parsePostgresError(op, err)
	}

	var regUser models.User

	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&regUser.ID,
		&regUser.Username,
		&regUser.Email,
		&regUser.TgUser,
		&regUser.VKID,
		&regUser.CreatedAt,
	); err != nil {
		return nil, parsePostgresError(op, err)
	}

	return &regUser, nil
}
