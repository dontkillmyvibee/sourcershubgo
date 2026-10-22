package domain

import (
	"time"
	"uuid"
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
