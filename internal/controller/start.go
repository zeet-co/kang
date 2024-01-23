package controller

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/storage/table"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
	"golang.org/x/sync/errgroup"
)

type StartEnvironmentOpts struct {
	ProjectBranchOverrides map[uuid.UUID]string
}

func (c *Controller) StartEnvironment(envID, teamID uuid.UUID, opts StartEnvironmentOpts) error {
	var env table.Environment

	err := c.db.DB.First(&env, envID).Error
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not find environment"))
	}

	group := ZeetGroupName
	subGroup := env.Name

	groupID, subGroupID, err := c.zeet.EnsureGroupsExist(group, subGroup, teamID)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not ensure group / subgroup"))
	}

	projects := make([]*v0.Repo, len(env.ProjectIDs))

	var wg sync.WaitGroup
	// var repos []*Repo
	eg := new(errgroup.Group)

	// Assuming you have a slice of inputs for GetRepo
	for i, id := range env.ProjectIDs {
		wg.Add(1)
		id := id // capture range variable
		i := i
		eg.Go(func() error {
			defer wg.Done()
			repo, err := c.zeet.GetRepo(context.Background(), id)
			if err != nil {
				return err
			}
			// Use a mutex or other synchronization method if needed
			projects[i] = repo
			return nil
		})
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Check if any goroutines returned an error
	if err := eg.Wait(); err != nil {
		return errors.WithStack(errors.Wrap(err, "could not fetch project information"))
	}

	for _, p := range projects {
		//TODO scale down resources on branch deployments (?)
		//TODO handle database linking

		newName := fmt.Sprintf("%s-%s", p.Name, env.Name)
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
