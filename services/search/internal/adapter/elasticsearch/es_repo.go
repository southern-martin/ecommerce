package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
)

const indexName = "products"

// ESSearchRepo implements domain.SearchRepository using Elasticsearch.
type ESSearchRepo struct {
	client *elasticsearch.Client
}

// NewESSearchRepo creates a new Elasticsearch-backed search repository.
// It connects to the given ES URL and ensures the "products" index exists
// with the correct mapping.
func NewESSearchRepo(esURL string) (*ESSearchRepo, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: failed to create client: %w", err)
	}

	// Verify connectivity.
	res, err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: ping failed: %w", err)
	}
	res.Body.Close()

	repo := &ESSearchRepo{client: client}

	if err := repo.ensureIndex(); err != nil {
		return nil, fmt.Errorf("elasticsearch: ensure index: %w", err)
	}

	return repo, nil
}

// ensureIndex creates the "products" index with proper mappings if it does
// not already exist.
func (r *ESSearchRepo) ensureIndex() error {
	// Check if index exists.
	res, err := r.client.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	res.Body.Close()

	if res.StatusCode == 200 {
		return nil // already exists
	}

	mapping := `{
  "mappings": {
    "properties": {
      "product_id":   {"type": "keyword"},
      "name":         {"type": "text", "analyzer": "standard", "fields": {"keyword": {"type": "keyword"}}},
      "slug":         {"type": "keyword"},
      "description":  {"type": "text", "analyzer": "standard"},
      "price_cents":  {"type": "long"},
      "currency":     {"type": "keyword"},
      "category_id":  {"type": "keyword"},
      "seller_id":    {"type": "keyword"},
      "image_url":    {"type": "keyword"},
      "rating":       {"type": "float"},
      "review_count": {"type": "integer"},
      "in_stock":     {"type": "boolean"},
      "tags":         {"type": "keyword"},
      "attributes":   {"type": "object", "enabled": false},
      "created_at":   {"type": "date"},
      "updated_at":   {"type": "date"},
      "suggest":      {"type": "completion"}
    }
  }
}`

	res, err = r.client.Indices.Create(
		indexName,
		r.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create index: %s", body)
	}

	log.Info().Str("index", indexName).Msg("elasticsearch index created")
	return nil
}

// esDocument is the internal representation stored in Elasticsearch.
type esDocument struct {
	ProductID   string            `json:"product_id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	PriceCents  int64             `json:"price_cents"`
	Currency    string            `json:"currency"`
	CategoryID  string            `json:"category_id"`
	SellerID    string            `json:"seller_id"`
	ImageURL    string            `json:"image_url"`
	Rating      float64           `json:"rating"`
	ReviewCount int               `json:"review_count"`
	InStock     bool              `json:"in_stock"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Suggest     suggestField      `json:"suggest"`
}

type suggestField struct {
	Input []string `json:"input"`
}

// Index upserts a product document into Elasticsearch using product_id as the
// document ID.
func (r *ESSearchRepo) Index(ctx context.Context, idx *domain.SearchIndex) error {
	doc := esDocument{
		ProductID:   idx.ProductID,
		Name:        idx.Name,
		Slug:        idx.Slug,
		Description: idx.Description,
		PriceCents:  idx.PriceCents,
		Currency:    idx.Currency,
		CategoryID:  idx.CategoryID,
		SellerID:    idx.SellerID,
		ImageURL:    idx.ImageURL,
		Rating:      idx.Rating,
		ReviewCount: idx.ReviewCount,
		InStock:     idx.InStock,
		Tags:        idx.Tags,
		Attributes:  idx.Attributes,
		CreatedAt:   idx.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   idx.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Suggest: suggestField{
			Input: buildSuggestInput(idx.Name),
		},
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: idx.ProductID,
		Body:       bytes.NewReader(body),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("elasticsearch: index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch: index error: %s", respBody)
	}

	log.Info().Str("product_id", idx.ProductID).Msg("product indexed in elasticsearch")
	return nil
}

// buildSuggestInput splits a product name into useful completion tokens.
func buildSuggestInput(name string) []string {
	inputs := []string{name}
	words := strings.Fields(name)
	if len(words) > 1 {
		for i := 1; i < len(words); i++ {
			inputs = append(inputs, strings.Join(words[i:], " "))
		}
	}
	return inputs
}

// Delete removes a product document from Elasticsearch by product_id.
func (r *ESSearchRepo) Delete(ctx context.Context, productID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: productID,
		Refresh:    "false",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("elasticsearch: delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elasticsearch: delete error: %s", respBody)
	}

	log.Info().Str("product_id", productID).Msg("product deleted from elasticsearch")
	return nil
}

