package tools

type ChatBot interface {
	Setup(token string) error
	TestConnection() error
	Query(string) string
}

type ChatBotInstance struct {
	ChatBot
}

func NewChatBot(bot ChatBot) *ChatBotInstance {
	return &ChatBotInstance{
		ChatBot: bot,
	}
}
