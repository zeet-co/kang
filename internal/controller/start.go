package controller

import (
	"context"
	"fmt"
	"strings"

	stdErrors "errors"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeet-co/kang/internal/parser"
	v0 "github.com/zeet-co/kang/internal/zeet/v0"
)

type StartEnvironmentOpts struct {
	Overrides  map[uuid.UUID]map[string]string
	EnvName    string
	ProjectIDs []uuid.UUID
	TeamID     uuid.UUID
}

type envOverride struct {
	key        string
	value      string
	isSymbolic bool
}

func (c *Controller) StartEnvironment(opts StartEnvironmentOpts) error {
	ctx := context.Background()

	teamName, err := c.zeet.GetTeamName(ctx, opts.TeamID)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not fetch team"))
	}

	group := c.groupName
	subGroup := opts.EnvName

	// Handle Group/ Sub Group creation
	groupID, subGroupID, err := c.zeet.EnsureGroupsExist(ctx, *teamName, group, subGroup, opts.TeamID)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not ensure group / subgroup"))
	}

	// Get all the projects that need to be duplicated
	projects, err := c.zeet.GetProjectsByID(ctx, opts.ProjectIDs)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "could not fetch project information"))
	}

	projectSuccessorMap := map[uuid.UUID]uuid.UUID{}

	errs := make([]error, len(projects))

	for i, p := range projects {
		// Duplicate each of the projects passed in via --ids
		//TODO scale down resources on branch deployments (?)
		//TODO handle database linking

		newName := fmt.Sprintf("%s-%s", p.Name, opts.EnvName)
		newProjectID, err := c.zeet.DuplicateProject(ctx, p.ID, groupID, subGroupID, newName)
		if err != nil {
			if err != v0.AlreadyExistsError {
				errs[i] = errors.WithStack(errors.Wrap(err, "could not duplicate project"))
				continue
			}
			fmt.Printf("Could not duplicate project %s, as it already exists. Checking for overrides..\n", newName)
			newProjectID, err = c.zeet.GetRepoByName(ctx, newName)
			if err != nil {
				errs[i] = errors.WithStack(errors.Wrap(err, "could not fetch project"))
				continue
			}
		}
		projectSuccessorMap[p.ID] = newProjectID
	}

	for i, p := range projects {
		newName := fmt.Sprintf("%s-%s", p.Name, opts.EnvName)
		if opts.Overrides[p.ID] != nil {
			fmt.Printf("Found overrides applying to %s: parsing now\n", newName)
			override := opts.Overrides[p.ID]
			errs[i] = c.applyOverrides(ctx, opts.TeamID, projectSuccessorMap, p.ID, override)
		}

		if errs[i] == nil {
			fmt.Printf("Done with project %s!\n", newName)
		}
	}

	return stdErrors.Join(errs...)
}

func overrideToUpdateInput(pID uuid.UUID, overrides map[string]string) (*v0.UpdateProjectInput, error) {
	out := v0.UpdateProjectInput{
		Id: pID,
	}

	err, anyFieldSet := parser.AssignValues(&out, overrides)

	if err != nil {
		return nil, err
	}

	if !anyFieldSet {
		return nil, nil
	}

	return &out, nil
}

