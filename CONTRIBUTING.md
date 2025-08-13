# Contributing to Askeladden

Takk for at du vil bidra til Askeladden! This guide will help you get started with contributing to our Norwegian Discord bot project.

## Development Environment

### Prerequisites
- Go 1.24+ (tested with 1.24.5)
- Git
- MySQL database (for local testing)
- Discord bot token (for testing)

### Quick Start
```bash
# Clone the repository
git clone https://github.com/helloellinor/askeladden.git
cd askeladden

# Download dependencies
go mod download

# Build the project
go build ./cmd/askeladden

# Run tests and validation
go fmt ./...
go vet ./cmd/askeladden ./internal/...
```

## Code Standards

Please read and follow our [Coding Standards](CODING_STANDARDS.md) before contributing. Key points:

- **Language**: Use Norwegian (Nynorsk) for comments and user-facing features
- **Import grouping**: Standard library, third-party, local packages
- **Documentation**: Document all exported functions
- **Error handling**: Always handle errors explicitly

## Development Workflow

### 1. Setting Up Local Development

Create configuration files for local testing:

```bash
# Copy the example configuration
cp config/config-beta.yaml config-beta.yaml
cp config/config-beta.yaml secrets-beta.yaml

# Edit secrets-beta.yaml to add your Discord token and database credentials
```

### 2. Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/ny-funksjon
   ```

2. **Make your changes** following the coding standards

3. **Test your changes**:
   ```bash
   # Format code
   go fmt ./...
   
   # Check for issues
   go vet ./cmd/askeladden ./internal/...
   
   # Build and verify
   go build ./cmd/askeladden
   ls -lh askeladden  # Should be ~12MB
   
   # Test help functionality
   go build ./tools/test_help
   ./test_help | wc -l  # Should output ~12 lines
   ```

4. **Run the bot locally** (optional):
   ```bash
   CONFIG_FILE=config-beta.yaml SECRETS_FILE=secrets-beta.yaml ./askeladden
   ```

### 3. Submitting Changes

1. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Legg til ny funksjon for X"
   ```

2. **Push to your branch**:
   ```bash
   git push origin feature/ny-funksjon
   ```

3. **Create a Pull Request** with:
   - Clear description of changes
   - Screenshots if UI changes
   - Test results

## Project Structure

```
askeladden/
├── cmd/askeladden/          # Main application entry point
├── internal/                # Core application logic
│   ├── bot/                 # Bot core and handlers
│   ├── commands/            # Discord commands
│   ├── config/              # Configuration management
│   ├── database/            # Database operations
│   ├── permissions/         # Role-based permissions
│   └── reactions/           # Reaction handlers
├── config/                  # Configuration files
├── docs/                    # Documentation
├── tools/                   # Utility tools and scripts
├── CODING_STANDARDS.md      # Code style guidelines
└── CONTRIBUTING.md          # This file
```

## Adding New Features

### Adding a New Command

1. **Create the command file**:
   ```go
   // internal/commands/ny_kommando.go
   package commands

   func init() {
       commands["ny-kommando"] = Command{
           name:        "ny-kommando",
           description: "Beskriving av ny kommando",
           emoji:       "🆕",
           handler:     NyKommando,
           aliases:     []string{},
           adminOnly:   false,
       }
   }

   func NyKommando(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
       // Implementation
   }
   ```

2. **Follow naming conventions**:
   - Norwegian command names and descriptions
   - Clear, descriptive function names
   - Appropriate emoji representation

3. **Add validation and error handling**:
   ```go
   func NyKommando(s *discordgo.Session, m *discordgo.MessageCreate, bot *bot.Bot) {
       // Validate input
       parts := strings.SplitN(m.Content, " ", 2)
       if len(parts) < 2 {
           embed := services.CreateBotEmbed(s, "Feil", "Du må oppgi ein parameter!", services.EmbedTypeError)
           s.ChannelMessageSendEmbed(m.ChannelID, embed)
           return
       }

       // Process and handle errors
       err := processCommand(parts[1])
       if err != nil {
           log.Printf("Feil i ny kommando: %v", err)
           embed := services.CreateBotEmbed(s, "Feil", "Kunne ikkje utføre kommando", services.EmbedTypeError)
           s.ChannelMessageSendEmbed(m.ChannelID, embed)
           return
       }

       // Success response
       embed := services.CreateBotEmbed(s, "Suksess", "Kommando utført!", services.EmbedTypeSuccess)
       s.ChannelMessageSendEmbed(m.ChannelID, embed)
   }
   ```

