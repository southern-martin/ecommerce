package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
)

// ManageRulesUseCase handles CRUD operations for TaxRules.
type ManageRulesUseCase struct {
	ruleRepo domain.TaxRuleRepository
}

// NewManageRulesUseCase creates a new ManageRulesUseCase.
func NewManageRulesUseCase(ruleRepo domain.TaxRuleRepository) *ManageRulesUseCase {
	return &ManageRulesUseCase{ruleRepo: ruleRepo}
}

// CreateRuleInput holds the input for creating a tax rule.
type CreateRuleInput struct {
	ZoneID    string
	TaxName   string
	Rate      float64
	Category  string
	Inclusive bool
	StartsAt  *time.Time
	ExpiresAt *time.Time
}

// CreateRule creates a new tax rule.
func (uc *ManageRulesUseCase) CreateRule(ctx context.Context, input CreateRuleInput) (*domain.TaxRule, error) {
	startsAt := time.Now()
	if input.StartsAt != nil {
		startsAt = *input.StartsAt
	}

	rule := &domain.TaxRule{
		ID:        uuid.New().String(),
		ZoneID:    input.ZoneID,
		TaxName:   input.TaxName,
		Rate:      input.Rate,
		Category:  input.Category,
		Inclusive:  input.Inclusive,
		StartsAt:  startsAt,
		ExpiresAt: input.ExpiresAt,
		IsActive:  true,
	}

	if err := uc.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// UpdateRuleInput holds the input for updating a tax rule.
type UpdateRuleInput struct {
	TaxName   *string
	Rate      *float64
	Category  *string
	Inclusive *bool
	IsActive  *bool
	ExpiresAt *time.Time
}

// UpdateRule updates an existing tax rule.
func (uc *ManageRulesUseCase) UpdateRule(ctx context.Context, id string, input UpdateRuleInput) (*domain.TaxRule, error) {
	rule, err := uc.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.TaxName != nil {
		rule.TaxName = *input.TaxName
	}
	if input.Rate != nil {
		rule.Rate = *input.Rate
	}
	if input.Category != nil {
		rule.Category = *input.Category
	}
	if input.Inclusive != nil {
		rule.Inclusive = *input.Inclusive
	}
	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}
	if input.ExpiresAt != nil {
		rule.ExpiresAt = input.ExpiresAt
	}

	if err := uc.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteRule deletes a tax rule by ID.
func (uc *ManageRulesUseCase) DeleteRule(ctx context.Context, id string) error {
	return uc.ruleRepo.Delete(ctx, id)
}

// ListRules lists all active tax rules.
func (uc *ManageRulesUseCase) ListRules(ctx context.Context) ([]*domain.TaxRule, error) {
	return uc.ruleRepo.ListActive(ctx)
}

// ListRulesByZone lists all tax rules for a specific zone.
func (uc *ManageRulesUseCase) ListRulesByZone(ctx context.Context, zoneID string) ([]*domain.TaxRule, error) {
	return uc.ruleRepo.ListByZone(ctx, zoneID)
}
