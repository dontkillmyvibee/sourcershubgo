package http

import (
	"net/http"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/response"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	userID, err := utils.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	resp := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(resp, http.StatusOK)
}
