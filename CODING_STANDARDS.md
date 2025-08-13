# Askeladden Coding Standards

This document establishes unified coding standards for the Askeladden Discord bot project to ensure code readability, maintainability, and consistency across the codebase.

## Language Guidelines

### Primary Language: Norwegian (Nynorsk)
Askeladden serves Norwegian language communities, and our code reflects this commitment:

- **Comments**: Write in Norwegian (Nynorsk) when possible
- **User-facing strings**: Always in Norwegian (Nynorsk)
- **Command names**: Use Norwegian terms (`spør`, `godkjenn`, `hjelp`)
- **Function names**: Prefer Norwegian when the function is domain-specific to Norwegian functionality
- **Interface/struct names**: Use English for technical interfaces, Norwegian for domain models

### Language Mixing Guidelines
```go
// ✅ Good: Norwegian comment for Norwegian domain function
// Spor handsamar spørsmål frå brukarar
func Spor(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
    // Implementation
}

// ✅ Good: English for technical interfaces
type DatabaseInterface interface {
    AddQuestion(question, authorID string) error
}

// ✅ Good: Norwegian for domain models
type Sporsmal struct {
    ID     int    `json:"id"`
    Tekst  string `json:"tekst"`
    Brukar string `json:"brukar"`
}
```

## File Organization

### Import Grouping
Always group imports in this order, separated by blank lines:

```go
import (
    // 1. Standard library
    "fmt"
    "log"
    "strings"

    // 2. Third-party packages
    "github.com/bwmarrin/discordgo"
    "github.com/go-sql-driver/mysql"

    // 3. Local packages
    "askeladden/internal/bot"
    "askeladden/internal/config"
)
```

### Package Structure
```
cmd/
├── askeladden/          # Main application entry point
internal/
├── bot/                 # Core bot logic and handlers
├── commands/            # Discord command implementations
├── config/              # Configuration management
├── database/            # Database operations and models
├── permissions/         # Role-based access control
└── reactions/           # Reaction handling
```

## Documentation Standards

### Package-Level Documentation
Every package must include a package comment:

```go
// Package commands implementerer alle Discord-kommandoar for Askeladden.
// Denne pakken handsamar kommandoregistrering og utføring.
package commands
```

### Function Documentation
Document all exported functions:

```go
// Spor handsamar spør-kommandoen som let brukarar leggje til spørsmål.
// Funksjonen validerer input og sender spørsmålet til godkjenning.
func Spor(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
    // Implementation
}
```

### Comment Style
- Use `//` for single-line comments
- Start comments with capital letters
- End comments with periods for complete sentences
- Keep line length under 80 characters when possible

## Code Style

### Error Handling
Always handle errors explicitly:

```go
// ✅ Good
embed := services.CreateBotEmbed(s, "Feil", "Kunne ikkje lagre spørsmål", services.EmbedTypeError)
response, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
if err != nil {
    log.Printf("Feil ved sending av melding: %v", err)
    return
}

// ❌ Bad
response, _ := s.ChannelMessageSendEmbed(m.ChannelID, embed)
```

### Variable Naming
- Use clear, descriptive names
- Prefer Norwegian for domain-specific variables
- Use English for technical/generic variables

```go
// ✅ Good
sporsmalsID := 42
authorID := m.Author.ID
channelID := m.ChannelID

// ❌ Bad
id := 42
a := m.Author.ID
c := m.ChannelID
```

### Function Structure
Keep functions focused and under 50 lines when possible:

```go
func ValditerSporsmål(tekst string) error {
    if strings.TrimSpace(tekst) == "" {
        return fmt.Errorf("spørsmål kan ikkje vere tomt")
    }
    if len(tekst) > 500 {
        return fmt.Errorf("spørsmål er for langt (maks 500 teikn)")
    }
    return nil
}
```

## Database Conventions

### Interface Design
- Use clear, English interface names for technical contracts
- Implement Norwegian method names where domain-appropriate

```go
type DatabaseInterface interface {
    AddQuestion(question, authorID, authorName, messageID, channelID string) (int64, error)
    ApproveQuestion(questionID int, approverID string) error
    GetPendingQuestion() (*Question, error)
}
```

### Model Naming
```go
// Norwegian for domain models
type Sporsmal struct {
    ID          int       `json:"id"`
    Tekst       string    `json:"tekst"`
    Brukar      string    `json:"brukar"`
    Godkjent    bool      `json:"godkjent"`
    OpprettaTid time.Time `json:"oppretta_tid"`
}

// English for technical models
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
```

## Discord Integration

### Embed Standards
Follow the embed guidelines in [docs/EMBEDS.md](docs/EMBEDS.md):

```go
// Use standardized embed creation
embed := services.CreateBotEmbed(
    s, 
    "✅ Godkjent", 
    "Spørsmålet ditt er godkjent!", 
    services.EmbedTypeSuccess,
)
```

### Message Handling
- Always validate input before processing
- Provide clear error messages in Norwegian
- Log errors with context for debugging

## Testing Guidelines

### Manual Testing
- Build with `go build ./cmd/askeladden`
- Test help functionality with `./tools/test_help`
- Verify binary size (~12MB expected)
- Run `go fmt ./...` before committing
- Run `go vet ./cmd/askeladden ./internal/...` to check for issues

### Command Testing
```bash
# Standard validation workflow
go mod download
go build ./cmd/askeladden
go fmt ./...
go vet ./cmd/askeladden ./internal/...
go build ./tools/test_help
./test_help | wc -l  # Should output ~12 lines
```

## Git Workflow

### Commit Messages
- Use Norwegian for feature commits
- Use English for technical/infrastructure changes
- Start with a verb: "Legg til", "Fiks", "Oppdater"

```
✅ Good: "Legg til validering for spørsmålslengd"
✅ Good: "Fix import grouping in handlers package"
❌ Bad: "updates"
❌ Bad: "misc changes"
```

### File Organization
- Exclude binaries and build artifacts
- Keep `tools/` separate for utility scripts
- Maintain clean `.gitignore`

## Configuration Management

### Environment Variables
- `CONFIG_FILE`: Path to main config (default: `config/config.yaml`)
- `SECRETS_FILE`: Path to secrets (default: `config/secrets.yaml`)

### YAML Structure
Use consistent formatting:

```yaml
discord:
  token: ""  # From secrets file
  prefix: "!"
  channels:
    log_channel_id: ""
    approval_channel_id: ""

database:
  host: "localhost"
  port: 3306
  name: "askeladden"
```

## Performance Guidelines

### Build Performance
- Expected build time: ~18 seconds initial, ~0.4 seconds subsequent
- Binary size should be ~11-12MB
- Module download: ~5 seconds

### Code Performance
- Use appropriate data structures
- Avoid unnecessary allocations in hot paths
- Cache frequently accessed data

## Review Checklist

Before submitting code:

- [ ] Go fmt applied (`go fmt ./...`)
- [ ] Go vet passes (`go vet ./cmd/askeladden ./internal/...`)
- [ ] Build succeeds (`go build ./cmd/askeladden`)
- [ ] Help tool works (`./test_help` outputs ~12 lines)
- [ ] Imports properly grouped
- [ ] Functions documented if exported
- [ ] Error handling implemented
- [ ] Norwegian language used appropriately
- [ ] Commit message follows conventions

---

*Denne standarden sikrar at Askeladden-koden er lesbar, konsis og enkel å vedlikehalde for det norske miljøet vårt.*