package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/southern-martin/ecommerce/services/product/internal/domain"
)

// Publisher implements domain.EventPublisher using NATS.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher creates a new NATS event publisher.
func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to NATS")
	return &Publisher{conn: nc}, nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

// ProductCreatedEvent is the payload for product.created events.
type ProductCreatedEvent struct {
	ID         string `json:"id"`
	SellerID   string `json:"seller_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	CategoryID string `json:"category_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

// ProductUpdatedEvent is the payload for product.updated events.
type ProductUpdatedEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

// ProductDeletedEvent is the payload for product.deleted events.
type ProductDeletedEvent struct {
	ID        string `json:"id"`
	DeletedAt string `json:"deleted_at"`
}

// StockUpdatedEvent is the payload for product.stock.updated events.
type StockUpdatedEvent struct {
	VariantID string `json:"variant_id"`
	NewStock  int    `json:"new_stock"`
	Delta     int    `json:"delta"`
	UpdatedAt string `json:"updated_at"`
}

func (p *Publisher) publish(subject string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, bytes)
}

// PublishProductCreated publishes a product.created event.
func (p *Publisher) PublishProductCreated(_ context.Context, product *domain.Product) error {
	event := ProductCreatedEvent{
		ID:         product.ID,
		SellerID:   product.SellerID,
		Name:       product.Name,
		Slug:       product.Slug,
		CategoryID: product.CategoryID,
		Status:     string(product.Status),
		CreatedAt:  product.CreatedAt.Format(time.RFC3339),
	}
	if err := p.publish("product.created", event); err != nil {
		log.Error().Err(err).Str("product_id", product.ID).Msg("Failed to publish product.created event")
		return err
	}
	log.Debug().Str("product_id", product.ID).Msg("Published product.created event")
	return nil
}

// PublishProductUpdated publishes a product.updated event.
func (p *Publisher) PublishProductUpdated(_ context.Context, product *domain.Product) error {
	event := ProductUpdatedEvent{
		ID:        product.ID,
		Name:      product.Name,
		Status:    string(product.Status),
		UpdatedAt: product.UpdatedAt.Format(time.RFC3339),
	}
	if err := p.publish("product.updated", event); err != nil {
		log.Error().Err(err).Str("product_id", product.ID).Msg("Failed to publish product.updated event")
		return err
	}
	log.Debug().Str("product_id", product.ID).Msg("Published product.updated event")
	return nil
}

// PublishProductDeleted publishes a product.deleted event.
func (p *Publisher) PublishProductDeleted(_ context.Context, productID string) error {
	event := ProductDeletedEvent{
		ID:        productID,
		DeletedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := p.publish("product.deleted", event); err != nil {
		log.Error().Err(err).Str("product_id", productID).Msg("Failed to publish product.deleted event")
		return err
	}
	log.Debug().Str("product_id", productID).Msg("Published product.deleted event")
	return nil
}

// PublishStockUpdated publishes a product.stock.updated event.
func (p *Publisher) PublishStockUpdated(_ context.Context, variantID string, newStock int, delta int) error {
	event := StockUpdatedEvent{
		VariantID: variantID,
		NewStock:  newStock,
		Delta:     delta,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := p.publish("product.stock.updated", event); err != nil {
		log.Error().Err(err).Str("variant_id", variantID).Msg("Failed to publish product.stock.updated event")
		return err
	}
	log.Debug().Str("variant_id", variantID).Int("delta", delta).Msg("Published product.stock.updated event")
	return nil
}
