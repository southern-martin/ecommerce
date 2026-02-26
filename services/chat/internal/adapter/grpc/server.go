package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/chat/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChatService defines the gRPC service interface.
type ChatService interface {
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
	GetConversation(ctx context.Context, req *GetConversationRequest) (*GetConversationResponse, error)
	GetUnreadCount(ctx context.Context, req *GetUnreadCountRequest) (*GetUnreadCountResponse, error)
}

// --- Request/Response types ---

type SendMessageRequest struct {
	ConversationID string
	SenderID       string
	SenderRole     string
	Content        string
	MessageType    string
}

type SendMessageResponse struct {
	ID             string
	ConversationID string
	SenderID       string
	Content        string
	MessageType    string
}

type GetConversationRequest struct {
	ConversationID string
}

type GetConversationResponse struct {
	ID       string
	Type     string
	BuyerID  string
	SellerID string
	OrderID  string
	Subject  string
	Status   string
}

type GetUnreadCountRequest struct {
	ConversationID string
	UserID         string
}

type GetUnreadCountResponse struct {
	Count int64
}

// Server implements the ChatService gRPC interface.
type Server struct {
	conversationUC *usecase.ConversationUseCase
	messageUC      *usecase.MessageUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(conversationUC *usecase.ConversationUseCase, messageUC *usecase.MessageUseCase) *Server {
	return &Server{
		conversationUC: conversationUC,
		messageUC:      messageUC,
	}
}

func (s *Server) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if req.ConversationID == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	message, err := s.messageUC.SendMessage(ctx, usecase.SendMessageRequest{
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		SenderRole:     req.SenderRole,
		Content:        req.Content,
		MessageType:    req.MessageType,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &SendMessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Content:        message.Content,
		MessageType:    string(message.MessageType),
	}, nil
}

func (s *Server) GetConversation(ctx context.Context, req *GetConversationRequest) (*GetConversationResponse, error) {
	if req.ConversationID == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}

	conv, err := s.conversationUC.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetConversationResponse{
		ID:       conv.ID,
		Type:     string(conv.Type),
		BuyerID:  conv.BuyerID,
		SellerID: conv.SellerID,
		OrderID:  conv.OrderID,
		Subject:  conv.Subject,
		Status:   string(conv.Status),
	}, nil
}

func (s *Server) GetUnreadCount(ctx context.Context, req *GetUnreadCountRequest) (*GetUnreadCountResponse, error) {
	if req.ConversationID == "" || req.UserID == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id and user_id are required")
	}

	count, err := s.messageUC.GetUnreadCount(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetUnreadCountResponse{Count: count}, nil
}

// --- gRPC ServiceDesc ---

func handlerSendMessage(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &SendMessageRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ChatService).SendMessage(ctx, req)
}

func handlerGetConversation(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetConversationRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ChatService).GetConversation(ctx, req)
}

func handlerGetUnreadCount(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetUnreadCountRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(ChatService).GetUnreadCount(ctx, req)
}

// ChatServiceDesc is the gRPC service descriptor.
var ChatServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "SendMessage", Handler: handlerSendMessage},
		{MethodName: "GetConversation", Handler: handlerGetConversation},
		{MethodName: "GetUnreadCount", Handler: handlerGetUnreadCount},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("chat_service.proto"),
}

// RegisterChatServiceServer registers the ChatService with a gRPC server.
func RegisterChatServiceServer(s *grpc.Server, srv ChatService) {
	s.RegisterService(&ChatServiceDesc, srv)
}
