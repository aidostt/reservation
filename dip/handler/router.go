package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitReservRoutes(rg *gin.RouterGroup) {
	reservations := rg.Group("reservations")
	{
		reservations.POST("/makereservation", h.MakeReservation)
		reservations.GET("/viewreservation", h.GetReservation)
		reservations.PATCH("/updatereservation", h.UpdateReservation)
		reservations.DELETE("/cancelreservation", h.DeleteReservationById)
		reservations.GET("/allreservations", h.GetAllReservByUserId)
		reservations.GET("/viewrestaurantinfo", h.GetRestaurantByReservationId)
		reservations.GET("/viewtableinfo", h.GetTableByReservationId)
	}
}

func (h *Handler) InitTableRoutes(rg *gin.RouterGroup) {
	reservations := rg.Group("tables")
	{
		reservations.GET("/viewtable", h.GetTable)
		reservations.GET("/alltablesofrestaurant", h.GetTablesByRestId) // restid
		reservations.POST("/addtable", h.AddTable)
		reservations.DELETE("/deletetable", h.DeleteTableById)
		reservations.GET("/allavailabletables", h.GetAvailableTables) // restid
		reservations.GET("/allreservedtables", h.GetReservedTables)   // restid
		reservations.PATCH("/updatetable", h.UpdateTableById)
	}
}

func (h *Handler) InitRestaurantRoutes(rg *gin.RouterGroup) {
	reservations := rg.Group("restaurants")
	{
		reservations.GET("/viewrestaurant", h.GetRestaurant)
		reservations.GET("/allRestaurants", h.GetAllRestaurants)
		reservations.POST("/addrestaurant", h.AddRestaurant)
		reservations.DELETE("/deleterestaurant", h.DeleteRestaurantById)
		reservations.PATCH("/updaterestaurant", h.UpdateRestById)
	}
}
