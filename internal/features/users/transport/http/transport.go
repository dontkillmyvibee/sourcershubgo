package http

import (
	"context"
	"net/http"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	PatchUser(ctx context.Context, id uuid.UUID, patch domain.UserPatch) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
