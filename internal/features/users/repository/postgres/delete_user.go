package postgres

import (
	"context"
	"fmt"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    UPDATE sourcershub.users
    SET
        deleted_at = now(),
        updated_at = now()
    WHERE id = $1
      AND deleted_at IS NULL
`

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id '%s' not found: %w", userID, errors.ErrNotFound)
	}

	return nil
}
