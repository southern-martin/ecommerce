package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks (shared with variant_test.go)
// ---------------------------------------------------------------------------

// --- ProductRepository mock ---

type mockProductRepo struct {
	createFn    func(ctx context.Context, p *domain.Product) error
	getByIDFn   func(ctx context.Context, id string) (*domain.Product, error)
	getBySlugFn func(ctx context.Context, slug string) (*domain.Product, error)
	listFn      func(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error)
	updateFn    func(ctx context.Context, p *domain.Product) error
	deleteFn    func(ctx context.Context, id string) error
}

func (m *mockProductRepo) Create(ctx context.Context, p *domain.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	return nil
}
func (m *mockProductRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockProductRepo) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, errors.New("not found")
}
func (m *mockProductRepo) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *mockProductRepo) Update(ctx context.Context, p *domain.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, p)
	}
	return nil
}
func (m *mockProductRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- CategoryRepository mock ---

type mockCategoryRepo struct {
	createFn  func(ctx context.Context, c *domain.Category) error
	getByIDFn func(ctx context.Context, id string) (*domain.Category, error)
	listFn    func(ctx context.Context) ([]*domain.Category, error)
	updateFn  func(ctx context.Context, c *domain.Category) error
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockCategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	if m.createFn != nil {
		return m.createFn(ctx, c)
	}
	return nil
}
func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockCategoryRepo) List(ctx context.Context) ([]*domain.Category, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}
func (m *mockCategoryRepo) Update(ctx context.Context, c *domain.Category) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, c)
	}
	return nil
}
func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- AttributeRepository mock ---

type mockAttributeRepo struct {
	createDefinitionFn     func(ctx context.Context, attr *domain.AttributeDefinition) error
	getDefinitionByIDFn    func(ctx context.Context, id string) (*domain.AttributeDefinition, error)
	listDefinitionsFn      func(ctx context.Context) ([]*domain.AttributeDefinition, error)
	updateDefinitionFn     func(ctx context.Context, attr *domain.AttributeDefinition) error
	deleteDefinitionFn     func(ctx context.Context, id string) error
	assignToCategoryFn     func(ctx context.Context, categoryID, attributeID string, sortOrder int) error
	removeFromCategoryFn   func(ctx context.Context, categoryID, attributeID string) error
	listByCategoryFn       func(ctx context.Context, categoryID string) ([]*domain.AttributeDefinition, error)
	setProductValuesFn     func(ctx context.Context, productID string, values []domain.ProductAttributeValue) error
	getProductValuesFn     func(ctx context.Context, productID string) ([]domain.ProductAttributeValue, error)
	setOptionValuesFn      func(ctx context.Context, attributeID string, opts []domain.AttributeOptionValue) error
	listOptionValuesFn     func(ctx context.Context, attributeID string) ([]domain.AttributeOptionValue, error)
}

