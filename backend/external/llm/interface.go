package llm

type LLM interface {
	Chat(prompt string, systemPrompt ...string) (string, error)
}
