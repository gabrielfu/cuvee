package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func SystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

func UserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

func AIMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

type LLM interface {
	Chat(messages []Message) (Message, error)
}
