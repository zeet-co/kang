package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

func (c *Controller) StopEnvironment(envName string) error {

	groupName := ZeetGroupName
	subGroup := envName

	fmt.Printf("Stopping environment located in %s/%s\n", groupName, subGroup)

	group, err := c.zeet.GetGroup(context.Background(), groupName)

	if err != nil {
		if err == v0.NotFoundError {
			// no group/subgroup = successfully deleted, exit
			fmt.Printf("Group doesn't exist; prior invocation must have succeeded, exiting\n")
			return nil
		}
		return pkgErrors.WithStack(err)
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
		fmt.Printf("SubGroup doesn't exist; prior invocation must have succeeded, exiting\n")
		return nil
	}

	fmt.Printf("Found %d projects in %s/%s; deleting now...\n", len(sg.Projects), groupName, subGroup)

	ids := make([]uuid.UUID, len(sg.Projects))

	for i, p := range sg.Projects {
		ids[i] = p.ID
	}

	errs := []error{}

	for _, id := range ids {
		fmt.Printf("Deleting project %s\n", id)
		if err := c.zeet.DeleteProject(context.Background(), id); err != nil {
			errs = append(errs, fmt.Errorf(fmt.Sprintf("failed to delete %s", err)))
		}
	}

	return errors.Join(errs...)
}
