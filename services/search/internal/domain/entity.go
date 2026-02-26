package domain

import "time"

// SearchResult represents a single search result item.
type SearchResult struct {
	ID          string
	ProductID   string
	Name        string
	Slug        string
	Description string
	PriceCents  int64
	Currency    string
	ImageURL    string
	SellerID    string
	CategoryID  string
	Rating      float64
	ReviewCount int
	InStock     bool
	Score       float64
	CreatedAt   time.Time
}

// SearchFilter holds the parameters for a search query.
type SearchFilter struct {
	Query      string
	CategoryID string
	MinPrice   int64
	MaxPrice   int64
	InStock    *bool
	SellerID   string
	SortBy     string
	SortOrder  string
	Page       int
	PageSize   int
}

// SearchIndex represents a product document in the search index.
type SearchIndex struct {
	ID          string
	ProductID   string
	Name        string
	Slug        string
	Description string
	PriceCents  int64
	Currency    string
	CategoryID  string
	SellerID    string
	ImageURL    string
	Rating      float64
	ReviewCount int
	InStock     bool
	Tags        []string
	Attributes  map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SearchSuggestion represents a search autocomplete suggestion.
type SearchSuggestion struct {
	Text      string
	Type      string
	ProductID string
}