// Search builds and executes an Elasticsearch query from the given filter,
// returning paginated results with relevance scores and a total count.
func (r *ESSearchRepo) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
	query := r.buildSearchQuery(filter)

	body, err := json.Marshal(query)
	if err != nil {
		return nil, 0, fmt.Errorf("elasticsearch: marshal query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(bytes.NewReader(body)),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("elasticsearch: search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		respBody, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("elasticsearch: search error: %s", respBody)
	}

	var esResp searchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, 0, fmt.Errorf("elasticsearch: decode response: %w", err)
	}

	total := esResp.Hits.Total.Value
	results := make([]domain.SearchResult, 0, len(esResp.Hits.Hits))
	for _, hit := range esResp.Hits.Hits {
		var doc esDocument
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			log.Warn().Err(err).Str("id", hit.ID).Msg("failed to unmarshal hit")
			continue
		}

		results = append(results, domain.SearchResult{
			ID:          hit.ID,
			ProductID:   doc.ProductID,
			Name:        doc.Name,
			Slug:        doc.Slug,
			Description: doc.Description,
			PriceCents:  doc.PriceCents,
			Currency:    doc.Currency,
			ImageURL:    doc.ImageURL,
			SellerID:    doc.SellerID,
			CategoryID:  doc.CategoryID,
			Rating:      doc.Rating,
			ReviewCount: doc.ReviewCount,
			InStock:     doc.InStock,
			Score:       hit.Score,
		})
	}

	return results, total, nil
}

// buildSearchQuery constructs the Elasticsearch query JSON from the filter.
func (r *ESSearchRepo) buildSearchQuery(filter domain.SearchFilter) map[string]interface{} {
	// Build the bool query.
	must := make([]map[string]interface{}, 0)
	filterClauses := make([]map[string]interface{}, 0)

	// Text query: multi-match on name (boosted) and description.
	if filter.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  filter.Query,
				"fields": []string{"name^3", "description"},
			},
		})
	}

	// Category filter.
	if filter.CategoryID != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{"category_id": filter.CategoryID},
		})
	}

	// Seller filter.
	if filter.SellerID != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{"seller_id": filter.SellerID},
		})
	}

	// In-stock filter.
	if filter.InStock != nil {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{"in_stock": *filter.InStock},
		})
	}

	// Price range filter.
	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		rangeQ := map[string]interface{}{}
		if filter.MinPrice > 0 {
			rangeQ["gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			rangeQ["lte"] = filter.MaxPrice
		}
		filterClauses = append(filterClauses, map[string]interface{}{
			"range": map[string]interface{}{"price_cents": rangeQ},
		})
	}

	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filterClauses) > 0 {
		boolQuery["filter"] = filterClauses
	}

	// If no must clauses and no text query, match all.
	var queryClause map[string]interface{}
	if len(boolQuery) > 0 {
		queryClause = map[string]interface{}{"bool": boolQuery}
	} else {
		queryClause = map[string]interface{}{"match_all": map[string]interface{}{}}
	}

	// Pagination.
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	from := (page - 1) * pageSize

	esQuery := map[string]interface{}{
		"query": queryClause,
		"from":  from,
		"size":  pageSize,
	}

	// Sorting.
	if filter.SortBy != "" {
		direction := "asc"
		if filter.SortOrder == "desc" {
			direction = "desc"
		}

		var sortField string
		switch filter.SortBy {
		case "price":
			sortField = "price_cents"
		case "rating":
			sortField = "rating"
		case "name":
			sortField = "name.keyword"
		case "created_at":
			sortField = "created_at"
		default:
			sortField = "created_at"
		}

		esQuery["sort"] = []map[string]interface{}{
			{sortField: map[string]interface{}{"order": direction}},
		}
	}

	return esQuery
}

// Suggest uses the completion suggester for autocomplete.
func (r *ESSearchRepo) Suggest(ctx context.Context, query string, limit int) ([]domain.SearchSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	suggestQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"product-suggest": map[string]interface{}{
				"prefix": query,
				"completion": map[string]interface{}{
					"field": "suggest",
					"size":  limit,
				},
			},
		},
		"_source": []string{"product_id", "name"},
	}

	body, err := json.Marshal(suggestQuery)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: marshal suggest query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: suggest request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		respBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch: suggest error: %s", respBody)
	}

	var esResp suggestResponse
	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, fmt.Errorf("elasticsearch: decode suggest response: %w", err)
	}

	options := esResp.Suggest["product-suggest"]
	suggestions := make([]domain.SearchSuggestion, 0)
	for _, entry := range options {
		for _, opt := range entry.Options {
			var doc esDocument
			if err := json.Unmarshal(opt.Source, &doc); err != nil {
				continue
			}
			suggestions = append(suggestions, domain.SearchSuggestion{
				Text:      opt.Text,
				Type:      "product",
				ProductID: doc.ProductID,
			})
		}
	}

	return suggestions, nil
}

// --- Elasticsearch response types ---

type searchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []searchHit `json:"hits"`
	} `json:"hits"`
}

type searchHit struct {
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

type suggestResponse struct {
	Suggest map[string][]suggestEntry `json:"suggest"`
}

type suggestEntry struct {
	Options []suggestOption `json:"options"`
}

type suggestOption struct {
	Text   string          `json:"text"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}
