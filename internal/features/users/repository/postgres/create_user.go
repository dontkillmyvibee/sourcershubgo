package postgres

import (
	"context"
	"fmt"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO sourcershub.users (first_name, second_name, middle_name, telegram_user)
	VALUES ($1, $2, $3, $4)
	RETURNING id, first_name, second_name, middle_name, telegram_user, created_at, updated_at, deleted_at
	`

	row := r.pool.QueryRow(ctx, query, user.FirstName, user.SecondName, user.MiddleName, user.TelegramUser)

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
		return domain.User{}, fmt.Errorf("scan user error: %w", err)
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
