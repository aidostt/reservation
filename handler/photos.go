package handler

import (
	"context"
	"dip/domain"
	"dip/internal/logger"
	proto_restaurant "github.com/aidostt/protos/gen/go/reservista/restaurant"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UploadPhotos(ctx context.Context, input *proto_restaurant.UploadPhotoRequest) (*proto_restaurant.StatusResponse, error) {
	if input.GetRestaurantID() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "restaurant id is required")
	}
	if input.GetUrls() == nil {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "image urls is required")
	}
	var photos []*domain.PhotoSql
	RestaurantUUID, err := uuid.FromString(input.GetRestaurantID())
	if err != nil {
		return nil, err
	}
	for _, url := range input.GetUrls() {
		photos = append(photos, &domain.PhotoSql{RestaurantID: RestaurantUUID, URl: url})
	}
	if err := h.service.Photos.Upload(ctx, photos); err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_restaurant.StatusResponse{Status: true}, nil
}

func (h *Handler) Delete(ctx context.Context, input *proto_restaurant.DeletePhotoRequest) (*proto_restaurant.StatusResponse, error) {
	if input.GetRestaurantID() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "restaurant id is required")
	}
	if input.GetUrl() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "image url is required")
	}

	if err := h.service.Photos.Delete(ctx, input.GetUrl(), input.GetRestaurantID()); err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_restaurant.StatusResponse{Status: true}, nil
}
