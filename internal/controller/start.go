package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type StartEnvironmentOpts struct {
	ProjectBranchOverrides map[uuid.UUID]string
	EnvName                string
	ProjectIDs             []uuid.UUID
	TeamID                 uuid.UUID
}

func (c *Controller) StartEnvironment(opts StartEnvironmentOpts) error {
	group := ZeetGroupName
	subGroup := opts.EnvName

	groupID, subGroupID, err := c.zeet.EnsureGroupsExist(group, subGroup, opts.TeamID)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not ensure group / subgroup"))
	}

	ctx := context.Background()

	projects, err := c.zeet.GetProjectsByID(ctx, opts.ProjectIDs)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not fetch project information"))
	}

	for _, p := range projects {
		//TODO scale down resources on branch deployments (?)
		//TODO handle database linking

		newName := fmt.Sprintf("%s-%s", p.Name, opts.EnvName)
		pID, err := c.zeet.DuplicateProject(context.Background(), p.ID, groupID, subGroupID, newName)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "could not duplicate project"))
		}

		if opts.ProjectBranchOverrides[p.ID] != "" {
			if err = c.zeet.UpdateProjectBranch(context.Background(), pID, opts.ProjectBranchOverrides[p.ID]); err != nil {
				return errors.WithStack(errors.Wrap(err, "could not apply branch override"))
			}
		}
	}

	return nil
}
