package domain

import "context"

// ReturnRepository defines the interface for return persistence.
type ReturnRepository interface {
	GetByID(ctx context.Context, id string) (*Return, error)
	GetByOrderID(ctx context.Context, orderID string) ([]Return, error)
	ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]Return, int64, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]Return, int64, error)
	Create(ctx context.Context, ret *Return) error
	Update(ctx context.Context, ret *Return) error
}

// DisputeRepository defines the interface for dispute persistence.
type DisputeRepository interface {
	GetByID(ctx context.Context, id string) (*Dispute, error)
	GetByOrderID(ctx context.Context, orderID string) ([]Dispute, error)
	ListAll(ctx context.Context, page, pageSize int) ([]Dispute, int64, error)
	ListByBuyer(ctx context.Context, buyerID string, page, pageSize int) ([]Dispute, int64, error)
	Create(ctx context.Context, dispute *Dispute) error
	Update(ctx context.Context, dispute *Dispute) error
}

// DisputeMessageRepository defines the interface for dispute message persistence.
type DisputeMessageRepository interface {
	GetByDisputeID(ctx context.Context, disputeID string) ([]DisputeMessage, error)
	Create(ctx context.Context, msg *DisputeMessage) error
}