func (m *mockAttributeRepo) CreateDefinition(ctx context.Context, attr *domain.AttributeDefinition) error {
	if m.createDefinitionFn != nil {
		return m.createDefinitionFn(ctx, attr)
	}
	return nil
}
func (m *mockAttributeRepo) GetDefinitionByID(ctx context.Context, id string) (*domain.AttributeDefinition, error) {
	if m.getDefinitionByIDFn != nil {
		return m.getDefinitionByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockAttributeRepo) ListDefinitions(ctx context.Context) ([]*domain.AttributeDefinition, error) {
	if m.listDefinitionsFn != nil {
		return m.listDefinitionsFn(ctx)
	}
	return nil, nil
}
func (m *mockAttributeRepo) UpdateDefinition(ctx context.Context, attr *domain.AttributeDefinition) error {
	if m.updateDefinitionFn != nil {
		return m.updateDefinitionFn(ctx, attr)
	}
	return nil
}
func (m *mockAttributeRepo) DeleteDefinition(ctx context.Context, id string) error {
	if m.deleteDefinitionFn != nil {
		return m.deleteDefinitionFn(ctx, id)
	}
	return nil
}
func (m *mockAttributeRepo) AssignToCategory(ctx context.Context, categoryID, attributeID string, sortOrder int) error {
	if m.assignToCategoryFn != nil {
		return m.assignToCategoryFn(ctx, categoryID, attributeID, sortOrder)
	}
	return nil
}
func (m *mockAttributeRepo) RemoveFromCategory(ctx context.Context, categoryID, attributeID string) error {
	if m.removeFromCategoryFn != nil {
		return m.removeFromCategoryFn(ctx, categoryID, attributeID)
	}
	return nil
}
func (m *mockAttributeRepo) ListByCategory(ctx context.Context, categoryID string) ([]*domain.AttributeDefinition, error) {
	if m.listByCategoryFn != nil {
		return m.listByCategoryFn(ctx, categoryID)
	}
	return nil, nil
}
func (m *mockAttributeRepo) SetProductValues(ctx context.Context, productID string, values []domain.ProductAttributeValue) error {
	if m.setProductValuesFn != nil {
		return m.setProductValuesFn(ctx, productID, values)
	}
	return nil
}
func (m *mockAttributeRepo) GetProductValues(ctx context.Context, productID string) ([]domain.ProductAttributeValue, error) {
	if m.getProductValuesFn != nil {
		return m.getProductValuesFn(ctx, productID)
	}
	return nil, nil
}
func (m *mockAttributeRepo) SetOptionValues(ctx context.Context, attributeID string, opts []domain.AttributeOptionValue) error {
	if m.setOptionValuesFn != nil {
		return m.setOptionValuesFn(ctx, attributeID, opts)
	}
	return nil
}
func (m *mockAttributeRepo) ListOptionValues(ctx context.Context, attributeID string) ([]domain.AttributeOptionValue, error) {
	if m.listOptionValuesFn != nil {
		return m.listOptionValuesFn(ctx, attributeID)
	}
	return nil, nil
}

// --- AttributeGroupRepository mock ---

type mockAttributeGroupRepo struct {
	createFn          func(ctx context.Context, group *domain.AttributeGroup) error
	getByIDFn         func(ctx context.Context, id string) (*domain.AttributeGroup, error)
	listFn            func(ctx context.Context) ([]*domain.AttributeGroup, error)
	updateFn          func(ctx context.Context, group *domain.AttributeGroup) error
	deleteFn          func(ctx context.Context, id string) error
	addAttributeFn    func(ctx context.Context, groupID, attributeID string, sortOrder int) error
	removeAttributeFn func(ctx context.Context, groupID, attributeID string) error
	listAttributesFn  func(ctx context.Context, groupID string) ([]*domain.AttributeDefinition, error)
}

func (m *mockAttributeGroupRepo) Create(ctx context.Context, group *domain.AttributeGroup) error {
	if m.createFn != nil {
		return m.createFn(ctx, group)
	}
	return nil
}
func (m *mockAttributeGroupRepo) GetByID(ctx context.Context, id string) (*domain.AttributeGroup, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockAttributeGroupRepo) List(ctx context.Context) ([]*domain.AttributeGroup, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}
func (m *mockAttributeGroupRepo) Update(ctx context.Context, group *domain.AttributeGroup) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, group)
	}
	return nil
}
func (m *mockAttributeGroupRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}
func (m *mockAttributeGroupRepo) AddAttribute(ctx context.Context, groupID, attributeID string, sortOrder int) error {
	if m.addAttributeFn != nil {
		return m.addAttributeFn(ctx, groupID, attributeID, sortOrder)
	}
	return nil
}
func (m *mockAttributeGroupRepo) RemoveAttribute(ctx context.Context, groupID, attributeID string) error {
	if m.removeAttributeFn != nil {
		return m.removeAttributeFn(ctx, groupID, attributeID)
	}
	return nil
}
func (m *mockAttributeGroupRepo) ListAttributes(ctx context.Context, groupID string) ([]*domain.AttributeDefinition, error) {
	if m.listAttributesFn != nil {
		return m.listAttributesFn(ctx, groupID)
	}
	return nil, nil
}

