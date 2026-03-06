package usecase

import (
	"context"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// --- mockCarrierRepo ---

type mockCarrierRepo struct {
	getAllFn     func(ctx context.Context) ([]domain.Carrier, error)
	getByCodeFn func(ctx context.Context, code string) (*domain.Carrier, error)
	createFn    func(ctx context.Context, carrier *domain.Carrier) error
	updateFn    func(ctx context.Context, carrier *domain.Carrier) error
}

func (m *mockCarrierRepo) GetAll(ctx context.Context) ([]domain.Carrier, error) {
	return m.getAllFn(ctx)
}

func (m *mockCarrierRepo) GetByCode(ctx context.Context, code string) (*domain.Carrier, error) {
	return m.getByCodeFn(ctx, code)
}

func (m *mockCarrierRepo) Create(ctx context.Context, carrier *domain.Carrier) error {
	return m.createFn(ctx, carrier)
}

func (m *mockCarrierRepo) Update(ctx context.Context, carrier *domain.Carrier) error {
	return m.updateFn(ctx, carrier)
}

// --- mockCarrierCredentialRepo ---

type mockCarrierCredentialRepo struct {
	getBySellerAndCarrierFn func(ctx context.Context, sellerID, carrierCode string) (*domain.CarrierCredential, error)
	listBySellerFn          func(ctx context.Context, sellerID string) ([]domain.CarrierCredential, error)
	createFn                func(ctx context.Context, cred *domain.CarrierCredential) error
	updateFn                func(ctx context.Context, cred *domain.CarrierCredential) error
	deleteFn                func(ctx context.Context, id string) error
}

func (m *mockCarrierCredentialRepo) GetBySellerAndCarrier(ctx context.Context, sellerID, carrierCode string) (*domain.CarrierCredential, error) {
	return m.getBySellerAndCarrierFn(ctx, sellerID, carrierCode)
}

func (m *mockCarrierCredentialRepo) ListBySeller(ctx context.Context, sellerID string) ([]domain.CarrierCredential, error) {
	return m.listBySellerFn(ctx, sellerID)
}

func (m *mockCarrierCredentialRepo) Create(ctx context.Context, cred *domain.CarrierCredential) error {
	return m.createFn(ctx, cred)
}

func (m *mockCarrierCredentialRepo) Update(ctx context.Context, cred *domain.CarrierCredential) error {
	return m.updateFn(ctx, cred)
}

func (m *mockCarrierCredentialRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

// --- mockShipmentRepo ---

type mockShipmentRepo struct {
	getByIDFn             func(ctx context.Context, id string) (*domain.Shipment, error)
	getByOrderIDFn        func(ctx context.Context, orderID string) ([]domain.Shipment, error)
	getByTrackingNumberFn func(ctx context.Context, trackingNumber string) (*domain.Shipment, error)
	listBySellerFn        func(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Shipment, int64, error)
	createFn              func(ctx context.Context, shipment *domain.Shipment) error
	updateFn              func(ctx context.Context, shipment *domain.Shipment) error
}

func (m *mockShipmentRepo) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockShipmentRepo) GetByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	return m.getByOrderIDFn(ctx, orderID)
}

func (m *mockShipmentRepo) GetByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	return m.getByTrackingNumberFn(ctx, trackingNumber)
}

func (m *mockShipmentRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Shipment, int64, error) {
	return m.listBySellerFn(ctx, sellerID, page, pageSize)
}

func (m *mockShipmentRepo) Create(ctx context.Context, shipment *domain.Shipment) error {
	return m.createFn(ctx, shipment)
}

func (m *mockShipmentRepo) Update(ctx context.Context, shipment *domain.Shipment) error {
	return m.updateFn(ctx, shipment)
}

// --- mockTrackingEventRepo ---

type mockTrackingEventRepo struct {
	getByShipmentIDFn func(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error)
	createFn          func(ctx context.Context, event *domain.TrackingEvent) error
}

func (m *mockTrackingEventRepo) GetByShipmentID(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error) {
	return m.getByShipmentIDFn(ctx, shipmentID)
}

func (m *mockTrackingEventRepo) Create(ctx context.Context, event *domain.TrackingEvent) error {
	return m.createFn(ctx, event)
}

// --- mockEventPublisher ---

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	return m.publishFn(ctx, subject, data)
}
