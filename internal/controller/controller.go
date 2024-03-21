package controller

import (
	"github.com/zeet-co/kang/internal/config"
	"github.com/zeet-co/kang/internal/zeet"
)

type Controller struct {
	zeet      *zeet.Client
	groupName string
}

func NewController(cfg *config.Config) (*Controller, error) {

	zeetCli := zeet.New(cfg.ZeetAPIKey)

	return &Controller{
		zeet:      zeetCli,
		groupName: cfg.ZeetGroupName,
	}, nil
}
