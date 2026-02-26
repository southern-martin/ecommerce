package domain

import "time"

// --- Return ---

// ReturnStatus represents the status of a return.
type ReturnStatus string

const (
	ReturnStatusRequested  ReturnStatus = "requested"
	ReturnStatusApproved   ReturnStatus = "approved"
	ReturnStatusRejected   ReturnStatus = "rejected"
	ReturnStatusShippedBack ReturnStatus = "shipped_back"
	ReturnStatusReceived   ReturnStatus = "received"
	ReturnStatusRefunded   ReturnStatus = "refunded"
)

// AllowedReturnTransitions defines valid return status transitions.
var AllowedReturnTransitions = map[ReturnStatus][]ReturnStatus{
	ReturnStatusRequested:   {ReturnStatusApproved, ReturnStatusRejected},
	ReturnStatusApproved:    {ReturnStatusShippedBack},
	ReturnStatusShippedBack: {ReturnStatusReceived},
	ReturnStatusReceived:    {ReturnStatusRefunded},
}

// CanReturnTransition checks if a return status transition is valid.
func CanReturnTransition(from, to ReturnStatus) bool {
	allowed, ok := AllowedReturnTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// ReturnReason represents the reason for a return.
type ReturnReason string

const (
	ReasonDefective      ReturnReason = "defective"
	ReasonWrongItem      ReturnReason = "wrong_item"
	ReasonNotAsDescribed ReturnReason = "not_as_described"
	ReasonChangedMind    ReturnReason = "changed_mind"
	ReasonOther          ReturnReason = "other"
)

// Return represents a return request.
type Return struct {
	ID                string
	OrderID           string
	BuyerID           string
	SellerID          string
	Status            ReturnStatus
	Reason            ReturnReason
	Description       string
	ImageURLs         []string
	Items             []ReturnItem
	RefundAmountCents int64
	RefundMethod      string // original_payment, wallet_credit
	ReturnTracking    string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ReturnItem represents an item in a return request.
type ReturnItem struct {
	ID          string
	ReturnID    string
	OrderItemID string
	ProductID   string
	VariantID   string
	Quantity    int
	Reason      string
}

// --- Dispute ---

// DisputeStatus represents the status of a dispute.
type DisputeStatus string

const (
	DisputeStatusOpen           DisputeStatus = "open"
	DisputeStatusUnderReview    DisputeStatus = "under_review"
	DisputeStatusResolvedBuyer  DisputeStatus = "resolved_buyer"
	DisputeStatusResolvedSeller DisputeStatus = "resolved_seller"
	DisputeStatusEscalated      DisputeStatus = "escalated"
)

// DisputeType represents the type of dispute.
type DisputeType string

const (
	DisputeTypeItemNotReceived    DisputeType = "item_not_received"
	DisputeTypeNotAsDescribed     DisputeType = "not_as_described"
	DisputeTypeUnauthorizedCharge DisputeType = "unauthorized_charge"
)

// Dispute represents a dispute.
type Dispute struct {
	ID          string
	OrderID     string
	ReturnID    string
	BuyerID     string
	SellerID    string
	Status      DisputeStatus
	Type        DisputeType
	Description string
	Messages    []DisputeMessage
	Resolution  string
	ResolvedBy  string
	CreatedAt   time.Time
	ResolvedAt  *time.Time
}

// DisputeMessage represents a message in a dispute thread.
type DisputeMessage struct {
	ID          string
	DisputeID   string
	SenderID    string
	SenderRole  string // buyer, seller, admin
	Message     string
	Attachments []string
	CreatedAt   time.Time
}
