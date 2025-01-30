package main

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

var _ MessageClient = (*DeepSeek)(nil)

type DeepSeek struct {
	apiKey string
	model  SupportedModel
}

func NewDeepSeek(apiKey string, model SupportedModel) *DeepSeek {
	return &DeepSeek{
		apiKey,
		model,
	}
}

func (d *DeepSeek) CreateMessage(ctx context.Context, system string, messages []Message) (string, error) {
	config := openai.DefaultConfig(d.apiKey)
	config.BaseURL = "https://api.deepseek.com"
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: string(d.model),
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system,
				},
			}, convertToOpenAIMessages(messages)...),
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
