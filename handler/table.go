package handler

import (
	"context"
	"dip/domain"
	"dip/internal/logger"
	"errors"
	proto_table "github.com/aidostt/protos/gen/go/reservista/table"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAllTables(ctx context.Context, input *proto_table.Empty) (*proto_table.TableListResponse, error) {
	tables, err := h.service.Tables.GetAll(ctx)
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}
	tablesResponse := make([]*proto_table.TableObject, len(tables))
	for index, table := range tables {
		tablesResponse[index] = &proto_table.TableObject{
			Id:            table.ID.String(),
			NumberOfSeats: int32(table.NumberOfSeats),
			IsReserved:    table.IsReserved,
			TableNumber:   int32(table.TableNumber),
			Restaurant: &proto_table.RestaurantObject{
				Id:      table.Restaurant.ID.String(),
				Name:    table.Restaurant.Name,
				Address: table.Restaurant.Address,
				Contact: table.Restaurant.Contact,
			},
		}
	}
	return &proto_table.TableListResponse{
		Tables: tablesResponse,
	}, nil
}

func (h *Handler) GetTablesByRestId(ctx context.Context, input *proto_table.IDRequest) (*proto_table.TableListResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	tables, err := h.service.Tables.GetAllByRestaurantId(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrNotFoundInDB):
			return nil, status.Error(codes.NotFound, "table not found")
		default:
			return nil, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}
	tablesResponse := make([]*proto_table.TableObject, len(tables))
	for index, table := range tables {
		tablesResponse[index] = &proto_table.TableObject{
			Id:            table.ID.String(),
			NumberOfSeats: int32(table.NumberOfSeats),
			IsReserved:    table.IsReserved,
			TableNumber:   int32(table.TableNumber),
			Restaurant: &proto_table.RestaurantObject{
				Id:      table.Restaurant.ID.String(),
				Name:    table.Restaurant.Name,
				Address: table.Restaurant.Address,
				Contact: table.Restaurant.Contact,
			},
		}
	}
	return &proto_table.TableListResponse{
		Tables: tablesResponse,
	}, nil
}

func (h *Handler) GetTable(ctx context.Context, input *proto_table.IDRequest) (*proto_table.TableObject, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	table, err := h.service.Tables.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrNotFoundInDB):
			return nil, status.Error(codes.NotFound, "table not found")
		default:
			return nil, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_table.TableObject{
		Id:            table.ID.String(),
		NumberOfSeats: int32(table.NumberOfSeats),
		IsReserved:    table.IsReserved,
		TableNumber:   int32(table.TableNumber),
		Restaurant: &proto_table.RestaurantObject{
			Id:      table.Restaurant.ID.String(),
			Name:    table.Restaurant.Name,
			Address: table.Restaurant.Address,
			Contact: table.Restaurant.Contact,
		},
	}, nil
}

func (h *Handler) AddTable(ctx context.Context, input *proto_table.AddTableRequest) (*proto_table.StatusResponse, error) {

	if input.NumberOfSeats == 0 {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "number of seats is required")
	}
	if input.TableNumber == 0 {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "table number is required")
	}
	if input.RestaurantID == "" {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "restaurant id is required")
	}

	if err := h.service.Tables.Create(ctx, &domain.TableInputSql{
		NumberOfSeats: uint(input.GetNumberOfSeats()),
		TableNumber:   uint(input.GetTableNumber()),
		IsReserved:    input.GetIsReserved(),
		RestaurantID:  input.GetRestaurantID(),
	}); err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_table.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_table.StatusResponse{Status: true}, nil
}

func (h *Handler) UpdateTableById(ctx context.Context, input *proto_table.UpdateTableRequest) (*proto_table.StatusResponse, error) {
	if input.GetId() == "" {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "id is required")
	}
	if input.GetNumberOfSeats() == 0 {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "number of seats is required")
	}
	if input.GetTableNumber() == 0 {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "table number is required")
	}
	err := h.service.Tables.UpdateById(ctx, &domain.UpdateTableInputSql{
		TableID:       input.GetId(),
		NumberOfSeats: uint(input.GetNumberOfSeats()),
		IsReserved:    input.GetIsReserved(),
		TableNumber:   uint(input.GetTableNumber()),
	})
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_table.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	return &proto_table.StatusResponse{Status: true}, nil
}

func (h *Handler) DeleteTableById(ctx context.Context, input *proto_table.IDRequest) (*proto_table.StatusResponse, error) {
	if input.GetId() == "" {
		return &proto_table.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "id is required")
	}

	err := h.service.Tables.Delete(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_table.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto_table.StatusResponse{Status: true}, nil
}

func (h *Handler) GetAvailableTables(ctx context.Context, input *proto_table.IDRequest) (*proto_table.TableListResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	tables, err := h.service.Tables.GetAvailable(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrNotFoundInDB):
			return nil, status.Error(codes.InvalidArgument, domain.ErrNotFoundInDB.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}

	tablesResponse := make([]*proto_table.TableObject, len(tables))
	for index, table := range tables {
		tablesResponse[index] = &proto_table.TableObject{
			Id:            table.ID.String(),
			NumberOfSeats: int32(table.NumberOfSeats),
			IsReserved:    table.IsReserved,
			TableNumber:   int32(table.TableNumber),
			Restaurant: &proto_table.RestaurantObject{
				Id:      table.Restaurant.ID.String(),
				Name:    table.Restaurant.Name,
				Address: table.Restaurant.Address,
				Contact: table.Restaurant.Contact,
			},
		}
	}
	return &proto_table.TableListResponse{
		Tables: tablesResponse,
	}, nil
}

func (h *Handler) GetReservedTables(ctx context.Context, input *proto_table.IDRequest) (*proto_table.TableListResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	tables, err := h.service.Tables.GetReserved(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrNotFoundInDB):
			return nil, status.Error(codes.InvalidArgument, domain.ErrNotFoundInDB.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error "+err.Error())
		}
	}
	tablesResponse := make([]*proto_table.TableObject, len(tables))
	for index, table := range tables {
		tablesResponse[index] = &proto_table.TableObject{
			Id:            table.ID.String(),
			NumberOfSeats: int32(table.NumberOfSeats),
			IsReserved:    table.IsReserved,
			TableNumber:   int32(table.TableNumber),
			Restaurant: &proto_table.RestaurantObject{
				Id:      table.Restaurant.ID.String(),
				Name:    table.Restaurant.Name,
				Address: table.Restaurant.Address,
				Contact: table.Restaurant.Contact,
			},
		}
	}
	return &proto_table.TableListResponse{
		Tables: tablesResponse,
	}, nil
}
