package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"askeladden/internal/config"
	"github.com/bwmarrin/discordgo"
)

// EmbedSession represents an ongoing embed creation session
type EmbedSession struct {
	UserID    string
	GuildID   string
	Step      int // 0=title, 1=description, 2=channel, 3=color(optional)
	Title     string
	Content   string
	ChannelID string
	Color     int
}

// EmbedBuilderService handles DM-based embed creation
type EmbedBuilderService struct {
	Config *config.Config
}

// NewEmbedBuilderService creates a new embed builder service
func NewEmbedBuilderService(cfg *config.Config) *EmbedBuilderService {
	return &EmbedBuilderService{
		Config: cfg,
	}
}

// StartSession starts a new embed creation session for a user
func (ebs *EmbedBuilderService) StartSession(userID, guildID string) {
	GlobalEmbedSessions.StartSession(userID, guildID)
}

// HandleDMMessage processes DM messages for embed creation
func (ebs *EmbedBuilderService) HandleDMMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	session, exists := GlobalEmbedSessions.GetSession(m.Author.ID)

	// If no session exists, ignore the DM
	if !exists {
		return false
	}

	// Handle cancel command
	if strings.ToLower(strings.TrimSpace(m.Content)) == "avbryt" {
		GlobalEmbedSessions.RemoveSession(m.Author.ID)
		embed := CreateBotEmbed(s, "❌ Avbrutt", "Embed-bygging avbrutt.", EmbedTypeError)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return true
	}

	switch session.Step {
	case 0: // Title
		session.Title = strings.TrimSpace(m.Content)
		if len(session.Title) > 256 {
			embed := CreateErrorEmbed("Tittel for lang", "Tittelen kan maksimalt vere 256 teikn. Prøv igjen.")
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return true
		}
		session.Step = 1
		embed := CreateBotEmbed(s, "✅ Tittel lagra",
			fmt.Sprintf("Tittel: **%s**\n\nNo send beskrivinga/innhaldet for embedden.", session.Title),
			EmbedTypeSuccess)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	case 1: // Description/Content
		session.Content = strings.TrimSpace(m.Content)
		if len(session.Content) > 4096 {
			embed := CreateErrorEmbed("Innhald for langt", "Innhaldet kan maksimalt vere 4096 teikn. Prøv igjen.")
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return true
		}
		session.Step = 2
		embed := CreateBotEmbed(s, "✅ Innhald lagra",
			"No send kanal-ID der embedden skal sendast.\n\nTips: Høgreklikk på kanalen og vel 'Copy Channel ID'",
			EmbedTypeSuccess)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	case 2: // Channel ID
		channelID := strings.TrimSpace(m.Content)

		// Verify the channel exists and the bot has access
		channel, err := s.Channel(channelID)
		if err != nil {
			embed := CreateErrorEmbed("Ugyldig kanal", "Kunne ikkje finne kanalen. Sjekk at kanal-ID er korrekt.")
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return true
		}

		// Check if the channel is in the same guild
		if channel.GuildID != session.GuildID {
			embed := CreateErrorEmbed("Feil server", "Kanalen må vere på same server der du køyrde kommandoen.")
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return true
		}

		session.ChannelID = channelID
		session.Step = 3

		// Show preview and ask for optional color
		previewEmbed := &discordgo.MessageEmbed{
			Title:       session.Title,
			Description: session.Content,
			Color:       session.Color,
		}

		embed := CreateBotEmbed(s, "✅ Kanal lagra",
			fmt.Sprintf("Kanalen: <#%s>\n\nFørehandsvising av embedden:", channelID),
			EmbedTypeSuccess)
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		s.ChannelMessageSendEmbed(m.ChannelID, previewEmbed)

		finalEmbed := CreateBotEmbed(s, "🎨 Valfritt: Farge",
			"Send ein hex-fargekode (t.d. #ff0000) for å endre fargen, eller send `ferdig` for å publisere embedden no.",
			EmbedTypeInfo)
		s.ChannelMessageSendEmbed(m.ChannelID, finalEmbed)

	case 3: // Optional color or finish
		content := strings.TrimSpace(strings.ToLower(m.Content))

		if content == "ferdig" {
			ebs.publishEmbed(s, m.ChannelID, session)
			GlobalEmbedSessions.RemoveSession(m.Author.ID)
			return true
		}

		// Try to parse hex color
		if strings.HasPrefix(content, "#") {
			color, err := ebs.parseHexColor(content)
			if err != nil {
				embed := CreateErrorEmbed("Ugyldig farge", "Bruk format som #ff0000. Send `ferdig` for å publisere med standard farge.")
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return true
			}
			session.Color = color

			// Show updated preview
			previewEmbed := &discordgo.MessageEmbed{
				Title:       session.Title,
				Description: session.Content,
				Color:       session.Color,
			}

			embed := CreateBotEmbed(s, "🎨 Farge oppdatert", "Oppdatert førehandsvising:", EmbedTypeSuccess)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			s.ChannelMessageSendEmbed(m.ChannelID, previewEmbed)

			finalEmbed := CreateBotEmbed(s, "✨ Klar til publisering", "Send `ferdig` for å publisere embedden no.", EmbedTypeInfo)
			s.ChannelMessageSendEmbed(m.ChannelID, finalEmbed)
		} else {
			embed := CreateErrorEmbed("Ukjend kommando", "Send ein hex-fargekode (t.d. #ff0000) eller `ferdig` for å publisere.")
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}

	return true
}

