// Package userclient provides a typed gRPC client for the user service.
// Request and response types mirror those defined in the user service's
// gRPC server (services/user/internal/adapter/grpc/server.go).
package userclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client provides typed access to the user gRPC service.
type Client struct {
	conn *grpc.ClientConn
}

// New creates a new user service gRPC client.
func New(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn}
}

// GetProfile fetches a user profile by user ID.
func (c *Client) GetProfile(ctx context.Context, userID string) (*UserProfileResponse, error) {
	req := &GetProfileRequest{UserID: userID}
	resp := &UserProfileResponse{}
	err := c.conn.Invoke(ctx, "/user.UserService/GetProfile", req, resp)
	if err != nil {
		return nil, fmt.Errorf("userclient.GetProfile: %w", err)
	}
	return resp, nil
}

// GetSellerProfile fetches a seller profile by seller ID.
func (c *Client) GetSellerProfile(ctx context.Context, sellerID string) (*SellerProfileResponse, error) {
	req := &GetSellerProfileRequest{SellerID: sellerID}
	resp := &SellerProfileResponse{}
	err := c.conn.Invoke(ctx, "/user.UserService/GetSellerProfile", req, resp)
	if err != nil {
		return nil, fmt.Errorf("userclient.GetSellerProfile: %w", err)
	}
	return resp, nil
}

// --- Request/Response types matching the server's types ---

// GetProfileRequest mirrors the server's GetProfileRequest.
type GetProfileRequest struct {
	UserID string `json:"user_id"`
}

// UserProfileResponse mirrors the server's UserProfileResponse.
type UserProfileResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Role        string `json:"role"`
}

// GetSellerProfileRequest mirrors the server's GetSellerProfileRequest.
type GetSellerProfileRequest struct {
	SellerID string `json:"seller_id"`
}

// SellerProfileResponse mirrors the server's SellerProfileResponse.
type SellerProfileResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	StoreName   string  `json:"store_name"`
	Description string  `json:"description"`
	LogoURL     string  `json:"logo_url"`
	Rating      float64 `json:"rating"`
	TotalSales  int     `json:"total_sales"`
	Status      string  `json:"status"`
}
