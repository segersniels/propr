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
	"github.com/segersniels/config"
	updater "github.com/segersniels/updater"
	"github.com/urfave/cli/v2"
)

var (
	AppVersion string
	AppName    string
)

type SupportedModel string

const (
	GPT4o             SupportedModel = "gpt-4o"
	GPT4oMini         SupportedModel = "gpt-4o-mini"
	GPTo1             SupportedModel = "o1"
	GPTo1Mini         SupportedModel = "o1-mini"
	Claude3Dot5Sonnet SupportedModel = "claude-3-5-sonnet-latest"
	Claude3Dot5Haiku  SupportedModel = "claude-3-5-haiku-latest"
	DeepSeekChat      SupportedModel = "deepseek-chat"
	DeepSeekReasoner  SupportedModel = "deepseek-reasoner"
)

var SupportedModels = []SupportedModel{
	GPT4o,
	GPT4oMini,
	GPTo1,
	GPTo1Mini,
	Claude3Dot5Sonnet,
	Claude3Dot5Haiku,
	DeepSeekChat,
	DeepSeekReasoner,
}

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

type MessageClient interface {
	CreateMessage(ctx context.Context, system string, messages []Message) (string, error)
}

type Config struct {
	Model       SupportedModel `json:"model"`
	Prompt      string         `json:"prompt"`
	Template    string         `json:"template"`
	PrettyPrint bool           `json:"pretty_print"`
}

var CONFIG = config.NewConfig("propr", Config{
	Model:       GPT4oMini,
	Prompt:      SYSTEM_MESSAGE,
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
	upd := updater.NewUpdater(AppName, AppVersion, "segersniels")
	err := upd.CheckIfNewVersionIsAvailable()
	if err != nil {
		log.Debug("Failed to check for latest release", "error", err)
	}

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
					&cli.BoolFlag{
						Name:  "empty",
						Usage: "Create an empty PR",
					},
					&cli.BoolFlag{
						Name:  "draft",
						Usage: "Create a draft PR",
					},
				},
				Action: func(ctx *cli.Context) error {
					draft := ctx.Bool("draft")
					branch := ctx.String("branch")
					if ctx.Bool("empty") {
						return propr.Create(branch, "", draft)
					}

					var description string
					for {
						response, err := propr.Generate(branch)
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

					return propr.Create(branch, description, draft)
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
							models := huh.NewOptions(SupportedModels...)
							form := huh.NewForm(
								huh.NewGroup(
									huh.NewSelect[SupportedModel]().Title("Model").Description("Configure the default model").Options(models...).Value(&CONFIG.Data.Model),
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
