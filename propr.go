package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v61/github"
)

type Propr struct {
	client MessageClient
	gh     *github.Client
	repo   *github.Repository
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

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	info, err := getRepositoryInformation()
	if err != nil {
		log.Fatal(err)
	}

	gh := github.NewClient(nil).WithAuthToken(token)
	repo, _, err := gh.Repositories.Get(context.Background(), info.Owner, info.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &Propr{
		client,
		gh,
		repo,
	}
}

func (p *Propr) Generate(target string) (string, error) {
	if target == "" {
		target = p.repo.GetDefaultBranch()
	}

	current, err := getCurrentBranch()
	if err != nil {
		return "", err
	}

	log.Debug("Fetching diff", "target", target)
	diff, err := getDiff(current, target)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("not enough changes found to generate")
	}

	var description string
	err = spinner.New().TitleStyle(lipgloss.NewStyle()).Title("Generating your pull request...").Action(func() {
		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		commits, err := getCommitMessages(current, target)
		if err != nil {
			log.Fatal(err)
		}

		// We provide as much context as possible to get the best possible description
		messages := []Message{
			{
				Role:    MessageRoleUser,
				Content: p.repo.GetURL(),
			},
			{
				Role:    MessageRoleAssistant,
				Content: "Thanks for providing the repository URL. What about the branch?",
			},
			{
				Role:    MessageRoleUser,
				Content: current,
			},
			{
				Role:    MessageRoleAssistant,
				Content: "Thanks for providing the branch. What about the commit messages?",
			},
			{
				Role:    MessageRoleUser,
				Content: strings.Join(commits, "\n"),
			},
			{
				Role:    MessageRoleAssistant,
				Content: "Thanks for providing the commit messages. Now the final step to generate a description is to see what's changed using the diff",
			},
			{
				Role:    MessageRoleUser,
				Content: prepareDiff(diff),
			},
		}

		response, err := p.client.CreateMessage(ctx, generateSystemMessageForDiff(CONFIG.Data.Prompt, CONFIG.Data.Template), messages)
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

func (p *Propr) Create(target string, description string) error {
	if target == "" {
		target = p.repo.GetDefaultBranch()
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	var title string
	err := huh.NewInput().Title("Provide a title for your pull request").Value(&title).Run()
	if err != nil {
		return nil
	}

	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}

	owner := p.repo.GetOwner().GetLogin()
	name := p.repo.GetName()

	log.Debug("Creating pull request", "head", branch, "base", target, "owner", owner, "name", name)
	pr, response, err := p.gh.PullRequests.Create(context.Background(), owner, name, &github.NewPullRequest{
		Head:  github.String(branch),
		Base:  github.String(target),
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
