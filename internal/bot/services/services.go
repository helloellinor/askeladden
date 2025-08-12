package services

import (
	"askeladden/internal/bot"
)

// BotServices holds all the services the bot can use.
type BotServices struct {
	Approval     *ApprovalService
	Dictionary   *DictionaryService
	EmbedBuilder *EmbedBuilderService
}

// New creates a new BotServices instance.
func New(b *bot.Bot) *BotServices {
	return &BotServices{
		Approval:     &ApprovalService{Bot: b},
		Dictionary:   NewDictionaryService(),
		EmbedBuilder: NewEmbedBuilderService(b.Config),
	}
}
