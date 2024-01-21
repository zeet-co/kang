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

var createEnvCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Environment by mapping existing Zeet Projects together",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.New(cmd)
		if err != nil {
			return err
		}

		kang, err := controller.NewController(cfg)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		projectIDs, _ := cmd.Flags().GetStringSlice("ids")

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
			return errors.New("Must specify at least 2 unique valid UUIDs")
		}

		return kang.CreateEnvironment(controller.CreateEnvironmentOptions{
			Name:       name,
			ProjectIDs: validIDs,
		})
	},
}

var destroyEnvCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a previously created Environment, preventing future instances from spawning",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New(cmd)
		if err != nil {
			return err
		}

		kang, err := controller.NewController(cfg)
		if err != nil {
			return err
		}

		inputID, _ := cmd.Flags().GetString("id")

		envID, err := uuid.Parse(inputID)
		if err != nil {
			return err
		}

		return kang.DestroyEnvironment(envID)
	},
}

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

		inputID, _ := cmd.Flags().GetString("id")
		overridesInput, _ := cmd.Flags().GetStringSlice("overrides")

		envID, err := uuid.Parse(inputID)
		if err != nil {
			return err
		}

		overrides := parseOverrides(overridesInput)

		return kang.StartEnvironment(envID, cfg.ZeetTeamID, controller.StartEnvironmentOpts{
			ProjectBranchOverrides: overrides,
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

		inputID, _ := cmd.Flags().GetString("id")

		envID, err := uuid.Parse(inputID)
		if err != nil {
			return err
		}

		return kang.StopEnvironment(envID)
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
	createEnvCmd.Flags().String("name", "", "Specify a name for your new environment")
	createEnvCmd.MarkFlagRequired("name")
	createEnvCmd.Flags().StringSlice("ids", []string{}, "Specify a comma-seperated list of Zeet Project IDs")
	createEnvCmd.MarkFlagRequired("ids")

	//TODO support name instead of ID for destroy, stop, start
	destroyEnvCmd.Flags().String("id", "", "Specify the ID of the environment you'd like to destroy")
	destroyEnvCmd.MarkFlagRequired("id")

	stopEnvironmentCmd.Flags().String("id", "", "Specify the ID of the environment you'd like to stop")
	stopEnvironmentCmd.MarkFlagRequired("id")

	startEnvironmentCmd.Flags().String("id", "", "Specify the ID of the environment you'd like to create an instance of")
	startEnvironmentCmd.MarkFlagRequired("id")
	startEnvironmentCmd.Flags().StringSlice("overrides", []string{}, "Specify the Project ID : Branch combos that you would like to override from the normal Production Branch of each Project. Format: project_id:branch,proj..")
}
