package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"gorm.io/gorm"
)

// ShipmentRepo implements domain.ShipmentRepository.
type ShipmentRepo struct {
	db *gorm.DB
}

// NewShipmentRepo creates a new ShipmentRepo.
func NewShipmentRepo(db *gorm.DB) *ShipmentRepo {
	return &ShipmentRepo{db: db}
}

func (r *ShipmentRepo) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	var model ShipmentModel
	if err := r.db.WithContext(ctx).Preload("Items").Preload("TrackingEvents").Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ShipmentRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	var models []ShipmentModel
	if err := r.db.WithContext(ctx).Preload("Items").Preload("TrackingEvents").Where("order_id = ?", orderID).Find(&models).Error; err != nil {
		return nil, err
	}
	shipments := make([]domain.Shipment, len(models))
	for i, m := range models {
		shipments[i] = *m.ToDomain()
	}
	return shipments, nil
}

func (r *ShipmentRepo) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	var model ShipmentModel
	if err := r.db.WithContext(ctx).Preload("Items").Preload("TrackingEvents").Where("tracking_number = ?", trackingNumber).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *ShipmentRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Shipment, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&ShipmentModel{}).Where("seller_id = ?", sellerID).Count(&total)

	var models []ShipmentModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Preload("Items").Where("seller_id = ?", sellerID).
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	shipments := make([]domain.Shipment, len(models))
	for i, m := range models {
		shipments[i] = *m.ToDomain()
	}
	return shipments, total, nil
}

func (r *ShipmentRepo) Create(ctx context.Context, shipment *domain.Shipment) error {
	model := ToShipmentModel(shipment)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *ShipmentRepo) Update(ctx context.Context, shipment *domain.Shipment) error {
	return r.db.WithContext(ctx).Model(&ShipmentModel{}).Where("id = ?", shipment.ID).Updates(map[string]interface{}{
		"carrier_code":    shipment.CarrierCode,
		"service_code":    shipment.ServiceCode,
		"tracking_number": shipment.TrackingNumber,
		"label_url":       shipment.LabelURL,
		"status":          string(shipment.Status),
		"weight_grams":    shipment.WeightGrams,
		"rate_cents":      shipment.RateCents,
	}).Error
}
