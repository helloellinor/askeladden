package commands

import (
	"fmt"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands["kontekst"] = Command{
		name:        "kontekst",
		description: "Vurder om eit ord er feil i ein spesifikk kontekst",
		emoji:       "🔍",
		handler:     handleKontekst,
		aliases:     []string{"context"},
		adminOnly:   false,
	}
}

func handleKontekst(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
	parts := strings.Fields(m.Content)
	if len(parts) < 3 {
		embed := services.CreateErrorEmbed("Manglande parameter", 
			"Bruk: `?kontekst <ord> <kontekst...>`\n"+
			"Døme: `?kontekst huse I huset bur det mange folk`\n\n"+
			"Dette vil sjekke om ordet 'huse' er korrekt i konteksten 'I huset bur det mange folk'.")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	word := parts[1]
	context := strings.Join(parts[2:], " ")

	// For now, this is a basic implementation that could be enhanced with AI or rule-based checking
	analysis := analyzeWordInContext(word, context)

	embed := &discordgo.MessageEmbed{
		Title: "🔍 Kontekstanalyse",
		Color: services.ColorInfo,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Ord",
				Value:  fmt.Sprintf("`%s`", word),
				Inline: true,
			},
			{
				Name:   "Kontekst",
				Value:  fmt.Sprintf("*%s*", context),
				Inline: false,
			},
			{
				Name:   "Vurdering",
				Value:  analysis,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Dette er ein grunnleggjande analyse. For meir avansert hjelp, bruk ?ord kommandoen eller spør i grammatikkkanalen.",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// analyzeWordInContext provides basic context analysis for words
func analyzeWordInContext(word, context string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	context = strings.ToLower(context)

	// Basic checks for common Norwegian grammar mistakes
	if strings.Contains(word, "og") && strings.Contains(context, "og") {
		return "✅ Ordet 'og' er vanlegvis korrekt som bindeord."
	}

	if strings.Contains(word, "å") && strings.Contains(context, "å") {
		return "✅ Infinitivsmerket 'å' ser ut til å vere brukt korrekt."
	}

	// Check for common bokmål vs nynorsk differences
	bokmålWords := map[string]string{
		"ikke": "ikkje",
		"jeg":  "eg",
		"det":  "det/den",
		"som":  "som",
		"en":   "ein",
		"et":   "eit",
	}

	if nynorsk, exists := bokmålWords[word]; exists {
		return fmt.Sprintf("⚠️ '%s' er bokmål. På nynorsk: '%s'", word, nynorsk)
	}

	// Check if word appears in context (basic validation)
	if !strings.Contains(context, word) {
		return "⚠️ Ordet er ikkje funnet i konteksten. Sjekk at du har skrive rett."
	}

	return fmt.Sprintf("ℹ️ Ordet '%s' treng nærmare vurdering. Bruk ?ord %s for ordbok-oppslag eller spør i grammatikkkanalen for hjelp.", word, word)
}