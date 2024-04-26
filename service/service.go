package service

import (
	"context"
	"dip/models"
	"dip/repository"
	"dip/service/reservation"
	"dip/service/restaurant"
	"dip/service/table"
)

type Restaurants interface {
	GetById(ctx context.Context, id string) (*models.RestaurantSql, error)
	GetAll(ctx context.Context) ([]*models.RestaurantSql, error)
	Create(ctx context.Context, res *models.RestaurantSql) error
	DeleteById(ctx context.Context, restId string) error
	UpdateById(ctx context.Context, upTable *models.UpdateRestaurantInputSql) error
}

type Tables interface {
	GetById(ctx context.Context, id string) (*models.TableStruct, error)
	GetAll(ctx context.Context) ([]*models.TableStruct, error)
	Create(ctx context.Context, res *models.TableInputSql) error
	UpdateById(ctx context.Context, upTable *models.UpdateTableInputSql) error
	MarkOccupied(ctx context.Context, tableId string) error
	MarkVacant(ctx context.Context, tableId string) error
	Delete(ctx context.Context, tableId string) error
	GetAvailable(ctx context.Context, restid string) ([]*models.TableStruct, error)
	GetReserved(ctx context.Context, restid string) ([]*models.TableStruct, error)
	GetAllByRestaurantId(ctx context.Context, restId string) ([]*models.TableStruct, error)
}

type Reservations interface {
	Create(ctx context.Context, reserv *models.ReservationInputSql) error // models.createReservation
	GetById(ctx context.Context, resId string) (*models.ReservationStruct, error)
	GetAllByUserId(ctx context.Context, userId string) ([]*models.ReservationStruct, error)
	Update(ctx context.Context, upReserv *models.UpdateReservationInputSql) error
	DeleteById(ctx context.Context, resId string) error
}

type Service struct {
	Restaurants
	Reservations
	Tables
}

type Dependencies struct {
	Repos       *repository.Repository
	Environment string
	Domain      string
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Restaurants:  restaurant.NewRestaurantService(deps.Repos.Restaurants),
		Tables:       table.NewTableService(deps.Repos.Tables),
		Reservations: reservation.NewreservationService(deps.Repos.Reservations),
	}
}
