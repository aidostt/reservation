package service

import (
	"context"
	"dip/models"
	"dip/repository"
	"dip/service/reservation"
	"dip/service/restaurant"
	"dip/service/table"
	"github.com/jackc/pgx/v5/pgtype"
)

type Restaurants interface {
	GetById(ctx context.Context, id pgtype.UUID) (*models.RestaurantSql, error)
	GetAll(ctx context.Context) ([]*models.RestaurantSql, error)
	Create(ctx context.Context, res *models.RestaurantSql) error
	DeleteById(ctx context.Context, restId pgtype.UUID) error
	UpdateById(ctx context.Context, upTable *models.UpdateRestaurantInputSql) error
}

type Tables interface {
	GetById(ctx context.Context, id pgtype.UUID) (*models.TableSql, error)
	GetAll(ctx context.Context) ([]*models.TableSql, error)
	Create(ctx context.Context, res *models.TableSql) error
	UpdateById(ctx context.Context, upTable *models.UpdateTableInputSql) error
	MarkOccupied(ctx context.Context, tableId pgtype.UUID) error
	MarkVacant(ctx context.Context, tableId pgtype.UUID) error
	Delete(ctx context.Context, tableId pgtype.UUID) error
	GetAvailable(ctx context.Context, restid pgtype.UUID) ([]*models.TableSql, error)
	GetReserved(ctx context.Context, restid pgtype.UUID) ([]*models.TableSql, error)
	GetAllByRestaurantId(ctx context.Context, restId pgtype.UUID) ([]*models.TableSql, error)
}

type Reservations interface {
	Create(ctx context.Context, reserv *models.ReservationSql) error // models.createReservation
	GetById(ctx context.Context, resId pgtype.UUID) (*models.ReservationSql, error)
	GetAllByUserId(ctx context.Context, userId pgtype.UUID) ([]*models.ReservationSql, error)
	Update(ctx context.Context, upReserv *models.UpdateReservationInputSql) error
	DeleteById(ctx context.Context, resId pgtype.UUID) error
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
