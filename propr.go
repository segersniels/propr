package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v61/github"
)

func printMarkdown(content string) error {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(content)
	if err != nil {
		return err
	}

	fmt.Print(out)
	return nil
}

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

	diff, err := getDiff(branch)
	if err != nil {
		return "", err
	}

	if diff == "" {
		return "", fmt.Errorf("no diff found")
	}

	var description string
	if err := spinner.New().TitleStyle(lipgloss.NewStyle()).Title("Generating your pull request...").Action(func() {
		if CONFIG.Data.Assistant.Enabled {
			if description, err = p.client.GetAssistantCompletion(diff); err != nil {
				log.Fatal(err)
			}
		} else {
			if description, err = p.client.GetChatCompletion(diff); err != nil {
				log.Fatal(err)
			}
		}
	}).Run(); err != nil {
		return "", err
	}

	return description, nil
}

func (p *Propr) Create(description string) error {
	info, err := getRepositoryInformation()
	if err != nil {
		return err
	}

	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(context.Background(), info.Owner, info.Name)
	if err != nil {
		return err
	}

	var title string
	err = huh.NewInput().Title("Provide a title for your pull request").Value(&title).Run()
	if err != nil {
		return nil
	}

	pr, response, err := client.PullRequests.Create(context.Background(), info.Owner, info.Name, &github.NewPullRequest{
		Base:  repository.DefaultBranch,
		Title: github.String(title),
		Body:  github.String(description),
	})

	if err != nil {
		return err
	} else if response.StatusCode != 201 {
		log.Fatal("failed to create pull request", "error", response.Status)
	}

	println("pull request created at %s", pr.URL)

	return nil
}
