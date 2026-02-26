package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"ecommerce/catalog-service/internal/domain"

	"github.com/google/uuid"
)

type productCategoryState struct {
	PrimaryCategoryID uuid.UUID
	CategoryIDs       []uuid.UUID
}

type Store struct {
	mu sync.RWMutex

	categories     map[uuid.UUID]domain.Category
	categoryBySlug map[string]uuid.UUID

	products          map[uuid.UUID]domain.Product
	productBySlug     map[string]uuid.UUID
	productCategories map[uuid.UUID]productCategoryState

	attributes          map[uuid.UUID]domain.CategoryAttribute
	attributeByCode     map[string]uuid.UUID
	attributesByCat     map[uuid.UUID][]uuid.UUID
	options             map[uuid.UUID]domain.AttributeOption
	optionByAttrValue   map[string]uuid.UUID
	optionsByAttribute  map[uuid.UUID][]uuid.UUID
	productAttrValueMap map[uuid.UUID]map[uuid.UUID]domain.ProductAttributeValue

	variantsByProduct map[uuid.UUID][]domain.ProductVariant
	variantByKey      map[string]uuid.UUID
	variantBySKU      map[string]uuid.UUID
}

type CategoryRepo struct {
	store *Store
}

type AttributeRepo struct {
	store *Store
}

type ProductRepo struct {
	store *Store
}

func NewStore() *Store {
	return &Store{
		categories:          map[uuid.UUID]domain.Category{},
		categoryBySlug:      map[string]uuid.UUID{},
		products:            map[uuid.UUID]domain.Product{},
		productBySlug:       map[string]uuid.UUID{},
		productCategories:   map[uuid.UUID]productCategoryState{},
		attributes:          map[uuid.UUID]domain.CategoryAttribute{},
		attributeByCode:     map[string]uuid.UUID{},
		attributesByCat:     map[uuid.UUID][]uuid.UUID{},
		options:             map[uuid.UUID]domain.AttributeOption{},
		optionByAttrValue:   map[string]uuid.UUID{},
		optionsByAttribute:  map[uuid.UUID][]uuid.UUID{},
		productAttrValueMap: map[uuid.UUID]map[uuid.UUID]domain.ProductAttributeValue{},
		variantsByProduct:   map[uuid.UUID][]domain.ProductVariant{},
		variantByKey:        map[string]uuid.UUID{},
		variantBySKU:        map[string]uuid.UUID{},
	}
}

func NewCategoryRepo(store *Store) *CategoryRepo {
	return &CategoryRepo{store: store}
}

func NewAttributeRepo(store *Store) *AttributeRepo {
	return &AttributeRepo{store: store}
}

func NewProductRepo(store *Store) *ProductRepo {
	return &ProductRepo{store: store}
}

func categorySlugKey(parentID *uuid.UUID, slug string) string {
	parent := "root"
	if parentID != nil {
		parent = parentID.String()
	}
	return parent + "|" + strings.ToLower(strings.TrimSpace(slug))
}

func attributeCodeKey(categoryID uuid.UUID, code string) string {
	return categoryID.String() + "|" + strings.ToLower(strings.TrimSpace(code))
}

func optionValueKey(attributeID uuid.UUID, value string) string {
	return attributeID.String() + "|" + strings.ToLower(strings.TrimSpace(value))
}

func variantCombinationKey(productID uuid.UUID, combinationKey string) string {
	return productID.String() + "|" + combinationKey
}

func variantSKUKey(productID uuid.UUID, sku string) string {
	return productID.String() + "|" + strings.TrimSpace(sku)
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	key := categorySlugKey(category.ParentID, category.Slug)
	if _, exists := r.store.categoryBySlug[key]; exists {
		return domain.ErrDuplicateSlugUnderParent
	}
	r.store.categories[category.ID] = category
	r.store.categoryBySlug[key] = category.ID
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	category, ok := r.store.categories[id]
	if !ok {
		return domain.Category{}, fmt.Errorf("%w: category %s", domain.ErrNotFound, id.String())
	}
	return category, nil
}

func (r *CategoryRepo) ExistsByParentAndSlug(ctx context.Context, parentID *uuid.UUID, slug string) (bool, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	_, exists := r.store.categoryBySlug[categorySlugKey(parentID, slug)]
	return exists, nil
}

func (r *CategoryRepo) ListChildren(ctx context.Context, parentID *uuid.UUID) ([]domain.Category, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	children := make([]domain.Category, 0)
	for _, category := range r.store.categories {
		switch {
		case category.ParentID == nil && parentID == nil:
			children = append(children, category)
		case category.ParentID != nil && parentID != nil && *category.ParentID == *parentID:
			children = append(children, category)
		}
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].SortOrder == children[j].SortOrder {
			return children[i].Name < children[j].Name
		}
		return children[i].SortOrder < children[j].SortOrder
	})
	return children, nil
}

