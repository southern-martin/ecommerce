package usecase

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

// ProgramUseCase handles affiliate program operations.
type ProgramUseCase struct {
	programRepo domain.AffiliateProgramRepository
}

// NewProgramUseCase creates a new ProgramUseCase.
func NewProgramUseCase(programRepo domain.AffiliateProgramRepository) *ProgramUseCase {
	return &ProgramUseCase{
		programRepo: programRepo,
	}
}

// GetProgram retrieves the active affiliate program.
func (uc *ProgramUseCase) GetProgram(ctx context.Context) (*domain.AffiliateProgram, error) {
	return uc.programRepo.Get(ctx)
}

// UpdateProgramRequest is the input for updating the affiliate program.
type UpdateProgramRequest struct {
	ID                 string
	CommissionRate     *float64
	MinPayoutCents     *int64
	CookieDays         *int
	ReferrerBonusCents *int64
	ReferredBonusCents *int64
	IsActive           *bool
}

// UpdateProgram updates the affiliate program settings.
func (uc *ProgramUseCase) UpdateProgram(ctx context.Context, req UpdateProgramRequest) (*domain.AffiliateProgram, error) {
	program, err := uc.programRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	if req.CommissionRate != nil {
		program.CommissionRate = *req.CommissionRate
	}
	if req.MinPayoutCents != nil {
		program.MinPayoutCents = *req.MinPayoutCents
	}
	if req.CookieDays != nil {
		program.CookieDays = *req.CookieDays
	}
	if req.ReferrerBonusCents != nil {
		program.ReferrerBonusCents = *req.ReferrerBonusCents
	}
	if req.ReferredBonusCents != nil {
		program.ReferredBonusCents = *req.ReferredBonusCents
	}
	if req.IsActive != nil {
		program.IsActive = *req.IsActive
	}

	if err := uc.programRepo.Update(ctx, program); err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}

	return program, nil
}
