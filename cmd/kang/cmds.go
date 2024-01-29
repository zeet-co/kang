package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zeet-co/kang/internal/config"

	"github.com/zeet-co/kang/internal/controller"
)

var startEnvironmentCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance of each Project in the Environment, using the specified Brnach overrides for any given Projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New(cmd)
		if err != nil {
			return err
		}

		kang, err := controller.NewController(cfg)
		if err != nil {
			return err
		}

		overridesInput, _ := cmd.Flags().GetStringSlice("overrides")
		projectIDs, _ := cmd.Flags().GetStringSlice("ids")
		envName, _ := cmd.Flags().GetString("name")

		dedupedIDs := map[uuid.UUID]interface{}{}

		for _, id := range projectIDs {
			if pID, err := uuid.Parse(id); err == nil && kang.CheckProjectExists(pID) {
				dedupedIDs[pID] = struct{}{}
			}
		}

		validIDs := []uuid.UUID{}
		for k := range dedupedIDs {
			validIDs = append(validIDs, k)
		}

		if len(validIDs) < 2 {
			return errors.New("Must specify at least 2 unique valid Project IDs")
		}

		overrides := parseOverrides(overridesInput)

		return kang.StartEnvironment(controller.StartEnvironmentOpts{
			TeamID:                 cfg.ZeetTeamID,
			ProjectBranchOverrides: overrides,
			EnvName:                envName,
			ProjectIDs:             validIDs,
		})
	},
}

var stopEnvironmentCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all running instances of each Project in the Environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New(cmd)
		if err != nil {
			return err
		}

		kang, err := controller.NewController(cfg)
		if err != nil {
			return err
		}

		envName, err := cmd.Flags().GetString("name")

		if err != nil {
			return err
		}

		return kang.StopEnvironment(envName)
	},
}

func parseOverrides(pairs []string) map[uuid.UUID]string {

	output := make(map[uuid.UUID]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			if id, err := uuid.Parse(kv[0]); err == nil {
				output[id] = kv[1]
			}
		} else {
			fmt.Println("Invalid key-value pair:", pair)
		}
	}
	return output
}

func init() {
	stopEnvironmentCmd.Flags().String("name", "", "Specify the name of the environment you'd like to stop")
	stopEnvironmentCmd.MarkFlagRequired("name")

	startEnvironmentCmd.Flags().String("name", "", "Specify the name of the environment you'd like to create an instance of")
	startEnvironmentCmd.MarkFlagRequired("name")

	startEnvironmentCmd.Flags().StringSlice("ids", []string{}, "Specify a comma-seperated list of Zeet Project IDs")
	startEnvironmentCmd.MarkFlagRequired("ids")

	startEnvironmentCmd.Flags().StringSlice("overrides", []string{}, "Specify the Project ID : Branch combos that you would like to override from the normal Production Branch of each Project. Format: project_id:branch,proj..")
}
