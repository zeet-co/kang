package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			TeamID:     cfg.ZeetTeamID,
			Overrides:  overrides,
			EnvName:    envName,
			ProjectIDs: validIDs,
		})
	},
}

var stopEnvironmentCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all running instances of each Project in the Environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
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

		return kang.StopEnvironment(ctx, envName, cfg.ZeetTeamID)
	},
}

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Comment high-level information on a given environment into a given Github Pull Request",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		cfg, err := config.New(cmd)
		if err != nil {
			return err
		}

		kang, err := controller.NewController(cfg)
		if err != nil {
			return err
		}

		repo, _ := cmd.Flags().GetString("repo")
		pr, _ := cmd.Flags().GetInt("pr")
		token, _ := cmd.Flags().GetString("token")
		envName, _ := cmd.Flags().GetString("env-name")

		return kang.CommentGithub(ctx, pr, repo, token, envName)
	},
}

func parseOverrides(stmts []string) map[uuid.UUID]map[string]string {
	// Each stmt is expected to be of format uuid:field:value,stmt..
	// TODO improve this format to be compatible with lists and dicts which may include the `,` character

	output := make(map[uuid.UUID]map[string]string)

	for _, stmt := range stmts {
		splitStmt := strings.SplitN(stmt, ":", 3)
		if len(splitStmt) == 3 || len(splitStmt) == 4 {
			if id, err := uuid.Parse(splitStmt[0]); err == nil {
				if output[id] == nil {
					output[id] = make(map[string]string)
				}
				output[id][splitStmt[1]] = splitStmt[2]
			}
		} else {
			fmt.Println("Invalid override; must be of format id:key:value. offending stmt:", stmt)
		}
	}
	return output
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "overrides":
		name = "override"
		break
	}
	return pflag.NormalizedName(name)
}

func init() {
	stopEnvironmentCmd.Flags().String("name", "", "Specify the name of the environment you'd like to stop")
	stopEnvironmentCmd.MarkFlagRequired("name")

	startEnvironmentCmd.Flags().String("name", "", "Specify the name of the environment you'd like to create an instance of")
	startEnvironmentCmd.MarkFlagRequired("name")

	startEnvironmentCmd.Flags().StringSlice("ids", []string{}, "Specify a comma-seperated list of Zeet Project IDs")
	startEnvironmentCmd.MarkFlagRequired("ids")

	startEnvironmentCmd.Flags().StringSlice("overrides", []string{}, "Specify the Project ID : field : value combos that you would like to override. Format: project_id:field:value,proj.. Example: 1c6ea878-f92e-435e-a849-7bccfe7c6e5a:branch:feature-1")
	startEnvironmentCmd.Flags().SetNormalizeFunc(aliasNormalizeFunc)

	commentCmd.Flags().String("repo", "", "Github Repo that will be commented on")
	commentCmd.MarkFlagRequired("repo")
	commentCmd.Flags().Int("pr", 0, "Specify the PR number that should have a comment added to it")
	commentCmd.MarkFlagRequired("pr")
	commentCmd.Flags().String("token", "", "Github Token that will has permission to comment on the specified Github repo & PR")
	commentCmd.MarkFlagRequired("token")

	commentCmd.Flags().String("env-name", "", "Specify the name of the environment you'd like to create the comment from")
	commentCmd.MarkFlagRequired("env-name")

}
