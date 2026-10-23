package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("validate json: %v: %w", err, errors.ErrInvalidArgument)
	}

	return nil
}