func (c *Controller) applyOverrides(ctx context.Context, teamID uuid.UUID, projectSuccessorMap map[uuid.UUID]uuid.UUID, oldProjectID uuid.UUID, override map[string]string) error {
	var changed bool
	newProjectID := projectSuccessorMap[oldProjectID]

	updateInput, err := overrideToUpdateInput(newProjectID, override)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not parse override %s", override)))
	}

	if updateInput != nil {
		fmt.Printf("Applying override stmt %s to %s\n", override, newProjectID)

		if err = c.zeet.UpdateProject(ctx, newProjectID, updateInput); err != nil {
			return errors.WithStack(errors.Wrap(err, "could not apply config overrides"))
		}
		changed = true
	}

	envs, envPresent := checkOverridesForEnvs(override)
	if envPresent {
		envsToSet := map[string]string{}
		symbolicEnvs := map[string]string{}

		for _, e := range envs {
			if e.isSymbolic {
				symbolicEnvs[e.key] = e.value
			} else {
				envsToSet[e.key] = e.value
			}
		}

		if err = c.addSymbolicEnvs(ctx, envsToSet, symbolicEnvs, projectSuccessorMap); err != nil {
			fmt.Printf("Failed to resolve references to other projects' outputs: %s\n", err)
		}

		if len(envsToSet) > 0 {
			fmt.Printf("Applying env overrides to %s: %v\n", newProjectID, envsToSet)
			finalEnvs := make(map[string]string)
			if existingRepo, err := c.zeet.GetRepoByID(ctx, oldProjectID); err != nil {
				fmt.Printf("Failed to fetch existing env vars for project %s, non-overriden env vars may be dropped: %s\n", oldProjectID, err)
			} else {
				//Because we're setting all the envs, we need to ensure we keep any existing env vars or else the override will wipe all non-override values
				finalEnvs = existingRepo.Envs
			}

			for k := range envsToSet {
				finalEnvs[k] = envsToSet[k]
			}

			if err = c.zeet.UpdateEnvs(ctx, newProjectID, finalEnvs); err != nil {
				return errors.WithStack(errors.Wrap(err, "could not apply env var override"))
			}
			changed = true
		}
	}

	clusterIDs, clusterPresent := checkOverridesForClusters(override)
	if clusterPresent {
		fmt.Printf("An override is trying to deploy this project to cluster %s\n", clusterIDs)
		clusters := make(map[uuid.UUID]*v0.Cluster)
		for _, clusterID := range clusterIDs {
			cluster, err := c.zeet.GetClusterByID(ctx, clusterID, teamID)
			if err != nil {
				//TODO handle 404
				return errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not fetch cluster %s", clusterID)))
			}
			clusters[cluster.ID] = cluster
		}

		newReplication := []v0.ReplicationInput{}
		for _, cluster := range clusters {
			newReplication = append(newReplication, v0.ReplicationInput{
				Region:    cluster.Region,
				Replicas:  1,
				ClusterID: &cluster.ID,
			})
		}

		updateObject := &v0.UpdateProjectInput{
			Id:          newProjectID,
			Replication: newReplication,
		}

		if err = c.zeet.UpdateProject(ctx, newProjectID, updateObject); err != nil {
			return errors.WithStack(errors.Wrap(err, "could not apply cluster overrides"))
		}
		changed = true
	}

	if changed {
		fmt.Printf("Project %s had config changes that may not be reflected in most recent deployments. Triggering re-build with cache\n", newProjectID)
		if err = c.zeet.RebuildProject(ctx, newProjectID); err != nil {
			return errors.WithStack(errors.Wrap(err, "could not trigger rebuild"))
		}
	}

	return nil
}

func checkOverridesForClusters(override map[string]string) ([]uuid.UUID, bool) {
	//TODO support cluster names instead of IDs as input
	out := []uuid.UUID{}
	isClusterPresent := false

	for k, v := range override {
		if k == "cluster" {

			clusterID, err := uuid.Parse(v)
			if err == nil {
				out = append(out, clusterID)

				isClusterPresent = true
			}
		}
	}

	return out, isClusterPresent
}

func checkOverridesForEnvs(override map[string]string) ([]envOverride, bool) {
	out := []envOverride{}
	isEnvPresent := false

	for k, v := range override {
		if strings.HasPrefix(k, "env.") {
			envKey := k[4:]
			out = append(out, envOverride{
				key:        envKey,
				value:      v,
				isSymbolic: isSymbolic(v),
			})

			isEnvPresent = true
		}
	}

	return out, isEnvPresent
}

// isSymbolic checks if the format of the env
func isSymbolic(v string) bool {
	split := strings.Split(v, ":")

	if len(split) == 3 {
		return false
	}

	if _, err := uuid.Parse(split[0]); err != nil {
		return false
	}

	return true

}

func (c *Controller) addSymbolicEnvs(ctx context.Context, out, symbolicEnvs map[string]string, projectSuccessorMap map[uuid.UUID]uuid.UUID) error {
	projectIDs := []uuid.UUID{}
	keyToProjectIDAndValue := map[string][]interface{}{}

	for k, v := range symbolicEnvs {
		s := strings.Split(v, ":")
		projectID, err := uuid.Parse(s[0])
		if err != nil {
			return err
		}

		if modernProjectID := projectSuccessorMap[projectID]; modernProjectID != uuid.Nil {
			projectID = modernProjectID
		}
		keyToProjectIDAndValue[k] = []interface{}{projectID, s[1]}
		projectIDs = append(projectIDs, projectID)
	}

	projects, err := c.zeet.GetProjectsByID(ctx, projectIDs)
	if err != nil {
		return err
	}

	projectsByID := make(map[uuid.UUID]*v0.Repo, len(projects))
	for _, p := range projects {
		projectsByID[p.ID] = p
	}

	for k := range symbolicEnvs {

		pID := keyToProjectIDAndValue[k][0].(uuid.UUID)
		field := keyToProjectIDAndValue[k][1].(string)
		p := projectsByID[pID]

		foundValue := parser.GetValue(*p, field)
		out[k] = foundValue
	}

	return nil
}
