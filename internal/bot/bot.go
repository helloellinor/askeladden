// Package bot implementerer kjernefunksjonaliteten til Askeladden Discord-boten.
// Denne pakken handsamar Discord-tilkobling, session-handsaming og grunnleggjande bot-operasjonar.
package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"askeladden/internal/config"
	"askeladden/internal/database"
)

// Bot represents the main bot structure.
type Bot struct {
	Session  *discordgo.Session
	Config   *config.Config
	Database *database.DB
}

// New creates a new Bot instance.
func New(cfg *config.Config, db *database.DB, session *discordgo.Session) *Bot {
	return &Bot{
		Session:  session,
		Config:   cfg,
		Database: db,
	}
}

// Start startar boten og opnar Discord-tilkoplinga.
func (b *Bot) Start() error {
	log.Println("[BOT] Prøver å kople til Discord...")
	// Open connection
	err := b.Session.Open()
	if err != nil {
		log.Printf("[BOT] Kunne ikkje opne Discord-sesjon: %v", err)
		return err
	}

	log.Println("[BOT] Discord-sesjon opna")
	log.Println("[BOT] Askeladden køyrer og er klar til å handtere meldingar.")
	return nil
}

// Stop stoppar boten og stenger alle tilkoplingar.
func (b *Bot) Stop() error {
	log.Println("[BOT] Askeladden loggar av.")
	// Log channel message will be sent from main.go before calling Stop()

	// Close database connection
	if b.Database != nil {
		b.Database.Close()
	}

	return b.Session.Close()
}

// Note: Direct field access is preferred in Go for simplicity
// Bot fields (Session, Config, Database) are exported for direct access
