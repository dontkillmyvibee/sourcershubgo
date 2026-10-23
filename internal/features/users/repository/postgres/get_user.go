package postgres

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
	coreerrors "github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT
        id, first_name, second_name, middle_name,
        telegram_user, created_at, updated_at, deleted_at
    FROM sourcershub.users
    WHERE id = $1
      AND deleted_at IS NULL
`

	row := r.pool.QueryRow(ctx, query, id)

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
				"user with id='%v' not found: %w",
				id,
				coreerrors.ErrNotFound,
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
