// Package orderclient provides a typed gRPC client for the order service.
// Request and response types mirror those defined in the order service's
// gRPC server (services/order/internal/adapter/grpc/server.go).
package orderclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client provides typed access to the order gRPC service.
type Client struct {
	conn *grpc.ClientConn
}

// New creates a new order service gRPC client.
func New(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

// GetOrder fetches an order by ID.
func (c *Client) GetOrder(ctx context.Context, orderID string) (*GetOrderResponse, error) {
	req := &GetOrderRequest{OrderID: orderID}
	resp := &GetOrderResponse{}
	err := c.conn.Invoke(ctx, "/order.OrderService/GetOrder", req, resp)
	if err != nil {
		return nil, fmt.Errorf("orderclient.GetOrder: %w", err)
	}
	return resp, nil
}

// UpdateOrderStatus updates an order's status.
func (c *Client) UpdateOrderStatus(ctx context.Context, orderID, newStatus string) (*UpdateOrderStatusResponse, error) {
	req := &UpdateOrderStatusRequest{OrderID: orderID, NewStatus: newStatus}
	resp := &UpdateOrderStatusResponse{}
	err := c.conn.Invoke(ctx, "/order.OrderService/UpdateOrderStatus", req, resp)
	if err != nil {
		return nil, fmt.Errorf("orderclient.UpdateOrderStatus: %w", err)
	}
	return resp, nil
}

// --- Request/Response types matching the server's types ---

// GetOrderRequest mirrors the server's GetOrderRequest.
type GetOrderRequest struct {
	OrderID string `json:"order_id"`
}

// GetOrderResponse mirrors the server's GetOrderResponse.
type GetOrderResponse struct {
	ID            string           `json:"id"`
	OrderNumber   string           `json:"order_number"`
	BuyerID       string           `json:"buyer_id"`
	Status        string           `json:"status"`
	SubtotalCents int64            `json:"subtotal_cents"`
	ShippingCents int64            `json:"shipping_cents"`
	TaxCents      int64            `json:"tax_cents"`
	DiscountCents int64            `json:"discount_cents"`
	TotalCents    int64            `json:"total_cents"`
	Currency      string           `json:"currency"`
	Items         []*OrderItemInfo `json:"items"`
}

// OrderItemInfo mirrors the server's OrderItemProto.
type OrderItemInfo struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	VariantID      string `json:"variant_id"`
	ProductName    string `json:"product_name"`
	Quantity       int32  `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	TotalCents     int64  `json:"total_cents"`
	SellerID       string `json:"seller_id"`
}

// UpdateOrderStatusRequest mirrors the server's UpdateOrderStatusRequest.
type UpdateOrderStatusRequest struct {
	OrderID   string `json:"order_id"`
	NewStatus string `json:"new_status"`
}

// UpdateOrderStatusResponse mirrors the server's UpdateOrderStatusResponse.
type UpdateOrderStatusResponse struct {
	ID          string `json:"id"`
	OrderNumber string `json:"order_number"`
	Status      string `json:"status"`
}