// --- OptionRepository mock ---

type mockOptionRepo struct {
	createOptionFn      func(ctx context.Context, option *domain.ProductOption) error
	updateOptionFn      func(ctx context.Context, option *domain.ProductOption) error
	deleteOptionFn      func(ctx context.Context, optionID string) error
	listByProductFn     func(ctx context.Context, productID string) ([]domain.ProductOption, error)
	createOptionValueFn func(ctx context.Context, value *domain.ProductOptionValue) error
	updateOptionValueFn func(ctx context.Context, value *domain.ProductOptionValue) error
	deleteOptionValueFn func(ctx context.Context, valueID string) error
}

func (m *mockOptionRepo) CreateOption(ctx context.Context, option *domain.ProductOption) error {
	if m.createOptionFn != nil {
		return m.createOptionFn(ctx, option)
	}
	return nil
}
func (m *mockOptionRepo) UpdateOption(ctx context.Context, option *domain.ProductOption) error {
	if m.updateOptionFn != nil {
		return m.updateOptionFn(ctx, option)
	}
	return nil
}
func (m *mockOptionRepo) DeleteOption(ctx context.Context, optionID string) error {
	if m.deleteOptionFn != nil {
		return m.deleteOptionFn(ctx, optionID)
	}
	return nil
}
func (m *mockOptionRepo) ListByProduct(ctx context.Context, productID string) ([]domain.ProductOption, error) {
	if m.listByProductFn != nil {
		return m.listByProductFn(ctx, productID)
	}
	return nil, nil
}
func (m *mockOptionRepo) CreateOptionValue(ctx context.Context, value *domain.ProductOptionValue) error {
	if m.createOptionValueFn != nil {
		return m.createOptionValueFn(ctx, value)
	}
	return nil
}
func (m *mockOptionRepo) UpdateOptionValue(ctx context.Context, value *domain.ProductOptionValue) error {
	if m.updateOptionValueFn != nil {
		return m.updateOptionValueFn(ctx, value)
	}
	return nil
}
func (m *mockOptionRepo) DeleteOptionValue(ctx context.Context, valueID string) error {
	if m.deleteOptionValueFn != nil {
		return m.deleteOptionValueFn(ctx, valueID)
	}
	return nil
}

// --- VariantRepository mock ---

type mockVariantRepo struct {
	createFn          func(ctx context.Context, v *domain.Variant) error
	getByIDFn         func(ctx context.Context, id string) (*domain.Variant, error)
	getBySKUFn        func(ctx context.Context, sku string) (*domain.Variant, error)
	listByProductFn   func(ctx context.Context, productID string) ([]domain.Variant, error)
	updateFn          func(ctx context.Context, v *domain.Variant) error
	deleteFn          func(ctx context.Context, id string) error
	bulkCreateFn      func(ctx context.Context, variants []domain.Variant) error
	updateStockFn     func(ctx context.Context, variantID string, delta int) error
	setOptionValuesFn func(ctx context.Context, variantID string, values []domain.VariantOptionValue) error
}

func (m *mockVariantRepo) Create(ctx context.Context, v *domain.Variant) error {
	if m.createFn != nil {
		return m.createFn(ctx, v)
	}
	return nil
}
func (m *mockVariantRepo) GetByID(ctx context.Context, id string) (*domain.Variant, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockVariantRepo) GetBySKU(ctx context.Context, sku string) (*domain.Variant, error) {
	if m.getBySKUFn != nil {
		return m.getBySKUFn(ctx, sku)
	}
	return nil, errors.New("not found")
}
func (m *mockVariantRepo) ListByProduct(ctx context.Context, productID string) ([]domain.Variant, error) {
	if m.listByProductFn != nil {
		return m.listByProductFn(ctx, productID)
	}
	return nil, nil
}
func (m *mockVariantRepo) Update(ctx context.Context, v *domain.Variant) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, v)
	}
	return nil
}
func (m *mockVariantRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}
func (m *mockVariantRepo) BulkCreate(ctx context.Context, variants []domain.Variant) error {
	if m.bulkCreateFn != nil {
		return m.bulkCreateFn(ctx, variants)
	}
	return nil
}
func (m *mockVariantRepo) UpdateStock(ctx context.Context, variantID string, delta int) error {
	if m.updateStockFn != nil {
		return m.updateStockFn(ctx, variantID, delta)
	}
	return nil
}
func (m *mockVariantRepo) SetOptionValues(ctx context.Context, variantID string, values []domain.VariantOptionValue) error {
	if m.setOptionValuesFn != nil {
		return m.setOptionValuesFn(ctx, variantID, values)
	}
	return nil
}

