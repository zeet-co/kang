package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/storage/table"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

func (c *Controller) StopEnvironment(envID uuid.UUID) error {

	var env table.Environment

	err := c.db.DB.First(&env, envID).Error
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	groupName := ZeetGroupName
	subGroup := env.Name

	group, err := c.zeet.GetGroup(context.Background(), groupName)

	if err != nil {
		if err == v0.NotFoundError {
			// no group/subgroup = successfully deleted, exit
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
		return nil
	}

	ids := make([]uuid.UUID, len(sg.Projects))

	for i, p := range sg.Projects {
		ids[i] = p.ID
	}

	errs := []error{}

	for _, id := range ids {
		if err := c.zeet.DeleteProject(context.Background(), id); err != nil {
			errs = append(errs, fmt.Errorf(fmt.Sprintf("failed to delete %s", err)))
		}
	}

	return errors.Join(errs...)
}
