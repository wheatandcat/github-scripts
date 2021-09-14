package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/machinebox/graphql"
)

// Config 設定ファイルのタイプ
type Config struct {
	GitHub GitHubConfig
}

// GitHubConfig GitHub設定ファイルのタイプ
type GitHubConfig struct {
	Token          string `toml:"token"`
	Owner          string `toml:"owner"`
	RepositoryName string `toml:"repositoryName"`
	IgnoreLabel    string `toml:"ignoreLabel"`
}

// ResponseData GASのResponse
type ResponseData struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Image        string   `json:"image"`
	Env          string   `json:"env"`
	Priority     string   `json:"priority"`
	Repositories []string `json:"repositories"`
}

// ResposeType Graphqlのタイプ
type ResposeType struct {
	Repository Repository `json:"repository"`
}

// Repository Repositoryのタイプ
type Repository struct {
	ID           string           `json:"id"`
	PullRequests PullRequestNodes `json:"pullRequests"`
}

type PullRequestNodes struct {
	Nodes []PullRequest `json:"nodes"`
}

type PullRequest struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	CreatedAt time.Time  `json:"createdAt"`
	Labels    LabelNodes `json:"labels"`
}

// LabelNodes ラベルリスト
type LabelNodes struct {
	Nodes []Label `json:"nodes"`
}

// Label ラベル
type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdatePullRequestInput struct {
	PullRequestId string `json:"pullRequestId"`
	State         string `json:"state"`
}

func include(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// 90日
const period = time.Hour * 24 * 30 * 3

func main() {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	is, err := config.GitHub.getPullRequests()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range is {
		input := UpdatePullRequestInput{
			PullRequestId: i.ID,
			State:         "CLOSED",
		}
		if err := config.GitHub.closeIssue(input); err != nil {
			fmt.Println("削除失敗:" + i.Title)
			log.Fatal(err)
		}

		fmt.Println("削除issue: " + i.Title)
	}
}

func (c *GitHubConfig) getPullRequests() ([]PullRequest, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    pullRequests(first: 100, states: [OPEN], orderBy: {field: CREATED_AT, direction: ASC}) {
      nodes {
        id
        title
        url
        createdAt
        labels(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
          nodes {
            id
            name
          }
        }
      }
    }
  }
}
`)

	req.Var("owner", c.Owner)
	req.Var("name", c.RepositoryName)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	prs := []PullRequest{}

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}

	for _, i := range respData.Repository.PullRequests.Nodes {
		if (time.Now().Add(-1 * period)).Before(i.CreatedAt) {
			// 指定のラベルは除外
			continue
		}

		labels := []string{}

		for _, l := range i.Labels.Nodes {
			labels = append(labels, l.Name)
		}

		if include(labels, c.IgnoreLabel) {
			// 指定のラベルは除外
			continue
		}

		prs = append(prs, i)
	}

	return prs, nil
}

func (c *GitHubConfig) closeIssue(input UpdatePullRequestInput) error {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
mutation UpdateRequestInput($input: UpdatePullRequestInput!) {
  updatePullRequest(input: $input) {
    pullRequest {
      id
    }
  }
}
`)

	req.Var("input", input)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		return err
	}

	return nil
}
