package domain

import (
	"errors"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	ExternalUID string    `json:"external_uid"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (u *User) Validate() error {
	// Validate ID
	if u.ID <= 0 {
		return errors.New("ID must be a positive integer")
	}

	// Validate ExternalUID
	if len(u.ExternalUID) > 128 {
		return errors.New("external UID must be 128 characters or less")
	}
	return nil
}
