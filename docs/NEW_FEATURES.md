# New Features Documentation

This update adds several new features to Askeladden for enhanced Norwegian language support and community management.

## New Commands

### Dictionary and Language Support
- **`?ord <word>`** - Look up words in ordbøkene.no with inflection (bøying) information
  - Example: `?ord hus` shows word definition and inflection forms
  - Supports aliases: `?ordbok`, `?lookup`

- **`?kontekst <word> <context>`** - Analyze words in context for grammar checking
  - Example: `?kontekst huse I huset bur det mange folk`
  - Provides basic bokmål vs nynorsk suggestions

### Administration
- **`?rolle <add/remove> <@user> <role>`** - Manage user roles (admin only)
  - Example: `?rolle add @user opplysar`
  - Supports roles: `opplysar`, `rettskrivar`

- **`?reglar`** - Display server rules and guidelines
  - Configurable content in `rules` section of config

- **`?embed`** - Create custom embeds via DM interface (opplysar only)
  - Interactive DM-based workflow for creating embeds
  - Supports custom colors, titles, and content

## New Features

### Welcome Messages
- Automatically welcomes new guild members
- Configurable welcome channel and message
- Template variables: `{user}`, `{username}`

### Enhanced Dictionary Integration
- Real-time word lookup from ordbøkene.no
- Extracts definition and bøying (inflection) information
- Handles Norwegian-specific grammar patterns

### Role Management
- Streamlined role assignment for moderators
- Supports existing permission system
- Integrated with approval workflows

### DM-Based Embed Creation
- Step-by-step embed creation process
- Preview functionality before publishing
- Custom color support with hex codes
- Logging of embed creation actions

## Configuration

Add these sections to your `config.yaml`:

```yaml
# Welcome messages
welcome:
  enabled: true
  channelID: "YOUR_WELCOME_CHANNEL_ID"
  message: "Velkommen til serveren, {user}! 🎉"

# Rules information  
rules:
  title: "📋 Serverreglar"
  content: "Your server rules here..."
```

## Technical Changes

### New Intents
- Added `IntentsGuildMembers` for welcome message functionality
- Enhanced DM handling for embed creation

### New Services
- `DictionaryService` - ordbøkene.no integration
- `EmbedBuilderService` - DM-based embed creation
- Global session management for embed workflows

### Command System Enhancement
- All new commands follow existing patterns
- Proper permission checking for admin commands
- Norwegian language interface maintained

## Usage Examples

### Dictionary Lookup
```
?ord hus
# Returns definition and inflection information

?kontekst huse I huset bur det mange folk  
# Analyzes grammar in context
```

### Role Management
```
?rolle add @johndoe opplysar
# Adds opplysar role to user

?rolle remove @johndoe rettskrivar
# Removes rettskrivar role from user
```

### Embed Creation
```
?embed
# Starts DM-based embed creation workflow
# Follow prompts in DM to create custom embeds
```

## Future Enhancements

The foundation is now in place for:
- Advanced context-based banned word detection
- AI-powered grammar checking
- Enhanced ordbøkene.no parsing
- More sophisticated role management