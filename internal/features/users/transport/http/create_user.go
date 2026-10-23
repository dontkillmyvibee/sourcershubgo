package http

import (
	"net/http"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/request"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FirstName    string  `json:"first_name" validate:"required,min=1,max=30"`
	SecondName   string  `json:"second_name" validate:"required,min=1,max=30"`
	MiddleName   *string `json:"middle_name,omitempty" validate:"omitempty,min=1,max=30"`
	TelegramUser *string `json:"telegram_user,omitempty" validate:"omitempty,min=1,max=30"`
}

type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("Create User")
	var req CreateUserRequest

	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	resp := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(resp, http.StatusOK)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FirstName, dto.SecondName, dto.MiddleName, dto.TelegramUser)
}
