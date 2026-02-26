package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

var (
	ErrInvalidUserID  = errors.New("user ID is required")
	ErrInvalidProduct = errors.New("product ID is required")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	ErrItemNotFound   = errors.New("item not found in cart")
)

// CartUseCase implements cart business logic.
type CartUseCase struct {
	repo      domain.CartRepository
	publisher domain.EventPublisher
	logger    zerolog.Logger
}

// NewCartUseCase creates a new CartUseCase.
func NewCartUseCase(repo domain.CartRepository, publisher domain.EventPublisher, logger zerolog.Logger) *CartUseCase {
	return &CartUseCase{
		repo:      repo,
		publisher: publisher,
		logger:    logger.With().Str("component", "cart_usecase").Logger(),
	}
}

// AddItem adds or increments an item in the cart. If the item already exists
// (matched by productID + variantID), the quantity is incremented.
func (uc *CartUseCase) AddItem(ctx context.Context, userID string, item domain.CartItem) (*domain.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if item.ProductID == "" {
		return nil, ErrInvalidProduct
	}
	if item.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart")
		return nil, err
	}

	idx := cart.FindItem(item.ProductID, item.VariantID)
	if idx >= 0 {
		cart.Items[idx].Quantity += item.Quantity
		// Update denormalized fields in case they changed
		cart.Items[idx].PriceCents = item.PriceCents
		cart.Items[idx].ProductName = item.ProductName
		cart.Items[idx].VariantName = item.VariantName
		cart.Items[idx].ImageURL = item.ImageURL
		cart.Items[idx].SKU = item.SKU
		cart.Items[idx].SellerID = item.SellerID
	} else {
		cart.Items = append(cart.Items, item)
	}

	cart.UpdatedAt = time.Now().UTC()

	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to save cart")
		return nil, err
	}

	uc.publishEvent(ctx, domain.EventCartItemAdded, domain.CartEvent{
		UserID:    userID,
		ProductID: item.ProductID,
		VariantID: item.VariantID,
		Quantity:  item.Quantity,
	})

	uc.logger.Info().Str("user_id", userID).Str("product_id", item.ProductID).Msg("item added to cart")
	return cart, nil
}

// RemoveItem removes an item from the cart by productID and variantID.
func (uc *CartUseCase) RemoveItem(ctx context.Context, userID, productID, variantID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if productID == "" {
		return nil, ErrInvalidProduct
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart")
		return nil, err
	}

	idx := cart.FindItem(productID, variantID)
	if idx < 0 {
		return nil, ErrItemNotFound
	}

	cart.Items = append(cart.Items[:idx], cart.Items[idx+1:]...)
	cart.UpdatedAt = time.Now().UTC()

	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to save cart")
		return nil, err
	}

	uc.publishEvent(ctx, domain.EventCartItemRemoved, domain.CartEvent{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
	})

	uc.logger.Info().Str("user_id", userID).Str("product_id", productID).Msg("item removed from cart")
	return cart, nil
}

// UpdateQuantity sets the quantity for a specific item in the cart.
func (uc *CartUseCase) UpdateQuantity(ctx context.Context, userID, productID, variantID string, quantity int) (*domain.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if productID == "" {
		return nil, ErrInvalidProduct
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart")
		return nil, err
	}

	idx := cart.FindItem(productID, variantID)
	if idx < 0 {
		return nil, ErrItemNotFound
	}

	cart.Items[idx].Quantity = quantity
	cart.UpdatedAt = time.Now().UTC()

	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to save cart")
		return nil, err
	}

	uc.publishEvent(ctx, domain.EventCartItemUpdated, domain.CartEvent{
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
		Quantity:  quantity,
	})

	uc.logger.Info().Str("user_id", userID).Str("product_id", productID).Int("quantity", quantity).Msg("cart item quantity updated")
	return cart, nil
}

// GetCart returns the user's cart.
func (uc *CartUseCase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart")
		return nil, err
	}

	return cart, nil
}

// ClearCart removes all items from the user's cart.
func (uc *CartUseCase) ClearCart(ctx context.Context, userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}

	if err := uc.repo.DeleteCart(ctx, userID); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to clear cart")
		return err
	}

	uc.publishEvent(ctx, domain.EventCartCleared, domain.CartEvent{
		UserID: userID,
	})

	uc.logger.Info().Str("user_id", userID).Msg("cart cleared")
	return nil
}

// MergeCart merges guest cart items into the authenticated user's cart.
// If an item already exists, its quantity is incremented.
func (uc *CartUseCase) MergeCart(ctx context.Context, userID string, guestItems []domain.CartItem) (*domain.Cart, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	cart, err := uc.repo.GetCart(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get cart for merge")
		return nil, err
	}

	for _, guestItem := range guestItems {
		if guestItem.ProductID == "" || guestItem.Quantity <= 0 {
			continue
		}

		idx := cart.FindItem(guestItem.ProductID, guestItem.VariantID)
		if idx >= 0 {
			cart.Items[idx].Quantity += guestItem.Quantity
			cart.Items[idx].PriceCents = guestItem.PriceCents
			cart.Items[idx].ProductName = guestItem.ProductName
			cart.Items[idx].VariantName = guestItem.VariantName
			cart.Items[idx].ImageURL = guestItem.ImageURL
			cart.Items[idx].SKU = guestItem.SKU
			cart.Items[idx].SellerID = guestItem.SellerID
		} else {
			cart.Items = append(cart.Items, guestItem)
		}
	}

	cart.UpdatedAt = time.Now().UTC()

	if err := uc.repo.SaveCart(ctx, cart); err != nil {
		uc.logger.Error().Err(err).Str("user_id", userID).Msg("failed to save merged cart")
		return nil, err
	}

	uc.logger.Info().Str("user_id", userID).Int("merged_items", len(guestItems)).Msg("cart merged")
	return cart, nil
}

// publishEvent publishes a domain event, logging any errors.
func (uc *CartUseCase) publishEvent(ctx context.Context, subject string, event domain.CartEvent) {
	if uc.publisher == nil {
		return
	}
	if err := uc.publisher.Publish(ctx, subject, event); err != nil {
		uc.logger.Error().Err(err).Str("subject", subject).Msg("failed to publish event")
	}
}
