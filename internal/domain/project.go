package domain

import (
	"errors"
	"net/url"
	"time"
)

type Project struct {
	ID             int           `json:"id"`
	Name           string        `json:"name,omitempty"`
	UserID         int           `json:"user_id,omitempty"`
	GitRepoURL     string        `json:"git_repo_url,omitempty"`
	Domain         string        `json:"domain,omitempty"`
	DeploymentName string        `json:"deployment_name,omitempty"`
	DockerfileDir  string        `json:"dockerfile_dir,omitempty"`
	Port           int           `json:"port,omitempty"`
	Environments   []Environment `json:"environments,omitempty"`
	Status         string        `json:"status,omitempty"`
	CreatedAt      time.Time     `json:"created_at,omitempty"`
	UpdatedAt      time.Time     `json:"updated_at,omitempty"`
}

func (p *Project) Validate() error {
	// Validate ID
	if p.ID < 0 {
		return errors.New("ID must be a positive integer")
	}
	// Validate Name
	if p.Name == "" {
		return errors.New("name is required")
	}

	// Validate UserID
	if p.UserID <= 0 {
		return errors.New("user ID must be a positive integer")
	}

	// Validate GitRepoURL
	if _, err := url.ParseRequestURI(p.GitRepoURL); err != nil {
		return errors.New("invalid git repository URL")
	}

	// Validate Domain
	if p.Domain == "" {
		return errors.New("domain is required")
	}

	// Validate DeploymentName
	if p.DeploymentName == "" {
		return errors.New("deployment name is required")
	}

	// Validate DockerfileDir
	if p.DockerfileDir == "" {
		return errors.New("dockerfile directory is required")
	}

	// Validate Port
	if p.Port <= 0 || p.Port > 65535 {
		return errors.New("port must be a positive integer between 1 and 65535")
	}

	// Validate Environments
	for _, env := range p.Environments {
		if err := env.Validate(); err != nil {
			return err
		}
	}

	// Validate Status
	if p.Status == "" {
		return errors.New("status is required")
	}

	return nil
}
