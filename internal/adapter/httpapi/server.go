package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ecommerce/catalog-service/internal/domain"
	"ecommerce/catalog-service/internal/port"
	"ecommerce/catalog-service/internal/usecase"

	"github.com/google/uuid"
)

type Server struct {
	categories  *usecase.CategoryService
	attributes  *usecase.AttributeService
	products    *usecase.ProductService
	variants    *usecase.VariantGenerationService
	productRepo port.ProductRepository
}

func NewServer(
	categories *usecase.CategoryService,
	attributes *usecase.AttributeService,
	products *usecase.ProductService,
	variants *usecase.VariantGenerationService,
	productRepo port.ProductRepository,
) *Server {
	return &Server{
		categories:  categories,
		attributes:  attributes,
		products:    products,
		variants:    variants,
		productRepo: productRepo,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /categories", s.handleCreateCategory)
	mux.HandleFunc("GET /categories/tree", s.handleGetCategoryTree)
	mux.HandleFunc("GET /categories/{categoryID}/children", s.handleListCategoryChildren)
	mux.HandleFunc("GET /categories/{categoryID}/products", s.handleListProductsByCategory)
	mux.HandleFunc("POST /categories/{categoryID}/attributes", s.handleCreateCategoryAttribute)
	mux.HandleFunc("GET /categories/{categoryID}/attributes", s.handleListCategoryAttributes)
	mux.HandleFunc("POST /attributes/{attributeID}/options", s.handleCreateAttributeOption)
	mux.HandleFunc("POST /products", s.handleCreateProduct)
	mux.HandleFunc("GET /products/{productID}", s.handleGetProduct)
	mux.HandleFunc("PUT /products/{productID}/attributes", s.handleSetProductAttributes)
	mux.HandleFunc("GET /products/{productID}/attributes", s.handleListProductAttributes)
	mux.HandleFunc("POST /products/{productID}/variants/generate", s.handleGenerateVariants)
	mux.HandleFunc("GET /products/{productID}/variants", s.handleListVariants)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type createCategoryRequest struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ParentID  *string `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}

func (s *Server) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var parentID *uuid.UUID
	if req.ParentID != nil && strings.TrimSpace(*req.ParentID) != "" {
		id, err := parseUUID(*req.ParentID, "parent_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		parentID = &id
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category, err := s.categories.CreateCategory(r.Context(), usecase.CreateCategoryInput{
		Name:      req.Name,
		Slug:      req.Slug,
		ParentID:  parentID,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toCategoryResponse(category))
}

func (s *Server) handleGetCategoryTree(w http.ResponseWriter, r *http.Request) {
	tree, err := s.categories.GetCategoryTree(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"categories": toCategoryTreeResponses(tree),
	})
}

func (s *Server) handleListCategoryChildren(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseUUID(r.PathValue("categoryID"), "categoryID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	children, err := s.categories.ListChildren(r.Context(), &categoryID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	res := make([]categoryResponse, 0, len(children))
	for _, child := range children {
		res = append(res, toCategoryResponse(child))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"parent_id": categoryID.String(),
		"children":  res,
	})
}

type createCategoryAttributeRequest struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Type          string `json:"type"`
	Required      bool   `json:"required"`
	IsVariantAxis bool   `json:"is_variant_axis"`
	IsFilterable  bool   `json:"is_filterable"`
	SortOrder     int    `json:"sort_order"`
}

func (s *Server) handleCreateCategoryAttribute(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseUUID(r.PathValue("categoryID"), "categoryID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req createCategoryAttributeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	attribute, err := s.attributes.CreateCategoryAttribute(r.Context(), usecase.CreateCategoryAttributeInput{
		CategoryID:    categoryID,
		Name:          req.Name,
		Code:          req.Code,
		Type:          req.Type,
		Required:      req.Required,
		IsVariantAxis: req.IsVariantAxis,
		IsFilterable:  req.IsFilterable,
		SortOrder:     req.SortOrder,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toAttributeResponse(attribute, nil))
}

type createAttributeOptionRequest struct {
	Value     string `json:"value"`
	Label     string `json:"label"`
	SortOrder int    `json:"sort_order"`
}

func (s *Server) handleCreateAttributeOption(w http.ResponseWriter, r *http.Request) {
	attributeID, err := parseUUID(r.PathValue("attributeID"), "attributeID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req createAttributeOptionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	option, err := s.attributes.AddAttributeOption(r.Context(), usecase.AddAttributeOptionInput{
		AttributeID: attributeID,
		Value:       req.Value,
		Label:       req.Label,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toOptionResponse(option))
}

func (s *Server) handleListCategoryAttributes(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseUUID(r.PathValue("categoryID"), "categoryID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	attributes, err := s.attributes.ListCategoryAttributes(r.Context(), categoryID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	responses := make([]attributeResponse, 0, len(attributes))
	for _, attribute := range attributes {
		options, err := s.attributes.ListAttributeOptions(r.Context(), attribute.ID)
		if err != nil {
			writeDomainError(w, err)
			return
		}
		responses = append(responses, toAttributeResponse(attribute, options))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"category_id": categoryID.String(),
		"attributes":  responses,
	})
}

type createProductRequest struct {
	Name                  string   `json:"name"`
	Slug                  string   `json:"slug"`
	Description           string   `json:"description"`
	PrimaryCategoryID     string   `json:"primary_category_id"`
	AdditionalCategoryIDs []string `json:"additional_category_ids"`
	BasePriceMinor        int64    `json:"base_price_minor"`
}

func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	primaryCategoryID, err := parseUUID(req.PrimaryCategoryID, "primary_category_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	additionalCategoryIDs, err := parseUUIDList(req.AdditionalCategoryIDs, "additional_category_ids")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := s.products.CreateProduct(r.Context(), usecase.CreateProductInput{
		Name:                  req.Name,
		Slug:                  req.Slug,
		Description:           req.Description,
		PrimaryCategoryID:     primaryCategoryID,
		AdditionalCategoryIDs: additionalCategoryIDs,
		BasePriceMinor:        req.BasePriceMinor,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toProductResponse(product))
}

func (s *Server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := parseUUID(r.PathValue("productID"), "productID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	product, err := s.products.GetProduct(r.Context(), productID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductResponse(product))
}

func (s *Server) handleListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseUUID(r.PathValue("categoryID"), "categoryID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	products, err := s.products.ListProductsByCategory(r.Context(), categoryID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	res := make([]productResponse, 0, len(products))
	for _, product := range products {
		res = append(res, toProductResponse(product))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"category_id": categoryID.String(),
		"count":       len(res),
		"products":    res,
	})
}

type setProductAttributesRequest struct {
	Values []setProductAttributeValueRequest `json:"values"`
}

type setProductAttributeValueRequest struct {
	AttributeID  string           `json:"attribute_id"`
	OptionID     *string          `json:"option_id,omitempty"`
	ValueText    *string          `json:"value_text,omitempty"`
	ValueNumber  *float64         `json:"value_number,omitempty"`
	ValueBoolean *bool            `json:"value_boolean,omitempty"`
	ValueJSON    *json.RawMessage `json:"value_json,omitempty"`
}

func (s *Server) handleSetProductAttributes(w http.ResponseWriter, r *http.Request) {
	productID, err := parseUUID(r.PathValue("productID"), "productID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req setProductAttributesRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	values := make([]domain.ProductAttributeValue, 0, len(req.Values))
	for idx, value := range req.Values {
		attributeID, err := parseUUID(value.AttributeID, fmt.Sprintf("values[%d].attribute_id", idx))
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		var optionID *uuid.UUID
		if value.OptionID != nil {
			parsed, err := parseUUID(*value.OptionID, fmt.Sprintf("values[%d].option_id", idx))
			if err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			optionID = &parsed
		}
		values = append(values, domain.ProductAttributeValue{
			ID:           uuid.New(),
			ProductID:    productID,
			AttributeID:  attributeID,
			OptionID:     optionID,
			ValueText:    value.ValueText,
			ValueNumber:  value.ValueNumber,
			ValueBoolean: value.ValueBoolean,
			ValueJSON:    value.ValueJSON,
		})
	}

	if err := s.products.SetProductAttributes(r.Context(), usecase.SetProductAttributesInput{
		ProductID: productID,
		Values:    values,
	}); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListProductAttributes(w http.ResponseWriter, r *http.Request) {
	productID, err := parseUUID(r.PathValue("productID"), "productID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	values, err := s.products.ListProductAttributeValues(r.Context(), productID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	res := make([]productAttributeValueResponse, 0, len(values))
	for _, value := range values {
		res = append(res, toProductAttributeValueResponse(value))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"product_id": productID.String(),
		"count":      len(res),
		"values":     res,
	})
}

type generateVariantsRequest struct {
	Axes            []variantAxisRequest `json:"axes"`
	BasePriceMinor  int64                `json:"base_price_minor"`
	InitialStockQty int64                `json:"initial_stock_qty"`
}

type variantAxisRequest struct {
	AttributeID string   `json:"attribute_id"`
	OptionIDs   []string `json:"option_ids"`
}

func (s *Server) handleGenerateVariants(w http.ResponseWriter, r *http.Request) {
	productID, err := parseUUID(r.PathValue("productID"), "productID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req generateVariantsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	axes := make([]usecase.VariantAxisInput, 0, len(req.Axes))
	for idx, axis := range req.Axes {
		attributeID, err := parseUUID(axis.AttributeID, fmt.Sprintf("axes[%d].attribute_id", idx))
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		optionIDs, err := parseUUIDList(axis.OptionIDs, fmt.Sprintf("axes[%d].option_ids", idx))
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		axes = append(axes, usecase.VariantAxisInput{
			AttributeID: attributeID,
			OptionIDs:   optionIDs,
		})
	}

	variants, err := s.variants.GenerateAndPersist(r.Context(), usecase.GenerateVariantsInput{
		ProductID:       productID,
		Axes:            axes,
		BasePriceMinor:  req.BasePriceMinor,
		InitialStockQty: req.InitialStockQty,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	res := make([]variantResponse, 0, len(variants))
	for _, variant := range variants {
		res = append(res, toVariantResponse(variant))
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"product_id": productID.String(),
		"count":      len(res),
		"variants":   res,
	})
}

func (s *Server) handleListVariants(w http.ResponseWriter, r *http.Request) {
	productID, err := parseUUID(r.PathValue("productID"), "productID")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	variants, err := s.productRepo.ListVariantsByProduct(r.Context(), productID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	res := make([]variantResponse, 0, len(variants))
	for _, variant := range variants {
		res = append(res, toVariantResponse(variant))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"product_id": productID.String(),
		"count":      len(res),
		"variants":   res,
	})
}

type categoryResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	ParentID  *string `json:"parent_id,omitempty"`
	Level     int     `json:"level"`
	Path      string  `json:"path"`
	SortOrder int     `json:"sort_order"`
	IsActive  bool    `json:"is_active"`
}

type categoryTreeNodeResponse struct {
	Category categoryResponse           `json:"category"`
	Children []categoryTreeNodeResponse `json:"children"`
}

func toCategoryResponse(c domain.Category) categoryResponse {
	var parentID *string
	if c.ParentID != nil {
		value := c.ParentID.String()
		parentID = &value
	}
	return categoryResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		ParentID:  parentID,
		Level:     c.Level,
		Path:      c.Path,
		SortOrder: c.SortOrder,
		IsActive:  c.IsActive,
	}
}

func toCategoryTreeResponses(nodes []usecase.CategoryTreeNode) []categoryTreeNodeResponse {
	out := make([]categoryTreeNodeResponse, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, categoryTreeNodeResponse{
			Category: toCategoryResponse(node.Category),
			Children: toCategoryTreeResponses(node.Children),
		})
	}
	return out
}

type optionResponse struct {
	ID          string `json:"id"`
	AttributeID string `json:"attribute_id"`
	Value       string `json:"value"`
	Label       string `json:"label"`
	SortOrder   int    `json:"sort_order"`
}

func toOptionResponse(o domain.AttributeOption) optionResponse {
	return optionResponse{
		ID:          o.ID.String(),
		AttributeID: o.AttributeID.String(),
		Value:       o.Value,
		Label:       o.Label,
		SortOrder:   o.SortOrder,
	}
}

type attributeResponse struct {
	ID            string           `json:"id"`
	CategoryID    string           `json:"category_id"`
	Name          string           `json:"name"`
	Code          string           `json:"code"`
	Type          string           `json:"type"`
	Required      bool             `json:"required"`
	IsVariantAxis bool             `json:"is_variant_axis"`
	IsFilterable  bool             `json:"is_filterable"`
	SortOrder     int              `json:"sort_order"`
	Options       []optionResponse `json:"options,omitempty"`
}

func toAttributeResponse(attribute domain.CategoryAttribute, options []domain.AttributeOption) attributeResponse {
	resp := attributeResponse{
		ID:            attribute.ID.String(),
		CategoryID:    attribute.CategoryID.String(),
		Name:          attribute.Name,
		Code:          attribute.Code,
		Type:          string(attribute.Type),
		Required:      attribute.Required,
		IsVariantAxis: attribute.IsVariantAxis,
		IsFilterable:  attribute.IsFilterable,
		SortOrder:     attribute.SortOrder,
	}
	if len(options) > 0 {
		resp.Options = make([]optionResponse, 0, len(options))
		for _, option := range options {
			resp.Options = append(resp.Options, toOptionResponse(option))
		}
	}
	return resp
}

type productResponse struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Slug              string   `json:"slug"`
	Description       string   `json:"description"`
	PrimaryCategoryID string   `json:"primary_category_id"`
	CategoryIDs       []string `json:"category_ids"`
	Status            string   `json:"status"`
	BasePriceMinor    int64    `json:"base_price_minor"`
}

func toProductResponse(product domain.Product) productResponse {
	categoryIDs := make([]string, 0, len(product.CategoryIDs))
	for _, id := range product.CategoryIDs {
		categoryIDs = append(categoryIDs, id.String())
	}
	return productResponse{
		ID:                product.ID.String(),
		Name:              product.Name,
		Slug:              product.Slug,
		Description:       product.Description,
		PrimaryCategoryID: product.PrimaryCategoryID.String(),
		CategoryIDs:       categoryIDs,
		Status:            string(product.Status),
		BasePriceMinor:    product.BasePriceMinor,
	}
}

type productAttributeValueResponse struct {
	ID           string           `json:"id"`
	ProductID    string           `json:"product_id"`
	AttributeID  string           `json:"attribute_id"`
	OptionID     *string          `json:"option_id,omitempty"`
	ValueText    *string          `json:"value_text,omitempty"`
	ValueNumber  *float64         `json:"value_number,omitempty"`
	ValueBoolean *bool            `json:"value_boolean,omitempty"`
	ValueJSON    *json.RawMessage `json:"value_json,omitempty"`
}

func toProductAttributeValueResponse(value domain.ProductAttributeValue) productAttributeValueResponse {
	var optionID *string
	if value.OptionID != nil {
		id := value.OptionID.String()
		optionID = &id
	}
	return productAttributeValueResponse{
		ID:           value.ID.String(),
		ProductID:    value.ProductID.String(),
		AttributeID:  value.AttributeID.String(),
		OptionID:     optionID,
		ValueText:    value.ValueText,
		ValueNumber:  value.ValueNumber,
		ValueBoolean: value.ValueBoolean,
		ValueJSON:    value.ValueJSON,
	}
}

type variantOptionValueResponse struct {
	AttributeID string `json:"attribute_id"`
	OptionID    string `json:"option_id"`
}

type variantResponse struct {
	ID             string                       `json:"id"`
	ProductID      string                       `json:"product_id"`
	SKU            string                       `json:"sku"`
	PriceMinor     int64                        `json:"price_minor"`
	StockQty       int64                        `json:"stock_qty"`
	Status         string                       `json:"status"`
	CombinationKey string                       `json:"combination_key"`
	Options        []variantOptionValueResponse `json:"options"`
}

func toVariantResponse(variant domain.ProductVariant) variantResponse {
	options := make([]variantOptionValueResponse, 0, len(variant.Options))
	for _, option := range variant.Options {
		options = append(options, variantOptionValueResponse{
			AttributeID: option.AttributeID.String(),
			OptionID:    option.OptionID.String(),
		})
	}
	return variantResponse{
		ID:             variant.ID.String(),
		ProductID:      variant.ProductID.String(),
		SKU:            variant.SKU,
		PriceMinor:     variant.PriceMinor,
		StockQty:       variant.StockQty,
		Status:         string(variant.Status),
		CombinationKey: variant.CombinationKey,
		Options:        options,
	}
}

func decodeJSON(r *http.Request, out any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	if decoder.More() {
		return errors.New("invalid request body: multiple JSON documents are not allowed")
	}
	return nil
}

func parseUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %w", field, err)
	}
	return id, nil
}

func parseUUIDList(values []string, field string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for idx, value := range values {
		id, err := parseUUID(value, fmt.Sprintf("%s[%d]", field, idx))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidCategory),
		errors.Is(err, domain.ErrInvalidCategoryHierarchy),
		errors.Is(err, domain.ErrDuplicateSlugUnderParent),
		errors.Is(err, domain.ErrInvalidAttribute),
		errors.Is(err, domain.ErrInvalidAttributeValue),
		errors.Is(err, domain.ErrInvalidProduct),
		errors.Is(err, domain.ErrInvalidVariantAxis),
		errors.Is(err, domain.ErrDuplicateVariantCombination):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
