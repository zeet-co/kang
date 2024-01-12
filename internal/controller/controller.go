package controller

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/storage"
	"github.com/zeet-co/kang/internal/zeet"
)

type Controller struct {
	db   *storage.DB
	zeet *zeet.Client
}

func NewController(zeetAPIKey string) (*Controller, error) {
	connStr := getConnStr()

	db, err := storage.NewDB(connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.AutoMigrate(); err != nil {
		return nil, errors.WithStack(err)
	}

	zeetCli := zeet.New(zeetAPIKey)

	return &Controller{
		db:   db,
		zeet: zeetCli,
	}, nil
}

func getConnStr() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbExtraOpts := os.Getenv("DB_EXTRA_OPTS")

	connStr := fmt.Sprintf("host=%s port=%s", dbHost, dbPort)

	if dbSSLMode != "" {
		connStr += fmt.Sprintf(" sslmode=%s", dbSSLMode)
	}

	if dbName != "" {
		connStr += fmt.Sprintf(" dbname=%s", dbName)
	}

	if dbUser != "" {
		connStr += fmt.Sprintf(" user=%s", dbUser)
	}

	if dbPass != "" {
		connStr += fmt.Sprintf(" password='%s'", dbPass)
	}

	if dbExtraOpts != "" {
		connStr += dbExtraOpts
	}

	return connStr
}
