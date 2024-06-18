package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/SerjRamone/vaultme/internal/models"
	pb "github.com/SerjRamone/vaultme/pkg/vaultme_v1"
)

const (
	userIDMetaKey = "USER_ID"
	maxFileSize   = 52428800 // 50 MB
	errForbidden  = "forbidden. error: %v"
)

// ItemServer - gRPC server for handling item requests
type ItemServer struct {
	pb.UnimplementedItemsServer
	log  *zap.Logger
	stor models.ItemStorage
}

// NewItemServer - returns new ItemServer instance
func NewItemServer(log *zap.Logger, storage models.ItemStorage) *ItemServer {
	return &ItemServer{
		log:  log,
		stor: storage,
	}
}

// GetItem returns single item
func (s *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	var resp pb.GetItemResponse
	uid, err := getUserIDFromMD(ctx)
	if err != nil {
		return &resp, status.Errorf(codes.Unauthenticated, fmt.Sprintf(errForbidden, err))
	}

	r, err := s.stor.GetItem(ctx, uid, req.GetId())
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("getting item error: %v", err))
	}

	item, err := convertItemToPB(r)
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("encode request error: %v", err))
	}

	resp.Item = item

	return &resp, nil
}

// CreateItem creates new item
func (s *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	var resp pb.CreateItemResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return &resp, status.Errorf(codes.Unauthenticated, fmt.Sprintf(errForbidden, err))
	}

	dat, err := convertItemDataTypeFromPB(req.Item.Data)
	if err != nil {
		return &resp, status.Errorf(codes.Internal,
			fmt.Sprintf("decode request error: %v", err))
	}

	idto, err := models.NewItemDTO(
		req.Item.GetName(),
		convertDataTypeFromPB(req.Item.GetType()),
		dat,
		convertMetaFromPB(req.Item.Meta),
	)
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("encode item error: %v", err))
	}

	if len(idto.Data) > maxFileSize {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("file too big. max size: %v", maxFileSize))
	}

	i, err := s.stor.CreateItem(ctx, userID, idto)
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("creating item error: %v", err))
	}

	resp.Id = i.ID
	return &resp, nil
}

// UpdateItem updates given item
func (s *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	var resp pb.UpdateItemResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return &resp, status.Errorf(codes.Unauthenticated, fmt.Sprintf(errForbidden, err))
	}

	i, err := convertItemFromPB(req.Item)
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("decode request error: %v", err))
	}

	i.Meta = convertMetaFromPB(req.Item.Meta)

	_, err = s.stor.UpdateItem(ctx, userID, i)
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("updating item error: %v", err))
	}

	resp.Id = i.ID
	return &resp, nil
}

// ListItems returns list of user's items
func (s *ItemServer) ListItems(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	var resp pb.ListItemResponse
	userID, err := getUserIDFromMD(ctx)
	if err != nil {
		return &resp, status.Errorf(codes.Unauthenticated, fmt.Sprintf(errForbidden, err))
	}

	items, err := s.stor.ListItems(ctx, userID, int(req.Limit), int(req.Offset))
	if err != nil {
		return &resp, status.Errorf(codes.Internal, fmt.Sprintf("get items batch error: %v", err))
	}

	for _, r := range items {
		i, err := convertItemToPB(r)
		if err != nil {
			return &resp, status.Errorf(codes.Internal, fmt.Sprintf("encode request error", err))
		}
		resp.Items = append(resp.Items, i)
	}

	return &resp, nil
}

// getUserIDFromMD returns user UUID from gRPC metadata
// checks that the field exists and is filled in
func getUserIDFromMD(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("empty metadata in gRPC request")
	}

	// slice of values
	ids := md.Get(userIDMetaKey)
	if len(ids) == 0 {
		return "", fmt.Errorf("meta with user's credentials are empty")
	}

	// get first value
	userID := ids[0]
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("user id is empty")
	}

	return userID, nil
}
