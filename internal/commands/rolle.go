package commands

import (
	"fmt"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["rolle"] = Command{
		name:        "rolle",
		description: "Administrer roller for brukarar (kun admin)",
		emoji:       "👑",
		handler:     handleRolle,
		aliases:     []string{"role"},
		adminOnly:   true,
	}
}

func handleRolle(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	parts := strings.Fields(m.Content)
	if len(parts) < 3 {
		embed := services.CreateErrorEmbed("Manglande parameter",
			"Bruk: `?rolle <add/remove> <brukar> <rolle>`\n"+
				"Tilgjengelege roller:\n"+
				"• `opplysar` - Kan godkjenne spørsmål og ord\n"+
				"• `rettskrivar` - Kan godkjenne rettskriving")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	action := strings.ToLower(parts[1])
	if action != "add" && action != "remove" {
		embed := services.CreateErrorEmbed("Ugyldig handling", "Bruk `add` eller `remove`")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Parse user mention
	userMention := parts[2]
	userID := ""

	// Extract user ID from mention
	if strings.HasPrefix(userMention, "<@") && strings.HasSuffix(userMention, ">") {
		userID = strings.TrimSuffix(strings.TrimPrefix(userMention, "<@"), ">")
		// Remove ! if present (for nickname mentions)
		userID = strings.TrimPrefix(userID, "!")
	} else {
		embed := services.CreateErrorEmbed("Ugyldig brukar", "Du må mentionere ein brukar med @brukarnamn")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	if len(parts) < 4 {
		embed := services.CreateErrorEmbed("Manglande rolle", "Du må spesifisere kva rolle som skal leggjast til/fjernast")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	roleName := strings.ToLower(parts[3])
	var roleID string

	switch roleName {
	case "opplysar":
		roleID = bot.Config.Approval.OpplysarRoleID
	case "rettskrivar":
		roleID = bot.Config.BannedWords.RettskrivarRoleID
	default:
		embed := services.CreateErrorEmbed("Ukjend rolle",
			"Tilgjengelege roller: `opplysar`, `rettskrivar`")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	if roleID == "" {
		embed := services.CreateErrorEmbed("Konfigurasjonsfeil",
			fmt.Sprintf("Rolle-ID for %s er ikkje konfigurert", roleName))
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Get user info for confirmation
	user, err := s.User(userID)
	if err != nil {
		embed := services.CreateErrorEmbed("Feil", "Kunne ikkje finne brukaren")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Perform the role action
	var actionErr error
	if action == "add" {
		actionErr = s.GuildMemberRoleAdd(m.GuildID, userID, roleID)
	} else {
		actionErr = s.GuildMemberRoleRemove(m.GuildID, userID, roleID)
	}

	if actionErr != nil {
		embed := services.CreateErrorEmbed("Feil",
			fmt.Sprintf("Kunne ikkje %s rolle: %v",
				map[string]string{"add": "legge til", "remove": "fjerne"}[action],
				actionErr))
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	// Success message
	actionText := map[string]string{
		"add":    "lagt til",
		"remove": "fjerna",
	}[action]

	embed := services.CreateSuccessEmbed("Rolle oppdatert",
		fmt.Sprintf("Rolla `%s` har blitt %s for %s", roleName, actionText, user.Username))
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