// --- EventPublisher mock ---

type mockEventPub struct {
	publishProductCreatedFn func(ctx context.Context, product *domain.Product) error
	publishProductUpdatedFn func(ctx context.Context, product *domain.Product) error
	publishProductDeletedFn func(ctx context.Context, productID string) error
	publishStockUpdatedFn   func(ctx context.Context, variantID string, newStock int, delta int) error
}

func (m *mockEventPub) PublishProductCreated(ctx context.Context, product *domain.Product) error {
	if m.publishProductCreatedFn != nil {
		return m.publishProductCreatedFn(ctx, product)
	}
	return nil
}
func (m *mockEventPub) PublishProductUpdated(ctx context.Context, product *domain.Product) error {
	if m.publishProductUpdatedFn != nil {
		return m.publishProductUpdatedFn(ctx, product)
	}
	return nil
}
func (m *mockEventPub) PublishProductDeleted(ctx context.Context, productID string) error {
	if m.publishProductDeletedFn != nil {
		return m.publishProductDeletedFn(ctx, productID)
	}
	return nil
}
func (m *mockEventPub) PublishStockUpdated(ctx context.Context, variantID string, newStock int, delta int) error {
	if m.publishStockUpdatedFn != nil {
		return m.publishStockUpdatedFn(ctx, variantID, newStock, delta)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newProductUseCase(
	productRepo *mockProductRepo,
	categoryRepo *mockCategoryRepo,
	attributeRepo *mockAttributeRepo,
	optionRepo *mockOptionRepo,
	variantRepo *mockVariantRepo,
	eventPub *mockEventPub,
) *ProductUseCase {
	return NewProductUseCase(productRepo, categoryRepo, attributeRepo, optionRepo, variantRepo, eventPub)
}

func defaultMocks() (*mockProductRepo, *mockCategoryRepo, *mockAttributeRepo, *mockOptionRepo, *mockVariantRepo, *mockEventPub) {
	return &mockProductRepo{},
		&mockCategoryRepo{},
		&mockAttributeRepo{},
		&mockOptionRepo{},
		&mockVariantRepo{},
		&mockEventPub{}
}

// ===========================================================================
// CreateProduct tests
// ===========================================================================

func TestCreateProduct_Success(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	var saved *domain.Product
	pRepo.createFn = func(_ context.Context, p *domain.Product) error {
		saved = p
		return nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	product, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID:       "seller-1",
		Name:           "Wireless Mouse",
		BasePriceCents: 2999,
	})

	require.NoError(t, err)
	require.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "seller-1", product.SellerID)
	assert.Equal(t, "Wireless Mouse", product.Name)
	assert.Equal(t, int64(2999), product.BasePriceCents)
	assert.Equal(t, domain.ProductStatusDraft, product.Status)
	assert.NotNil(t, saved)
}

func TestCreateProduct_ValidationEmptyName(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)

	_, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID: "seller-1",
		Name:     "",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product name is required")
}

func TestCreateProduct_ValidationEmptySellerID(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)

	_, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID: "",
		Name:     "Widget",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "seller ID is required")
}

func TestCreateProduct_ValidationNegativePrice(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)

	_, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID:       "seller-1",
		Name:           "Widget",
		BasePriceCents: -100,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "base price must be non-negative")
}

