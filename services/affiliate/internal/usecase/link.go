package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/domain"
)

const (
	codeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength  = 8
)

// LinkUseCase handles affiliate link operations.
type LinkUseCase struct {
	linkRepo  domain.AffiliateLinkRepository
	publisher domain.EventPublisher
}

// NewLinkUseCase creates a new LinkUseCase.
func NewLinkUseCase(linkRepo domain.AffiliateLinkRepository, publisher domain.EventPublisher) *LinkUseCase {
	return &LinkUseCase{
		linkRepo:  linkRepo,
		publisher: publisher,
	}
}

// CreateLinkRequest is the input for creating an affiliate link.
type CreateLinkRequest struct {
	UserID    string
	TargetURL string
}

// CreateLink creates a new affiliate link with a unique referral code.
func (uc *LinkUseCase) CreateLink(ctx context.Context, req CreateLinkRequest) (*domain.AffiliateLink, error) {
	code, err := generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate referral code: %w", err)
	}

	link := &domain.AffiliateLink{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Code:      code,
		TargetURL: req.TargetURL,
	}

	if err := uc.linkRepo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("failed to create affiliate link: %w", err)
	}

	return link, nil
}

// GetLink retrieves an affiliate link by ID.
func (uc *LinkUseCase) GetLink(ctx context.Context, id string) (*domain.AffiliateLink, error) {
	return uc.linkRepo.GetByID(ctx, id)
}

// ListUserLinks lists affiliate links for a user with pagination.
func (uc *LinkUseCase) ListUserLinks(ctx context.Context, userID string, page, pageSize int) ([]domain.AffiliateLink, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.linkRepo.ListByUser(ctx, userID, page, pageSize)
}

// TrackClick increments the click count for an affiliate link and publishes an event.
func (uc *LinkUseCase) TrackClick(ctx context.Context, code string) (*domain.AffiliateLink, error) {
	link, err := uc.linkRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("affiliate link not found: %w", err)
	}

	if err := uc.linkRepo.IncrementClicks(ctx, link.ID); err != nil {
		return nil, fmt.Errorf("failed to increment clicks: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "affiliate.click.tracked", map[string]interface{}{
		"link_id": link.ID,
		"user_id": link.UserID,
		"code":    link.Code,
	})

	link.ClickCount++
	return link, nil
}

// generateCode generates a random 8-character alphanumeric string using crypto/rand.
func generateCode() (string, error) {
	result := make([]byte, codeLength)
	charsetLen := big.NewInt(int64(len(codeCharset)))
	for i := range result {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = codeCharset[n.Int64()]
	}
	return string(result), nil
}
