package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/storage/table"
)

type StartEnvironmentOpts struct {
	ProjectBranchOverrides map[uuid.UUID]string
}

func (c *Controller) StartEnvironment(envID, teamID uuid.UUID, opts StartEnvironmentOpts) error {
	var env table.Environment

	err := c.db.DB.First(&env, envID).Error
	if err != nil {
		return errors.WithStack(err)
	}

	group := "kang"
	subGroup := env.Name

	groupID, subGroupID, err := c.zeet.EnsureGroupsExist(group, subGroup, teamID)
	if err != nil {
		return errors.WithStack(err)
	}

	for i, id := range env.ProjectIDs {
		//TODO scale down resources on branch deployments (?)
		//TODO handle database linking

		newName := fmt.Sprintf("kang-%s_%d", subGroup, i)
		pID, err := c.zeet.DuplicateProject(context.Background(), id, groupID, subGroupID, newName)
		if err != nil {
			return errors.WithStack(err)
		}

		if opts.ProjectBranchOverrides[id] != "" {
			if err = c.zeet.UpdateProjectBranch(context.Background(), pID, opts.ProjectBranchOverrides[id]); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}
