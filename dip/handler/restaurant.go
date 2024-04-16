package handler

import (
	"dip/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllRestaurants(c *gin.Context) {
	res, err := h.service.Restaurants.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetRestaurant(c *gin.Context) {
	var getRestInp models.GetByIdInputSql
	if err := c.BindJSON(&getRestInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res, err := h.service.Restaurants.GetById(c.Request.Context(), getRestInp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) AddRestaurant(c *gin.Context) {
	var tableInp models.RestaurantInputSql
	if err := c.BindJSON(&tableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	table := models.RestaurantSql{Name: tableInp.Name,
		Address: tableInp.Address, Contact: tableInp.Contact}

	if err := h.service.Restaurants.Create(c.Request.Context(), &table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal servar error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "restaurant made"})
}

func (h *Handler) DeleteRestaurantById(c *gin.Context) {
	var delReserv models.DeleteInputSql
	if err := c.BindJSON(&delReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := h.service.Restaurants.DeleteById(c.Request.Context(), delReserv.DeleteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"Restaurant Deleted": nil})
}

func (h *Handler) UpdateRestById(c *gin.Context) {
	var upRest models.UpdateRestaurantInputSql
	if err := c.BindJSON(&upRest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := h.service.Restaurants.UpdateById(c.Request.Context(), &upRest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated restaurant": nil})
}
