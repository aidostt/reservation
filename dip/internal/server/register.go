package server

import (
	"dip/handler"
)

func (s *Server) RegisterServers(h *handler.Handler) {
	reservation.RegisterReservationServer(s.GrpcServer, h)
	restaurant.RegisterRestaurantServer(s.GrpcServer, h)
	table.RegisterTableServer(s.GrpcServer, h)
}
