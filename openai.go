package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	openai "github.com/sashabaranov/go-openai"
)

var FILES_TO_IGNORE = []string{
	"package-lock.json",
	"yarn.lock",
	"npm-debug.log",
	"yarn-debug.log",
	"yarn-error.log",
	".pnpm-debug.log",
	"Cargo.lock",
	"Gemfile.lock",
	"mix.lock",
	"Pipfile.lock",
	"composer.lock",
	"go.sum",
}

func splitDiffIntoChunks(diff string) []string {
	split := strings.Split(diff, "diff --git")[1:]
	for i, chunk := range split {
		split[i] = strings.TrimSpace(chunk)
	}

	return split
}

func removeLockFiles(chunks []string) []string {
	var wg sync.WaitGroup

	filtered := make(chan string)

	for _, chunk := range chunks {
		wg.Add(1)

		go func(chunk string) {
			defer wg.Done()
			shouldIgnore := false
			header := strings.Split(chunk, "\n")[0]

			// Check if the first line contains any of the files to ignore
			for _, file := range FILES_TO_IGNORE {
				if strings.Contains(header, file) {
					log.Debug("Ignoring", "file", file)
					shouldIgnore = true
				}
			}

			if !shouldIgnore {
				log.Debug("Adding", "header", header)
				filtered <- chunk
			}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(filtered)
	}()

	var result []string
	for chunk := range filtered {
		result = append(result, chunk)
	}

	return result
}

// Split the diff in chunks and remove any lock files to save on tokens
func prepareDiff(diff string) string {
	chunks := splitDiffIntoChunks(diff)

	return strings.Join(removeLockFiles(chunks), "\n")
}

func generateSystemMessageForDiff(systemMessage string, template string) string {
	return fmt.Sprintf("%s\n\nFollow this exact template to write your description:\n\n```\n%s\n```", systemMessage, template)
}

func generateUserMessage(diff string, template string) string {
	return fmt.Sprintf("Use the following template to write your description, don't deviate from the template:\n\n```\n%s\n```\n\nThe diff:\n\n```\n%s\n```", template, prepareDiff(diff))
}

type OpenAI struct {
	ApiKey string
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		ApiKey: apiKey,
	}
}

func (o *OpenAI) GetChatCompletion(diff string) (string, error) {
	client := openai.NewClient(o.ApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: CONFIG.Data.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: generateSystemMessageForDiff(CONFIG.Data.Prompt, CONFIG.Data.Template),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prepareDiff(diff),
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	log.Debug("Run status", "usage", resp.Usage)

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) GetAssistantCompletion(diff string) (string, error) {
	config := openai.DefaultConfig(o.ApiKey)
	config.AssistantVersion = "v2" // `gpt-4o` requires the v2 assistant API
	client := openai.NewClientWithConfig(config)

	ctx := context.Background()
	resp, err := client.CreateThreadAndRun(
		ctx,
		openai.CreateThreadAndRunRequest{
			RunRequest: openai.RunRequest{
				AssistantID: CONFIG.Data.Assistant.Id,
			},
			Thread: openai.ThreadRequest{
				Messages: []openai.ThreadMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: generateUserMessage(diff, CONFIG.Data.Template),
					},
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	for {
		run, err := client.RetrieveRun(ctx, resp.ThreadID, resp.ID)
		if err != nil {
			return "", err
		}

		log.Debug("Run status", "status", run.Status)

		if run.Status == openai.RunStatusFailed || run.Status == openai.RunStatusCancelled || run.Status == openai.RunStatusExpired {
			return "", fmt.Errorf("run failed: %v", run.LastError)
		}

		if run.Status == openai.RunStatusCompleted {
			amountOfMessages := 1
			messages, err := client.ListMessage(ctx, run.ThreadID, &amountOfMessages, nil, nil, nil)

			if err != nil {
				return "", err
			}

			if len(messages.Messages) == 0 {
				return "", fmt.Errorf("no messages found")
			}

			log.Debug("Run status", "usage", run.Usage)

			return messages.Messages[0].Content[0].Text.Value, nil
		}
	}
}
