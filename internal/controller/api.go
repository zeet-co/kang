package controller

import (
	"context"

	"github.com/google/uuid"
)

func (c *Controller) CheckProjectExists(id uuid.UUID) bool {
	_, err := c.zeet.GetRepo(context.Background(), id)
	return err == nil
}
