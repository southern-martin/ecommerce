package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
)

// CreateFlashSaleInput represents the input for creating a flash sale.
type CreateFlashSaleInput struct {
	Name     string
	StartsAt time.Time
	EndsAt   time.Time
	Items    []CreateFlashSaleItemInput
}

// CreateFlashSaleItemInput represents a single item in the flash sale creation request.
type CreateFlashSaleItemInput struct {
	ProductID      string
	VariantID      string
	SalePriceCents int64
	QuantityLimit  int
}

// FlashSaleUseCase handles flash sale business logic.
type FlashSaleUseCase struct {
	flashSaleRepo     domain.FlashSaleRepository
	flashSaleItemRepo domain.FlashSaleItemRepository
	publisher         domain.EventPublisher
}

// NewFlashSaleUseCase creates a new FlashSaleUseCase instance.
func NewFlashSaleUseCase(
	flashSaleRepo domain.FlashSaleRepository,
	flashSaleItemRepo domain.FlashSaleItemRepository,
	publisher domain.EventPublisher,
) *FlashSaleUseCase {
	return &FlashSaleUseCase{
		flashSaleRepo:     flashSaleRepo,
		flashSaleItemRepo: flashSaleItemRepo,
		publisher:         publisher,
	}
}

// CreateFlashSale creates a new flash sale with items.
func (uc *FlashSaleUseCase) CreateFlashSale(ctx context.Context, input CreateFlashSaleInput) (*domain.FlashSale, error) {
	if input.Name == "" {
		return nil, errors.New("flash sale name is required")
	}
	if input.StartsAt.IsZero() {
		return nil, errors.New("starts_at is required")
	}
	if input.EndsAt.IsZero() {
		return nil, errors.New("ends_at is required")
	}
	if input.EndsAt.Before(input.StartsAt) {
		return nil, errors.New("ends_at must be after starts_at")
	}

	flashSale := domain.NewFlashSale(input.Name, input.StartsAt, input.EndsAt)

	// Add items
	for _, item := range input.Items {
		fsItem := domain.NewFlashSaleItem(
			flashSale.ID,
			item.ProductID,
			item.VariantID,
			item.SalePriceCents,
			item.QuantityLimit,
		)
		flashSale.Items = append(flashSale.Items, *fsItem)
	}

	if err := uc.flashSaleRepo.Create(ctx, flashSale); err != nil {
		return nil, err
	}

	// Publish flash_sale.started event if sale is currently active
	now := time.Now()
	if flashSale.IsActive && now.After(flashSale.StartsAt) && now.Before(flashSale.EndsAt) {
		event := domain.FlashSaleEvent{
			FlashSaleID: flashSale.ID,
			Name:        flashSale.Name,
		}
		_ = uc.publisher.Publish(ctx, domain.EventFlashSaleStarted, event)
	}

	return flashSale, nil
}

// GetFlashSale retrieves a flash sale by ID.
func (uc *FlashSaleUseCase) GetFlashSale(ctx context.Context, id string) (*domain.FlashSale, error) {
	if id == "" {
		return nil, errors.New("flash sale id is required")
	}
	return uc.flashSaleRepo.GetByID(ctx, id)
}

// ListFlashSales retrieves a paginated list of all flash sales.
func (uc *FlashSaleUseCase) ListFlashSales(ctx context.Context, page, pageSize int) ([]*domain.FlashSale, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.flashSaleRepo.ListAll(ctx, page, pageSize)
}

// ListActiveFlashSales retrieves all currently active flash sales.
func (uc *FlashSaleUseCase) ListActiveFlashSales(ctx context.Context) ([]*domain.FlashSale, error) {
	return uc.flashSaleRepo.ListActive(ctx)
}

// UpdateFlashSale updates an existing flash sale and publishes relevant events.
func (uc *FlashSaleUseCase) UpdateFlashSale(ctx context.Context, flashSale *domain.FlashSale) error {
	if err := uc.flashSaleRepo.Update(ctx, flashSale); err != nil {
		return err
	}

	event := domain.FlashSaleEvent{
		FlashSaleID: flashSale.ID,
		Name:        flashSale.Name,
	}

	if !flashSale.IsActive {
		_ = uc.publisher.Publish(ctx, domain.EventFlashSaleEnded, event)
	} else {
		now := time.Now()
		if now.After(flashSale.StartsAt) && now.Before(flashSale.EndsAt) {
			_ = uc.publisher.Publish(ctx, domain.EventFlashSaleStarted, event)
		}
	}

	return nil
}
