package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

func (c *Controller) StopEnvironment(ctx context.Context, envName string, teamID uuid.UUID) error {

	teamName, err := c.zeet.GetTeamName(ctx, teamID)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	groupName := fmt.Sprintf("%s/%s", *teamName, c.groupName)
	subGroup := envName

	fmt.Printf("Stopping environment located in %s/%s\n", groupName, subGroup)

	errs := []error{}

	group, err := c.zeet.GetGroup(ctx, groupName)

	if err != nil {
		if err != v0.NotFoundError {
			return pkgErrors.WithStack(err)
		}
		// no group/subgroup = successfully deleted, exit
		fmt.Printf("Group %s doesn't exist\n", groupName)
	}

	sg := findSubGroup(group, subGroup)

	if sg == nil {
		// no subgroup = successfully deleted, exit
		fmt.Printf("SubGroup %s doesn't exist\n", subGroup)
		return nil
	}

	ids, err := c.getProjectsInSubGroup(ctx, group.Name, sg)
	if err != nil {
		return err
	}

	if ids != nil {
		for _, id := range ids {
			fmt.Printf("Deleting project %s\n", id)
			if err := c.zeet.DeleteProject(ctx, id); err != nil {
				errs = append(errs, fmt.Errorf(fmt.Sprintf("failed to delete %s", err)))
			}
		}

	}

	if compositeErr := errors.Join(errs...); compositeErr != nil {
		return err
	}

	c.zeet.DeleteSubGroup(ctx, sg.ID)

	return nil
}

func (c *Controller) getProjectsInSubGroup(ctx context.Context, groupName string, sg *v0.SubGroup) ([]uuid.UUID, error) {

	fmt.Printf("Found %d projects in %s/%s\n", len(sg.Projects), groupName, sg.Name)

	ids := make([]uuid.UUID, len(sg.Projects))

	for i, p := range sg.Projects {
		ids[i] = p.ID
	}

	return ids, nil
}

func findSubGroup(group *v0.GetSubGroupsForGroupResp, subGroup string) (sg *v0.SubGroup) {
	for _, gSg := range group.SubGroups {
		if gSg.Name == subGroup {
			sg = &gSg
			break
		}
	}

	return sg
}
