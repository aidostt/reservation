package handler

import (
	"dip/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MakeReservation(c *gin.Context) {
	var reservInp models.ReservationInputSql
	if err := c.BindJSON(&reservInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	reservation := models.ReservationSql{UserID: reservInp.UserID, TableID: reservInp.TableID,
		ReservationTime: reservInp.ReservationTime}

	err := h.service.Tables.MarkOccupied(c, reservInp.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request error": err.Error()})
		c.Abort()
		return
	}

	if err = h.service.Reservations.Create(c.Request.Context(), &reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal servar error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "reservation made"})

}

func (h *Handler) GetReservation(c *gin.Context) {
	var getReservInp models.GetByIdInputSql
	if err := c.BindJSON(&getReservInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res, err := h.service.Reservations.GetById(c.Request.Context(), getReservInp.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"not found": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteReservationById(c *gin.Context) {
	var delReserv models.DeleteInputSql
	if err := c.BindJSON(&delReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	reserv, err := h.service.Reservations.GetById(c.Request.Context(), delReserv.DeleteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	err = h.service.Reservations.DeleteById(c.Request.Context(), delReserv.DeleteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	err = h.service.Tables.MarkVacant(c, reserv.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad request error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"Reservation Deleted": nil})
}

func (h *Handler) GetAllReservByUserId(c *gin.Context) {
	var allReserv models.GetAllInputSql
	if err := c.BindJSON(&allReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res, err := h.service.Reservations.GetAllByUserId(c.Request.Context(), allReserv.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)
}

// rewrite
func (h *Handler) UpdateReservation(c *gin.Context) {
	var upReservInp models.UpdateReservationInputSql
	if err := c.BindJSON(&upReservInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := h.service.Reservations.Update(c.Request.Context(), &upReservInp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"not found": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetRestaurantByReservationId(c *gin.Context) {
	var idReserv models.GetByIdInputSql
	if err := c.BindJSON(&idReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	reserv, err := h.service.Reservations.GetById(c.Request.Context(), idReserv.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	table, err := h.service.Tables.GetById(c.Request.Context(), reserv.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	res, err := h.service.Restaurants.GetById(c.Request.Context(), table.RestaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetTableByReservationId(c *gin.Context) {
	var idReserv models.GetByIdInputSql
	if err := c.BindJSON(&idReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	reserv, err := h.service.Reservations.GetById(c.Request.Context(), idReserv.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	table, err := h.service.Tables.GetById(c.Request.Context(), reserv.TableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, table)
}
