package postgres

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
	errors2 "github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE sourcershub.users
	SET 
	    first_name = $1,
	    second_name = $2,
	    middle_name = $3,
	    telegram_user = $4,
	    updated_at = now()
	WHERE id = $5 
	  AND deleted_at IS NULL
	RETURNING id,first_name,second_name,middle_name,telegram_user,created_at,updated_at,deleted_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.SecondName,
		user.MiddleName,
		user.TelegramUser,
		id,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.FirstName,
		&userModel.SecondName,
		&userModel.MiddleName,
		&userModel.TelegramUser,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
		&userModel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id '%v' concurrently accessed or not found: %w",
				id,
				errors2.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.FirstName,
		userModel.SecondName,
		userModel.MiddleName,
		userModel.TelegramUser,
		userModel.CreatedAt,
		userModel.UpdatedAt,
		userModel.DeletedAt,
	)

	return userDomain, nil
}
