package models

import (
	"github.com/gofrs/uuid"
)

type RestaurantSql struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Contact string    `json:"contact"`
}

type TableSql struct {
	ID            uuid.UUID `json:"id"`
	NumberOfSeats uint      `json:"numberOfSeats"`
	IsReserved    bool      `json:"isReserved"`
	TableNumber   uint      `json:"tableNumber"`
	RestaurantID  uuid.UUID `json:"restaurantId"`
}

type TableStruct struct {
	ID            uuid.UUID `json:"id"`
	NumberOfSeats uint      `json:"numberOfSeats"`
	IsReserved    bool      `json:"isReserved"`
	TableNumber   uint      `json:"tableNumber"`
	Restaurant    RestaurantSql
}

type ReservationSql struct {
	ID              string `json:"id"`
	UserID          string `json:"userId"`
	TableID         string `json:"tableId"`
	ReservationTime string `json:"reservationTime"`
}

type ReservationStruct struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"userId"`
	Table           TableStruct
	ReservationTime string `json:"reservationTime"`
}
