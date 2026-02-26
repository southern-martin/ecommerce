package usecase

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// ManageReturnUseCase handles return management operations.
type ManageReturnUseCase struct {
	returnRepo domain.ReturnRepository
	publisher  domain.EventPublisher
}

// NewManageReturnUseCase creates a new ManageReturnUseCase.
func NewManageReturnUseCase(returnRepo domain.ReturnRepository, publisher domain.EventPublisher) *ManageReturnUseCase {
	return &ManageReturnUseCase{
		returnRepo: returnRepo,
		publisher:  publisher,
	}
}

// GetReturn retrieves a return by ID.
func (uc *ManageReturnUseCase) GetReturn(ctx context.Context, id string) (*domain.Return, error) {
	return uc.returnRepo.GetByID(ctx, id)
}

// ListBuyerReturns lists returns for a buyer.
func (uc *ManageReturnUseCase) ListBuyerReturns(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Return, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.returnRepo.ListByBuyer(ctx, buyerID, page, pageSize)
}

// ListSellerReturns lists returns for a seller.
func (uc *ManageReturnUseCase) ListSellerReturns(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Return, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.returnRepo.ListBySeller(ctx, sellerID, page, pageSize)
}

// ApproveReturn approves a return request.
func (uc *ManageReturnUseCase) ApproveReturn(ctx context.Context, id, sellerID string, refundAmountCents int64) (*domain.Return, error) {
	ret, err := uc.returnRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("return not found: %w", err)
	}

	if ret.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: return belongs to different seller")
	}

	if !domain.CanReturnTransition(ret.Status, domain.ReturnStatusApproved) {
		return nil, fmt.Errorf("cannot approve return in status %s", ret.Status)
	}

	ret.Status = domain.ReturnStatusApproved
	if refundAmountCents > 0 {
		ret.RefundAmountCents = refundAmountCents
	}

	if err := uc.returnRepo.Update(ctx, ret); err != nil {
		return nil, fmt.Errorf("failed to update return: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "return.approved", map[string]interface{}{
		"return_id":          ret.ID,
		"order_id":           ret.OrderID,
		"refund_amount_cents": ret.RefundAmountCents,
	})

	return ret, nil
}

// RejectReturn rejects a return request.
func (uc *ManageReturnUseCase) RejectReturn(ctx context.Context, id, sellerID string) (*domain.Return, error) {
	ret, err := uc.returnRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("return not found: %w", err)
	}

	if ret.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: return belongs to different seller")
	}

	if !domain.CanReturnTransition(ret.Status, domain.ReturnStatusRejected) {
		return nil, fmt.Errorf("cannot reject return in status %s", ret.Status)
	}

	ret.Status = domain.ReturnStatusRejected

	if err := uc.returnRepo.Update(ctx, ret); err != nil {
		return nil, fmt.Errorf("failed to update return: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "return.rejected", map[string]interface{}{
		"return_id": ret.ID,
		"order_id":  ret.OrderID,
	})

	return ret, nil
}

// UpdateReturnStatus updates a return's status with validation.
func (uc *ManageReturnUseCase) UpdateReturnStatus(ctx context.Context, id, sellerID string, newStatus domain.ReturnStatus) (*domain.Return, error) {
	ret, err := uc.returnRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("return not found: %w", err)
	}

	if ret.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: return belongs to different seller")
	}

	if !domain.CanReturnTransition(ret.Status, newStatus) {
		return nil, fmt.Errorf("invalid transition from %s to %s", ret.Status, newStatus)
	}

	ret.Status = newStatus

	if err := uc.returnRepo.Update(ctx, ret); err != nil {
		return nil, fmt.Errorf("failed to update return: %w", err)
	}

	// Publish completed event if refunded
	if newStatus == domain.ReturnStatusRefunded {
		_ = uc.publisher.Publish(ctx, "return.completed", map[string]interface{}{
			"return_id":          ret.ID,
			"order_id":           ret.OrderID,
			"refund_amount_cents": ret.RefundAmountCents,
		})
	}

	return ret, nil
}
