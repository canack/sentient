package tools

import "github.com/canack/sentient/pkg/models"

type ChatBot interface {
	Setup(token string) error
	TestConnection() error
	Query(string) models.ResponseMessage
}

type ChatBotInstance struct {
	ChatBot
}

func NewChatBot(bot ChatBot) *ChatBotInstance {
	return &ChatBotInstance{
		ChatBot: bot,
	}
}
