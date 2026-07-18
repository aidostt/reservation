package server

import (
	"dip/internal/delivery"
	reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	restaurant "github.com/aidostt/protos/gen/go/reservista/restaurant"
	table "github.com/aidostt/protos/gen/go/reservista/table"
)

func (s *Server) RegisterServers(h *delivery.Handler) {
	reservation.RegisterReservationServer(s.GrpcServer, h)
	restaurant.RegisterRestaurantServer(s.GrpcServer, h)
	table.RegisterTableServer(s.GrpcServer, h)
}
