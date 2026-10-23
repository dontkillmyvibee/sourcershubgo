package postgres

import (
	"time"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

type UserModel struct {
	ID           uuid.UUID  `db:"id"`
	FirstName    string     `db:"first_name"`
	SecondName   string     `db:"second_name"`
	MiddleName   *string    `db:"middle_name"`
	TelegramUser *string    `db:"telegram_user"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.FirstName,
			user.SecondName,
			user.MiddleName,
			user.TelegramUser,
			user.CreatedAt,
			user.UpdatedAt,
			user.DeletedAt,
		)
	}

	return userDomains
}
