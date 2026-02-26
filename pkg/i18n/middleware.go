package i18n

import (
	"github.com/gin-gonic/gin"
)

const (
	// contextKeyLanguage is the gin context key used to store the resolved language.
	contextKeyLanguage = "i18n_language"
)

// GinMiddleware returns a gin.HandlerFunc that extracts the Accept-Language
// header from incoming requests, matches it to the closest supported language
// using the bundle's BCP47 matcher, and stores the resolved language code
// in the gin context for downstream handlers.
func GinMiddleware(bundle *Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		lang := "en"

		if acceptLang != "" {
			lang = bundle.matchLanguage(acceptLang)
		}

		c.Set(contextKeyLanguage, lang)
		c.Next()
	}
}

// GetLanguage retrieves the current language from the gin context. If no
// language has been set (i.e., the middleware was not applied), it defaults
// to "en".
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get(contextKeyLanguage); exists {
		if s, ok := lang.(string); ok && s != "" {
			return s
		}
	}
	return "en"
}
