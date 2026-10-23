package postgres

import (
	"context"
	"fmt"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT
        id, first_name, second_name, middle_name,
        telegram_user, created_at, updated_at, deleted_at
    FROM sourcershub.users
    WHERE deleted_at IS NULL
    ORDER BY id ASC
    LIMIT $1
    OFFSET $2
`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
