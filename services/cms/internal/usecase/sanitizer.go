package usecase

import "github.com/microcosm-cc/bluemonday"

var (
	// htmlPolicy allows safe HTML subset (links, formatting, images)
	htmlPolicy = bluemonday.UGCPolicy()
	// strictPolicy strips ALL HTML — for plain text fields
	strictPolicy = bluemonday.StrictPolicy()
)

// sanitizeHTML sanitizes user-provided HTML content, allowing safe formatting tags.
func sanitizeHTML(input string) string {
	return htmlPolicy.Sanitize(input)
}

// sanitizeText strips all HTML tags from text.
func sanitizeText(input string) string {
	return strictPolicy.Sanitize(input)
}
