package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// CreateReturnUseCase handles return creation.
type CreateReturnUseCase struct {
	returnRepo domain.ReturnRepository
	publisher  domain.EventPublisher
}

// NewCreateReturnUseCase creates a new CreateReturnUseCase.
func NewCreateReturnUseCase(returnRepo domain.ReturnRepository, publisher domain.EventPublisher) *CreateReturnUseCase {
	return &CreateReturnUseCase{
		returnRepo: returnRepo,
		publisher:  publisher,
	}
}

// CreateReturnRequest is the input for creating a return.
type CreateReturnRequest struct {
	OrderID           string
	BuyerID           string
	SellerID          string
	Reason            string
	Description       string
	ImageURLs         []string
	RefundAmountCents int64
	RefundMethod      string
	Items             []CreateReturnItemRequest
}

// CreateReturnItemRequest is the input for a return item.
type CreateReturnItemRequest struct {
	OrderItemID string
	ProductID   string
	VariantID   string
	Quantity    int
	Reason      string
}

// Execute creates a new return request.
func (uc *CreateReturnUseCase) Execute(ctx context.Context, req CreateReturnRequest) (*domain.Return, error) {
	if req.OrderID == "" || req.BuyerID == "" || req.SellerID == "" {
		return nil, fmt.Errorf("order_id, buyer_id, and seller_id are required")
	}
	if req.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	returnID := uuid.New().String()

	items := make([]domain.ReturnItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.ReturnItem{
			ID:          uuid.New().String(),
			ReturnID:    returnID,
			OrderItemID: item.OrderItemID,
			ProductID:   item.ProductID,
			VariantID:   item.VariantID,
			Quantity:    item.Quantity,
			Reason:      item.Reason,
		}
	}

	ret := &domain.Return{
		ID:                returnID,
		OrderID:           req.OrderID,
		BuyerID:           req.BuyerID,
		SellerID:          req.SellerID,
		Status:            domain.ReturnStatusRequested,
		Reason:            domain.ReturnReason(req.Reason),
		Description:       req.Description,
		ImageURLs:         req.ImageURLs,
		Items:             items,
		RefundAmountCents: req.RefundAmountCents,
		RefundMethod:      req.RefundMethod,
	}

	if err := uc.returnRepo.Create(ctx, ret); err != nil {
		return nil, fmt.Errorf("failed to create return: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "return.requested", map[string]interface{}{
		"return_id":  ret.ID,
		"order_id":   ret.OrderID,
		"buyer_id":   ret.BuyerID,
		"seller_id":  ret.SellerID,
		"reason":     string(ret.Reason),
		"item_count": len(ret.Items),
	})

	return ret, nil
}
