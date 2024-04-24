package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v61/github"
)

type Propr struct {
	client OpenAI
}

func NewPropr() *Propr {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API is not set")
	}

	return &Propr{
		client: *NewOpenAI(apiKey),
	}
}

func (p *Propr) Generate() (string, error) {
	branch, err := getDefaultBranch()
	if err != nil {
		return "", err
	}

	log.Debug("Fetching diff", "branch", branch)
	diff, err := getDiff(branch)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("no diff found")
	}

	var description string
	err = spinner.New().TitleStyle(lipgloss.NewStyle()).Title("Generating your pull request...").Action(func() {
		if CONFIG.Data.Assistant.Enabled && CONFIG.Data.Assistant.Id != "" {
			log.Debug("Using assistant completion")
			response, err := p.client.GetAssistantCompletion(diff)
			if err != nil {
				log.Fatal(err)
			}

			description = response
		} else {
			log.Debug("Using chat completion")
			response, err := p.client.GetChatCompletion(diff)
			if err != nil {
				log.Fatal(err)
			}

			description = response
		}
	}).Run()

	if err != nil {
		return "", err
	}

	return description, nil
}

func (p *Propr) Create(description string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	info, err := getRepositoryInformation()
	if err != nil {
		return err
	}

	client := github.NewClient(nil).WithAuthToken(token)
	repository, _, err := client.Repositories.Get(context.Background(), info.Owner, info.Name)
	if err != nil {
		return err
	}

	var title string
	err = huh.NewInput().Title("Provide a title for your pull request").Value(&title).Run()
	if err != nil {
		return nil
	}

	base := repository.GetDefaultBranch()
	head, err := getCurrentBranch()
	if err != nil {
		return err
	}

	log.Debug("Creating pull request", "head", head, "base", base)
	pr, response, err := client.PullRequests.Create(context.Background(), info.Owner, info.Name, &github.NewPullRequest{
		Head:  github.String(head),
		Base:  github.String(base),
		Title: github.String(title),
		Body:  github.String(description),
	})

	if err != nil {
		return err
	} else if response.StatusCode != 201 {
		log.Fatal("Failed to create pull request", "error", response.Status)
	}

	fmt.Printf("pull request created at %s\n", pr.GetHTMLURL())

	return nil
}
