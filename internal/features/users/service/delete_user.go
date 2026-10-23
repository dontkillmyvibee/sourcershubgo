package service

import (
	"context"
	"fmt"
	"uuid"
)

func (s *UsersService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