// publishEmbed publishes the completed embed to the target channel
func (ebs *EmbedBuilderService) publishEmbed(s *discordgo.Session, dmChannelID string, session *EmbedSession) {
	// Create the final embed
	finalEmbed := &discordgo.MessageEmbed{
		Title:       session.Title,
		Description: session.Content,
		Color:       session.Color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Laga av Askeladden",
		},
	}

	// Send to target channel
	_, err := s.ChannelMessageSendEmbed(session.ChannelID, finalEmbed)
	if err != nil {
		embed := CreateErrorEmbed("Publiseringsfeil", fmt.Sprintf("Kunne ikkje sende embedden: %v", err))
		s.ChannelMessageSendEmbed(dmChannelID, embed)
		return
	}

	// Confirm success
	embed := CreateSuccessEmbed("✅ Publisert!",
		fmt.Sprintf("Embedden har blitt sendt til <#%s>!", session.ChannelID))
	s.ChannelMessageSendEmbed(dmChannelID, embed)

	// Log the action if log channel is configured
	if ebs.Config != nil && ebs.Config.Discord.LogChannelID != "" {
		user, _ := s.User(session.UserID)
		username := "Ukjend brukar"
		if user != nil {
			username = user.Username
		}

		logEmbed := CreateBotEmbed(s, "📝 Embed publisert",
			fmt.Sprintf("**Brukar:** %s\n**Kanal:** <#%s>\n**Tittel:** %s",
				username, session.ChannelID, session.Title),
			EmbedTypeInfo)
		s.ChannelMessageSendEmbed(ebs.Config.Discord.LogChannelID, logEmbed)
	}
}

// parseHexColor converts a hex color string to an integer
func (ebs *EmbedBuilderService) parseHexColor(hexColor string) (int, error) {
	// Remove # if present
	hexColor = strings.TrimPrefix(hexColor, "#")

	// Validate hex format
	if !regexp.MustCompile(`^[0-9a-fA-F]{6}$`).MatchString(hexColor) {
		return 0, fmt.Errorf("ugyldig hex-format")
	}

	// Parse to integer
	color, err := strconv.ParseInt(hexColor, 16, 64)
	if err != nil {
		return 0, err
	}

	return int(color), nil
}

// HasActiveSession checks if a user has an active embed creation session
func (ebs *EmbedBuilderService) HasActiveSession(userID string) bool {
	return GlobalEmbedSessions.HasActiveSession(userID)
}
