package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type RestaurantSql struct {
	ID      pgtype.UUID `json:"id"`
	Name    string      `json:"name"`
	Address string      `json:"address"`
	Contact string      `json:"contact"`
}

type TableSql struct {
	ID            pgtype.UUID `json:"id"`
	NumberOfSeats uint        `json:"numberOfSeats"`
	IsReserved    bool        `json:"isReserved"`
	TableNumber   uint        `json:"tableNumber"`
	RestaurantID  pgtype.UUID `json:"restaurantId"`
}

type ReservationSql struct {
	ID              pgtype.UUID `json:"id"`
	UserID          pgtype.UUID `json:"userId"`
	TableID         pgtype.UUID `json:"tableId"`
	ReservationTime string      `json:"reservationTime"`
}
