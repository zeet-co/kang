package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Config struct {
	ZeetAPIKey    string
	ZeetTeamID    uuid.UUID
	ZeetGroupName string

	DBConnectionString string
}

const DefaultZeetGroupName = "kang"

func New(cmd *cobra.Command) (*Config, error) {

	var (
		apiKey    string
		teamID    uuid.UUID = uuid.Nil
		groupName string    = DefaultZeetGroupName
	)

	if envAPIKey := os.Getenv("ZEET_API_KEY"); envAPIKey != "" {
		apiKey = envAPIKey
	}

	if cliApiKey, _ := cmd.Flags().GetString("api-key"); cliApiKey != "" {
		apiKey = cliApiKey
	}

	if apiKey == "" {
		return nil, errors.New("Missing Zeet API Key. Set via env var ZEET_API_KEY or CLI --api-key")
	}

	teamID, err := getTeamID(cmd)

	if err != nil {
		return nil, err
	}

	if teamID == uuid.Nil {
		return nil, errors.New("Missing Zeet Team ID. Set via env var ZEET_TEAM_ID or CLI --team-id")
	}

	if envGroupName := os.Getenv("ZEET_GROUP_NAME"); envGroupName != "" {
		groupName = envGroupName
	}

	if cliGroupName, _ := cmd.Flags().GetString("group-name"); cliGroupName != "" {
		groupName = cliGroupName
	}

	return &Config{
		ZeetAPIKey:    apiKey,
		ZeetTeamID:    teamID,
		ZeetGroupName: groupName,

		DBConnectionString: getConnStr(),
	}, nil
}

func getTeamID(cmd *cobra.Command) (uuid.UUID, error) {

	var teamIDString string

	if envTeamID := os.Getenv("ZEET_TEAM_ID"); envTeamID != "" {
		teamIDString = envTeamID
	}

	if teamIDString == "" {
		teamID, err := cmd.Flags().GetString("team-id")
		if err != nil {
			return uuid.Nil, err
		}
		teamIDString = teamID
	}

	if teamIDString == "" {
		return uuid.Nil, nil
	}

	res, err := uuid.Parse(teamIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
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
