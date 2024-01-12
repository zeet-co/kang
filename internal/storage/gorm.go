package storage

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zeet-co/kang/internal/storage/table"
)

type DB struct {
	DB *gorm.DB
}

func NewDB(connStr string) (*DB, error) {

	db, err := gorm.Open(postgres.Open(connStr))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &DB{
		db,
	}, nil
}

func (d *DB) AutoMigrate() error {
	return d.DB.AutoMigrate(
		&table.Environment{},
	)
}