### Adding Database Operations

1. **Add to interface** in `internal/database/database.go`:
   ```go
   type DatabaseIface interface {
       // Existing methods...
       NyOperasjon(param string) error
   }
   ```

2. **Implement the method**:
   ```go
   func (db *DB) NyOperasjon(param string) error {
       query := `INSERT INTO tabell (kolonne) VALUES (?)`
       _, err := db.conn.Exec(query, param)
       return err
   }
   ```

3. **Add appropriate error handling and logging**

### Adding Configuration Options

1. **Update config struct** in `internal/config/config.go`:
   ```go
   type Config struct {
       // Existing fields...
       NyInnstilling string `yaml:"ny_innstilling"`
   }
   ```

2. **Update example config** in `config/config-beta.yaml`:
   ```yaml
   ny_innstilling: "standardverdi"
   ```

3. **Document the new option** in README.md

## Testing

### Manual Testing Process

1. **Build verification**:
   ```bash
   rm -f askeladden test_help
   go build ./cmd/askeladden
   go build ./tools/test_help
   ```

2. **Code quality checks**:
   ```bash
   go fmt ./...
   go vet ./cmd/askeladden ./internal/...
   ```

3. **Functional testing**:
   ```bash
   ./test_help  # Should output help text
   # Test bot locally with your changes
   ```

### Testing Commands

Test new commands by:
1. Running the bot locally with test configuration
2. Using a test Discord server
3. Verifying proper error handling
4. Testing edge cases and invalid input

## Deployment

### Beta Testing
Use the beta environment for testing:
```bash
# Ensure you have config-beta.yaml and secrets-beta.yaml in root
./run-beta.sh
```

### Production Deployment
Production deployment is handled by maintainers using:
```bash
./build-and-deploy.sh
```

This script:
- Builds for Linux (`GOOS=linux GOARCH=amd64`)
- Deploys to `heim.bitraf.no`
- Starts in tmux session
- Logs to `askeladden.log`

## Code Review Process

### Before Submitting
- [ ] Code follows [Coding Standards](CODING_STANDARDS.md)
- [ ] All functions are documented
- [ ] Error handling is implemented
- [ ] Norwegian language is used appropriately
- [ ] Build succeeds and tests pass
- [ ] Imports are properly grouped

### Review Criteria
- **Functionality**: Does it work as intended?
- **Code quality**: Follows standards and best practices?
- **Documentation**: Adequate comments and documentation?
- **Testing**: Properly tested and validated?
- **Norwegian language**: Appropriate use of Norwegian vs English?

## Common Issues

### Build Issues
- **Module download fails**: Check internet connection and Go proxy settings
- **Build timeout**: Allow up to 60 seconds for initial build
- **Binary size wrong**: Expected ~12MB, check for missing dependencies

### Runtime Issues
- **Config file not found**: Check `CONFIG_FILE` and `SECRETS_FILE` environment variables
- **Database connection fails**: Verify database credentials and connectivity
- **Discord permissions**: Ensure bot has necessary permissions in test server

### Code Issues
- **Import organization**: Use `go fmt` and group imports properly
- **Go vet warnings**: Address all warnings before submitting
- **Language mixing**: Follow guidelines for Norwegian vs English usage

## Getting Help

- **Issues**: Use GitHub issues for bug reports and feature requests
- **Questions**: Ask in project discussions or Discord
- **Code review**: Tag maintainers in pull requests

## Recognition

Contributors who follow these guidelines and provide quality code will be acknowledged in:
- Git commit history
- Release notes
- Project documentation

Takk for at du bidrar til Askeladden! 🇳🇴