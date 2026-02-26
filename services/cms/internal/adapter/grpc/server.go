package grpc

import (
	"context"
	"fmt"

	"github.com/southern-martin/ecommerce/services/cms/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CMSService defines the gRPC service interface.
type CMSService interface {
	ListBanners(ctx context.Context, req *ListBannersRequest) (*ListBannersResponse, error)
	GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error)
}

// --- Request/Response types ---

type ListBannersRequest struct {
	Position string
}

type ListBannersResponse struct {
	Banners []*BannerProto
}

type BannerProto struct {
	ID       string
	Title    string
	ImageURL string
	LinkURL  string
	Position string
}

type GetPageRequest struct {
	Slug string
}

type GetPageResponse struct {
	ID              string
	Title           string
	Slug            string
	ContentHTML     string
	MetaTitle       string
	MetaDescription string
	Status          string
}

// Server implements the CMSService gRPC interface.
type Server struct {
	bannerUC *usecase.BannerUseCase
	pageUC   *usecase.PageUseCase
}

// NewServer creates a new gRPC Server.
func NewServer(bannerUC *usecase.BannerUseCase, pageUC *usecase.PageUseCase) *Server {
	return &Server{
		bannerUC: bannerUC,
		pageUC:   pageUC,
	}
}

func (s *Server) ListBanners(ctx context.Context, req *ListBannersRequest) (*ListBannersResponse, error) {
	banners, err := s.bannerUC.ListActiveBanners(ctx, req.Position)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &ListBannersResponse{}
	for _, b := range banners {
		resp.Banners = append(resp.Banners, &BannerProto{
			ID:       b.ID,
			Title:    b.Title,
			ImageURL: b.ImageURL,
			LinkURL:  b.LinkURL,
			Position: b.Position,
		})
	}
	return resp, nil
}

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	page, err := s.pageUC.GetPageBySlug(ctx, req.Slug)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &GetPageResponse{
		ID:              page.ID,
		Title:           page.Title,
		Slug:            page.Slug,
		ContentHTML:     page.ContentHTML,
		MetaTitle:       page.MetaTitle,
		MetaDescription: page.MetaDescription,
		Status:          string(page.Status),
	}, nil
}

// --- gRPC ServiceDesc ---

func handlerListBanners(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &ListBannersRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(CMSService).ListBanners(ctx, req)
}

func handlerGetPage(srv interface{}, ctx context.Context, dec func(interface{}) error, _ grpc.UnaryServerInterceptor) (interface{}, error) {
	req := &GetPageRequest{}
	if err := dec(req); err != nil {
		return nil, err
	}
	return srv.(CMSService).GetPage(ctx, req)
}

// CMSServiceDesc is the gRPC service descriptor.
var CMSServiceDesc = grpc.ServiceDesc{
	ServiceName: "cms.CMSService",
	HandlerType: (*CMSService)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ListBanners", Handler: handlerListBanners},
		{MethodName: "GetPage", Handler: handlerGetPage},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fmt.Sprintf("cms_service.proto"),
}

// RegisterCMSServiceServer registers the CMSService with a gRPC server.
func RegisterCMSServiceServer(s *grpc.Server, srv CMSService) {
	s.RegisterService(&CMSServiceDesc, srv)
}
