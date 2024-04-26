package repository

import (
	"context"
	"dip/models"
	"dip/repository/reservation"
	"dip/repository/restaurant"
	"dip/repository/table"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Restaurants interface {
	GetById(ctx context.Context, id uuid.UUID) (*models.RestaurantSql, error)
	GetAll(ctx context.Context) ([]*models.RestaurantSql, error)
	Create(ctx context.Context, res *models.RestaurantSql) error
	Delete(ctx context.Context, restId uuid.UUID) error
	UpdateById(ctx context.Context, upTable *models.RestaurantSql) error
}

type Tables interface {
	GetById(ctx context.Context, id uuid.UUID) (*models.TableStruct, error)
	GetAll(ctx context.Context) ([]*models.TableStruct, error)
	Create(ctx context.Context, res *models.TableSql) error
	UpdateById(ctx context.Context, upTable *models.TableSql) error
	SetStatusById(ctx context.Context, upTable *models.StatusTableInputSql) error
	Delete(ctx context.Context, tableId uuid.UUID) error
	GetReserved(ctx context.Context, restid uuid.UUID) ([]*models.TableStruct, error)
	GetAvailable(ctx context.Context, restid uuid.UUID) ([]*models.TableStruct, error)
	GetAllByRestaurantId(ctx context.Context, restId uuid.UUID) ([]*models.TableStruct, error)
}

type Reservations interface {
	Create(ctx context.Context, reserv *models.ReservationSql) error // models.createReservation
	GetById(ctx context.Context, resId uuid.UUID) (*models.ReservationStruct, error)
	GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]*models.ReservationStruct, error)
	Update(ctx context.Context, upReserv *models.ReservationSql) error
	Delete(ctx context.Context, resId uuid.UUID) error
}

type Repository struct {
	Restaurants
	Tables
	Reservations
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Restaurants:  restaurant.NewRestaurantRepo(db),
		Tables:       table.NewRestaurantRepo(db),
		Reservations: reservation.NewReservationRepo(db),
	}
}