func (r *CategoryRepo) ListAll(ctx context.Context) ([]domain.Category, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]domain.Category, 0, len(r.store.categories))
	for _, category := range r.store.categories {
		out = append(out, category)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].Name < out[j].Name
		}
		return out[i].Path < out[j].Path
	})
	return out, nil
}

func (r *AttributeRepo) Create(ctx context.Context, attribute domain.CategoryAttribute) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, ok := r.store.categories[attribute.CategoryID]; !ok {
		return fmt.Errorf("%w: category %s", domain.ErrNotFound, attribute.CategoryID.String())
	}
	codeKey := attributeCodeKey(attribute.CategoryID, attribute.Code)
	if _, exists := r.store.attributeByCode[codeKey]; exists {
		return fmt.Errorf("%w: duplicate category attribute code", domain.ErrInvalidAttribute)
	}
	r.store.attributes[attribute.ID] = attribute
	r.store.attributeByCode[codeKey] = attribute.ID
	r.store.attributesByCat[attribute.CategoryID] = append(r.store.attributesByCat[attribute.CategoryID], attribute.ID)
	return nil
}

func (r *AttributeRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.CategoryAttribute, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	attribute, ok := r.store.attributes[id]
	if !ok {
		return domain.CategoryAttribute{}, fmt.Errorf("%w: attribute %s", domain.ErrNotFound, id.String())
	}
	return attribute, nil
}

func (r *AttributeRepo) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.CategoryAttribute, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	ids := append([]uuid.UUID(nil), r.store.attributesByCat[categoryID]...)
	attrs := make([]domain.CategoryAttribute, 0, len(ids))
	for _, id := range ids {
		attrs = append(attrs, r.store.attributes[id])
	}
	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i].SortOrder == attrs[j].SortOrder {
			return attrs[i].Code < attrs[j].Code
		}
		return attrs[i].SortOrder < attrs[j].SortOrder
	})
	return attrs, nil
}

func (r *AttributeRepo) ExistsByCategoryAndCode(ctx context.Context, categoryID uuid.UUID, code string) (bool, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	_, exists := r.store.attributeByCode[attributeCodeKey(categoryID, code)]
	return exists, nil
}

func (r *AttributeRepo) CreateOption(ctx context.Context, option domain.AttributeOption) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, ok := r.store.attributes[option.AttributeID]; !ok {
		return fmt.Errorf("%w: attribute %s", domain.ErrNotFound, option.AttributeID.String())
	}
	key := optionValueKey(option.AttributeID, option.Value)
	if _, exists := r.store.optionByAttrValue[key]; exists {
		return fmt.Errorf("%w: duplicate option value", domain.ErrInvalidAttribute)
	}
	r.store.options[option.ID] = option
	r.store.optionByAttrValue[key] = option.ID
	r.store.optionsByAttribute[option.AttributeID] = append(r.store.optionsByAttribute[option.AttributeID], option.ID)
	return nil
}

func (r *AttributeRepo) ListOptionsByAttribute(ctx context.Context, attributeID uuid.UUID) ([]domain.AttributeOption, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	ids := append([]uuid.UUID(nil), r.store.optionsByAttribute[attributeID]...)
	options := make([]domain.AttributeOption, 0, len(ids))
	for _, id := range ids {
		options = append(options, r.store.options[id])
	}
	sort.Slice(options, func(i, j int) bool {
		if options[i].SortOrder == options[j].SortOrder {
			return options[i].Value < options[j].Value
		}
		return options[i].SortOrder < options[j].SortOrder
	})
	return options, nil
}

func (r *AttributeRepo) OptionBelongsToAttribute(ctx context.Context, optionID, attributeID uuid.UUID) (bool, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	option, ok := r.store.options[optionID]
	if !ok {
		return false, nil
	}
	return option.AttributeID == attributeID, nil
}

func (r *ProductRepo) Create(ctx context.Context, product domain.Product) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, ok := r.store.categories[product.PrimaryCategoryID]; !ok {
		return fmt.Errorf("%w: category %s", domain.ErrNotFound, product.PrimaryCategoryID.String())
	}
	slug := strings.ToLower(strings.TrimSpace(product.Slug))
	if _, exists := r.store.productBySlug[slug]; exists {
		return fmt.Errorf("%w: duplicate product slug", domain.ErrInvalidProduct)
	}
	r.store.products[product.ID] = product
	r.store.productBySlug[slug] = product.ID
	r.store.productCategories[product.ID] = productCategoryState{
		PrimaryCategoryID: product.PrimaryCategoryID,
		CategoryIDs:       append([]uuid.UUID(nil), product.CategoryIDs...),
	}
	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	product, ok := r.store.products[id]
	if !ok {
		return domain.Product{}, fmt.Errorf("%w: product %s", domain.ErrNotFound, id.String())
	}
	return product, nil
}

