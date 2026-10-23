package service

import (
	"context"
	"fmt"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit,
	offset *int,
) ([]domain.User, error) {
	if err := domain.ValidateLimitAndOffset(limit, offset); err != nil {
		return nil, err
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}
