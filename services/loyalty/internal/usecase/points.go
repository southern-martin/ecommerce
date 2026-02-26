package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/domain"
)

// PointsUseCase handles points operations.
type PointsUseCase struct {
	membershipRepo domain.MembershipRepository
	txRepo         domain.PointsTransactionRepository
	membershipUC   *MembershipUseCase
	publisher      domain.EventPublisher
}

// NewPointsUseCase creates a new PointsUseCase.
func NewPointsUseCase(
	membershipRepo domain.MembershipRepository,
	txRepo domain.PointsTransactionRepository,
	membershipUC *MembershipUseCase,
	publisher domain.EventPublisher,
) *PointsUseCase {
	return &PointsUseCase{
		membershipRepo: membershipRepo,
		txRepo:         txRepo,
		membershipUC:   membershipUC,
		publisher:      publisher,
	}
}

// EarnPointsRequest is the input for earning points.
type EarnPointsRequest struct {
	UserID      string
	Points      int64
	Source      domain.PointsSource
	ReferenceID string
	Description string
}

// EarnPoints adds points to a user's balance and lifetime total.
func (uc *PointsUseCase) EarnPoints(ctx context.Context, req EarnPointsRequest) (*domain.PointsTransaction, error) {
	membership, err := uc.membershipRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	tx := &domain.PointsTransaction{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Type:        domain.TransactionEarn,
		Points:      req.Points,
		Source:      req.Source,
		ReferenceID: req.ReferenceID,
		Description: req.Description,
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	newBalance := membership.PointsBalance + req.Points
	newLifetime := membership.LifetimePoints + req.Points
	if err := uc.membershipRepo.UpdatePoints(ctx, req.UserID, newBalance, newLifetime); err != nil {
		return nil, fmt.Errorf("failed to update points: %w", err)
	}

	// Check for tier upgrade
	_ = uc.membershipUC.CheckAndUpgradeTier(ctx, req.UserID)

	_ = uc.publisher.Publish(ctx, "loyalty.points.earned", map[string]interface{}{
		"user_id":      req.UserID,
		"points":       req.Points,
		"source":       string(req.Source),
		"reference_id": req.ReferenceID,
	})

	return tx, nil
}

// RedeemPointsRequest is the input for redeeming points.
type RedeemPointsRequest struct {
	UserID      string
	Points      int64
	OrderID     string
	Description string
}

// RedeemPoints subtracts points from a user's balance.
func (uc *PointsUseCase) RedeemPoints(ctx context.Context, req RedeemPointsRequest) (*domain.PointsTransaction, error) {
	membership, err := uc.membershipRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	if membership.PointsBalance < req.Points {
		return nil, fmt.Errorf("insufficient points: have %d, need %d", membership.PointsBalance, req.Points)
	}

	tx := &domain.PointsTransaction{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Type:        domain.TransactionRedeem,
		Points:      req.Points,
		Source:      domain.SourceOrder,
		ReferenceID: req.OrderID,
		Description: req.Description,
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	newBalance := membership.PointsBalance - req.Points
	if err := uc.membershipRepo.UpdatePoints(ctx, req.UserID, newBalance, membership.LifetimePoints); err != nil {
		return nil, fmt.Errorf("failed to update points: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "loyalty.points.redeemed", map[string]interface{}{
		"user_id":  req.UserID,
		"points":   req.Points,
		"order_id": req.OrderID,
	})

	return tx, nil
}

// GetBalance retrieves a user's points balance.
func (uc *PointsUseCase) GetBalance(ctx context.Context, userID string) (int64, error) {
	membership, err := uc.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return membership.PointsBalance, nil
}

// ListTransactions lists a user's points transactions with pagination.
func (uc *PointsUseCase) ListTransactions(ctx context.Context, userID string, page, pageSize int) ([]domain.PointsTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.txRepo.ListByUser(ctx, userID, page, pageSize)
}
