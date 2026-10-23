package utils

import (
	"fmt"
	"net/http"
	"uuid"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil(), fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			errors.ErrInvalidArgument,
		)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil(), fmt.Errorf(
			"path value='%s' by key='%s' not a valid uuid: %w",
			pathValue,
			key,
			err,
			errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
