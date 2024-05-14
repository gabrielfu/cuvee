package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAILLM struct {
	client *openai.Client
	model  string
}

func translateRole(role Role) string {
	switch role {
	case RoleSystem:
		return openai.ChatMessageRoleSystem
	case RoleUser:
		return openai.ChatMessageRoleUser
	case RoleAssistant:
		return openai.ChatMessageRoleAssistant
	default:
		return openai.ChatMessageRoleUser
	}
}

func NewOpenAILLM(apiKey, model string) *OpenAILLM {
	client := openai.NewClient(apiKey)
	return &OpenAILLM{client: client, model: model}
}

func (o *OpenAILLM) Chat(messages []Message) (Message, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    translateRole(m.Role),
			Content: m.Content,
		}
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    o.model,
			Messages: openaiMessages,
		},
	)
	if err != nil {
		return Message{}, err
	}
	return AIMessage(resp.Choices[0].Message.Content), nil
}
