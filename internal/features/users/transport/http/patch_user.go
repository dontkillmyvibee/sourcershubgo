package http

import (
	"fmt"
	"net/http"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/request"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/response"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/types"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FirstName    types.Nullable[string] `json:"first_name"`
	SecondName   types.Nullable[string] `json:"second_name"`
	MiddleName   types.Nullable[string] `json:"middle_name"`
	TelegramUser types.Nullable[string] `json:"telegram_user"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FirstName.Set {
		if r.FirstName.Value == nil {
			return fmt.Errorf("`first_name` can`t be null")
		}
		if firstNameLen := len([]rune(*r.FirstName.Value)); firstNameLen < 1 || firstNameLen > 30 {
			return fmt.Errorf("`first_name` length must be between 1 and 30")
		}
	}

	if r.SecondName.Set {
		if r.SecondName.Value == nil {
			return fmt.Errorf("`second_name` can`t be null")
		}
		if secondNameLen := len([]rune(*r.SecondName.Value)); secondNameLen < 1 || secondNameLen > 30 {
			return fmt.Errorf("`second_name` length must be between 1 and 30")
		}
	}

	if r.MiddleName.Set {
		if r.MiddleName.Value != nil {
			if middleNameLen := len([]rune(*r.MiddleName.Value)); middleNameLen < 1 || middleNameLen > 30 {
				return fmt.Errorf("`middle_name` length must be between 1 and 30")
			}
		}
	}

	if r.TelegramUser.Set {
		if r.TelegramUser.Value != nil {
			if telegramUserLen := len([]rune(*r.TelegramUser.Value)); telegramUserLen < 5 || telegramUserLen > 32 {
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	userID, err := utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	var req PatchUserRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	resp := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(resp, http.StatusOK)
}

func userPatchFromRequest(req PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FirstName:    req.FirstName.ToDomain(),
		SecondName:   req.SecondName.ToDomain(),
		MiddleName:   req.MiddleName.ToDomain(),
		TelegramUser: req.TelegramUser.ToDomain(),
	}
}
