package controller

import (
	"github.com/google/uuid"
	"github.com/zeet-co/kang/internal/storage/table"
)

type CreateEnvironmentOptions struct {
	Name       string
	ProjectIDs []uuid.UUID
}

func (c *Controller) CreateEnvironment(opts CreateEnvironmentOptions) error {

	//TODO check that each project ID exists in Zeet from the vantage of the authenticated user

	env := table.Environment{
		Name:       opts.Name,
		ProjectIDs: opts.ProjectIDs,
	}

	return c.db.DB.Save(&env).Error
}
