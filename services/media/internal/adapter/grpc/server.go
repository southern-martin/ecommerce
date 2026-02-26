package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/media/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MediaService defines the gRPC service interface.
type MediaService interface {
	GetMedia(ctx context.Context, req *GetMediaRequest) (*GetMediaResponse, error)
}

// --- Request/Response types ---

type GetMediaRequest struct {
	MediaID string
}

type GetMediaResponse struct {
	ID           string
	OwnerID      string
	OwnerType    string
	FileName     string
	OriginalName string
	ContentType  string
	SizeBytes    int64
	URL          string
	ThumbnailURL string
	Width        int32
	Height       int32
	Status       string
}

// Server implements the MediaService gRPC interface.
type Server struct {
	mediaUC *usecase.MediaUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(mediaUC *usecase.MediaUseCase) *Server {
	return &Server{
		mediaUC: mediaUC,
	}
}

func (s *Server) GetMedia(ctx context.Context, req *GetMediaRequest) (*GetMediaResponse, error) {
	if req.MediaID == "" {
		return nil, status.Error(codes.InvalidArgument, "media_id is required")
	}

	media, err := s.mediaUC.GetMedia(ctx, req.MediaID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetMediaResponse{
		ID:           media.ID,
		OwnerID:      media.OwnerID,
		OwnerType:    media.OwnerType,
		FileName:     media.FileName,
		OriginalName: media.OriginalName,
		ContentType:  media.ContentType,
		SizeBytes:    media.SizeBytes,
		URL:          media.URL,
		ThumbnailURL: media.ThumbnailURL,
		Width:        int32(media.Width),
		Height:       int32(media.Height),
		Status:       string(media.Status),
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerGetMedia(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetMediaRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(MediaService).GetMedia(ctx, req)
}

// MediaServiceDesc is the gRPC service descriptor.
var MediaServiceDesc = grpc.ServiceDesc{
	ServiceName: "media.MediaService",
	HandlerType: (*MediaService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetMedia", Handler: handlerGetMedia},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("media_service.proto"),
}

// RegisterMediaServiceServer registers the MediaService with a gRPC server.
func RegisterMediaServiceServer(s *grpc.Server, srv MediaService) {
	s.RegisterService(&MediaServiceDesc, srv)
}
