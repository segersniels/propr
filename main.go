package main

import (
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

	app := &cli.App{
		Name:    AppName,
		Usage:   "Generate your PRs from the command line with AI",
		Version: AppVersion,
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Creates a PR with a generated description",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "branch",
						Aliases: []string{"b"},
						Usage:   "Select a branch from a list",
					},
					&cli.BoolFlag{
						Name:    "model",
						Aliases: []string{"m"},
						Usage:   "Select a model from a list",
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

					propr, err := NewPropr(ctx)
					if err != nil {
						log.Fatal(err)
					}

					if ctx.Bool("empty") {
						return propr.Create("", "", draft)
					}

					var description string
					for {
						response, err := propr.Generate("")
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

					return propr.Create("", description, draft)
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
					&cli.BoolFlag{
						Name:    "branch",
						Aliases: []string{"b"},
						Usage:   "Select a branch from a list",
					},
					&cli.BoolFlag{
						Name:    "model",
						Aliases: []string{"m"},
						Usage:   "Select a model from a list",
					},
				},
				Action: func(ctx *cli.Context) error {
					propr, err := NewPropr(ctx)
					if err != nil {
						log.Fatal(err)
					}

					description, err := propr.Generate("")
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
