package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type RestaurantSql struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Contact string    `json:"contact"`
	Photos  []PhotoSql
}

type PhotoSql struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurantID"`
	URl          string    `json:"url"`
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
	Confirmed       bool   `json:"confirmed"`
}

type ReservationStruct struct {
	ID              uuid.UUID `json:"id"`
	UserID          string    `json:"userId"`
	Table           TableStruct
	ReservationTime string    `json:"reservationTime"`
	ReservationDate time.Time `json:"reservationDate"`
	Confirmed       bool      `json:"confirmed"`
}
