package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (c *Controller) CommentGithub(ctx context.Context, prNumber int, repo, token, envName string) error {

	projectIDs, err := c.getProjectsInSubGroup(ctx, ZeetGroupName, envName)
	if err != nil {
		return err
	}

	client, err := newGitHubAPIClient(ctx, token, "https://api.github.com", nil)

	if err != nil {
		return err
	}

	owner, repo, err := splitGitHubProject(repo)
	if err != nil {
		return err
	}

	body, err := c.genGithubCommentBody(ctx, projectIDs)
	if err != nil {
		return err
	}

	comment, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{Body: github.String(body)})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully commented: %s\n", comment.GetHTMLURL())

	return nil
}

func newGitHubAPIClient(ctx context.Context, token string, apiURL string, tlsConfig *tls.Config) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = tlsConfig
	client := &http.Client{Transport: transport}
	httpCtx := context.WithValue(ctx, oauth2.HTTPClient, client)

	tc := oauth2.NewClient(httpCtx, ts)

	// Handle default GitHub API client
	if apiURL == "" || apiURL == "https://api.github.com" {
		return github.NewClient(tc), nil
	}

	// Handle GitHub Enterprise API client

	// GitHub Enterprise v3 client needs a base URL and upload URL
	// So we need to parse the API URL and add the necessary parts
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing API URL")
	}

	// Add trailing slash
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	// Check if "/api/v3/" already exists in the URL path.
	if strings.HasSuffix(u.Path, "/api/v3/") {
		return nil, fmt.Errorf("the /api/v3/ suffix should not be included in the --github-api-url")
	}

	// Add api to path if it doesn't exist
	if !strings.HasSuffix(u.Path, "/api/") {
		u.Path += "api/"
	}

	apiURL = u.String()

	v3client, err := github.NewEnterpriseClient(apiURL+"v3/", apiURL+"uploads/", tc)
	if err != nil {
		return nil, err
	}

	return v3client, nil
}

// splitGitHubProject parses a GitHub project string into its owner and repo parts.
func splitGitHubProject(project string) (string, string, error) {
	parts := strings.SplitN(project, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid GitHub repository name: %s, expecting owner/repo", project)
	}
	return parts[0], parts[1], nil
}

func (c *Controller) genGithubCommentBody(ctx context.Context, projectIDs []uuid.UUID) (string, error) {

	projects, err := c.zeet.GetProjectsByID(ctx, projectIDs)
	if err != nil {
		return "", err
	}

	body := fmt.Sprintf(`
We've created a new ephemeral environment for this pull request over on Zeet!
We're deploying %d projects into your new environment. The table below includes the Dashboard page where you can modify settings, as well as the Preview Endpoint you can use to test before merging:
| Project Name  | Dashboard Link | Preview Endpoint |
| ------------- | ------------- | ------------- |
`, len(projects))

	for _, p := range projects {
		projectName := p.Name

		dashboardLink := fmt.Sprintf("https://zeet.co/%s/%s/%s/%s/deployments/%s", p.Owner, p.GroupName, p.SubGroupName, p.Name, p.ProductionDeployment.ID)
		previewLink := "https://" + p.ProductionDeployment.Endpoints[0]

		line := fmt.Sprintf("| %s | [Go to Project in Zeet Dashboard](%s) | [Preview Project](%s) | \n", projectName, dashboardLink, previewLink)
		body += line
	}

	return body, nil
}