func TestCreateProduct_DefaultCurrencyUSD(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	pRepo.createFn = func(_ context.Context, _ *domain.Product) error { return nil }

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	product, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID:       "seller-1",
		Name:           "Widget",
		BasePriceCents: 500,
		Currency:       "", // empty => should default to USD
	})

	require.NoError(t, err)
	assert.Equal(t, "USD", product.Currency)
}

func TestCreateProduct_DefaultTypeSimple(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	pRepo.createFn = func(_ context.Context, _ *domain.Product) error { return nil }

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	product, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID: "seller-1",
		Name:     "Widget",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ProductTypeSimple, product.ProductType)
	assert.False(t, product.HasVariants)
}

func TestCreateProduct_ConfigurableHasVariants(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	pRepo.createFn = func(_ context.Context, _ *domain.Product) error { return nil }

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	product, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID:    "seller-1",
		Name:        "T-Shirt",
		ProductType: domain.ProductTypeConfigurable,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ProductTypeConfigurable, product.ProductType)
	assert.True(t, product.HasVariants)
}

func TestCreateProduct_InvalidCategory(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	cRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Category, error) {
		return nil, errors.New("category does not exist")
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	_, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID:   "seller-1",
		Name:       "Widget",
		CategoryID: "bad-cat-id",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "category not found")
}

