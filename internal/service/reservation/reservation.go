package reservation

import (
	"context"
	"dip/internal/domain"
	repo "dip/internal/repository"
	"time"

	"github.com/gofrs/uuid"
)

// maxActiveReservationsPerUser caps how many not-yet-ended reservations a single
// user may hold at once, to blunt hoarding and abuse.
const maxActiveReservationsPerUser = 5

type ReservationService struct {
	repo repo.Reservations
	// turnDuration is the restaurant's seating length; the reservation occupies
	// its table for [start, start+turnDuration). It is policy, not client input.
	turnDuration time.Duration
}

func NewReservationService(repo repo.Reservations, turnDuration time.Duration) *ReservationService {
	return &ReservationService{repo: repo, turnDuration: turnDuration}
}

func (s *ReservationService) Create(ctx context.Context, reservation *domain.ReservationInputSql) (string, error) {
	active, err := s.repo.CountActiveByUser(ctx, reservation.UserID)
	if err != nil {
		return "", err
	}
	if active >= maxActiveReservationsPerUser {
		return "", domain.ErrTooManyActiveReservations
	}

	newReservation := domain.ReservationSql{
		UserID:    reservation.UserID,
		TableID:   reservation.TableID,
		StartAt:   reservation.StartAt,
		EndsAt:    reservation.StartAt.Add(s.turnDuration),
		PartySize: reservation.PartySize,
		Confirmed: reservation.Confirmed,
	}
	if err := s.repo.Create(ctx, &newReservation); err != nil {
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
	return s.repo.Update(ctx, &domain.ReservationSql{
		ID:        upReserv.ReservationID,
		TableID:   upReserv.TableID,
		StartAt:   upReserv.StartAt,
		EndsAt:    upReserv.StartAt.Add(s.turnDuration),
		PartySize: upReserv.PartySize,
		Confirmed: upReserv.Confirmed,
	})
}
