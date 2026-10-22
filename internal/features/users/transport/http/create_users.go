package http

import "net/http"

type CreateUserRequest struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MiddleName   *string `json:"middle_name"`
	TelegramUser *string `json:"telegram_user"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
