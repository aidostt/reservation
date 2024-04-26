package handler

import (
	"dip/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// rewrited
func (h *Handler) GetAllTables(c *gin.Context) {
	tables, err := h.service.Tables.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, tables)
}

// rewrited
func (h *Handler) GetTablesByRestId(c *gin.Context) {
	var getTableInp models.GetByIdInputSql
	if err := c.BindJSON(&getTableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	tables, err := h.service.Tables.GetAllByRestaurantId(c.Request.Context(), getTableInp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, tables)
}

// rewrited
func (h *Handler) GetTable(c *gin.Context) {
	var getTableInp models.GetByIdInputSql
	if err := c.BindJSON(&getTableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	table, err := h.service.Tables.GetById(c.Request.Context(), getTableInp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *Handler) AddTable(c *gin.Context) {
	var tableInp models.TableInputSql
	if err := c.BindJSON(&tableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// table := models.TableSql{NumberOfSeats: tableInp.NumberOfSeats,
	// 	IsReserved: false, RestaurantID: tableInp.RestaurantID}

	if err := h.service.Tables.Create(c.Request.Context(), &tableInp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"internal servar error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "table made"})
}

func (h *Handler) UpdateTableById(c *gin.Context) {
	var upTable models.UpdateTableInputSql
	if err := c.BindJSON(&upTable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := h.service.Tables.UpdateById(c.Request.Context(), &upTable)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated table": nil})
}

func (h *Handler) DeleteTableById(c *gin.Context) {
	var delReserv models.DeleteInputSql
	if err := c.BindJSON(&delReserv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err := h.service.Tables.Delete(c.Request.Context(), delReserv.DeleteId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"Table Deleted": nil})
}

// rewrited
func (h *Handler) GetAvailableTables(c *gin.Context) {
	var getTableInp models.GetByIdInputSql
	if err := c.BindJSON(&getTableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	tables, err := h.service.Tables.GetAvailable(c.Request.Context(), getTableInp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, tables)
}

// rewrited
func (h *Handler) GetReservedTables(c *gin.Context) {
	var getTableInp models.GetByIdInputSql
	if err := c.BindJSON(&getTableInp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	tables, err := h.service.Tables.GetReserved(c.Request.Context(), getTableInp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bad reques": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, tables)
}
