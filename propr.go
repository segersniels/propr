package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v61/github"
)

type Propr struct {
	client MessageClient
}

func NewPropr() *Propr {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API is not set")
	}

	// Depending on the user selected model, we need to set the corresponding API key
	var client MessageClient
	switch CONFIG.Data.Model {
	case Claude3Dot5Sonnet:
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			log.Fatal("ANTHROPIC_API_KEY is not set")
		}

		client = NewAnthropic(apiKey, CONFIG.Data.Model)
	default:
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY is not set")
		}

		client = NewOpenAI(apiKey, CONFIG.Data.Model)
	}

	return &Propr{
		client,
	}
}

func (p *Propr) Generate(branch string) (string, error) {
	if branch == "" {
		head, err := getDefaultBranch()
		if err != nil {
			return "", err
		}

		branch = head
	}

	log.Debug("Fetching diff", "base", branch)
	diff, err := getDiff(branch)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("not enough changes found to generate")
	}

	var description string
	err = spinner.New().TitleStyle(lipgloss.NewStyle()).Title("Generating your pull request...").Action(func() {
		log.Debug("Using chat completion")

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := p.client.CreateMessage(ctx, generateSystemMessageForDiff(CONFIG.Data.Prompt, CONFIG.Data.Template), prepareDiff(diff))
		if err != nil {
			log.Fatal(err)
		}

		description = response
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

	fmt.Printf("Pull request created at %s\n", pr.GetHTMLURL())

	return nil
}
