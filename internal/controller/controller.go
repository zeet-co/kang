package controller

import (
	"github.com/zeet-co/kang/internal/config"
	"github.com/zeet-co/kang/internal/zeet"
)

const ZeetGroupName = "kang"

type Controller struct {
	zeet *zeet.Client
}

func NewController(cfg *config.Config) (*Controller, error) {

	zeetCli := zeet.New(cfg.ZeetAPIKey)

	return &Controller{
		zeet: zeetCli,
	}, nil
}
