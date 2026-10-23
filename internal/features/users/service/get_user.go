package service

import (
	"context"
	"fmt"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	return user, nil
}
