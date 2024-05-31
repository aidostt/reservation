package reservation

import (
	"context"
	"dip/domain"
	repo "dip/repository"
	"errors"
	"github.com/gofrs/uuid"
)

type ReservationService struct {
	repo repo.Reservations
}

func NewReservationService(repo repo.Reservations) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Create(ctx context.Context, reservation *domain.ReservationInputSql) (string, error) {
	newReservation := domain.ReservationSql{
		UserID:          reservation.UserID,
		TableID:         reservation.TableID,
		ReservationTime: reservation.ReservationTime,
		Confirmed:       reservation.Confirmed,
	}
	err := s.repo.Create(ctx, &newReservation)
	if err != nil {
		return "", err
	}
	return newReservation.ID, nil
}

func (s *ReservationService) DeleteById(ctx context.Context, reservationId string) error {
	newReservId, err := uuid.FromString(reservationId)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, newReservId)
}

func (s *ReservationService) GetById(ctx context.Context, reservationId string) (*domain.ReservationStruct, error) {
	newReservId, err := uuid.FromString(reservationId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetById(ctx, newReservId)
}

func (s *ReservationService) GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error) {
	return s.repo.GetAllByUserId(ctx, userId)
}

func (s *ReservationService) GetAllByRestaurantId(ctx context.Context, restaurantId string) ([]*domain.ReservationStruct, error) {
	return s.repo.GetAllByRestaurantId(ctx, restaurantId)
}

func (s *ReservationService) Update(ctx context.Context, upReserv *domain.UpdateReservationInputSql) error {
	newReservation := domain.ReservationSql{
		ID:              upReserv.ReservationID,
		TableID:         upReserv.TableID,
		ReservationTime: upReserv.ReservationTime,
		Confirmed:       upReserv.Confirmed,
	}

	return s.repo.Update(ctx, &newReservation)
}

func (s *ReservationService) TableOccupied(ctx context.Context, tableID, reservationTime string) (bool, error) {
	tableUUID, err := uuid.FromString(tableID)
	if err != nil {
		return true, err
	}
	err = s.repo.TableOccupied(ctx, tableUUID, reservationTime)
	if err != nil {
		if errors.Is(err, domain.ErrTableOccupied) {
			return true, nil
		}
		return true, err
	}
	return false, nil
}
