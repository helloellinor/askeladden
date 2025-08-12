# Askeladden Discord Bot

Askeladden is a Discord bot for Norwegian language communities, written in Go, with features for grammar correction, daily questions, starboard, and community engagement.

**ALWAYS reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

## Working Effectively

### Prerequisites and Dependencies
- Go 1.24+ is required (project uses Go 1.24.1)
- MySQL database for data storage
- Discord bot token for production use
- All dependencies are managed via `go.mod` and downloaded automatically during build

### Build and Development Commands
- **Initial build**: `go build ./cmd/askeladden` 
  - First build takes ~13 seconds (downloads dependencies). NEVER CANCEL. Set timeout to 60+ seconds.
  - Subsequent builds take ~0.08 seconds
  - Produces binary named `askeladden`
- **Beta build**: `go build -o askeladden-beta ./cmd/askeladden`
- **Format code**: `go fmt ./...` - takes ~0.04 seconds, formats all Go files
- **Code analysis**: `go vet ./...` - takes ~2 seconds, may show warnings about tools/ directory (has multiple main functions, this is expected)
- **Clean build cache**: `go clean -cache` then rebuild (for timing tests)
- **Clean build**: `go clean` then rebuild

### Configuration Requirements
- Create `config/config.yaml` with Discord channel IDs, role IDs, and bot settings
- Create `config/secrets.yaml` with database credentials and Discord bot token
- See `config/config-beta.yaml` for configuration structure example
- Bot expects these files in `config/` directory relative to binary location
- Can override paths with environment variables: `CONFIG_FILE` and `SECRETS_FILE`

### Testing
- **No automated tests exist** - this project has no `*_test.go` files
- **Manual testing tools available**:
  - `go build -o test_help ./tools/test_help && ./test_help` - shows bot command help text
  - `go build -o test_match ./tools/test_match && ./test_match` - tests command matching (may crash without full bot context)
- **Validation approach**: Build successfully, test configuration loading, verify Discord connection (will fail without valid token)
- **Expected configuration test output**: "Could not connect to the database: dial tcp [::1]:3306: connect: connection refused"

### Creating Test Configuration Files
```bash
mkdir -p config
cat > config/config.yaml << 'EOF'
discord:
  prefix: "!"
  logChannelID: "123456789"
  defaultChannelID: "123456789"
approval:
  queueChannelID: "123456789"
  opplysarRoleID: "123456789"
bannedwords:
  approvalChannelID: "123456789"
  rettskrivarRoleID: "123456789"
grammar:
  channelID: "123456789"
starboard:
  channelID: "123456789"
  threshold: 3
  emoji: "⭐"
database:
  host: "localhost"
  port: 3306
  dbname: "test"
scheduler:
  enabled: false
  timezone: "Europe/Oslo"
reactions:
  question: "❓"
environment: "test"
EOF

cat > config/secrets.yaml << 'EOF'
database:
  user: "test"
  password: "test"
discord:
  token: "test_token"
EOF
```

### Running the Application
- **Production**: `./askeladden` (requires valid config files)
- **Beta mode**: Use `./run-beta.sh` script (requires `config-beta.yaml` and `secrets-beta.yaml`)
- **Testing startup**: Create minimal config files to test configuration loading
- Bot will fail at database connection or Discord authentication without valid credentials
- **Expected behavior**: Should load configuration, attempt database connection, then Discord connection

## Validation and Quality Assurance

### Pre-commit Checklist
- **ALWAYS run**: `go fmt ./...` before committing changes
- **ALWAYS run**: `go vet ./...` and fix any NEW warnings (ignore existing tools/ warnings)
- **ALWAYS build**: `go build ./cmd/askeladden` to ensure compilation succeeds
- **Test configuration loading**: Run `./askeladden` briefly to verify config parsing works

### Manual Testing Scenarios
- **Command help validation**: Build and run `go build -o test_help ./tools/test_help && ./test_help` to verify command help text
- **Build verification**: Ensure binary is created and has execute permissions
- **Configuration testing**: Create minimal config files and run `timeout 3s ./askeladden` to verify config parsing works (should fail at database connection)
- **Beta workflow testing**: Copy `config/config-beta.yaml` to `config-beta.yaml`, create `secrets-beta.yaml`, run `timeout 8s ./run-beta.sh`
- **Development workflow**: Make small code change → `go fmt ./...` → `go vet ./...` → `go build ./cmd/askeladden` → test startup
- **Integration testing**: If you have valid Discord credentials, test basic bot startup and shutdown

### Development Workflow
1. Make code changes
2. Run `go fmt ./...` 
3. Run `go vet ./...` and address any NEW issues
4. Build with `go build ./cmd/askeladden`
5. Test configuration loading
6. Commit changes

## Repository Structure and Key Files

### Essential Directories
- `cmd/askeladden/` - Main application entry point (`main.go`, `scheduler.go`)
- `internal/bot/` - Core bot implementation and Discord handling
- `internal/commands/` - Bot command implementations (12 commands total)
- `internal/database/` - Database connection and operations
- `internal/config/` - Configuration loading and management
- `config/` - Configuration files (not in git, created by user)
- `tools/` - Utility tools for testing and database operations

