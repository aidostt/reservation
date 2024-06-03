package handler

import (
	"context"
	"dip/domain"
	"dip/internal/logger"
	"errors"
	proto_restaurant "github.com/aidostt/protos/gen/go/reservista/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAllRestaurants(ctx context.Context, input *proto_restaurant.Empty) (*proto_restaurant.RestaurantListResponse, error) {
	restaurants, err := h.service.Restaurants.GetAll(ctx)
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}
	restaurantResponse := make([]*proto_restaurant.RestaurantObject, len(restaurants))
	for index, restaurant := range restaurants {
		photoUrls := make([]string, len(restaurant.Photos))
		for i, photo := range restaurant.Photos {
			photoUrls[i] = photo.URl
		}
		restaurantResponse[index] = &proto_restaurant.RestaurantObject{
			Id:        restaurant.ID.String(),
			Name:      restaurant.Name,
			Address:   restaurant.Address,
			Contact:   restaurant.Contact,
			ImageUrls: photoUrls,
		}
	}
	return &proto_restaurant.RestaurantListResponse{
		Restaurants: restaurantResponse,
	}, nil
}

func (h *Handler) SearchRestaurants(ctx context.Context, input *proto_restaurant.SearchRequest) (*proto_restaurant.SearchResponse, error) {
	restaurants, total, err := h.service.Restaurants.Search(ctx, input.Query, int(input.Limit), int(input.Offset))
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}
	restaurantResponse := make([]*proto_restaurant.RestaurantObject, len(restaurants))
	for index, restaurant := range restaurants {
		photoUrls := make([]string, len(restaurant.Photos))
		for i, photo := range restaurant.Photos {
			photoUrls[i] = photo.URl
		}
		restaurantResponse[index] = &proto_restaurant.RestaurantObject{
			Id:        restaurant.ID.String(),
			Name:      restaurant.Name,
			Address:   restaurant.Address,
			Contact:   restaurant.Contact,
			ImageUrls: photoUrls,
		}
	}
	limit := int(input.Limit)
	if limit == 0 {
		limit = 10 // default limit
	}
	totalPages := (total + limit - 1) / limit

	return &proto_restaurant.SearchResponse{
		Restaurants: restaurantResponse,
		TotalPages:  int32(totalPages),
	}, nil
}

func (h *Handler) GetRestaurantSuggestions(ctx context.Context, input *proto_restaurant.SuggestionRequest) (*proto_restaurant.RestaurantListResponse, error) {
	restaurants, err := h.service.Restaurants.GetSuggestions(ctx, input.Query)
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	restaurantResponse := make([]*proto_restaurant.RestaurantObject, len(restaurants))
	for index, restaurant := range restaurants {
		restaurantResponse[index] = &proto_restaurant.RestaurantObject{
			Id:      restaurant.ID.String(),
			Name:    restaurant.Name,
			Address: restaurant.Address,
			Contact: restaurant.Contact,
		}
	}
	return &proto_restaurant.RestaurantListResponse{
		Restaurants: restaurantResponse,
	}, nil
}

func (h *Handler) GetRestaurant(ctx context.Context, input *proto_restaurant.IDRequest) (*proto_restaurant.RestaurantObject, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	res, err := h.service.Restaurants.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrNotFoundInDB):
			return nil, status.Error(codes.NotFound, "restaurant not found")
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}

	photoUrls := make([]string, len(res.Photos))
	for i, photo := range res.Photos {
		photoUrls[i] = photo.URl
	}

	return &proto_restaurant.RestaurantObject{
		Id:        res.ID.String(),
		Name:      res.Name,
		Address:   res.Address,
		Contact:   res.Contact,
		ImageUrls: photoUrls,
	}, nil
}

func (h *Handler) AddRestaurant(ctx context.Context, input *proto_restaurant.RestaurantObject) (*proto_restaurant.StatusResponse, error) {
	if input.GetName() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "name is required")
	}
	if input.GetAddress() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "address is required")
	}
	if input.GetContact() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "contact is required")
	}

	table := domain.RestaurantSql{
		Name:    input.GetName(),
		Address: input.GetAddress(),
		Contact: input.GetContact(),
	}
	if err := h.service.Restaurants.Create(ctx, &table); err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_restaurant.StatusResponse{Status: true}, nil
}

func (h *Handler) DeleteRestaurantById(ctx context.Context, input *proto_restaurant.IDRequest) (*proto_restaurant.StatusResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := h.service.Restaurants.DeleteById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_restaurant.StatusResponse{Status: true}, nil
}

func (h *Handler) UpdateRestById(ctx context.Context, input *proto_restaurant.RestaurantObject) (*proto_restaurant.StatusResponse, error) {
	if input.GetId() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "id is required")
	}
	if input.GetName() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "name is required")
	}
	if input.GetAddress() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "address is required")
	}
	if input.GetContact() == "" {
		return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "contact is required")
	}

	err := h.service.Restaurants.UpdateById(ctx, &domain.UpdateRestaurantInputSql{
		RestaurantId: input.GetId(),
		Name:         input.GetName(),
		Contact:      input.GetContact(),
		Address:      input.GetAddress(),
	})
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_restaurant.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_restaurant.StatusResponse{Status: true}, nil
}
