package domain

import "github.com/gofrs/uuid"

type GetByIdInputSql struct {
	ID uuid.UUID `json:"id"`
}

type ReservationInputSql struct {
	UserID          string `json:"userId"`
	TableID         string `json:"tableId"`
	ReservationTime string `json:"reservationTime"`
	Confirmed       bool   `json:"confirmed"`
}

type UpdateReservationInputSql struct {
	ReservationID   string `json:"reservationId"`
	TableID         string `json:"tableId"`
	ReservationTime string `json:"reservationTime"`
	Confirmed       bool   `json:"confirmed"`
}

type DeleteInputSql struct {
	DeleteId string `json:"id"`
}

type StatusTableInputSql struct {
	TableID    uuid.UUID `json:"tableId"`
	IsReserved bool      `json:"isReserved"`
}

type UpdateTableInputSql struct {
	TableID       string `json:"tableId"`
	NumberOfSeats uint   `json:"numberOfSeats"`
	TableNumber   uint   `json:"tableNumber"`
	IsReserved    bool   `json:"isReserved"`
}

type TableInputSql struct {
	NumberOfSeats uint   `json:"numberOfSeats"`
	IsReserved    bool   `json:"isReserved"`
	TableNumber   uint   `json:"tableNumber"`
	RestaurantID  string `json:"restaurantId"`
}
type RestaurantInputSql struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type UpdateRestaurantInputSql struct {
	RestaurantId string `json:"restaurantId"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Contact      string `json:"contact"`
}
