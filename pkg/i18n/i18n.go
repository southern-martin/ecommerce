package i18n

import (
	"fmt"
	"sync"

	"golang.org/x/text/language"
)

// Supported language tags.
var supportedLanguages = []language.Tag{
	language.English,            // en
	language.Spanish,            // es
	language.French,             // fr
	language.Chinese,            // zh
	language.Arabic,             // ar
	language.German,             // de
	language.Japanese,           // ja
	language.Portuguese,         // pt
	language.Korean,             // ko
	language.Italian,            // it
}

// Bundle holds translations for multiple languages and provides thread-safe
// lookup with BCP47 language tag matching via golang.org/x/text/language.
type Bundle struct {
	mu           sync.RWMutex
	translations map[string]map[string]string
	matcher      language.Matcher
}

// NewBundle creates a new, empty Bundle with a language matcher configured
// for the supported languages.
func NewBundle() *Bundle {
	return &Bundle{
		translations: make(map[string]map[string]string),
		matcher:      language.NewMatcher(supportedLanguages),
	}
}

// LoadLanguage registers translations for a single language. If translations
// already exist for the language, the new entries are merged in, with new
// values overwriting existing ones for the same key.
func (b *Bundle) LoadLanguage(lang string, translations map[string]string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.translations[lang] == nil {
		b.translations[lang] = make(map[string]string, len(translations))
	}
	for k, v := range translations {
		b.translations[lang][k] = v
	}
}

// T translates the given key for the specified language. If the key is not
// found in the requested language, it falls back to English. If the key is
// not found at all, the key itself is returned.
//
// Optional args are passed to fmt.Sprintf for parameterized messages.
func (b *Bundle) T(lang, key string, args ...interface{}) string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Match the requested language to the closest supported language.
	matched := b.matchLanguage(lang)

	// Try the matched language.
	if msgs, ok := b.translations[matched]; ok {
		if msg, ok := msgs[key]; ok {
			return b.format(msg, args...)
		}
	}

	// Fall back to English.
	if matched != "en" {
		if msgs, ok := b.translations["en"]; ok {
			if msg, ok := msgs[key]; ok {
				return b.format(msg, args...)
			}
		}
	}

	// Return the key itself as a last resort.
	return key
}

// SetupDefaults loads the embedded default translations for all supported
// languages. This should be called once at application startup.
func (b *Bundle) SetupDefaults() {
	for lang, translations := range defaultTranslations {
		b.LoadLanguage(lang, translations)
	}
}

// matchLanguage uses BCP47 matching to find the best supported language
// for the given language string. Returns the language code as a string.
func (b *Bundle) matchLanguage(lang string) string {
	tag, err := language.Parse(lang)
	if err != nil {
		return "en"
	}
	_, idx, _ := b.matcher.Match(tag)
	if idx < 0 || idx >= len(supportedLanguages) {
		return "en"
	}
	base, _ := supportedLanguages[idx].Base()
	return base.String()
}

// format applies fmt.Sprintf if args are provided, otherwise returns the
// message unchanged.
func (b *Bundle) format(msg string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}