func (r *ProductRepo) SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID, primaryCategoryID uuid.UUID) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, ok := r.store.products[productID]; !ok {
		return fmt.Errorf("%w: product %s", domain.ErrNotFound, productID.String())
	}
	for _, categoryID := range categoryIDs {
		if _, ok := r.store.categories[categoryID]; !ok {
			return fmt.Errorf("%w: category %s", domain.ErrNotFound, categoryID.String())
		}
	}
	r.store.productCategories[productID] = productCategoryState{
		PrimaryCategoryID: primaryCategoryID,
		CategoryIDs:       append([]uuid.UUID(nil), categoryIDs...),
	}
	return nil
}

func (r *ProductRepo) UpsertAttributeValues(ctx context.Context, productID uuid.UUID, values []domain.ProductAttributeValue) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if _, ok := r.store.products[productID]; !ok {
		return fmt.Errorf("%w: product %s", domain.ErrNotFound, productID.String())
	}
	if _, ok := r.store.productAttrValueMap[productID]; !ok {
		r.store.productAttrValueMap[productID] = map[uuid.UUID]domain.ProductAttributeValue{}
	}
	for _, value := range values {
		if _, ok := r.store.attributes[value.AttributeID]; !ok {
			return fmt.Errorf("%w: attribute %s", domain.ErrNotFound, value.AttributeID.String())
		}
		if value.ID == uuid.Nil {
			value.ID = uuid.New()
		}
		value.ProductID = productID
		r.store.productAttrValueMap[productID][value.AttributeID] = value
	}
	return nil
}

func (r *ProductRepo) ListAttributeValues(ctx context.Context, productID uuid.UUID) ([]domain.ProductAttributeValue, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	if _, ok := r.store.products[productID]; !ok {
		return nil, fmt.Errorf("%w: product %s", domain.ErrNotFound, productID.String())
	}
	valueMap := r.store.productAttrValueMap[productID]
	out := make([]domain.ProductAttributeValue, 0, len(valueMap))
	for _, value := range valueMap {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].AttributeID.String() < out[j].AttributeID.String()
	})
	return out, nil
}

func (r *ProductRepo) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]domain.Product, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	if _, ok := r.store.categories[categoryID]; !ok {
		return nil, fmt.Errorf("%w: category %s", domain.ErrNotFound, categoryID.String())
	}
	products := make([]domain.Product, 0)
	for productID, state := range r.store.productCategories {
		for _, id := range state.CategoryIDs {
			if id == categoryID {
				products = append(products, r.store.products[productID])
				break
			}
		}
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})
	return products, nil
}

func (r *ProductRepo) CreateVariants(ctx context.Context, variants []domain.ProductVariant) error {
	_ = ctx
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if len(variants) == 0 {
		return nil
	}
	productID := variants[0].ProductID
	if _, ok := r.store.products[productID]; !ok {
		return fmt.Errorf("%w: product %s", domain.ErrNotFound, productID.String())
	}

	pendingKeyIndex := map[string]struct{}{}
	pendingSKUIndex := map[string]struct{}{}
	for _, variant := range variants {
		if variant.ProductID != productID {
			return fmt.Errorf("%w: all variants must belong to the same product", domain.ErrInvalidVariantAxis)
		}
		if err := variant.Validate(); err != nil {
			return err
		}
		key := variantCombinationKey(variant.ProductID, variant.CombinationKey)
		if _, exists := r.store.variantByKey[key]; exists {
			return fmt.Errorf("%w: combination already exists", domain.ErrDuplicateVariantCombination)
		}
		if _, exists := pendingKeyIndex[key]; exists {
			return fmt.Errorf("%w: duplicate combination in batch", domain.ErrDuplicateVariantCombination)
		}
		pendingKeyIndex[key] = struct{}{}

		if strings.TrimSpace(variant.SKU) != "" {
			skuKey := variantSKUKey(variant.ProductID, variant.SKU)
			if _, exists := r.store.variantBySKU[skuKey]; exists {
				return fmt.Errorf("%w: variant sku already exists", domain.ErrInvalidVariantAxis)
			}
			if _, exists := pendingSKUIndex[skuKey]; exists {
				return fmt.Errorf("%w: duplicate variant sku in batch", domain.ErrInvalidVariantAxis)
			}
			pendingSKUIndex[skuKey] = struct{}{}
		}
	}

	for _, variant := range variants {
		r.store.variantsByProduct[variant.ProductID] = append(r.store.variantsByProduct[variant.ProductID], variant)
		r.store.variantByKey[variantCombinationKey(variant.ProductID, variant.CombinationKey)] = variant.ID
		if strings.TrimSpace(variant.SKU) != "" {
			r.store.variantBySKU[variantSKUKey(variant.ProductID, variant.SKU)] = variant.ID
		}
	}
	return nil
}

func (r *ProductRepo) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	_ = ctx
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	if _, ok := r.store.products[productID]; !ok {
		return nil, fmt.Errorf("%w: product %s", domain.ErrNotFound, productID.String())
	}
	variants := append([]domain.ProductVariant(nil), r.store.variantsByProduct[productID]...)
	sort.Slice(variants, func(i, j int) bool {
		return variants[i].CombinationKey < variants[j].CombinationKey
	})
	return variants, nil
}
