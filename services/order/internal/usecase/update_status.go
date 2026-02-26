package usecase

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
)

// UpdateOrderStatusUseCase handles updating the status of seller orders.
type UpdateOrderStatusUseCase struct {
	orderRepo       domain.OrderRepository
	sellerOrderRepo domain.SellerOrderRepository
	publisher       domain.EventPublisher
}

// NewUpdateOrderStatusUseCase creates a new UpdateOrderStatusUseCase instance.
func NewUpdateOrderStatusUseCase(
	orderRepo domain.OrderRepository,
	sellerOrderRepo domain.SellerOrderRepository,
	publisher domain.EventPublisher,
) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo:       orderRepo,
		sellerOrderRepo: sellerOrderRepo,
		publisher:       publisher,
	}
}

// Execute updates the status of a seller order, enforcing the state machine.
func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, sellerOrderID string, newStatus domain.OrderStatus) (*domain.SellerOrder, error) {
	sellerOrder, err := uc.sellerOrderRepo.GetByID(ctx, sellerOrderID)
	if err != nil {
		return nil, err
	}

	if !domain.CanTransition(sellerOrder.Status, newStatus) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", sellerOrder.Status, newStatus)
	}

	if err := uc.sellerOrderRepo.UpdateStatus(ctx, sellerOrderID, newStatus); err != nil {
		return nil, err
	}

	sellerOrder.Status = newStatus

	// Publish status change event
	order, err := uc.orderRepo.GetByID(ctx, sellerOrder.OrderID)
	if err == nil {
		statusEvent := domain.OrderStatusEvent{
			OrderID:     order.ID,
			OrderNumber: order.OrderNumber,
			BuyerID:     order.BuyerID,
			Status:      newStatus,
		}
		subject := statusToEventSubject(newStatus)
		if subject != "" {
			_ = uc.publisher.Publish(ctx, subject, statusEvent)
		}
	}

	return sellerOrder, nil
}

// UpdateOrderStatus updates the status of an order directly (used by gRPC / inter-service).
func (uc *UpdateOrderStatusUseCase) UpdateOrderStatus(ctx context.Context, orderID string, newStatus domain.OrderStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if !domain.CanTransition(order.Status, newStatus) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, newStatus)
	}

	if err := uc.orderRepo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		return nil, err
	}

	order.Status = newStatus

	// Publish status change event
	statusEvent := domain.OrderStatusEvent{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		BuyerID:     order.BuyerID,
		Status:      newStatus,
	}
	subject := statusToEventSubject(newStatus)
	if subject != "" {
		_ = uc.publisher.Publish(ctx, subject, statusEvent)
	}

	return order, nil
}

// statusToEventSubject maps an order status to the corresponding event subject.
func statusToEventSubject(status domain.OrderStatus) string {
	switch status {
	case domain.OrderStatusConfirmed:
		return domain.EventOrderConfirmed
	case domain.OrderStatusCancelled:
		return domain.EventOrderCancelled
	case domain.OrderStatusShipped:
		return domain.EventOrderShipped
	case domain.OrderStatusDelivered:
		return domain.EventOrderDelivered
	case domain.OrderStatusCompleted:
		return domain.EventOrderCompleted
	default:
		return ""
	}
}
