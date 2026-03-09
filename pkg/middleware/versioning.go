// Package middleware provides HTTP middleware for the ecommerce platform.
//
// API Versioning Strategy
//
// This middleware implements API versioning with the following approach:
//
//   - URL-based versioning is the primary mechanism (/api/v1, /api/v2, etc.).
//     Clients specify the version directly in the URL path.
//
//   - Accept header versioning is a secondary mechanism. Clients may send
//     Accept: application/vnd.ecommerce.v2+json to request a specific version.
//     The URL path version takes precedence when both are present.
//
//   - Deprecation headers follow RFC 8594 (Sunset header). When a version is
//     deprecated, responses include Deprecation, Sunset, and Link headers to
//     guide clients toward the successor version.
//
//   - The resolved API version is stored in the Gin context under the key
//     "api_version" (e.g. "v1", "v2") for use by downstream handlers.
package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyAPIVersion is the context key where the resolved API version is stored.
	ContextKeyAPIVersion = "api_version"

	// acceptHeaderPrefix is the vendor media-type prefix used for header-based versioning.
	acceptHeaderPrefix = "application/vnd.ecommerce."
)

// acceptVersionRe matches "application/vnd.ecommerce.v2+json" and captures "v2".
var acceptVersionRe = regexp.MustCompile(`application/vnd\.ecommerce\.(v\d+)\+json`)

// pathVersionRe matches "/api/v1/..." and captures "v1".
var pathVersionRe = regexp.MustCompile(`/api/(v\d+)(?:/|$)`)

// APIVersion extracts the API version from the URL path or Accept header and
// stores it in the Gin context under ContextKeyAPIVersion.
//
// Resolution order:
//  1. URL path segment (e.g. /api/v1/auth/login -> "v1")
//  2. Accept header (e.g. application/vnd.ecommerce.v2+json -> "v2")
//  3. Default: "v1"
func APIVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := ""

		// 1. Try URL path.
		if matches := pathVersionRe.FindStringSubmatch(c.Request.URL.Path); len(matches) > 1 {
			version = matches[1]
		}

		// 2. Fall back to Accept header.
		if version == "" {
			accept := c.GetHeader("Accept")
			if strings.Contains(accept, acceptHeaderPrefix) {
				if matches := acceptVersionRe.FindStringSubmatch(accept); len(matches) > 1 {
					version = matches[1]
				}
			}
		}

		// 3. Default to v1.
		if version == "" {
			version = "v1"
		}

		c.Set(ContextKeyAPIVersion, version)
		c.Next()
	}
}

// DeprecateVersion marks an API version as deprecated by setting standard HTTP
// headers on every response. This middleware should be applied to route groups
// that serve a deprecated version (e.g. the v1 group once v2 is available).
//
// Parameters:
//   - sunset: the date after which the version may be removed, in RFC 7231
//     format (e.g. "Sat, 01 Mar 2025 00:00:00 GMT").
//   - alternative: the base path of the successor version (e.g. "/api/v2").
//
// Headers set:
//
//	Deprecation: true
//	Sunset: <sunset>
//	Link: </api/v2/...>; rel="successor-version"
//	X-API-Deprecated: true
func DeprecateVersion(sunset string, alternative string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Deprecation", "true")
		c.Header("Sunset", sunset)
		c.Header("X-API-Deprecated", "true")

		// Build the successor link using the current request path rewritten to
		// the alternative version prefix.
		successorPath := alternative
		if matches := pathVersionRe.FindStringSubmatch(c.Request.URL.Path); len(matches) > 1 {
			// Replace /api/vN with the alternative prefix in the original path.
			oldPrefix := "/api/" + matches[1]
			successorPath = strings.Replace(c.Request.URL.Path, oldPrefix, alternative, 1)
		}
		c.Header("Link", fmt.Sprintf("<%s>; rel=\"successor-version\"", successorPath))

		c.Next()
	}
}

// VersionNotFound returns a 404 JSON error for requests targeting an
// unsupported API version. Apply this as a catch-all handler for version
// prefixes that are not (or no longer) served.
//
// The supportedVersions parameter lists the versions currently available
// (e.g. []string{"v1"}).
func VersionNotFound(supportedVersions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":              "API version not supported",
			"supported_versions": supportedVersions,
		})
	}
}
