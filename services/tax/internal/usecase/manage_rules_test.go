package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Tests for ManageRulesUseCase.CreateRule
// ---------------------------------------------------------------------------

func TestManageRules_CreateRule(t *testing.T) {
	ctx := context.Background()

	t.Run("success defaults StartsAt to now if nil", func(t *testing.T) {
		before := time.Now()
		var captured *domain.TaxRule

		ruleRepo := &mockTaxRuleRepo{
			createFn: func(_ context.Context, rule *domain.TaxRule) error {
				captured = rule
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		input := CreateRuleInput{
			ZoneID:   "zone-1",
			TaxName:  "GST",
			Rate:     0.10,
			Category: "electronics",
			Inclusive: true,
			StartsAt: nil, // should default to now
		}

		rule, err := uc.CreateRule(ctx, input)
		after := time.Now()

		require.NoError(t, err)
		require.NotNil(t, rule)
		assert.NotEmpty(t, rule.ID)
		assert.Equal(t, "zone-1", rule.ZoneID)
		assert.Equal(t, "GST", rule.TaxName)
		assert.Equal(t, 0.10, rule.Rate)
		assert.Equal(t, "electronics", rule.Category)
		assert.True(t, rule.Inclusive)
		assert.True(t, rule.IsActive)
		assert.Nil(t, rule.ExpiresAt)
		// StartsAt should be between before and after
		assert.False(t, rule.StartsAt.Before(before), "StartsAt should not be before test start")
		assert.False(t, rule.StartsAt.After(after), "StartsAt should not be after test end")

		// Verify mock was called
		require.NotNil(t, captured)
		assert.Equal(t, rule.ID, captured.ID)
	})

	t.Run("success with explicit StartsAt", func(t *testing.T) {
		customStart := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

		ruleRepo := &mockTaxRuleRepo{
			createFn: func(_ context.Context, _ *domain.TaxRule) error {
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		input := CreateRuleInput{
			ZoneID:   "zone-2",
			TaxName:  "VAT",
			Rate:     0.20,
			Category: "",
			Inclusive: false,
			StartsAt: &customStart,
		}

		rule, err := uc.CreateRule(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, customStart, rule.StartsAt)
	})

	t.Run("success with all fields including ExpiresAt", func(t *testing.T) {
		customStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		customExpiry := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

		ruleRepo := &mockTaxRuleRepo{
			createFn: func(_ context.Context, _ *domain.TaxRule) error {
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		input := CreateRuleInput{
			ZoneID:    "zone-3",
			TaxName:   "Holiday Tax",
			Rate:      0.05,
			Category:  "gifts",
			Inclusive:  false,
			StartsAt:  &customStart,
			ExpiresAt: &customExpiry,
		}

		rule, err := uc.CreateRule(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, rule.ExpiresAt)
		assert.Equal(t, customExpiry, *rule.ExpiresAt)
		assert.Equal(t, customStart, rule.StartsAt)
		assert.True(t, rule.IsActive)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			createFn: func(_ context.Context, _ *domain.TaxRule) error {
				return errors.New("db connection failed")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		input := CreateRuleInput{
			ZoneID:  "zone-1",
			TaxName: "GST",
			Rate:    0.10,
		}

		rule, err := uc.CreateRule(ctx, input)
		require.Error(t, err)
		assert.Nil(t, rule)
		assert.Contains(t, err.Error(), "db connection failed")
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageRulesUseCase.UpdateRule
// ---------------------------------------------------------------------------

func TestManageRules_UpdateRule(t *testing.T) {
	ctx := context.Background()

	existingRule := func() *domain.TaxRule {
		return &domain.TaxRule{
			ID:       "rule-1",
			ZoneID:   "zone-1",
			TaxName:  "GST",
			Rate:     0.10,
			Category: "electronics",
			Inclusive: true,
			IsActive: true,
		}
	}

	t.Run("success with all fields updated", func(t *testing.T) {
		var updatedRule *domain.TaxRule

		ruleRepo := &mockTaxRuleRepo{
			getByIDFn: func(_ context.Context, id string) (*domain.TaxRule, error) {
				assert.Equal(t, "rule-1", id)
				return existingRule(), nil
			},
			updateFn: func(_ context.Context, rule *domain.TaxRule) error {
				updatedRule = rule
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		newName := "Updated GST"
		newRate := 0.15
		newCategory := "all"
		newInclusive := false
		newActive := false
		newExpiry := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		input := UpdateRuleInput{
			TaxName:   &newName,
			Rate:      &newRate,
			Category:  &newCategory,
			Inclusive:  &newInclusive,
			IsActive:  &newActive,
			ExpiresAt: &newExpiry,
		}

		result, err := uc.UpdateRule(ctx, "rule-1", input)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Updated GST", result.TaxName)
		assert.Equal(t, 0.15, result.Rate)
		assert.Equal(t, "all", result.Category)
		assert.False(t, result.Inclusive)
		assert.False(t, result.IsActive)
		require.NotNil(t, result.ExpiresAt)
		assert.Equal(t, newExpiry, *result.ExpiresAt)

		// Ensure update was called
		require.NotNil(t, updatedRule)
	})

	t.Run("partial update only rate", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			getByIDFn: func(_ context.Context, _ string) (*domain.TaxRule, error) {
				return existingRule(), nil
			},
			updateFn: func(_ context.Context, _ *domain.TaxRule) error {
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		newRate := 0.12
		input := UpdateRuleInput{
			Rate: &newRate,
		}

		result, err := uc.UpdateRule(ctx, "rule-1", input)
		require.NoError(t, err)
		// Rate should be updated
		assert.Equal(t, 0.12, result.Rate)
		// All other fields should remain unchanged
		assert.Equal(t, "GST", result.TaxName)
		assert.Equal(t, "electronics", result.Category)
		assert.True(t, result.Inclusive)
		assert.True(t, result.IsActive)
	})

	t.Run("not found returns error", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			getByIDFn: func(_ context.Context, _ string) (*domain.TaxRule, error) {
				return nil, errors.New("rule not found")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		newRate := 0.12
		input := UpdateRuleInput{Rate: &newRate}

		result, err := uc.UpdateRule(ctx, "nonexistent", input)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "rule not found")
	})

	t.Run("repo update error propagates", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			getByIDFn: func(_ context.Context, _ string) (*domain.TaxRule, error) {
				return existingRule(), nil
			},
			updateFn: func(_ context.Context, _ *domain.TaxRule) error {
				return errors.New("update failed")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		newRate := 0.12
		input := UpdateRuleInput{Rate: &newRate}

		result, err := uc.UpdateRule(ctx, "rule-1", input)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update failed")
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageRulesUseCase.DeleteRule
// ---------------------------------------------------------------------------

func TestManageRules_DeleteRule(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		var deletedID string
		ruleRepo := &mockTaxRuleRepo{
			deleteFn: func(_ context.Context, id string) error {
				deletedID = id
				return nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		err := uc.DeleteRule(ctx, "rule-42")
		require.NoError(t, err)
		assert.Equal(t, "rule-42", deletedID)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			deleteFn: func(_ context.Context, _ string) error {
				return errors.New("delete failed")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		err := uc.DeleteRule(ctx, "rule-42")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageRulesUseCase.ListRules
// ---------------------------------------------------------------------------

func TestManageRules_ListRules(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns active rules", func(t *testing.T) {
		expected := []*domain.TaxRule{
			{ID: "r1", TaxName: "GST", Rate: 0.10, IsActive: true},
			{ID: "r2", TaxName: "VAT", Rate: 0.20, IsActive: true},
		}
		ruleRepo := &mockTaxRuleRepo{
			listActiveFn: func(_ context.Context) ([]*domain.TaxRule, error) {
				return expected, nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		rules, err := uc.ListRules(ctx)
		require.NoError(t, err)
		assert.Len(t, rules, 2)
		assert.Equal(t, "r1", rules[0].ID)
		assert.Equal(t, "r2", rules[1].ID)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			listActiveFn: func(_ context.Context) ([]*domain.TaxRule, error) {
				return nil, errors.New("list failed")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		rules, err := uc.ListRules(ctx)
		require.Error(t, err)
		assert.Nil(t, rules)
	})
}

// ---------------------------------------------------------------------------
// Tests for ManageRulesUseCase.ListRulesByZone
// ---------------------------------------------------------------------------

func TestManageRules_ListRulesByZone(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns rules for zone", func(t *testing.T) {
		expected := []*domain.TaxRule{
			{ID: "r1", ZoneID: "zone-au", TaxName: "GST", Rate: 0.10},
		}
		ruleRepo := &mockTaxRuleRepo{
			listByZoneFn: func(_ context.Context, zoneID string) ([]*domain.TaxRule, error) {
				assert.Equal(t, "zone-au", zoneID)
				return expected, nil
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		rules, err := uc.ListRulesByZone(ctx, "zone-au")
		require.NoError(t, err)
		assert.Len(t, rules, 1)
		assert.Equal(t, "zone-au", rules[0].ZoneID)
	})

	t.Run("repo error propagates", func(t *testing.T) {
		ruleRepo := &mockTaxRuleRepo{
			listByZoneFn: func(_ context.Context, _ string) ([]*domain.TaxRule, error) {
				return nil, errors.New("zone rules query failed")
			},
		}
		uc := NewManageRulesUseCase(ruleRepo)

		rules, err := uc.ListRulesByZone(ctx, "zone-au")
		require.Error(t, err)
		assert.Nil(t, rules)
	})
}
