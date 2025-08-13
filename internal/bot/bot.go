// Package bot implementerer kjernefunksjonaliteten til Askeladden Discord-boten.
// Denne pakken handsamar Discord-tilkobling, session-handsaming og grunnleggjande bot-operasjonar.
package bot

import (
	"github.com/bwmarrin/discordgo"

	"askeladden/internal/config"
	"askeladden/internal/database"
	"askeladden/internal/logging"
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
	logger := logging.GetLogger("BOT")
	logger.Info("Attempting to connect to Discord...")

	// Open connection
	err := b.Session.Open()
	if err != nil {
		logger.Error("Could not open Discord session: %v", err)
		return err
	}

	logger.Info("Discord session opened")
	logger.Info("Askeladden is running and ready to handle messages")
	return nil
}

// Stop stoppar boten og stenger alle tilkoplingar.
func (b *Bot) Stop() error {
	logger := logging.GetLogger("BOT")
	logger.Info("Askeladden is logging off")
	// Log channel message will be sent from main.go before calling Stop()

	// Close database connection
	if b.Database != nil {
		b.Database.Close()
	}

	return b.Session.Close()
}

// Note: Direct field access is preferred in Go for simplicity
// Bot fields (Session, Config, Database) are exported for direct access
