package domain

import (
	"errors"
	"time"
)

type Environment struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	EnvKey    string    `json:"name"`
	EnvValue  string    `json:"value"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (e *Environment) Validate() error {
	// Validate ProjectID
	if e.ProjectID <= 0 {
		return errors.New("project ID must be a positive integer")
	}

	// Validate EnvKey
	if e.EnvKey == "" {
		return errors.New("environment variable key is required")
	}

	// Validate EnvValue
	if e.EnvValue == "" {
		return errors.New("environment variable value is required")
	}

	return nil
}
