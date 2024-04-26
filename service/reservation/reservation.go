package reservation

import (
	"context"
	"dip/models"
	repo "dip/repository"
	"github.com/gofrs/uuid"
)

type ReservationService struct {
	repo repo.Reservations
}

func NewreservationService(repo repo.Reservations) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Create(ctx context.Context, reservation *models.ReservationInputSql) error {
	newReservation := models.ReservationSql{
		UserID:          reservation.UserID,
		TableID:         reservation.TableID,
		ReservationTime: reservation.ReservationTime,
	}

	return s.repo.Create(ctx, &newReservation)
}

func (s *ReservationService) DeleteById(ctx context.Context, reservationId string) error {
	newReservId, err := uuid.FromString(reservationId)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, newReservId)
}

func (s *ReservationService) GetById(ctx context.Context, reservationId string) (*models.ReservationStruct, error) {
	newReservId, err := uuid.FromString(reservationId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetById(ctx, newReservId)
}

func (s *ReservationService) GetAllByUserId(ctx context.Context, userId string) ([]*models.ReservationStruct, error) {
	newUserId, err := uuid.FromString(userId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllByUserId(ctx, newUserId)
}

func (s *ReservationService) Update(ctx context.Context, upReserv *models.UpdateReservationInputSql) error {
	newReservation := models.ReservationSql{
		ID:              upReserv.ReservationID,
		TableID:         upReserv.TableID,
		ReservationTime: upReserv.ReservationTime,
	}

	return s.repo.Update(ctx, &newReservation)
}
