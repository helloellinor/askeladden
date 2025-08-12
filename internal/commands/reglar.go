package commands

import (
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["reglar"] = Command{
		name:        "reglar",
		description: "Syn reglar og retningslinjer for serveren",
		emoji:       "📋",
		handler:     handleReglar,
		aliases:     []string{"rules", "regler"},
		adminOnly:   false,
	}
}

func handleReglar(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	title := bot.Config.Rules.Title
	if title == "" {
		title = "📋 Serverreglar"
	}

	content := bot.Config.Rules.Content
	if content == "" {
		content = "Ver snill og følg Discords retningslinjer og vær respektfull mot andre medlemmar.\n\n" +
			"For meir detaljerte reglar, sjå pinned meldingar eller spør ein moderator."
	}

	embed := services.CreateBotEmbed(s, title, content, services.EmbedTypeInfo)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
