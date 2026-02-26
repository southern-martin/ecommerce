package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/southern-martin/ecommerce/services/shipping/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShippingService defines the gRPC service interface.
type ShippingService interface {
	GetShippingRates(ctx context.Context, req *GetRatesRequest) (*GetRatesResponse, error)
	CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*CreateShipmentResponse, error)
	GetShipment(ctx context.Context, req *GetShipmentRequest) (*GetShipmentResponse, error)
}

// --- Request/Response types ---

type GetRatesRequest struct {
	OriginCountry      string
	DestinationCountry string
	WeightGrams        int32
	Currency           string
}

type GetRatesResponse struct {
	Rates []*ShippingRateProto
}

type ShippingRateProto struct {
	CarrierCode      string
	ServiceName      string
	RateCents        int64
	Currency         string
	EstimatedDaysMin int32
	EstimatedDaysMax int32
}

type CreateShipmentRequest struct {
	OrderID     string
	SellerID    string
	CarrierCode string
	ServiceCode string
	WeightGrams int32
}

type CreateShipmentResponse struct {
	ID          string
	OrderID     string
	Status      string
	CarrierCode string
}

type GetShipmentRequest struct {
	ShipmentID string
}

type GetShipmentResponse struct {
	ID             string
	OrderID        string
	SellerID       string
	CarrierCode    string
	TrackingNumber string
	Status         string
	LabelURL       string
}

// Server implements the ShippingService gRPC interface.
type Server struct {
	rateUC     *usecase.RateUseCase
	shipmentUC *usecase.ShipmentUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(rateUC *usecase.RateUseCase, shipmentUC *usecase.ShipmentUseCase) *Server {
	return &Server{
		rateUC:     rateUC,
		shipmentUC: shipmentUC,
	}
}

func (s *Server) GetShippingRates(ctx context.Context, req *GetRatesRequest) (*GetRatesResponse, error) {
	if req.WeightGrams <= 0 {
		return nil, status.Error(codes.InvalidArgument, "weight_grams must be positive")
	}

	rates, err := s.rateUC.GetShippingRates(ctx, usecase.GetShippingRatesRequest{
		Destination: domain.Address{Country: req.DestinationCountry},
		Origin:      domain.Address{Country: req.OriginCountry},
		WeightGrams: int(req.WeightGrams),
		Currency:    req.Currency,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &GetRatesResponse{}
	for _, r := range rates {
		resp.Rates = append(resp.Rates, &ShippingRateProto{
			CarrierCode:      r.CarrierCode,
			ServiceName:      r.ServiceName,
			RateCents:        r.RateCents,
			Currency:         r.Currency,
			EstimatedDaysMin: int32(r.EstimatedDaysMin),
			EstimatedDaysMax: int32(r.EstimatedDaysMax),
		})
	}
	return resp, nil
}

func (s *Server) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*CreateShipmentResponse, error) {
	if req.OrderID == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	shipment, err := s.shipmentUC.CreateShipment(ctx, usecase.CreateShipmentRequest{
		OrderID:     req.OrderID,
		SellerID:    req.SellerID,
		CarrierCode: req.CarrierCode,
		ServiceCode: req.ServiceCode,
		WeightGrams: int(req.WeightGrams),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateShipmentResponse{
		ID:          shipment.ID,
		OrderID:     shipment.OrderID,
		Status:      string(shipment.Status),
		CarrierCode: shipment.CarrierCode,
	}, nil
}

func (s *Server) GetShipment(ctx context.Context, req *GetShipmentRequest) (*GetShipmentResponse, error) {
	if req.ShipmentID == "" {
		return nil, status.Error(codes.InvalidArgument, "shipment_id is required")
	}

	shipment, err := s.shipmentUC.GetShipment(ctx, req.ShipmentID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetShipmentResponse{
		ID:             shipment.ID,
		OrderID:        shipment.OrderID,
		SellerID:       shipment.SellerID,
		CarrierCode:    shipment.CarrierCode,
		TrackingNumber: shipment.TrackingNumber,
		Status:         string(shipment.Status),
		LabelURL:       shipment.LabelURL,
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerGetShippingRates(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetRatesRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ShippingService).GetShippingRates(ctx, req)
}

func handlerCreateShipment(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &CreateShipmentRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ShippingService).CreateShipment(ctx, req)
}

func handlerGetShipment(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetShipmentRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ShippingService).GetShipment(ctx, req)
}

// ShippingServiceDesc is the gRPC service descriptor.
var ShippingServiceDesc = grpc.ServiceDesc{
	ServiceName: "shipping.ShippingService",
	HandlerType: (*ShippingService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetShippingRates", Handler: handlerGetShippingRates},
		{MethodName: "CreateShipment", Handler: handlerCreateShipment},
		{MethodName: "GetShipment", Handler: handlerGetShipment},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("shipping_service.proto"),
}

// RegisterShippingServiceServer registers the ShippingService with a gRPC server.
func RegisterShippingServiceServer(s *grpc.Server, srv ShippingService) {
	s.RegisterService(&ShippingServiceDesc, srv)
}
