package repository

import (
	"context"
	"dip/domain"
	"dip/repository/photos"
	"dip/repository/reservation"
	"dip/repository/restaurant"
	"dip/repository/table"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Restaurants interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.RestaurantSql, error)
	GetAll(ctx context.Context) ([]*domain.RestaurantSql, error)
	Create(ctx context.Context, res *domain.RestaurantSql) error
	Delete(ctx context.Context, restId uuid.UUID) error
	UpdateById(ctx context.Context, upTable *domain.RestaurantSql) error
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.RestaurantSql, int, error)
	GetSuggestions(ctx context.Context, query string) ([]*domain.RestaurantSql, error)
}

type Tables interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.TableStruct, error)
	GetAll(ctx context.Context) ([]*domain.TableStruct, error)
	Create(ctx context.Context, res *domain.TableSql) error
	UpdateById(ctx context.Context, upTable *domain.TableSql) error
	SetStatusById(ctx context.Context, upTable *domain.StatusTableInputSql) error
	Delete(ctx context.Context, tableId uuid.UUID) error
	GetReserved(ctx context.Context, restid uuid.UUID) ([]*domain.TableStruct, error)
	GetAvailable(ctx context.Context, restid uuid.UUID) ([]*domain.TableStruct, error)
	GetAllByRestaurantId(ctx context.Context, restId uuid.UUID) ([]*domain.TableStruct, error)
}

type Reservations interface {
	Create(ctx context.Context, reserv *domain.ReservationSql) error // domain.createReservation
	GetById(ctx context.Context, resId uuid.UUID) (*domain.ReservationStruct, error)
	GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error)
	GetAllByRestaurantId(ctx context.Context, reservationId string) ([]*domain.ReservationStruct, error)
	Update(ctx context.Context, upReserv *domain.ReservationSql) error
	Delete(ctx context.Context, resId uuid.UUID) error
	TableOccupied(context.Context, uuid.UUID, string) error
}

type Photos interface {
	Upload(ctx context.Context, photos []*domain.PhotoSql) error
	Delete(ctx context.Context, url string, restaurantID uuid.UUID) error
}

type Repository struct {
	Restaurants
	Tables
	Reservations
	Photos
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Restaurants:  restaurant.NewRestaurantRepo(db),
		Tables:       table.NewRestaurantRepo(db),
		Reservations: reservation.NewReservationRepo(db),
		Photos:       photos.NewPhotoRepo(db),
	}
}