func TestCreateProduct_WithAttributes(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()
	pRepo.createFn = func(_ context.Context, _ *domain.Product) error { return nil }

	var savedValues []domain.ProductAttributeValue
	aRepo.setProductValuesFn = func(_ context.Context, _ string, values []domain.ProductAttributeValue) error {
		savedValues = values
		return nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	_, err := uc.CreateProduct(context.Background(), CreateProductInput{
		SellerID: "seller-1",
		Name:     "Widget",
		Attributes: []AttributeValueInput{
			{AttributeID: "attr-1", Value: "Red"},
			{AttributeID: "attr-2", Value: "Large"},
		},
	})

	require.NoError(t, err)
	require.Len(t, savedValues, 2)
	assert.Equal(t, "attr-1", savedValues[0].AttributeID)
	assert.Equal(t, "Red", savedValues[0].Value)
}

// ===========================================================================
// UpdateProduct tests
// ===========================================================================

func TestUpdateProduct_Success(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	existing := &domain.Product{
		ID:             "prod-1",
		SellerID:       "seller-1",
		Name:           "Old Name",
		Slug:           "old-name-abc12345",
		BasePriceCents: 1000,
		Currency:       "USD",
		Status:         domain.ProductStatusDraft,
		ProductType:    domain.ProductTypeSimple,
	}
	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return existing, nil
	}
	var updated *domain.Product
	pRepo.updateFn = func(_ context.Context, p *domain.Product) error {
		updated = p
		return nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	newName := "New Name"
	newPrice := int64(2500)
	product, err := uc.UpdateProduct(context.Background(), "prod-1", "seller-1", UpdateProductInput{
		Name:           &newName,
		BasePriceCents: &newPrice,
	})

	require.NoError(t, err)
	assert.Equal(t, "New Name", product.Name)
	assert.Equal(t, int64(2500), product.BasePriceCents)
	assert.NotNil(t, updated)
}

func TestUpdateProduct_WrongSeller(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return &domain.Product{ID: "prod-1", SellerID: "seller-1"}, nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	_, err := uc.UpdateProduct(context.Background(), "prod-1", "seller-WRONG", UpdateProductInput{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestUpdateProduct_NotFound(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return nil, errors.New("record not found")
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	_, err := uc.UpdateProduct(context.Background(), "nonexistent", "seller-1", UpdateProductInput{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

func TestUpdateProduct_PartialUpdate(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	existing := &domain.Product{
		ID:             "prod-1",
		SellerID:       "seller-1",
		Name:           "Widget",
		Description:    "Original description",
		BasePriceCents: 1000,
		Currency:       "USD",
		Status:         domain.ProductStatusDraft,
	}
	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return existing, nil
	}
	pRepo.updateFn = func(_ context.Context, _ *domain.Product) error { return nil }

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	newDesc := "Updated description"
	product, err := uc.UpdateProduct(context.Background(), "prod-1", "seller-1", UpdateProductInput{
		Description: &newDesc,
		// Name is nil so should remain unchanged
	})

	require.NoError(t, err)
	assert.Equal(t, "Widget", product.Name)                  // unchanged
	assert.Equal(t, "Updated description", product.Description) // updated
	assert.Equal(t, int64(1000), product.BasePriceCents)        // unchanged
}

func TestUpdateProduct_SlugRegenerationOnNameChange(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	existing := &domain.Product{
		ID:       "prod-1",
		SellerID: "seller-1",
		Name:     "Old Name",
		Slug:     "old-name-abc12345",
	}
	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return existing, nil
	}
	pRepo.updateFn = func(_ context.Context, _ *domain.Product) error { return nil }

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	newName := "Brand New Name"
	product, err := uc.UpdateProduct(context.Background(), "prod-1", "seller-1", UpdateProductInput{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(product.Slug, "brand-new-name-"), "slug should start with 'brand-new-name-'")
	assert.NotEqual(t, "old-name-abc12345", product.Slug)
}

// ===========================================================================
// DeleteProduct tests
// ===========================================================================

func TestDeleteProduct_Success(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return &domain.Product{ID: "prod-1", SellerID: "seller-1"}, nil
	}
	var deletedID string
	pRepo.deleteFn = func(_ context.Context, id string) error {
		deletedID = id
		return nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	err := uc.DeleteProduct(context.Background(), "prod-1", "seller-1")

	require.NoError(t, err)
	assert.Equal(t, "prod-1", deletedID)
}

func TestDeleteProduct_WrongSeller(t *testing.T) {
	pRepo, cRepo, aRepo, oRepo, vRepo, pub := defaultMocks()

	pRepo.getByIDFn = func(_ context.Context, _ string) (*domain.Product, error) {
		return &domain.Product{ID: "prod-1", SellerID: "seller-1"}, nil
	}

	uc := newProductUseCase(pRepo, cRepo, aRepo, oRepo, vRepo, pub)
	err := uc.DeleteProduct(context.Background(), "prod-1", "seller-WRONG")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

// ===========================================================================
// generateSlug tests
// ===========================================================================

func TestGenerateSlug_FormatValidation(t *testing.T) {
	slug := generateSlug("Hello World Product")

	// Should be lowercase, hyphenated, with UUID suffix
	assert.True(t, strings.HasPrefix(slug, "hello-world-product-"), "slug should start with 'hello-world-product-'")
	// The suffix should be 8 hex chars
	parts := strings.Split(slug, "hello-world-product-")
	require.Len(t, parts, 2)
	assert.Len(t, parts[1], 8, "UUID suffix should be 8 characters")
}

func TestGenerateSlug_SpecialCharsHandling(t *testing.T) {
	slug := generateSlug("Product @#$% Name!!!")

	// Special characters should be stripped; result should be alphanumeric + hyphens + suffix
	assert.True(t, strings.HasPrefix(slug, "product-name-"), "slug should strip special chars: got %s", slug)
	// Verify no special chars remain (except hyphens)
	base := slug[:len(slug)-9] // strip -<8 hex chars>
	for _, c := range base {
		assert.True(t, (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-',
			"unexpected char %q in slug base %q", string(c), base)
	}
}

func TestGenerateSlug_Uniqueness(t *testing.T) {
	slug1 := generateSlug("Same Name")
	slug2 := generateSlug("Same Name")

	// UUID suffixes should make them different
	assert.NotEqual(t, slug1, slug2, "two calls to generateSlug with the same input should produce different slugs")
	// But prefixes should match
	assert.True(t, strings.HasPrefix(slug1, "same-name-"))
	assert.True(t, strings.HasPrefix(slug2, "same-name-"))
}
