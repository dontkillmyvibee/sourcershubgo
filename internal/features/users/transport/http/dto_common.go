package http

import (
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

type UserDTOResponse struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	SecondName   string    `json:"second_name"`
	MiddleName   *string   `json:"middle_name,omitempty"`
	TelegramUser *string   `json:"telegram_user,omitempty"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		SecondName:   user.SecondName,
		MiddleName:   user.MiddleName,
		TelegramUser: user.TelegramUser,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
