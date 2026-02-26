package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/southern-martin/ecommerce/services/tax/internal/usecase"
	"google.golang.org/grpc"
)

// TaxService defines the gRPC service interface.
type TaxService interface {
	CalculateTax(ctx context.Context, req *CalculateTaxRequest) (*CalculateTaxResponse, error)
	GetTaxRules(ctx context.Context, req *GetTaxRulesRequest) (*GetTaxRulesResponse, error)
}

// --- Request / Response types for gRPC ---

// CalculateTaxRequest is the gRPC request for tax calculation.
type CalculateTaxRequest struct {
	Items           []*TaxItemMessage
	ShippingAddress *TaxAddressMessage
}

// TaxItemMessage represents a tax item in gRPC messages.
type TaxItemMessage struct {
	ProductID  string
	VariantID  string
	Category   string
	PriceCents int64
	Quantity   int32
}

// TaxAddressMessage represents a tax address in gRPC messages.
type TaxAddressMessage struct {
	CountryCode string
	StateCode   string
	City        string
	PostalCode  string
}

// CalculateTaxResponse is the gRPC response for tax calculation.
type CalculateTaxResponse struct {
	SubtotalCents  int64
	TaxAmountCents int64
	Breakdown      []*TaxBreakdownMessage
}

// TaxBreakdownMessage represents a tax breakdown item in gRPC messages.
type TaxBreakdownMessage struct {
	TaxName      string
	Rate         float64
	AmountCents  int64
	Jurisdiction string
}

// GetTaxRulesRequest is the gRPC request for getting tax rules.
type GetTaxRulesRequest struct {
	CountryCode string
	StateCode   string
}

// GetTaxRulesResponse is the gRPC response for getting tax rules.
type GetTaxRulesResponse struct {
	Rules []*TaxRuleMessage
}

// TaxRuleMessage represents a tax rule in gRPC messages.
type TaxRuleMessage struct {
	ID        string
	ZoneID    string
	TaxName   string
	Rate      float64
	Category  string
	Inclusive bool
	IsActive  bool
}

// Server implements the TaxService gRPC interface.
type Server struct {
	calculateTax *usecase.CalculateTaxUseCase
	manageRules  *usecase.ManageRulesUseCase
	manageZones  *usecase.ManageZonesUseCase
}

// NewServer creates a new gRPC server.
func NewServer(
	calculateTax *usecase.CalculateTaxUseCase,
	manageRules *usecase.ManageRulesUseCase,
	manageZones *usecase.ManageZonesUseCase,
) *Server {
	return &Server{
		calculateTax: calculateTax,
		manageRules:  manageRules,
		manageZones:  manageZones,
	}
}

// CalculateTax handles the gRPC CalculateTax call.
func (s *Server) CalculateTax(ctx context.Context, req *CalculateTaxRequest) (*CalculateTaxResponse, error) {
	if req.ShippingAddress == nil {
		return nil, fmt.Errorf("shipping address is required")
	}

	domainReq := &domain.TaxCalculationRequest{
		Items: make([]domain.TaxItem, len(req.Items)),
		ShippingAddress: domain.TaxAddress{
			CountryCode: req.ShippingAddress.CountryCode,
			StateCode:   req.ShippingAddress.StateCode,
			City:        req.ShippingAddress.City,
			PostalCode:  req.ShippingAddress.PostalCode,
		},
	}

	for i, item := range req.Items {
		domainReq.Items[i] = domain.TaxItem{
			ProductID:  item.ProductID,
			VariantID:  item.VariantID,
			Category:   item.Category,
			PriceCents: item.PriceCents,
			Quantity:   int(item.Quantity),
		}
	}

	result, err := s.calculateTax.Execute(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	breakdown := make([]*TaxBreakdownMessage, len(result.Breakdown))
	for i, b := range result.Breakdown {
		breakdown[i] = &TaxBreakdownMessage{
			TaxName:      b.TaxName,
			Rate:         b.Rate,
			AmountCents:  b.AmountCents,
			Jurisdiction: b.Jurisdiction,
		}
	}

	return &CalculateTaxResponse{
		SubtotalCents:  result.SubtotalCents,
		TaxAmountCents: result.TaxAmountCents,
		Breakdown:      breakdown,
	}, nil
}

// GetTaxRules handles the gRPC GetTaxRules call.
func (s *Server) GetTaxRules(ctx context.Context, req *GetTaxRulesRequest) (*GetTaxRulesResponse, error) {
	zone, err := s.manageZones.GetZoneByLocation(ctx, req.CountryCode, req.StateCode)
	if err != nil {
		return &GetTaxRulesResponse{Rules: []*TaxRuleMessage{}}, nil
	}

	rules, err := s.manageRules.ListRulesByZone(ctx, zone.ID)
	if err != nil {
		return nil, err
	}

	ruleMessages := make([]*TaxRuleMessage, len(rules))
	for i, r := range rules {
		ruleMessages[i] = &TaxRuleMessage{
			ID:        r.ID,
			ZoneID:    r.ZoneID,
			TaxName:   r.TaxName,
			Rate:      r.Rate,
			Category:  r.Category,
			Inclusive: r.Inclusive,
			IsActive:  r.IsActive,
		}
	}

	return &GetTaxRulesResponse{Rules: ruleMessages}, nil
}

// ServiceDesc returns the gRPC ServiceDesc for manual registration.
var ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tax.TaxService",
	HandlerType: (*TaxService)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CalculateTax",
			Handler:    calculateTaxHandler,
		},
		{
			MethodName: "GetTaxRules",
			Handler:    getTaxRulesHandler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

func calculateTaxHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(CalculateTaxRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaxService).CalculateTax(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tax.TaxService/CalculateTax",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaxService).CalculateTax(ctx, req.(*CalculateTaxRequest))
	}
	return interceptor(ctx, req, info, handler)
}

func getTaxRulesHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	req := new(GetTaxRulesRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaxService).GetTaxRules(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tax.TaxService/GetTaxRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaxService).GetTaxRules(ctx, req.(*GetTaxRulesRequest))
	}
	return interceptor(ctx, req, info, handler)
}
