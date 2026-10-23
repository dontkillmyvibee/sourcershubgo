package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"uuid"

	coreerrors "github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
)

var telegramUsernameRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9_]{3,30}[a-zA-Z0-9]$`,
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	SecondName   string
	MiddleName   *string
	TelegramUser *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewUser(
	id uuid.UUID,
	firstName,
	secondName string,
	middleName,
	telegramUser *string,
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) User {
	return User{
		ID:           id,
		FirstName:    firstName,
		SecondName:   secondName,
		MiddleName:   middleName,
		TelegramUser: telegramUser,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
	}
}

func NewUserUninitialized(
	firstName string,
	secondName string,
	middleName *string,
	telegramUser *string,
) User {

	user := User{
		ID:           UninitializedID,
		FirstName:    firstName,
		SecondName:   secondName,
		MiddleName:   middleName,
		TelegramUser: telegramUser,
	}

	user.Normalize()

	return user
}

func (u *User) Normalize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.SecondName = strings.TrimSpace(u.SecondName)

	if u.MiddleName != nil {
		value := strings.TrimSpace(*u.MiddleName)
		u.MiddleName = &value
	}

	if u.TelegramUser != nil {
		value := strings.TrimSpace(*u.TelegramUser)
		value = strings.TrimPrefix(value, "@")
		u.TelegramUser = &value
	}
}

func validateRequiredName(fieldName, value string) error {
	if length := len([]rune(value)); length < 1 || length > 30 {
		return fmt.Errorf(
			"invalid `%s` length, got: %d: %w",
			fieldName,
			length,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) Validate() error {
	if err := validateRequiredName("first_name", u.FirstName); err != nil {
		return err
	}

	if err := validateRequiredName("second_name", u.SecondName); err != nil {
		return err
	}

	if u.MiddleName != nil {
		if err := validateRequiredName("middle_name", *u.MiddleName); err != nil {
			return err
		}
	}

	if u.TelegramUser != nil {
		tgUser := *u.TelegramUser
		if telegramUserLength := len([]rune(tgUser)); telegramUserLength < 5 {
			return fmt.Errorf(
				"invalid `TelegramUser` length, got %d: %w",
				telegramUserLength,
				coreerrors.ErrInvalidArgument,
			)
		}

		if !telegramUsernameRegex.MatchString(tgUser) {
			return fmt.Errorf(
				"invalid `TelegramUser` format, got: %s: %w",
				tgUser,
				coreerrors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func ValidateLimitAndOffset(limit, offset *int) error {
	if limit != nil && *limit < 0 {
		return fmt.Errorf(
			"limit must be non-negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return fmt.Errorf(
			"offset must be non-negative: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserPatch struct {
	FirstName    Nullable[string]
	SecondName   Nullable[string]
	MiddleName   Nullable[string]
	TelegramUser Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FirstName.Set && p.FirstName.Value == nil {
		return fmt.Errorf("`first_name` can`t be patched to null: %w", coreerrors.ErrInvalidArgument)
	}

	if p.SecondName.Set && p.SecondName.Value == nil {
		return fmt.Errorf("`second_name` can`t be patched to null: %w", coreerrors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("invalid patch: %w", err)
	}

	tmp := *u

	if patch.FirstName.Set {
		tmp.FirstName = *patch.FirstName.Value
	}
	if patch.SecondName.Set {
		tmp.SecondName = *patch.SecondName.Value
	}
	if patch.MiddleName.Set {
		tmp.MiddleName = patch.MiddleName.Value
	}
	if patch.TelegramUser.Set {
		tmp.TelegramUser = patch.TelegramUser.Value
	}

	tmp.Normalize()

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
