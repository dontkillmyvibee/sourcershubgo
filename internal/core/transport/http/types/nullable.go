package types

import (
	"encoding/json/v2"

	"github.com/dontkillmyvibee/sourcershubgo/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	n.Set = true

	if string(bytes) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
