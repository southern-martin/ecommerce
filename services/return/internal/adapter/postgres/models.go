package postgres

import (
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/return/internal/domain"
)

// ReturnModel is the GORM model for the returns table.
type ReturnModel struct {
	ID                string         `gorm:"type:uuid;primaryKey"`
	OrderID           string         `gorm:"type:uuid;index;not null"`
	BuyerID           string         `gorm:"type:uuid;index;not null"`
	SellerID          string         `gorm:"type:uuid;index;not null"`
	Status            string         `gorm:"type:varchar(20);not null;default:'requested'"`
	Reason            string         `gorm:"type:varchar(50);not null"`
	Description       string         `gorm:"type:text"`
	ImageURLs         pq.StringArray `gorm:"type:text[]"`
	RefundAmountCents int64          `gorm:"default:0"`
	RefundMethod      string         `gorm:"type:varchar(20)"`
	ReturnTracking    string         `gorm:"type:varchar(100)"`
	Items             []ReturnItemModel `gorm:"foreignKey:ReturnID;constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
}

func (ReturnModel) TableName() string { return "returns" }

func (m *ReturnModel) ToDomain() *domain.Return {
	ret := &domain.Return{
		ID:                m.ID,
		OrderID:           m.OrderID,
		BuyerID:           m.BuyerID,
		SellerID:          m.SellerID,
		Status:            domain.ReturnStatus(m.Status),
		Reason:            domain.ReturnReason(m.Reason),
		Description:       m.Description,
		ImageURLs:         []string(m.ImageURLs),
		RefundAmountCents: m.RefundAmountCents,
		RefundMethod:      m.RefundMethod,
		ReturnTracking:    m.ReturnTracking,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
	for _, item := range m.Items {
		ret.Items = append(ret.Items, *item.ToDomain())
	}
	return ret
}

func ToReturnModel(r *domain.Return) *ReturnModel {
	model := &ReturnModel{
		ID:                r.ID,
		OrderID:           r.OrderID,
		BuyerID:           r.BuyerID,
		SellerID:          r.SellerID,
		Status:            string(r.Status),
		Reason:            string(r.Reason),
		Description:       r.Description,
		ImageURLs:         pq.StringArray(r.ImageURLs),
		RefundAmountCents: r.RefundAmountCents,
		RefundMethod:      r.RefundMethod,
		ReturnTracking:    r.ReturnTracking,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
	for _, item := range r.Items {
		model.Items = append(model.Items, *ToReturnItemModel(&item))
	}
	return model
}

// ReturnItemModel is the GORM model for the return_items table.
type ReturnItemModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	ReturnID    string `gorm:"type:uuid;index;not null"`
	OrderItemID string `gorm:"type:uuid;not null"`
	ProductID   string `gorm:"type:uuid;not null"`
	VariantID   string `gorm:"type:uuid"`
	Quantity    int    `gorm:"not null;default:1"`
	Reason      string `gorm:"type:varchar(50)"`
}

func (ReturnItemModel) TableName() string { return "return_items" }

func (m *ReturnItemModel) ToDomain() *domain.ReturnItem {
	return &domain.ReturnItem{
		ID:          m.ID,
		ReturnID:    m.ReturnID,
		OrderItemID: m.OrderItemID,
		ProductID:   m.ProductID,
		VariantID:   m.VariantID,
		Quantity:    m.Quantity,
		Reason:      m.Reason,
	}
}

func ToReturnItemModel(i *domain.ReturnItem) *ReturnItemModel {
	return &ReturnItemModel{
		ID:          i.ID,
		ReturnID:    i.ReturnID,
		OrderItemID: i.OrderItemID,
		ProductID:   i.ProductID,
		VariantID:   i.VariantID,
		Quantity:    i.Quantity,
		Reason:      i.Reason,
	}
}

// DisputeModel is the GORM model for the disputes table.
type DisputeModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	OrderID     string `gorm:"type:uuid;index;not null"`
	ReturnID    string `gorm:"type:uuid;index"`
	BuyerID     string `gorm:"type:uuid;index;not null"`
	SellerID    string `gorm:"type:uuid;index;not null"`
	Status      string `gorm:"type:varchar(20);not null;default:'open'"`
	Type        string `gorm:"type:varchar(30);not null"`
	Description string `gorm:"type:text;not null"`
	Resolution  string `gorm:"type:text"`
	ResolvedBy  string `gorm:"type:uuid"`
	Messages    []DisputeMessageModel `gorm:"foreignKey:DisputeID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	ResolvedAt  *time.Time `gorm:""`
}

func (DisputeModel) TableName() string { return "disputes" }

func (m *DisputeModel) ToDomain() *domain.Dispute {
	d := &domain.Dispute{
		ID:          m.ID,
		OrderID:     m.OrderID,
		ReturnID:    m.ReturnID,
		BuyerID:     m.BuyerID,
		SellerID:    m.SellerID,
		Status:      domain.DisputeStatus(m.Status),
		Type:        domain.DisputeType(m.Type),
		Description: m.Description,
		Resolution:  m.Resolution,
		ResolvedBy:  m.ResolvedBy,
		CreatedAt:   m.CreatedAt,
		ResolvedAt:  m.ResolvedAt,
	}
	for _, msg := range m.Messages {
		d.Messages = append(d.Messages, *msg.ToDomain())
	}
	return d
}

func ToDisputeModel(d *domain.Dispute) *DisputeModel {
	return &DisputeModel{
		ID:          d.ID,
		OrderID:     d.OrderID,
		ReturnID:    d.ReturnID,
		BuyerID:     d.BuyerID,
		SellerID:    d.SellerID,
		Status:      string(d.Status),
		Type:        string(d.Type),
		Description: d.Description,
		Resolution:  d.Resolution,
		ResolvedBy:  d.ResolvedBy,
		CreatedAt:   d.CreatedAt,
		ResolvedAt:  d.ResolvedAt,
	}
}

// DisputeMessageModel is the GORM model for the dispute_messages table.
type DisputeMessageModel struct {
	ID          string         `gorm:"type:uuid;primaryKey"`
	DisputeID   string         `gorm:"type:uuid;index;not null"`
	SenderID    string         `gorm:"type:uuid;not null"`
	SenderRole  string         `gorm:"type:varchar(10);not null"`
	Message     string         `gorm:"type:text;not null"`
	Attachments pq.StringArray `gorm:"type:text[]"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
}

func (DisputeMessageModel) TableName() string { return "dispute_messages" }

func (m *DisputeMessageModel) ToDomain() *domain.DisputeMessage {
	return &domain.DisputeMessage{
		ID:          m.ID,
		DisputeID:   m.DisputeID,
		SenderID:    m.SenderID,
		SenderRole:  m.SenderRole,
		Message:     m.Message,
		Attachments: []string(m.Attachments),
		CreatedAt:   m.CreatedAt,
	}
}

func ToDisputeMessageModel(msg *domain.DisputeMessage) *DisputeMessageModel {
	return &DisputeMessageModel{
		ID:          msg.ID,
		DisputeID:   msg.DisputeID,
		SenderID:    msg.SenderID,
		SenderRole:  msg.SenderRole,
		Message:     msg.Message,
		Attachments: pq.StringArray(msg.Attachments),
		CreatedAt:   msg.CreatedAt,
	}
}
