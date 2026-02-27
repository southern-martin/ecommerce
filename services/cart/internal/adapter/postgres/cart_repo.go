package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

// CartModel is the GORM model for durable cart persistence.
type CartModel struct {
	UserID    string    `gorm:"type:varchar(36);primaryKey"`
	CartData  string    `gorm:"type:jsonb;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (CartModel) TableName() string {
	return "carts"
}

// postgresCartRepo implements domain.CartRepository using PostgreSQL.
type postgresCartRepo struct {
	db *gorm.DB
}

// NewPostgresCartRepository creates a new Postgres-backed cart repository.
func NewPostgresCartRepository(db *gorm.DB) domain.CartRepository {
	return &postgresCartRepo{db: db}
}

func (r *postgresCartRepo) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	var model CartModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return &domain.Cart{
			UserID:    userID,
			Items:     []domain.CartItem{},
			UpdatedAt: time.Now().UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres get cart: %w", err)
	}

	var cart domain.Cart
	if err := json.Unmarshal([]byte(model.CartData), &cart); err != nil {
		return nil, fmt.Errorf("unmarshal cart from postgres: %w", err)
	}

	if cart.Items == nil {
		cart.Items = []domain.CartItem{}
	}

	return &cart, nil
}

func (r *postgresCartRepo) SaveCart(ctx context.Context, cart *domain.Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("marshal cart for postgres: %w", err)
	}

	model := CartModel{
		UserID:    cart.UserID,
		CartData:  string(data),
		UpdatedAt: cart.UpdatedAt,
	}

	// Upsert: insert or update on conflict
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"cart_data", "updated_at"}),
	}).Create(&model)

	if result.Error != nil {
		return fmt.Errorf("postgres save cart: %w", result.Error)
	}

	return nil
}

func (r *postgresCartRepo) DeleteCart(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&CartModel{})
	if result.Error != nil {
		return fmt.Errorf("postgres delete cart: %w", result.Error)
	}
	return nil
}
