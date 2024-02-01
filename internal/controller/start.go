package controller

import (
	"context"
	"fmt"
	"strings"

	stdErrors "errors"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

type StartEnvironmentOpts struct {
	Overrides  map[uuid.UUID]map[string]string
	EnvName    string
	ProjectIDs []uuid.UUID
	TeamID     uuid.UUID
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

	errs := make([]error, len(projects))

	for i, p := range projects {
		//TODO scale down resources on branch deployments (?)
		//TODO handle database linking

		newName := fmt.Sprintf("%s-%s", p.Name, opts.EnvName)
		newProjectID, err := c.zeet.DuplicateProject(ctx, p.ID, groupID, subGroupID, newName)
		if err != nil {
			if err != v0.AlreadyExistsError {
				errs[i] = errors.WithStack(errors.Wrap(err, "could not duplicate project"))
				continue
			}
			fmt.Printf("Could not duplicate project %s, as it already exists. Checking for overrides..\n", newName)
			newProjectID, err = c.zeet.GetRepoByName(ctx, newName)
			if err != nil {
				errs[i] = errors.WithStack(errors.Wrap(err, "could not fetch project"))
				continue
			}

		}

		if opts.Overrides[p.ID] != nil {
			fmt.Printf("Found overrides applying to %s: parsing now\n", newName)
			override := opts.Overrides[p.ID]
			errs[i] = c.applyOverrides(ctx, newProjectID, override)
		}
		fmt.Printf("Done with project %s!\n", newName)
	}

	return stdErrors.Join(errs...)
}

func overrideToUpdateInput(pID uuid.UUID, overrides map[string]string) (*v0.UpdateProjectInput, error) {
	out := v0.UpdateProjectInput{
		Id: pID,
	}

	//TODO handle ref to symbolic value in another project
	err, anyFieldSet := assignValues(&out, overrides)

	if err != nil {
		return nil, err
	}

	if !anyFieldSet {
		return nil, nil
	}

	return &out, nil
}

func (c *Controller) applyOverrides(ctx context.Context, newProjectID uuid.UUID, override map[string]string) error {
	updateInput, err := overrideToUpdateInput(newProjectID, override)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not parse override %s", override)))
	}

	if updateInput != nil {
		fmt.Printf("Applying override stmt %s to %s\n", override, newProjectID)
		// fmt.Printf("%#v\n", updateInput)

		if err = c.zeet.UpdateProject(ctx, newProjectID, updateInput); err != nil {
			return errors.WithStack(errors.Wrap(err, "could not apply config overrides"))
		}
	}

	envs, envPresent := checkOverridesForEnvs(override)
	if envPresent {
		fmt.Printf("Applying env overrides to %s: %s\n", newProjectID, envs)
		if err = c.zeet.UpdateEnvs(ctx, newProjectID, envs); err != nil {
			return errors.WithStack(errors.Wrap(err, "could not apply env var override"))
		}
	}

	return nil
}

func checkOverridesForEnvs(override map[string]string) (map[string]string, bool) {
	out := make(map[string]string)
	isEnvPresent := false

	for k, v := range override {
		if strings.HasPrefix(k, "env.") {
			out[k[4:]] = v
			isEnvPresent = true
		}
	}

	return out, isEnvPresent
}
