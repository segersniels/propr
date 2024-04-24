package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/sashabaranov/go-openai"
	"github.com/segersniels/propr/config"
	"github.com/urfave/cli/v2"
)

var AppVersion string
var AppName string

type AssistantConfig struct {
	Enabled bool   `json:"enabled"`
	Id      string `json:"id"`
}

type Config struct {
	Model     string          `json:"model"`
	Prompt    string          `json:"prompt"`
	Template  string          `json:"template"`
	Assistant AssistantConfig `json:"assistant"`
}

var CONFIG = config.NewConfig("propr", Config{
	Model: openai.GPT4TurboPreview,
	Prompt: `You will be asked to write a concise GitHub PR description based on a provided git diff.
Analyze the code changes and provide a concise explanation of the changes, their context and why they were made.
Don't reference file names or directories directly, instead give a general explanation of the changes made.
Do not treat imports and requires as changes or new features. If the provided message is not a diff respond with an appropriate message.
Don't surround your description in backticks but still write GitHub supported markdown.`,
	Template: "# Description",
	Assistant: AssistantConfig{
		Enabled: false,
		Id:      "",
	},
})

func main() {
	propr := NewPropr()
	app := &cli.App{
		Name:    AppName,
		Usage:   "Generate your PRs from the command line with AI",
		Version: AppVersion,
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initializes propr with a base configuration",
				Action: func(ctx *cli.Context) error {
					form := huh.NewForm(
						huh.NewGroup(
							huh.NewConfirm().Title("Assistant").Description("Do you want to use an OpenAI assistant to control your prompt?").Value(&CONFIG.Data.Assistant.Enabled),
						),
						huh.NewGroup(
							huh.NewInput().Title("Assistant").Description("Provide the assistant's id").Value(&CONFIG.Data.Assistant.Id),
						).WithHideFunc(func() bool {
							return !CONFIG.Data.Assistant.Enabled
						}),
						huh.NewGroup(
							huh.NewSelect[string]().Title("Model").Description("Configure the default model").Options(huh.NewOption(openai.GPT4TurboPreview, openai.GPT4TurboPreview), huh.NewOption(openai.GPT3Dot5Turbo, openai.GPT3Dot5Turbo)).Value(&CONFIG.Data.Model),
							huh.NewText().Title("Prompt").Description("Configure the default prompt").CharLimit(99999).Value(&CONFIG.Data.Prompt),
							huh.NewText().Title("Template").Description("Configure the default template").Value(&CONFIG.Data.Template),
						).WithHideFunc(func() bool {
							return CONFIG.Data.Assistant.Enabled
						}),
					)

					err := form.Run()
					if err != nil {
						return err
					}

					return CONFIG.Save()
				},
			},
			{
				Name:  "create",
				Usage: "Creates a PR with a generated description",
				Action: func(ctx *cli.Context) error {
					var description string
					for {
						response, err := propr.Generate()
						if err != nil {
							log.Fatal(err)
						}

						printMarkdown(response)

						var confirmation bool
						err = huh.NewConfirm().Title("Do you want to create this pull request?").Value(&confirmation).Run()
						if err != nil {
							return nil
						}

						if confirmation {
							description = response
							break
						}
					}

					return propr.Create(description)
				},
			},
			{
				Name:  "generate",
				Usage: "Generates a PR description and outputs it",
				Action: func(ctx *cli.Context) error {
					description, err := propr.Generate()
					if err != nil {
						log.Fatal(err)
					}

					printMarkdown(description)
					return nil
				},
			},
			{
				Name:  "config",
				Usage: "Configure propr to your liking",
				Subcommands: []*cli.Command{
					{
						Name:  "ls",
						Usage: "List the current configuration",
						Action: func(ctx *cli.Context) error {
							data, err := json.MarshalIndent(CONFIG.Data, "", "  ")
							if err != nil {
								return err
							}

							fmt.Println(string(data))
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
