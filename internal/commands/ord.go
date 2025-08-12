package commands

import (
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["ord"] = Command{
		name:        "ord",
		description: "Slå opp eit ord i ordbøkene.no",
		emoji:       "📚",
		handler:     handleOrd,
		aliases:     []string{"ordbok", "lookup"},
		adminOnly:   false,
	}
}

func handleOrd(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// Extract the word from the message
	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		embed := services.CreateErrorEmbed("Manglande ord", "Bruk: `?ord <ord>` for å slå opp eit ord i ordbøkene.no")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Get the word to look up (everything after the command)
	word := strings.Join(parts[1:], " ")

	// Create dictionary service and look up the word
	dictService := services.NewDictionaryService()
	wordInfo, err := dictService.LookupWord(word)
	if err != nil {
		embed := services.CreateErrorEmbed("Feil ved oppslag", "Det oppstod ein feil under oppslag av ordet.")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Create and send the embed
	embed := dictService.CreateWordLookupEmbed(wordInfo)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
