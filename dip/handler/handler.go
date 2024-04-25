package handler

import (
	"dip/service"
	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	proto_restaurant "github.com/aidostt/protos/gen/go/reservista/restaurant"
	proto_table "github.com/aidostt/protos/gen/go/reservista/table"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Init() *gin.Engine {

	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handler := NewHandler(h.service)
	api := router.Group("/api")
	{
		handler.InitReservRoutes(api)
		handler.InitRestaurantRoutes(api)
		handler.InitTableRoutes(api)
	}
}
