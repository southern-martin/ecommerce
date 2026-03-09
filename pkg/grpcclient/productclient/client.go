// Package productclient provides a typed gRPC client for the product service.
// Request and response types mirror those defined in the product service's
// gRPC server (services/product/internal/adapter/grpc/server.go).
package productclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client provides typed access to the product gRPC service.
type Client struct {
	conn *grpc.ClientConn
}

// New creates a new product service gRPC client.
func New(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

// GetProduct fetches a product by ID.
func (c *Client) GetProduct(ctx context.Context, productID string) (*GetProductResponse, error) {
	req := &GetProductRequest{ProductID: productID}
	resp := &GetProductResponse{}
	err := c.conn.Invoke(ctx, "/product.ProductService/GetProduct", req, resp)
	if err != nil {
		return nil, fmt.Errorf("productclient.GetProduct: %w", err)
	}
	return resp, nil
}

// GetVariant fetches a variant by ID.
func (c *Client) GetVariant(ctx context.Context, variantID string) (*GetVariantResponse, error) {
	req := &GetVariantRequest{VariantID: variantID}
	resp := &GetVariantResponse{}
	err := c.conn.Invoke(ctx, "/product.ProductService/GetVariant", req, resp)
	if err != nil {
		return nil, fmt.Errorf("productclient.GetVariant: %w", err)
	}
	return resp, nil
}

// UpdateStock updates variant stock by delta.
func (c *Client) UpdateStock(ctx context.Context, variantID string, delta int32) (*UpdateStockResponse, error) {
	req := &UpdateStockRequest{VariantID: variantID, Delta: delta}
	resp := &UpdateStockResponse{}
	err := c.conn.Invoke(ctx, "/product.ProductService/UpdateStock", req, resp)
	if err != nil {
		return nil, fmt.Errorf("productclient.UpdateStock: %w", err)
	}
	return resp, nil
}

// ListVariantsByProduct fetches all variants for a product.
func (c *Client) ListVariantsByProduct(ctx context.Context, productID string) (*ListVariantsByProductResponse, error) {
	req := &ListVariantsByProductRequest{ProductID: productID}
	resp := &ListVariantsByProductResponse{}
	err := c.conn.Invoke(ctx, "/product.ProductService/ListVariantsByProduct", req, resp)
	if err != nil {
		return nil, fmt.Errorf("productclient.ListVariantsByProduct: %w", err)
	}
	return resp, nil
}

// --- Request/Response types matching the server's types ---

// GetProductRequest mirrors the server's GetProductRequest.
type GetProductRequest struct {
	ProductID string `json:"product_id"`
}

// GetProductResponse mirrors the server's GetProductResponse.
type GetProductResponse struct {
	ID             string `json:"id"`
	SellerID       string `json:"seller_id"`
	CategoryID     string `json:"category_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	BasePriceCents int64  `json:"base_price_cents"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
	HasVariants    bool   `json:"has_variants"`
}

// GetVariantRequest mirrors the server's GetVariantRequest.
type GetVariantRequest struct {
	VariantID string `json:"variant_id"`
}

// GetVariantResponse mirrors the server's GetVariantResponse.
type GetVariantResponse struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Stock      int    `json:"stock"`
	IsActive   bool   `json:"is_active"`
}

// UpdateStockRequest mirrors the server's UpdateStockRequest.
type UpdateStockRequest struct {
	VariantID string `json:"variant_id"`
	Delta     int32  `json:"delta"`
}

// UpdateStockResponse mirrors the server's UpdateStockResponse.
type UpdateStockResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ListVariantsByProductRequest mirrors the server's ListVariantsByProductRequest.
type ListVariantsByProductRequest struct {
	ProductID string `json:"product_id"`
}

// VariantInfo mirrors the server's VariantInfo.
type VariantInfo struct {
	ID         string `json:"id"`
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Stock      int    `json:"stock"`
	IsActive   bool   `json:"is_active"`
}

// ListVariantsByProductResponse mirrors the server's ListVariantsByProductResponse.
type ListVariantsByProductResponse struct {
	Variants []VariantInfo `json:"variants"`
}
