package coppermind

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *github.Client
}

func NewClient() *GithubClient {
	GITHUB_STAR_KEY := os.Getenv("GITHUB_STAR_KEY")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_STAR_KEY},
	)

	return &GithubClient{
		client: github.NewClient(oauth2.NewClient(ctx, ts)),
	}
}

func (gh *GithubClient) ListStarred() ([]*github.StarredRepository, error) {
	repos, _, err := gh.client.Activity.ListStarred(context.Background(), "", nil)

	if err != nil {
		return nil, err
	}

	return repos, nil
}

type StarredRepository struct {
	FullName    string
	Description string
}

func StoreGithubStars(db *CoppermindDb) error {
	client := NewClient()
	starred, err := client.ListStarred()
	if err != nil {
		return err
	}

	for _, repo := range starred {
		data := StarredRepository{}

		name := repo.Repository.FullName
		description := repo.Repository.Description

		if name != nil {
			data.FullName = *name
		}

		if description != nil {
			data.Description = *description
		}

	}

	return nil
}
