package reservation

import (
	"context"
	"dip/models"
	repo "dip/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type ReservationService struct {
	repo repo.Reservations
}

func NewreservationService(repo repo.Reservations) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Create(ctx context.Context, reservation *models.ReservationSql) error {
	return s.repo.Create(ctx, reservation)
}

func (s *ReservationService) DeleteById(ctx context.Context, reservationId pgtype.UUID) error {
	return s.repo.Delete(ctx, reservationId)
}

func (s *ReservationService) GetById(ctx context.Context, reservationId pgtype.UUID) (*models.ReservationStruct, error) {
	return s.repo.GetById(ctx, reservationId)
}

func (s *ReservationService) GetAllByUserId(ctx context.Context, userId pgtype.UUID) ([]*models.ReservationStruct, error) {
	return s.repo.GetAllByUserId(ctx, userId)
}

func (s *ReservationService) Update(ctx context.Context, upReserv *models.UpdateReservationInputSql) error {
	return s.repo.Update(ctx, upReserv)
}
