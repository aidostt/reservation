package handler

import (
	"context"
	"dip/internal/logger"
	"dip/models"
	"errors"

	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) MakeReservation(ctx context.Context, input *proto_reservation.ReservationSQLRequest) (*proto_reservation.StatusResponse, error) {
	if input.UserID == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "user id is required")
	}
	if input.TableID == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "table id is required")
	}
	if input.ReservationTime == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "reservation time is required")
	}

	err := h.service.Tables.MarkOccupied(ctx, input.GetTableID())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}

	}
	if err = h.service.Reservations.Create(ctx, &models.ReservationInputSql{
		UserID:          input.GetUserID(),
		TableID:         input.GetTableID(),
		ReservationTime: input.GetReservationTime(),
	}); err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}
	return &proto_reservation.StatusResponse{Status: true}, nil
}

func (h *Handler) GetReservation(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.ReservationObject, error) {
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	reservation, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, errors.New("not found in db")):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}

	}

	return &proto_reservation.ReservationObject{
		Id:     reservation.ID.String(),
		UserID: reservation.UserID.String(),
		Table: &proto_reservation.TableObject{
			Id:            reservation.Table.ID.String(),
			NumberOfSeats: int32(reservation.Table.NumberOfSeats),
			IsReserved:    reservation.Table.IsReserved,
			TableNumber:   int32(reservation.Table.TableNumber),
			Restaurant: &proto_reservation.RestaurantObject{
				Id:      reservation.Table.Restaurant.ID.String(),
				Name:    reservation.Table.Restaurant.Name,
				Address: reservation.Table.Restaurant.Address,
				Contact: reservation.Table.Restaurant.Contact,
			},
		},
		ReservationTime: reservation.ReservationTime,
	}, nil
}

func (h *Handler) DeleteReservationById(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.StatusResponse, error) {
	if input.Id == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "id is required")
	}
	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}

	err = h.service.Reservations.DeleteById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}

	err = h.service.Tables.MarkVacant(ctx, reserv.Table.ID.String())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}
	}
	return &proto_reservation.StatusResponse{Status: true}, nil
}

func (h *Handler) GetAllReservationByUserId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.ReservationListResponse, error) {
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	reservations, err := h.service.Reservations.GetAllByUserId(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	reservResp := make([]*proto_reservation.ReservationObject, len(reservations))
	for i, res := range reservations {
		reservResp[i] = &proto_reservation.ReservationObject{
			Id:     res.ID.String(),
			UserID: res.UserID.String(),
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
			ReservationTime: res.ReservationTime,
		}
	}

	return &proto_reservation.ReservationListResponse{
		Reservations: reservResp,
	}, nil
}

func (h *Handler) UpdateReservation(ctx context.Context, input *proto_reservation.ReservationSQLRequest) (*proto_reservation.StatusResponse, error) {
	if input.UserID == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "user id is required")
	}
	if input.TableID == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "table id is required")
	}
	if input.ReservationTime == "" {
		return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.InvalidArgument, "reservation time is required")
	}

	err := h.service.Reservations.Update(ctx, &models.UpdateReservationInputSql{
		ReservationID:   input.GetUserID(),
		TableID:         input.GetTableID(),
		ReservationTime: input.GetReservationTime(),
	})
	if err != nil {
		logger.Error(err)
		switch {
		case errors.Is(err, errors.New("not found in db")):
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.NotFound, "user not found")
		default:
			return &proto_reservation.StatusResponse{Status: false}, status.Error(codes.Internal, "internal error")
		}

	}

	return &proto_reservation.StatusResponse{Status: true}, nil
}

func (h *Handler) GetRestaurantByReservationId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.RestaurantObject, error) {
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto_reservation.RestaurantObject{
		Id:      reserv.Table.Restaurant.ID.String(),
		Name:    reserv.Table.Restaurant.Name,
		Address: reserv.Table.Restaurant.Address,
		Contact: reserv.Table.Restaurant.Contact,
	}, nil
}

func (h *Handler) GetTableByReservationId(ctx context.Context, input *proto_reservation.IDRequest) (*proto_reservation.TableObject, error) {
	if input.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	reserv, err := h.service.Reservations.GetById(ctx, input.GetId())
	if err != nil {
		logger.Error(err)
		switch {
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &proto_reservation.TableObject{
		Id:            reserv.Table.ID.String(),
		NumberOfSeats: int32(reserv.Table.NumberOfSeats),
		IsReserved:    reserv.Table.IsReserved,
		TableNumber:   int32(reserv.Table.TableNumber),
		Restaurant: &proto_reservation.RestaurantObject{
			Id:      reserv.Table.Restaurant.ID.String(),
			Name:    reserv.Table.Restaurant.Name,
			Address: reserv.Table.Restaurant.Address,
			Contact: reserv.Table.Restaurant.Contact,
		},
	}, nil
}
