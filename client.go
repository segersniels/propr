package main

import "context"

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

type SupportedModel string

const (
	GPT4o             SupportedModel = "gpt-4o"
	GPT4oMini         SupportedModel = "gpt-4o-mini"
	GPTo1             SupportedModel = "o1"
	GPTo1Mini         SupportedModel = "o1-mini"
	GPTo3Mini         SupportedModel = "o3-mini"
	Claude3Dot7Sonnet SupportedModel = "claude-3-7-sonnet-latest"
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
	GPTo3Mini,
	Claude3Dot7Sonnet,
	Claude3Dot5Sonnet,
	Claude3Dot5Haiku,
	DeepSeekChat,
	DeepSeekReasoner,
}
