package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// DisputeUseCase handles dispute operations.
type DisputeUseCase struct {
	disputeRepo    domain.DisputeRepository
	messageRepo    domain.DisputeMessageRepository
	publisher      domain.EventPublisher
}

// NewDisputeUseCase creates a new DisputeUseCase.
func NewDisputeUseCase(
	disputeRepo domain.DisputeRepository,
	messageRepo domain.DisputeMessageRepository,
	publisher domain.EventPublisher,
) *DisputeUseCase {
	return &DisputeUseCase{
		disputeRepo: disputeRepo,
		messageRepo: messageRepo,
		publisher:   publisher,
	}
}

// CreateDisputeRequest is the input for creating a dispute.
type CreateDisputeRequest struct {
	OrderID     string
	ReturnID    string
	BuyerID     string
	SellerID    string
	Type        string
	Description string
}

// CreateDispute creates a new dispute.
func (uc *DisputeUseCase) CreateDispute(ctx context.Context, req CreateDisputeRequest) (*domain.Dispute, error) {
	if req.OrderID == "" || req.BuyerID == "" || req.SellerID == "" {
		return nil, fmt.Errorf("order_id, buyer_id, and seller_id are required")
	}
	if req.Type == "" || req.Description == "" {
		return nil, fmt.Errorf("type and description are required")
	}

	dispute := &domain.Dispute{
		ID:          uuid.New().String(),
		OrderID:     req.OrderID,
		ReturnID:    req.ReturnID,
		BuyerID:     req.BuyerID,
		SellerID:    req.SellerID,
		Status:      domain.DisputeStatusOpen,
		Type:        domain.DisputeType(req.Type),
		Description: req.Description,
	}

	if err := uc.disputeRepo.Create(ctx, dispute); err != nil {
		return nil, fmt.Errorf("failed to create dispute: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "dispute.opened", map[string]interface{}{
		"dispute_id": dispute.ID,
		"order_id":   dispute.OrderID,
		"type":       string(dispute.Type),
		"buyer_id":   dispute.BuyerID,
		"seller_id":  dispute.SellerID,
	})

	return dispute, nil
}

// GetDispute retrieves a dispute by ID.
func (uc *DisputeUseCase) GetDispute(ctx context.Context, id string) (*domain.Dispute, error) {
	return uc.disputeRepo.GetByID(ctx, id)
}

// ListAllDisputes lists all disputes (admin).
func (uc *DisputeUseCase) ListAllDisputes(ctx context.Context, page, pageSize int) ([]domain.Dispute, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.disputeRepo.ListAll(ctx, page, pageSize)
}

// ListBuyerDisputes lists disputes for a buyer.
func (uc *DisputeUseCase) ListBuyerDisputes(ctx context.Context, buyerID string, page, pageSize int) ([]domain.Dispute, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.disputeRepo.ListByBuyer(ctx, buyerID, page, pageSize)
}

// AddMessageRequest is the input for adding a dispute message.
type AddMessageRequest struct {
	DisputeID   string
	SenderID    string
	SenderRole  string
	Message     string
	Attachments []string
}

// AddMessage adds a message to a dispute.
func (uc *DisputeUseCase) AddMessage(ctx context.Context, req AddMessageRequest) (*domain.DisputeMessage, error) {
	// Verify dispute exists
	dispute, err := uc.disputeRepo.GetByID(ctx, req.DisputeID)
	if err != nil {
		return nil, fmt.Errorf("dispute not found: %w", err)
	}

	// Don't allow messages on resolved disputes
	if dispute.Status == domain.DisputeStatusResolvedBuyer || dispute.Status == domain.DisputeStatusResolvedSeller {
		return nil, fmt.Errorf("cannot add message to resolved dispute")
	}

	// If seller or buyer adds a message to an open dispute, move to under_review
	if dispute.Status == domain.DisputeStatusOpen {
		dispute.Status = domain.DisputeStatusUnderReview
		_ = uc.disputeRepo.Update(ctx, dispute)
	}

	msg := &domain.DisputeMessage{
		ID:          uuid.New().String(),
		DisputeID:   req.DisputeID,
		SenderID:    req.SenderID,
		SenderRole:  req.SenderRole,
		Message:     req.Message,
		Attachments: req.Attachments,
	}

	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return msg, nil
}

// ResolveDisputeRequest is the input for resolving a dispute.
type ResolveDisputeRequest struct {
	DisputeID  string
	Resolution string
	ResolvedBy string
	Status     string // resolved_buyer or resolved_seller
}

// ResolveDispute resolves a dispute (admin only).
func (uc *DisputeUseCase) ResolveDispute(ctx context.Context, req ResolveDisputeRequest) (*domain.Dispute, error) {
	dispute, err := uc.disputeRepo.GetByID(ctx, req.DisputeID)
	if err != nil {
		return nil, fmt.Errorf("dispute not found: %w", err)
	}

	if dispute.Status == domain.DisputeStatusResolvedBuyer || dispute.Status == domain.DisputeStatusResolvedSeller {
		return nil, fmt.Errorf("dispute is already resolved")
	}

	resolvedStatus := domain.DisputeStatus(req.Status)
	if resolvedStatus != domain.DisputeStatusResolvedBuyer && resolvedStatus != domain.DisputeStatusResolvedSeller {
		return nil, fmt.Errorf("invalid resolution status: must be resolved_buyer or resolved_seller")
	}

	now := time.Now()
	dispute.Status = resolvedStatus
	dispute.Resolution = req.Resolution
	dispute.ResolvedBy = req.ResolvedBy
	dispute.ResolvedAt = &now

	if err := uc.disputeRepo.Update(ctx, dispute); err != nil {
		return nil, fmt.Errorf("failed to resolve dispute: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "dispute.resolved", map[string]interface{}{
		"dispute_id": dispute.ID,
		"order_id":   dispute.OrderID,
		"status":     string(dispute.Status),
		"resolution": dispute.Resolution,
	})

	return dispute, nil
}
