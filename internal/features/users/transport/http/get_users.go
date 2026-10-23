package http

import (
	"fmt"
	"net/http"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/response"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query params",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	resp := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(resp, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
