package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// DictionaryService handles word lookups from ordbøkene.no
type DictionaryService struct {
	httpClient *http.Client
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WordInfo contains information about a word from ordbøkene.no
type WordInfo struct {
	Word       string
	Definition string
	Bøying     string // Inflection information
	URL        string
	Found      bool
	Error      string
}

// LookupWord searches for a word on ordbøkene.no and extracts bøying information
func (ds *DictionaryService) LookupWord(word string) (*WordInfo, error) {
	// Clean the word input
	cleanWord := strings.TrimSpace(strings.ToLower(word))
	if cleanWord == "" {
		return &WordInfo{
			Word:  word,
			Found: false,
			Error: "Tomt ord",
		}, nil
	}

	// Construct the URL - using the search endpoint first
	searchURL := fmt.Sprintf("https://ordbokene.no/nob/nn/%s", url.QueryEscape(cleanWord))

	log.Printf("[DICTIONARY] Looking up word: %s at %s", cleanWord, searchURL)

	// Make HTTP request
	resp, err := ds.httpClient.Get(searchURL)
	if err != nil {
		return &WordInfo{
			Word:  word,
			Found: false,
			Error: fmt.Sprintf("Feil ved oppslag: %v", err),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &WordInfo{
			Word:  word,
			Found: false,
			Error: fmt.Sprintf("HTTP-feil: %d", resp.StatusCode),
			URL:   searchURL,
		}, nil
	}

	// Read the response body
	body := make([]byte, 50000) // Limit to prevent excessive memory usage
	n, _ := resp.Body.Read(body)
	content := string(body[:n])

	// Parse the content to extract bøying information
	wordInfo := &WordInfo{
		Word:  cleanWord,
		URL:   searchURL,
		Found: true,
	}

	// Extract definition (look for the first definition in the content)
	if def := ds.extractDefinition(content); def != "" {
		wordInfo.Definition = def
	}

	// Extract bøying (inflection) information
	if bøying := ds.extractBøying(content); bøying != "" {
		wordInfo.Bøying = bøying
	} else {
		wordInfo.Bøying = "Ingen bøying funnet"
	}

	// If we couldn't find any useful information, mark as not found
	if wordInfo.Definition == "" && wordInfo.Bøying == "Ingen bøying funnet" {
		wordInfo.Found = false
		wordInfo.Error = "Ordet vart ikkje funnet i ordbøkene"
	}

	return wordInfo, nil
}

// extractDefinition attempts to extract the first definition from the HTML content
func (ds *DictionaryService) extractDefinition(content string) string {
	// Look for definition patterns in Norwegian dictionary format
	// This is a simplified extraction - in practice, you'd want more robust HTML parsing

	// Pattern for definition sections
	defPattern := regexp.MustCompile(`(?i)<div[^>]*class="[^"]*definition[^"]*"[^>]*>(.*?)</div>`)
	matches := defPattern.FindStringSubmatch(content)

	if len(matches) > 1 {
		def := ds.cleanHTML(matches[1])
		if len(def) > 200 {
			def = def[:200] + "..."
		}
		return def
	}

	// Alternative pattern for meaning/betydning
	meaningPattern := regexp.MustCompile(`(?i)<div[^>]*class="[^"]*meaning[^"]*"[^>]*>(.*?)</div>`)
	matches = meaningPattern.FindStringSubmatch(content)

	if len(matches) > 1 {
		def := ds.cleanHTML(matches[1])
		if len(def) > 200 {
			def = def[:200] + "..."
		}
		return def
	}

	return ""
}

// extractBøying attempts to extract inflection information from the HTML content
func (ds *DictionaryService) extractBøying(content string) string {
	// Look for bøying/inflection patterns
	bøyingPatterns := []string{
		`(?i)<div[^>]*class="[^"]*bøying[^"]*"[^>]*>(.*?)</div>`,
		`(?i)<div[^>]*class="[^"]*inflection[^"]*"[^>]*>(.*?)</div>`,
		`(?i)<div[^>]*class="[^"]*conjugation[^"]*"[^>]*>(.*?)</div>`,
		`(?i)<span[^>]*class="[^"]*bøying[^"]*"[^>]*>(.*?)</span>`,
	}

	for _, pattern := range bøyingPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)

		if len(matches) > 1 {
			bøying := ds.cleanHTML(matches[1])
			if bøying != "" {
				return bøying
			}
		}
	}

	// Look for word forms in tables or lists
	formPattern := regexp.MustCompile(`(?i)(<table[^>]*>.*?</table>|<ul[^>]*>.*?</ul>|<ol[^>]*>.*?</ol>)`)
	matches := formPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 && strings.Contains(strings.ToLower(match[1]), "bøy") {
			forms := ds.extractWordForms(match[1])
			if forms != "" {
				return forms
			}
		}
	}

	return ""
}

// extractWordForms extracts word forms from table/list HTML
func (ds *DictionaryService) extractWordForms(htmlContent string) string {
	// Extract text content from table cells or list items
	cellPattern := regexp.MustCompile(`(?i)<(?:td|li)[^>]*>(.*?)</(?:td|li)>`)
	matches := cellPattern.FindAllStringSubmatch(htmlContent, -1)

	var forms []string
	for _, match := range matches {
		if len(match) > 1 {
			form := ds.cleanHTML(match[1])
			if form != "" && len(form) < 50 { // Reasonable length for word forms
				forms = append(forms, form)
			}
		}
	}

	if len(forms) > 0 {
		// Limit to first 5 forms to avoid too long messages
		if len(forms) > 5 {
			forms = forms[:5]
		}
		return strings.Join(forms, ", ")
	}

	return ""
}

// cleanHTML removes HTML tags and cleans up text content
func (ds *DictionaryService) cleanHTML(html string) string {
	// Remove HTML tags
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	text := tagPattern.ReplaceAllString(html, "")

	// Clean up whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Remove common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")

	return text
}

// CreateWordLookupEmbed creates a Discord embed for word lookup results
func (ds *DictionaryService) CreateWordLookupEmbed(wordInfo *WordInfo) *discordgo.MessageEmbed {
	if !wordInfo.Found {
		return &discordgo.MessageEmbed{
			Title:       "📚 Ordoppslag",
			Description: fmt.Sprintf("**Ord:** %s\n\n❌ %s", wordInfo.Word, wordInfo.Error),
			Color:       ColorError,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Kilde: ordbøkene.no",
			},
		}
	}

	fields := []*discordgo.MessageEmbedField{}

	if wordInfo.Definition != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "📖 Tyding",
			Value:  wordInfo.Definition,
			Inline: false,
		})
	}

	if wordInfo.Bøying != "" && wordInfo.Bøying != "Ingen bøying funnet" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "🔄 Bøying",
			Value:  wordInfo.Bøying,
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Ordoppslag",
		Description: fmt.Sprintf("**Ord:** %s", wordInfo.Word),
		Color:       ColorSuccess,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Kilde: ordbøkene.no",
		},
	}

	if wordInfo.URL != "" {
		embed.Description += fmt.Sprintf("\n\n[Sjå på ordbøkene.no](%s)", wordInfo.URL)
	}

	return embed
}
