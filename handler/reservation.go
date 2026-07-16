package handler

import (
	"context"
	"dip/domain"
	"dip/internal/logger"
	"errors"
	"time"

	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toReservationObject maps a domain reservation to its transport representation.
func toReservationObject(res *domain.ReservationStruct) *proto_reservation.ReservationObject {
	return &proto_reservation.ReservationObject{
		Id:     res.ID.String(),
		UserID: res.UserID,
		Table: &proto_reservation.TableObject{
			Id:            res.Table.ID.String(),
			NumberOfSeats: int32(res.Table.NumberOfSeats),
			IsReserved:    res.Table.IsReserved,
			TableNumber:   int32(res.Table.TableNumber),
			Restaurant: &proto_reservation.RestaurantObject{
				Id:      res.Table.Restaurant.ID.String(),
				Name:    res.Table.Restaurant.Name,
				Address: res.Table.Restaurant.Address,
				Contact: res.Table.Restaurant.Contact,
			},
		},
		StartAt:   timestamppb.New(res.StartAt),
		PartySize: int32(res.PartySize),
		Confirmed: res.Confirmed,
	}
}

func (h *Handler) MakeReservation(ctx context.Context, input *proto_reservation.ReservationSQLRequest) (*proto_reservation.IDRequest, error) {
	if input.GetUserID() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if input.GetTableID() == "" {
		return nil, status.Error(codes.InvalidArgument, "table id is required")
	}
	if input.GetStartAt() == nil {
		return nil, status.Error(codes.InvalidArgument, "start time is required")
	}
	if input.GetPartySize() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "party size must be positive")
	}
	startAt := input.GetStartAt().AsTime()
	if startAt.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "start time must be in the future")
	}

	// The party must fit the table.
	table, err := h.service.Tables.GetById(ctx, input.GetTableID())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return nil, status.Error(codes.InvalidArgument, "invalid table id")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	if input.GetPartySize() > int32(table.NumberOfSeats) {
		return nil, status.Error(codes.InvalidArgument, "party size exceeds table capacity")
	}

	id, err := h.service.Reservations.Create(ctx, &domain.ReservationInputSql{
		UserID:    input.GetUserID(),
		TableID:   input.GetTableID(),
		StartAt:   startAt,
		PartySize: int(input.GetPartySize()),
		Confirmed: false,
	})
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrTableOccupied) {
			return nil, status.Error(codes.AlreadyExists, domain.ErrTableOccupied.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	// Flip the table's convenience flag after the reservation is persisted, so a
	// rejected overlap never leaves a stale flag. Availability is authoritative
	// via the reservation intervals, so a failure here is not fatal.
	if err := h.service.Tables.MarkOccupied(ctx, input.GetTableID()); err != nil {
		logger.Error(err)
	}

	return &proto_reservation.IDRequest{Id: id}, nil
}

func (h *Handler) GetReservation(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.ReservationObject, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reservation, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return nil, status.Error(codes.NotFound, "reservation not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return toReservationObject(reservation), nil
}

func (h *Handler) DeleteReservationById(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.StatusResponse, error) {
	if input.GetId() == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "id is required")
	}
	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.NotFound, domain.ErrNotFoundInDB.Error())
		}
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
	}

	if err = h.service.Reservations.DeleteById(ctx, input.GetId()); err != nil {
		logger.Error(err)
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
	}

	if err = h.service.Tables.MarkVacant(ctx, reserv.Table.ID.String()); err != nil {
		logger.Error(err)
	}
	return &proto_reservation.StatusResponse{Status: true}, nil
}

func (h *Handler) GetAllReservationByUserId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.ReservationListResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reservations, err := h.service.Reservations.GetAllByUserId(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return reservationListResponse(reservations), nil
}

func (h *Handler) GetAllReservationByRestaurantId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.ReservationListResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reservations, err := h.service.Reservations.GetAllByRestaurantId(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return reservationListResponse(reservations), nil
}

func reservationListResponse(reservations []*domain.ReservationStruct) *proto_reservation.ReservationListResponse {
	objects := make([]*proto_reservation.ReservationObject, len(reservations))
	for i, res := range reservations {
		objects[i] = toReservationObject(res)
	}
	return &proto_reservation.ReservationListResponse{Reservations: objects}
}

func (h *Handler) UpdateReservation(ctx context.Context, input *proto_reservation.UpdateReservationRequest) (*proto_reservation.StatusResponse, error) {
	if input.GetReservationID() == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "reservation id is required")
	}
	if input.GetTableID() == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "table id is required")
	}
	if input.GetStartAt() == nil {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "start time is required")
	}
	if input.GetPartySize() <= 0 {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "party size must be positive")
	}

	err := h.service.Reservations.Update(ctx, &domain.UpdateReservationInputSql{
		ReservationID: input.GetReservationID(),
		TableID:       input.GetTableID(),
		StartAt:       input.GetStartAt().AsTime(),
		PartySize:     int(input.GetPartySize()),
		Confirmed:     false,
	})
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, domain.ErrTableOccupied):
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.AlreadyExists, domain.ErrTableOccupied.Error())
		case errors.Is(err, domain.ErrNotFoundInDB):
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, domain.ErrNotFoundInDB.Error())
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}
	return &proto_reservation.StatusResponse{Status: true}, nil
}

func (h *Handler) GetRestaurantByReservationId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.RestaurantObject, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return nil, status.Error(codes.NotFound, "reservation not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &proto_reservation.RestaurantObject{
		Id:      reserv.Table.Restaurant.ID.String(),
		Name:    reserv.Table.Restaurant.Name,
		Address: reserv.Table.Restaurant.Address,
		Contact: reserv.Table.Restaurant.Contact,
	}, nil
}

func (h *Handler) GetTableByReservationId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.TableObject, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return nil, status.Error(codes.NotFound, "reservation not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return toReservationObject(reserv).GetTable(), nil
}

func (h *Handler) ConfirmReservation(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.StatusResponse, error) {
	if input.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	reservation, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.NotFound, "reservation not found")
		}
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
	}
	if reservation.Confirmed {
		return &proto_reservation.StatusResponse{Status: true}, nil
	}

	err = h.service.Reservations.Update(ctx, &domain.UpdateReservationInputSql{
		ReservationID: reservation.ID.String(),
		TableID:       reservation.Table.ID.String(),
		StartAt:       reservation.StartAt,
		PartySize:     reservation.PartySize,
		Confirmed:     true,
	})
	if err != nil {
		logger.Error(err)
		if errors.Is(err, domain.ErrNotFoundInDB) {
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, domain.ErrNotFoundInDB.Error())
		}
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
	}
	return &proto_reservation.StatusResponse{Status: true}, nil
}
