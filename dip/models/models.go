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

type TableStruct struct {
	ID            pgtype.UUID `json:"id"`
	NumberOfSeats uint        `json:"numberOfSeats"`
	IsReserved    bool        `json:"isReserved"`
	TableNumber   uint        `json:"tableNumber"`
	Restaurant    RestaurantSql
}

type ReservationSql struct {
	ID              pgtype.UUID `json:"id"`
	UserID          pgtype.UUID `json:"userId"`
	TableID         pgtype.UUID `json:"tableId"`
	RestaurantID    pgtype.UUID `json:"restaurantId"`
	ReservationTime string      `json:"reservationTime"`
}

type ReservationStruct struct {
	ID              pgtype.UUID `json:"id"`
	UserID          pgtype.UUID `json:"userId"`
	Table           TableStruct
	ReservationTime string `json:"reservationTime"`
}
