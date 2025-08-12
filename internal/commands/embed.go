package commands

import (
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["embed"] = Command{
		name:        "embed",
		description: "Opprett tilpassa embed-meldingar (kun opplysar)",
		emoji:       "✨",
		handler:     handleEmbed,
		aliases:     []string{"lag-embed"},
		adminOnly:   true, // Only opplysar role
	}
}

func handleEmbed(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// Start DM conversation with the user
	dmChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		embed := services.CreateErrorEmbed("Feil", "Kunne ikkje opprette DM-samtale. Sjekk at du har tillatt DMs frå servermedlemmar.")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Start an embed creation session using the global session manager
	services.GlobalEmbedSessions.StartSession(m.Author.ID, m.GuildID)

	// Send instructions to the user via DM
	instructionEmbed := services.CreateBotEmbed(s, "✨ Embed-byggjar",
		"Velkommen til embed-byggjaren! Følg desse stega:\n\n"+
			"1️⃣ Send tittelen på embedden\n"+
			"2️⃣ Send innhaldet/beskrivinga\n"+
			"3️⃣ Send kanal-ID der embedden skal sendast\n"+
			"4️⃣ (Valfritt) Send farge som hex-kode (t.d. #ff0000)\n\n"+
			"Send `avbryt` når som helst for å avbryte.",
		services.EmbedTypeInfo)

	_, err = s.ChannelMessageSendEmbed(dmChannel.ID, instructionEmbed)
	if err != nil {
		embed := services.CreateErrorEmbed("Feil", "Kunne ikkje sende DM-melding. Sjekk at du har tillatt DMs.")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Confirm in the original channel that DM was sent
	confirmEmbed := services.CreateSuccessEmbed("Embed-byggjar", "Sjekk DMs for instruksjonar om å bygge embedden!")
	s.ChannelMessageSendEmbed(m.ChannelID, confirmEmbed)
}
