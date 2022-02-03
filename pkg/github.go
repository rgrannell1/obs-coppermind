package coppermind

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type StarredRepository struct {
	FullName    string
	Description string
	Login       string
	Url         string
	Language    string
	Topics      []string
}

type GithubClient struct {
	client *github.Client
}

/*
 * Construct a github client
 *
 */
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

func CountGithubStars(client *GithubClient) (int, error) {
	_, res, err := client.client.Activity.ListStarred(context.Background(), "", &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})

	if err != nil {
		return 0, err
	}

	return res.LastPage, nil
}

/*
 * Store github stars in the diatom database
 *
 */
func StoreGithubStars(db *CoppermindDb) error {
	client := NewClient()
	stars, err := CountGithubStars(client)
	if err != nil {
		return err
	}

	changed, err := db.GithubChanged(stars)
	if err != nil {
		return err
	}

	if !changed {
		return nil
	}

	page := 0

	//enumerate through pages
	for {
		page++

		// list results from this page
		starred, _, err := client.client.Activity.ListStarred(context.Background(), "", &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				PerPage: 10,
				Page:    page,
			},
		})
		if err != nil {
			return err
		}

		// break when all are enumerated
		if len(starred) == 0 {
			break
		}

		// store each repository
		for _, repo := range starred {
			data := StarredRepository{}

			name := repo.Repository.FullName
			description := repo.Repository.Description
			data.Login = *repo.Repository.Owner.Login
			data.Url = *repo.Repository.HTMLURL
			data.Topics = repo.Repository.Topics

			language := repo.Repository.Language
			if language != nil {
				data.Language = *language
			}

			if name != nil {
				data.FullName = *name
			}

			if description != nil {
				data.Description = *description
			}

			// -- TODO

			err := db.InsertStar(&data)
			if err != nil {
				return err
			}
		}
	}

	err = db.UpdateStarCount(stars)
	if err != nil {
		return err
	}

	return nil
}
