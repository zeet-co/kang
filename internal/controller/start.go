package controller

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
	"golang.org/x/sync/errgroup"
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

	projects := make([]*v0.Repo, len(opts.ProjectIDs))

	var wg sync.WaitGroup
	// var repos []*Repo
	eg := new(errgroup.Group)

	// Assuming you have a slice of inputs for GetRepo
	for i, id := range opts.ProjectIDs {
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