### Important Files
- `README.md` - Project documentation with setup and deployment instructions
- `go.mod`/`go.sum` - Go module definition and dependency management
- `build-and-deploy.sh` - Production deployment script for remote server deployment
- `run-beta.sh` - Beta environment script for testing
- `.gitignore` - Excludes config files and binaries from git

### Deployment
- **Production deployment**: Use `./build-and-deploy.sh` (deploys to `heim.bitraf.no`)
- **Manual deployment**: Build binary, copy config files, start with tmux
- **Local development**: Create local config files and run binary directly

## Common Tasks and Outputs

### Repository Root Contents
```
.git/
.gitignore
README.md
build-and-deploy.sh
cmd/
config/
go.mod
go.sum
internal/
run-beta.sh
tools/
```

### Available Bot Commands (from test_help output)
```
🗑️ tøm-db - Tømmer databasen for alle spørsmål
🤐 kjeften - Si Askeladden han må teie for dei upratsame
🏓 ping - Sjekk om boten svarar
👉 poke - Utløys dagens spørsmål for hand (kun admin)
❓ spør - Legg til eit spørsmål for daglege spørsmål
🔧 config - Vis gjeldende bot-konfigurasjon
✅ godkjenn - Godkjenn eit spørsmål for hand (kun for opplysarar)
👋 hei - Sei hei til boten
❓ hjelp - Syn denne hjelpemeldinga
📊 info - Syn opplysingar om boten
👋 loggav - Loggar av boten og avsluttar programmet
```

### Build Timing Expectations
- **First build**: 13 seconds (dependency download). NEVER CANCEL. Set timeout to 60+ seconds.
- **Subsequent builds**: 0.08 seconds
- **Go fmt**: 0.04 seconds
- **Go vet**: 2 seconds (includes tools directory warnings)
- **Tool builds**: 0.4 seconds each
- **Configuration loading test**: Instant (fails at database connection as expected)

## Troubleshooting

### Common Issues
- **"Could not load configuration"**: Create `config/config.yaml` and `config/secrets.yaml` files
- **Database connection failed**: Expected without local MySQL; verify config syntax instead
- **Discord authentication failed**: Expected without valid bot token; verify config loading works
- **Tools build errors**: Normal for `go test ./...` due to multiple main functions in tools/
- **Missing binary**: Run `go build ./cmd/askeladden` to create executable

### Working Around Limitations
- **No tests**: Focus on build validation and configuration testing
- **Requires external services**: Create minimal configs for syntax validation
- **Production deployment**: Use provided deployment script for `heim.bitraf.no` server
- **Beta testing**: Use beta configuration files and separate Discord bot instance

## Complete Development Workflow Validation

### End-to-End Testing Script
Run this complete workflow to validate everything works:

```bash
# 1. Clean start
go clean -cache
rm -f askeladden askeladden-beta test_help

# 2. Build (should take ~13 seconds first time)
time go build ./cmd/askeladden

# 3. Format and vet
time go fmt ./...
time go vet ./...  # Ignore tools/ warnings

# 4. Create test configuration
mkdir -p config
cat > config/config.yaml << 'EOF'
discord:
  prefix: "!"
  logChannelID: "123456789"
  defaultChannelID: "123456789"
approval:
  queueChannelID: "123456789"
  opplysarRoleID: "123456789"
bannedwords:
  approvalChannelID: "123456789"
  rettskrivarRoleID: "123456789"
grammar:
  channelID: "123456789"
starboard:
  channelID: "123456789"
  threshold: 3
  emoji: "⭐"
database:
  host: "localhost"
  port: 3306
  dbname: "test"
scheduler:
  enabled: false
  timezone: "Europe/Oslo"
reactions:
  question: "❓"
environment: "test"
EOF

cat > config/secrets.yaml << 'EOF'
database:
  user: "test"
  password: "test"
discord:
  token: "test_token"
EOF

# 5. Test configuration loading
timeout 3s ./askeladden  # Should fail at database connection

# 6. Test help tool
go build -o test_help ./tools/test_help
./test_help  # Should show Norwegian command list

# 7. Test beta workflow (if beta configs exist)
cp config/config-beta.yaml config-beta.yaml
cp config/secrets.yaml secrets-beta.yaml
timeout 8s ./run-beta.sh

# 8. Clean up
rm -f config/config.yaml config/secrets.yaml config-beta.yaml secrets-beta.yaml test_help
```

### Expected Results
- Build: No errors, ~13 seconds first time, ~0.08 seconds subsequent
- Format: No output, ~0.04 seconds
- Vet: Tools warnings (expected), ~2 seconds
- Config test: "Could not connect to the database: dial tcp [::1]:3306: connect: connection refused"
- Help tool: List of 11 Norwegian bot commands with emojis
- Beta script: Shows beta configuration details, fails at database connection, restores files

## Key Architecture Notes
- **Discord integration**: Uses `github.com/bwmarrin/discordgo` library
- **Database**: MySQL with custom connection management
- **Configuration**: YAML-based with environment variable overrides
- **Commands**: Modular command system with Norwegian language interface
- **Features**: Grammar correction, daily questions, starboard, role-based permissions
- **Deployment**: Remote server deployment with tmux session management