package main

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

var _ MessageClient = (*OpenAI)(nil)

type OpenAI struct {
	apiKey string
	model  SupportedModel
}

func NewOpenAI(apiKey string, model SupportedModel) *OpenAI {
	return &OpenAI{
		apiKey,
		model,
	}
}

func convertToOpenAIMessages(messages []Message) []openai.ChatCompletionMessage {
	var msgs []openai.ChatCompletionMessage
	for _, m := range messages {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	return msgs
}

func (o *OpenAI) CreateMessage(ctx context.Context, system string, messages []Message) (string, error) {
	client := openai.NewClient(o.apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: string(o.model),
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
