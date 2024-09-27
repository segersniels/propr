package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
	"github.com/segersniels/config"
	"github.com/urfave/cli/v2"
)

var AppVersion string
var AppName string

const (
	GPT4o             = "gpt-4o"
	GPT4oMini         = "gpt-4o-mini"
	GPT4Turbo         = "gpt-4-turbo"
	GPT3Dot5Turbo     = "gpt-3.5-turbo"
	Claude3Dot5Sonnet = "claude-3-5-sonnet-20240620"
)

const (
	MessageRoleSystem    = "system"
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

type MessageClient interface {
	CreateMessage(ctx context.Context, system string, prompt string) (string, error)
}

type Config struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Template    string `json:"template"`
	PrettyPrint bool   `json:"pretty_print"`
}

var CONFIG = config.NewConfig("propr", Config{
	Model: openai.GPT4o,
	Prompt: `You will be asked to write a concise GitHub PR description based on a provided git diff.
Analyze the code changes and provide a concise explanation of the changes, their context and why they were made.
Don't reference file names or directories directly, instead give a general explanation of the changes made.
Do not treat imports and requires as changes or new features. If the provided message is not a diff respond with an appropriate message.
Don't surround your description in backticks but still write GitHub supported markdown.`,
	Template:    "# Description",
	PrettyPrint: true,
})

func printMarkdown(content string, pretty bool) error {
	// Remove leading and trailing backticks
	if strings.HasPrefix(content, "```") {
		content = strings.Trim(content, "`")
	}

	if !pretty {
		println(content)
		return nil
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(content)
	if err != nil {
		return err
	}

	println(out)
	return nil
}

func main() {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		log.SetLevel(log.DebugLevel)
	}

	propr := NewPropr()
	app := &cli.App{
		Name:    AppName,
		Usage:   "Generate your PRs from the command line with AI",
		Version: AppVersion,
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Creates a PR with a generated description",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "branch",
						Usage:       "The branch to compare your changes against",
						DefaultText: "HEAD",
					},
				},
				Action: func(ctx *cli.Context) error {
					var description string
					for {
						response, err := propr.Generate(ctx.String("branch"))
						if err != nil {
							log.Fatal(err)
						}

						err = printMarkdown(response, CONFIG.Data.PrettyPrint)
						if err != nil {
							log.Fatal(err)
						}

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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "plain",
						Usage: "Output the generated description without any formatting",
					},
					&cli.StringFlag{
						Name:        "branch",
						Usage:       "The branch to compare your changes against",
						DefaultText: "HEAD",
					},
				},
				Action: func(ctx *cli.Context) error {
					description, err := propr.Generate(ctx.String("branch"))
					if err != nil {
						log.Fatal(err)
					}

					pretty := CONFIG.Data.PrettyPrint
					if ctx.Bool("plain") {
						pretty = false
					}

					return printMarkdown(description, pretty)
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
					{
						Name:  "init",
						Usage: "Initializes propr with a base configuration",
						Action: func(ctx *cli.Context) error {
							models := huh.NewOptions(GPT4oMini, GPT4o, GPT4Turbo, GPT3Dot5Turbo, Claude3Dot5Sonnet)
							form := huh.NewForm(
								huh.NewGroup(
									huh.NewSelect[string]().Title("Model").Description("Configure the default model").Options(models...).Value(&CONFIG.Data.Model),
									huh.NewText().Title("Prompt").Description("Configure the default prompt").CharLimit(99999).Value(&CONFIG.Data.Prompt),
								),
								huh.NewGroup(
									huh.NewText().Title("Template").Description("Configure the default template").Value(&CONFIG.Data.Template),
								),
								huh.NewGroup(
									huh.NewConfirm().Title("Pretty Print").Description("Do you want to pretty print the generated output?").Value(&CONFIG.Data.PrettyPrint),
								),
							)

							err := form.Run()
							if err != nil {
								return err
							}

							return CONFIG.Save()
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
