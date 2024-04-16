package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type GetByIdInputSql struct {
	ID pgtype.UUID `json:"id"`
}

type UpdateReservationInputSql struct {
	ReservationID   pgtype.UUID `json:"reservationId"`
	TableID         pgtype.UUID `json:"tableId"`
	ReservationTime string      `json:"reservationTime"`
}

type ReservationInputSql struct {
	UserID          pgtype.UUID `json:"userId"`
	TableID         pgtype.UUID `json:"tableId"`
	ReservationTime string      `json:"reservationTime"`
}

type DeleteInputSql struct {
	DeleteId pgtype.UUID `json:"id"`
}

type GetAllInputSql struct {
	UserId pgtype.UUID `json:"userId"`
}

type StatusTableInputSql struct {
	TableID    pgtype.UUID `json:"tableId"`
	IsReserved bool        `json:"isReserved"`
}

type UpdateTableInputSql struct {
	TableID       pgtype.UUID `json:"tableId"`
	NumberOfSeats uint        `json:"numberOfSeats"`
	TableNumber   uint        `json:"tableNumber"`
	IsReserved    bool        `json:"isReserved"`
}

type TableInputSql struct {
	NumberOfSeats uint        `json:"numberOfSeats"`
	IsReserved    bool        `json:"isReserved"`
	RestaurantID  pgtype.UUID `json:"restaurantId"`
}
type RestaurantInputSql struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type UpdateRestaurantInputSql struct {
	RestaurantId pgtype.UUID `json:"restaurantId"`
	Name         string      `json:"name"`
	Address      string      `json:"address"`
	Contact      string      `json:"contact"`
}
