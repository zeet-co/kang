package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "kang", Short: "CLI to manage multi-Project PRs with Zeet"}

func init() {
	rootCmd.AddCommand(createEnvCmd)
	rootCmd.AddCommand(destroyEnvCmd)
	rootCmd.AddCommand(startInstance)
}

func main() {

	rootCmd.PersistentFlags().String("api-key", "", "Input your Zeet API Key. For more information see https://docs.zeet.co/graphql/")
	rootCmd.MarkPersistentFlagRequired("api-key")

	rootCmd.PersistentFlags().String("team-id", "", "Input your Zeet Team ID. For more information see https://docs.zeet.co/graphql/")
	rootCmd.MarkPersistentFlagRequired("team-id")

	rand.Seed(time.Now().UnixNano())
	cobra.CheckErr(rootCmd.Execute())
}

func getTeamID(cmd *cobra.Command) (uuid.UUID, error) {
	teamID, err := cmd.Flags().GetString("team-id")
	if err != nil {
		return uuid.Nil, err
	}

	res, err := uuid.Parse(teamID)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}
