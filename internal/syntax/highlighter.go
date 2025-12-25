package syntax

import (
	"hash/fnv"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	_ "github.com/javanhut/vem/internal/syntax/lexers"
)

// Token represents a syntax-highlighted token with its text, type, and color.
type Token struct {
	Text  string
	Type  chroma.TokenType
	Style *chroma.Style
}

// HighlightedLine represents a cached highlighted line.
type HighlightedLine struct {
	Tokens []Token
	Hash   uint64
}

// Highlighter provides syntax highlighting for a specific file/language.
type Highlighter struct {
	lexer     chroma.Lexer
	style     *chroma.Style
	cache     map[int]*HighlightedLine
	formatter *chroma.Formatter
	enabled   bool
}

// NewHighlighter creates a new highlighter for the given file path.
// It auto-detects the language from the file extension.
func NewHighlighter(filePath string) *Highlighter {
	// Try to match lexer by filename
	lexer := lexers.Match(filePath)

	// If no match, try by extension
	if lexer == nil {
		ext := filepath.Ext(filePath)
		if ext != "" {
			lexer = lexers.Get(ext[1:]) // Remove the leading dot
		}
	}

	// Fallback to plain text
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Get the style (theme)
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	return &Highlighter{
		lexer:   chroma.Coalesce(lexer),
		style:   style,
		cache:   make(map[int]*HighlightedLine),
		enabled: true,
	}
}

// NewPlainHighlighter creates a highlighter without syntax highlighting (plain text).
func NewPlainHighlighter() *Highlighter {
	return &Highlighter{
		lexer:   lexers.Fallback,
		style:   styles.Fallback,
		cache:   make(map[int]*HighlightedLine),
		enabled: false,
	}
}

// HighlightLine tokenizes and highlights a single line of text.
// It uses caching to avoid re-tokenizing unchanged lines.
func (h *Highlighter) HighlightLine(lineNum int, text string) []Token {
	// If highlighting is disabled, return plain text token
	if !h.enabled {
		return []Token{{Text: text, Type: chroma.Text, Style: h.style}}
	}

	// Check cache first
	textHash := hashString(text)
	if cached, ok := h.cache[lineNum]; ok {
		if cached.Hash == textHash {
			return cached.Tokens
		}
	}

	// Tokenize the line
	tokens := make([]Token, 0)

	iterator, err := h.lexer.Tokenise(nil, text)
	if err != nil {
		// On error, return plain text
		return []Token{{Text: text, Type: chroma.Text, Style: h.style}}
	}

	// Convert chroma tokens to our Token type
	for _, token := range iterator.Tokens() {
		if token.Value == "" {
			continue
		}
		tokens = append(tokens, Token{
			Text:  token.Value,
			Type:  token.Type,
			Style: h.style,
		})
	}

	// If no tokens were generated, return the whole line as plain text
	if len(tokens) == 0 {
		tokens = []Token{{Text: text, Type: chroma.Text, Style: h.style}}
	}

	// Cache the result
	h.cache[lineNum] = &HighlightedLine{
		Tokens: tokens,
		Hash:   textHash,
	}

	return tokens
}

// InvalidateLine removes a line from the cache (called when line is edited).
func (h *Highlighter) InvalidateLine(lineNum int) {
	delete(h.cache, lineNum)
}

// InvalidateAll clears the entire cache (called on major changes).
func (h *Highlighter) InvalidateAll() {
	h.cache = make(map[int]*HighlightedLine)
}

// SetTheme changes the color theme.
func (h *Highlighter) SetTheme(themeName string) {
	style := styles.Get(themeName)
	if style != nil {
		h.style = style
		h.InvalidateAll() // Re-highlight with new theme
	}
}

// GetThemeName returns the current theme name.
func (h *Highlighter) GetThemeName() string {
	return h.style.Name
}

// SetEnabled enables or disables syntax highlighting.
func (h *Highlighter) SetEnabled(enabled bool) {
	if h.enabled != enabled {
		h.enabled = enabled
		h.InvalidateAll()
	}
}

// IsEnabled returns whether syntax highlighting is enabled.
func (h *Highlighter) IsEnabled() bool {
	return h.enabled
}

// GetLanguage returns the name of the detected language.
func (h *Highlighter) GetLanguage() string {
	if h.lexer != nil {
		config := h.lexer.Config()
		if config != nil {
			return config.Name
		}
	}
	return "Plain Text"
}

// ListAvailableThemes returns a list of all available color themes.
func ListAvailableThemes() []string {
	themes := styles.Names()

	// Sort alphabetically for consistency
	result := make([]string, len(themes))
	copy(result, themes)

	return result
}

// hashString computes a hash of the string for cache comparison.
func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// ShouldHighlight determines if a file should have syntax highlighting based on its path.
func ShouldHighlight(filePath string) bool {
	if filePath == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	// Skip very large files or binary files
	// (This is a simple heuristic; you might want to add file size checking)
	nonTextExtensions := []string{
		".bin", ".exe", ".so", ".dylib", ".dll",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico",
		".mp3", ".mp4", ".avi", ".mov", ".mkv",
		".zip", ".tar", ".gz", ".7z", ".rar",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
	}

	for _, badExt := range nonTextExtensions {
		if ext == badExt {
			return false
		}
	}

	return true
}
