package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

func (c *Controller) StopEnvironment(ctx context.Context, envName string) error {

	groupName := ZeetGroupName
	subGroup := envName

	fmt.Printf("Stopping environment located in %s/%s\n", groupName, subGroup)

	errs := []error{}

	ids, err := c.getProjectsInSubGroup(ctx, groupName, subGroup)
	if err != nil {
		return err
	}

	if ids == nil {
		fmt.Println("No Projects found; exiting")
		return nil
	}

	for _, id := range ids {
		fmt.Printf("Deleting project %s\n", id)
		if err := c.zeet.DeleteProject(context.Background(), id); err != nil {
			errs = append(errs, fmt.Errorf(fmt.Sprintf("failed to delete %s", err)))
		}
	}

	return errors.Join(errs...)
}

func (c *Controller) getProjectsInSubGroup(ctx context.Context, groupName, subGroup string) ([]uuid.UUID, error) {
	group, err := c.zeet.GetGroup(context.Background(), groupName)

	if err != nil {
		if err == v0.NotFoundError {
			// no group/subgroup = successfully deleted, exit
			fmt.Printf("Group %s doesn't exist\n", groupName)
			return nil, nil
		}
		return nil, pkgErrors.WithStack(err)
	}

	var (
		sg *v0.SubGroup
	)

	for _, gSg := range group.SubGroups {
		if gSg.Name == subGroup {
			sg = &gSg
			break
		}
	}

	if sg == nil {
		// no subgroup = successfully deleted, exit
		fmt.Printf("SubGroup %s doesn't exist\n", subGroup)
		return nil, nil
	}

	fmt.Printf("Found %d projects in %s/%s\n", len(sg.Projects), groupName, subGroup)

	ids := make([]uuid.UUID, len(sg.Projects))

	for i, p := range sg.Projects {
		ids[i] = p.ID
	}

	return ids, nil
}
