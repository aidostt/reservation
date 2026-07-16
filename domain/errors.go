package domain

import "errors"

var (
	ErrTableOccupied             = errors.New("table occupied")
	ErrNotFoundInDB              = errors.New("missing field in db")
	ErrTooManyActiveReservations = errors.New("too many active reservations")
)
