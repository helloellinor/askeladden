// Package commands implementerer alle Discord-kommandoar for Askeladden.
// Denne pakken handsamar kommandoregistrering, utføring og validering av alle
// brukarkommandoar som boten støttar.
package commands

import (
	"fmt"
	"log"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

// Command defines the structure for a command
type Command struct {
	name        string
	description string
	emoji       string
	handler     func(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot)
	aliases     []string
	adminOnly   bool
}

// commands holds all the registered commands
var commands = make(map[string]Command)

// MatchAndRunCommand finn og utfører ein kommando basert på namn eller alias.
// Input inneheld kommandoen med prefix (t.d. "!spør"), og funksjonen fjernar
// prefikset og søkjer etter samsvarande kommando eller alias.
func MatchAndRunCommand(input string, s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	// `input` is the command with prefix, e.g., "?spør"
	// Remove prefix to get the actual command
	commandWithoutPrefix := strings.TrimPrefix(input, bot.Config.Discord.Prefix)

	log.Printf("[DEBUG] MatchAndRunCommand: input='%s', prefix='%s', command='%s'", input, bot.Config.Discord.Prefix, commandWithoutPrefix)

	// Debug: list all registered commands
	log.Printf("[DEBUG] Registered commands: %v", getCommandNames())

	// Try to find command by name without prefix
	if cmd, exists := commands[commandWithoutPrefix]; exists {
		log.Printf("[DEBUG] Found command '%s', executing", commandWithoutPrefix)
		cmd.handler(s, m, bot)
		return
	}

	log.Printf("[DEBUG] Command '%s' not found, checking aliases", commandWithoutPrefix)
	// Check aliases
	for _, cmd := range commands {
		for _, alias := range cmd.aliases {
			if alias == commandWithoutPrefix {
				log.Printf("[DEBUG] Found alias '%s', executing", alias)
				cmd.handler(s, m, bot)
				return
			}
		}
	}

	log.Printf("[DEBUG] No command or alias found for '%s'", commandWithoutPrefix)
}

// Helper function to get command names for debugging
func getCommandNames() []string {
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	return names
}

// IsAdminCommand sjekkar om ein kommando er berre for administratorar.
// Funksjonen fjernar prefix frå kommandonamnet og søkjer i kommandoregisteret.
func IsAdminCommand(commandName string) bool {
	// Remove prefix from command name for lookup
	commandWithoutPrefix := strings.TrimPrefix(commandName, "!")
	commandWithoutPrefix = strings.TrimPrefix(commandWithoutPrefix, "?")

	if cmd, exists := commands[commandWithoutPrefix]; exists {
		return cmd.adminOnly
	}
	return false
}

// GetHelpText genererer hjelpetekst for alle kommandoar.
// Returnerer ein formatert streng med alle tilgjengelege kommandoar og deira beskriving.
func GetHelpText() string {
	var helpText strings.Builder
	helpText.WriteString("**Askeladden - Kommandoer:**\n")

	for _, cmd := range commands {
		helpText.WriteString(fmt.Sprintf("%s `%s` - %s\n", cmd.emoji, cmd.name, cmd.description))
	}

	return strings.TrimSpace(helpText.String())
}

// ListCommands lagar ei liste over kommandoar basert på admin-status.
// Returnerer ein Discord-embed med kommandoar gruppert etter tilgang.
func ListCommands(isAdmin bool) *discordgo.MessageEmbed {
	var generalCommands, adminCommands strings.Builder

	for _, cmd := range commands {
		commandLine := fmt.Sprintf("%s `%s` - %s\n", cmd.emoji, cmd.name, cmd.description)
		if cmd.adminOnly {
			adminCommands.WriteString(commandLine)
		} else {
			generalCommands.WriteString(commandLine)
		}
	}

	embed := &discordgo.MessageEmbed{
		Title: "Askeladden - Kommandoer",
		Color: services.ColorInfo, // Blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Generelle kommandoer",
				Value: generalCommands.String(),
			},
		},
	}

	if isAdmin && adminCommands.Len() > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Admin-kommandoer",
			Value: adminCommands.String(),
		})
	}

	return embed
}
