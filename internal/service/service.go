package service

import (
	"context"
	"dip/internal/domain"
	"dip/internal/repository"
	"dip/internal/service/photos"
	"dip/internal/service/reservation"
	"dip/internal/service/restaurant"
	"dip/internal/service/table"
	"time"
)

type Restaurants interface {
	GetById(ctx context.Context, id string) (*domain.RestaurantSql, error)
	GetAll(ctx context.Context) ([]*domain.RestaurantSql, error)
	Create(ctx context.Context, res *domain.RestaurantSql) error
	DeleteById(ctx context.Context, restId string) error
	UpdateById(ctx context.Context, upTable *domain.UpdateRestaurantInputSql) error
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.RestaurantSql, int, error)
	GetSuggestions(ctx context.Context, query string) ([]*domain.RestaurantSql, error)
}

type Tables interface {
	GetById(ctx context.Context, id string) (*domain.TableStruct, error)
	GetAll(ctx context.Context) ([]*domain.TableStruct, error)
	Create(ctx context.Context, res *domain.TableInputSql) error
	UpdateById(ctx context.Context, upTable *domain.UpdateTableInputSql) error
	Delete(ctx context.Context, tableId string) error
	GetAvailable(ctx context.Context, restid string, startAt time.Time) ([]*domain.TableStruct, error)
	GetReserved(ctx context.Context, restid string, startAt time.Time) ([]*domain.TableStruct, error)
	GetAllByRestaurantId(ctx context.Context, restId string) ([]*domain.TableStruct, error)
}

type Reservations interface {
	Create(ctx context.Context, reserv *domain.ReservationInputSql) (string, error)
	GetById(ctx context.Context, resId string) (*domain.ReservationStruct, error)
	GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error)
	GetAllByRestaurantId(ctx context.Context, restaurantId string) ([]*domain.ReservationStruct, error)
	Update(ctx context.Context, upReserv *domain.UpdateReservationInputSql) error
	DeleteById(ctx context.Context, resId string) error
}

type Photos interface {
	Upload(ctx context.Context, photos []*domain.PhotoSql) error
	Delete(ctx context.Context, url, restaurantID string) error
}

type Service struct {
	Restaurants
	Reservations
	Tables
	Photos
}

type Dependencies struct {
	Repos        *repository.Repository
	Environment  string
	Domain       string
	TurnDuration time.Duration
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Restaurants:  restaurant.NewRestaurantService(deps.Repos.Restaurants),
		Tables:       table.NewTableService(deps.Repos.Tables, deps.TurnDuration),
		Reservations: reservation.NewReservationService(deps.Repos.Reservations, deps.TurnDuration),
		Photos:       photos.NewPhotosService(deps.Repos.Photos),
	}
}
