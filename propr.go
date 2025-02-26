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

func NewMessageClient() MessageClient {
	var client MessageClient

	switch CONFIG.Data.Model {
	case Claude3Dot7Sonnet, Claude3Dot5Haiku, Claude3Dot5Sonnet:
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			log.Fatal("ANTHROPIC_API_KEY is not set")
		}

		client = NewAnthropic(apiKey, CONFIG.Data.Model)
	case DeepSeekChat, DeepSeekReasoner:
		apiKey := os.Getenv("DEEPSEEK_API_KEY")
		if apiKey == "" {
			log.Fatal("DEEPSEEK_API_KEY is not set")
		}

		client = NewDeepSeek(apiKey, CONFIG.Data.Model)
	default:
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY is not set")
		}

		client = NewOpenAI(apiKey, CONFIG.Data.Model)
	}

	log.Debug("Client initialized for", "model", CONFIG.Data.Model)

	return client
}

type GitHub struct {
	client *github.Client
	repo   *github.Repository
}

func NewGitHub() *GitHub {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	gh := github.NewClient(nil).WithAuthToken(token)
	info, err := getRepositoryInformation()
	if err != nil {
		log.Fatal(err)
	}

	repo, _, err := gh.Repositories.Get(context.Background(), info.Owner, info.Name)
	if err != nil {
		log.Fatal(err)
	}

	return &GitHub{
		gh,
		repo,
	}
}

type Propr struct{}

func NewPropr() *Propr {
	return &Propr{}
}

func (p *Propr) Generate(target string) (string, error) {
	gh := NewGitHub()
	if target == "" {
		target = gh.repo.GetDefaultBranch()
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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		commits, err := getCommitMessages(current, target)
		if err != nil {
			log.Fatal(err)
		}

		// We provide as much context as possible to get the best possible description
		messages := []Message{
			{
				Role:    MessageRoleUser,
				Content: gh.repo.GetURL(),
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

		log.Debug("Constructed messages to send to the provider", "messages", messages)

		client := NewMessageClient()
		systemMessage := generateSystemMessageForDiff(CONFIG.Data.Prompt, CONFIG.Data.Template)
		log.Debug("Constructed system instructions", "message", systemMessage)

		response, err := client.CreateMessage(ctx, systemMessage, messages)
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

func (p *Propr) Create(target string, description string, draft bool) error {
	gh := NewGitHub()
	if target == "" {
		target = gh.repo.GetDefaultBranch()
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

	owner := gh.repo.GetOwner().GetLogin()
	name := gh.repo.GetName()

	log.Debug("Creating pull request", "head", branch, "base", target, "owner", owner, "name", name)
	pr, response, err := gh.client.PullRequests.Create(context.Background(), owner, name, &github.NewPullRequest{
		Head:  github.String(branch),
		Base:  github.String(target),
		Title: github.String(title),
		Body:  github.String(description),
		Draft: github.Bool(draft),
	})

	if err != nil {
		return err
	} else if response.StatusCode != 201 {
		log.Fatal("Failed to create pull request", "error", response.Status)
	}

	fmt.Printf("Pull request created at %s\n", pr.GetHTMLURL())

	return nil
}
