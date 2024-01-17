package controller

import (
	"github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/config"
	"github.com/zeet-co/kang/internal/storage"
	"github.com/zeet-co/kang/internal/zeet"
)

type Controller struct {
	db   *storage.DB
	zeet *zeet.Client
}

func NewController(cfg *config.Config) (*Controller, error) {

	db, err := storage.NewDB(cfg.DBConnectionString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.AutoMigrate(); err != nil {
		return nil, errors.WithStack(err)
	}

	zeetCli := zeet.New(cfg.ZeetAPIKey)

	return &Controller{
		db:   db,
		zeet: zeetCli,
	}, nil
}
