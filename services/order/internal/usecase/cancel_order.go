package usecase

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

// CancelOrderUseCase handles cancelling an order.
type CancelOrderUseCase struct {
	orderRepo       domain.OrderRepository
	sellerOrderRepo domain.SellerOrderRepository
	publisher       domain.EventPublisher
}

// NewCancelOrderUseCase creates a new CancelOrderUseCase instance.
func NewCancelOrderUseCase(
	orderRepo domain.OrderRepository,
	sellerOrderRepo domain.SellerOrderRepository,
	publisher domain.EventPublisher,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		orderRepo:       orderRepo,
		sellerOrderRepo: sellerOrderRepo,
		publisher:       publisher,
	}
}

// Execute cancels an order and all its seller orders.
func (uc *CancelOrderUseCase) Execute(ctx context.Context, orderID string, buyerID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Verify the buyer owns this order
	if order.BuyerID != buyerID {
		return nil, fmt.Errorf("order does not belong to this buyer")
	}

	if !domain.CanTransition(order.Status, domain.OrderStatusCancelled) {
		return nil, fmt.Errorf("order cannot be cancelled from status %s", order.Status)
	}

	// Cancel the main order
	if err := uc.orderRepo.UpdateStatus(ctx, orderID, domain.OrderStatusCancelled); err != nil {
		return nil, err
	}

	// Cancel all seller orders
	sellerOrders, err := uc.sellerOrderRepo.ListByOrder(ctx, orderID)
	if err == nil {
		for _, so := range sellerOrders {
			if domain.CanTransition(so.Status, domain.OrderStatusCancelled) {
				_ = uc.sellerOrderRepo.UpdateStatus(ctx, so.ID, domain.OrderStatusCancelled)
			}
		}
	}

	order.Status = domain.OrderStatusCancelled

	// Publish order.cancelled event
	event := domain.OrderStatusEvent{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		BuyerID:     order.BuyerID,
		Status:      domain.OrderStatusCancelled,
	}
	_ = uc.publisher.Publish(ctx, domain.EventOrderCancelled, event)

	return order, nil
}
