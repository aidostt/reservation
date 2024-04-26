package handler

import (
	"dip/service"
	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	proto_restaurant "github.com/aidostt/protos/gen/go/reservista/restaurant"
	proto_table "github.com/aidostt/protos/gen/go/reservista/table"
)

type Handler struct {
	proto_reservation.UnimplementedReservationServer
	proto_restaurant.UnimplementedRestaurantServer
	proto_table.UnimplementedTableServer
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
