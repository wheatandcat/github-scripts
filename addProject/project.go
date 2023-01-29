package addproject

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v29/github"
	"github.com/machinebox/graphql"
)

type GitHubClient struct {
	Token   string
	GraphQL *graphql.Client
}

type AddProjectV2ItemByIdResult struct {
	AddProjectV2ItemById AddProjectV2ItemById `json:"addProjectV2ItemById"`
}

type AddProjectV2ItemById struct {
	Item AddProjectV2Item `json:"item"`
}

type AddProjectV2Item struct {
	ID string `json:"id"`
}

func processIssuesEvent(ctx context.Context, event *github.IssuesEvent) error {
	log.Println("GetAction:", event.GetAction())
	if event.GetAction() != "opened" {
		return nil
	}

	checkRepo := []string{"github-scripts"}
	if !contains(checkRepo, event.Repo.GetName()) {
		return nil
	}

	token := os.Getenv("ORGANIZATIONS_GITHUB")

	log.Println("token: ", token)

	c := GitHubClient{
		Token:   token,
		GraphQL: graphql.NewClient("https://api.github.com/graphql"),
	}

	item, err := c.addProject(ctx, event.Issue)
	if err != nil {
		return err
	}

	log.Println("item: ", item.AddProjectV2ItemById.Item.ID)

	return nil
}

func (c *GitHubClient) addProject(ctx context.Context, issue *github.Issue) (AddProjectV2ItemByIdResult, error) {
	req := graphql.NewRequest(`
		mutation ($projectID: ID!, $nodeID: ID!) {
			addProjectV2ItemById(input: {projectId: $projectID, contentId: $nodeID}) {
				item {
				id
				}
			}
		}
	`)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	req.Var("projectID", os.Getenv("GITHUB_PROJECT_ID"))
	req.Var("nodeID", issue.GetNodeID())

	var res AddProjectV2ItemByIdResult
	if err := c.GraphQL.Run(ctx, req, &res); err != nil {
		return res, err
	}

	return res, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
