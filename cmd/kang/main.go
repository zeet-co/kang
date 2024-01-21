package main

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "kang", Short: "CLI to manage multi-Project PRs with Zeet"}

func init() {
	rootCmd.AddCommand(createEnvCmd)
	rootCmd.AddCommand(destroyEnvCmd)
	rootCmd.AddCommand(startEnvironmentCmd)
	rootCmd.AddCommand(stopEnvironmentCmd)
}

func main() {

	rootCmd.PersistentFlags().String("api-key", "", "Input your Zeet API Key. For more information see https://docs.zeet.co/graphql/")

	rootCmd.PersistentFlags().String("team-id", "", "Input your Zeet Team ID. For more information see https://docs.zeet.co/graphql/")

	rand.Seed(time.Now().UnixNano())
	cobra.CheckErr(rootCmd.Execute())
}
