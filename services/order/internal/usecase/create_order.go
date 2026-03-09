package usecase

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

// ProductServiceClient defines the interface for cross-service product calls.
type ProductServiceClient interface {
	UpdateStock(ctx context.Context, variantID string, delta int32) error
}

// CreateOrderInput represents the input for creating a new order.
type CreateOrderInput struct {
	BuyerID         string
	BuyerEmail      string
	Currency        string
	ShippingAddress domain.Address
	Items           []CreateOrderItemInput
}

// CreateOrderItemInput represents a single item in the order creation request.
type CreateOrderItemInput struct {
	ProductID      string
	VariantID      string
	ProductName    string
	VariantName    string
	SKU            string
	Quantity       int
	UnitPriceCents int64
	SellerID       string
	ImageURL       string
}

// CreateOrderUseCase handles the creation of new orders.
type CreateOrderUseCase struct {
	orderRepo       domain.OrderRepository
	sellerOrderRepo domain.SellerOrderRepository
	publisher       domain.EventPublisher
	productClient   ProductServiceClient
}

// NewCreateOrderUseCase creates a new CreateOrderUseCase instance.
func NewCreateOrderUseCase(
	orderRepo domain.OrderRepository,
	sellerOrderRepo domain.SellerOrderRepository,
	publisher domain.EventPublisher,
	productClient ProductServiceClient,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:       orderRepo,
		sellerOrderRepo: sellerOrderRepo,
		publisher:       publisher,
		productClient:   productClient,
	}
}

// Execute creates a new order, splits it by seller, persists it, and publishes an event.
func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	if input.BuyerID == "" {
		return nil, errors.New("buyer_id is required")
	}
	if len(input.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}
	if input.Currency == "" {
		input.Currency = "USD"
	}

	// Convert input items to domain items
	var items []domain.OrderItem
	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than 0")
		}
		if item.UnitPriceCents <= 0 {
			return nil, errors.New("item unit price must be greater than 0")
		}
		if item.SellerID == "" {
			return nil, errors.New("seller_id is required for each item")
		}
		items = append(items, domain.OrderItem{
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			ProductName:    item.ProductName,
			VariantName:    item.VariantName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
			SellerID:       item.SellerID,
			ImageURL:       item.ImageURL,
		})
	}

	// Create the order with seller splitting
	order := domain.NewOrder(input.BuyerID, input.Currency, input.ShippingAddress, items)

	// Persist the order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Persist seller orders
	for i := range order.SellerOrders {
		if err := uc.sellerOrderRepo.Create(ctx, &order.SellerOrders[i]); err != nil {
			return nil, err
		}
	}

	// Reserve inventory by decrementing stock for each item.
	if uc.productClient != nil {
		for _, item := range order.Items {
			if err := uc.productClient.UpdateStock(ctx, item.VariantID, -int32(item.Quantity)); err != nil {
				log.Warn().Err(err).
					Str("variant_id", item.VariantID).
					Int("quantity", item.Quantity).
					Msg("failed to reserve stock for order item")
			}
		}
	}

	// Publish order.created event
	var eventItems []domain.ItemEvent
	for _, item := range order.Items {
		eventItems = append(eventItems, domain.ItemEvent{
			ProductID:      item.ProductID,
			VariantID:      item.VariantID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
			SellerID:       item.SellerID,
		})
	}
	event := domain.OrderCreatedEvent{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		BuyerID:     order.BuyerID,
		BuyerEmail:  input.BuyerEmail,
		TotalCents:  order.TotalCents,
		Currency:    order.Currency,
		Items:       eventItems,
	}
	if pubErr := uc.publisher.Publish(ctx, domain.EventOrderCreated, event); pubErr != nil {
		log.Warn().Err(pubErr).Str("event", domain.EventOrderCreated).Msg("failed to publish event")
	}

	return order, nil
}
