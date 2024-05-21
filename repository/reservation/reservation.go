package reservation

import (
	"context"
	"dip/domain"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepo struct {
	db *pgxpool.Pool
}

func NewReservationRepo(db *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{
		db: db,
	}
}

func (r *ReservationRepo) Create(ctx context.Context, reservation *domain.ReservationSql) error {
	query := `
	INSERT INTO reservations (userid, tableid, reservationtime, reservationdate) 
VALUES ($1, $2, $3, CURRENT_DATE)
RETURNING id;`
	args := []any{reservation.UserID, reservation.TableID, reservation.ReservationTime}

	err := r.db.QueryRow(ctx, query, args...).Scan(&reservation.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if strings.Contains(pgErr.ConstraintName, "reservationTime") {
				return domain.ErrDuplicateKeyErr
			}
		}
		return err
	}
	return nil
}

func (r *ReservationRepo) Delete(ctx context.Context, reservationId uuid.UUID) error {
	query := `DELETE FROM reservations WHERE id = $1`

	_, err := r.db.Exec(ctx, query, reservationId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReservationRepo) GetById(ctx context.Context, resId uuid.UUID) (*domain.ReservationStruct, error) {
	query := `Select reservations.id, reservations.userid, restables.id, restables.numberofseats,
restables.isreserved, restables.tablenumber,  restaurants.*, reservations.reservationtime
from reservations 
join restables on reservations.tableid = restables.id 
join restaurants on restables.restaurantid = restaurants.id 
where reservations.id = $1`
	var reservation domain.ReservationStruct
	err := r.db.QueryRow(ctx, query, resId).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.Table.ID,
		&reservation.Table.NumberOfSeats,
		&reservation.Table.IsReserved,
		&reservation.Table.TableNumber,
		&reservation.Table.Restaurant.ID,
		&reservation.Table.Restaurant.Name,
		&reservation.Table.Restaurant.Address,
		&reservation.Table.Restaurant.Contact,
		&reservation.ReservationTime,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *ReservationRepo) GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error) {
	query := `Select reservations.id, reservations.userid, restables.id, restables.numberofseats,
restables.isreserved, restables.tablenumber,  restaurants.* , reservations.reservationtime
from reservations 
join restables on reservations.tableid = restables.id 
join restaurants on restables.restaurantid = restaurants.id 
where reservations.userId = $1`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	reservations := make([]*domain.ReservationStruct, 0)
	for rows.Next() {
		reservation := new(domain.ReservationStruct)
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.Table.ID,
			&reservation.Table.NumberOfSeats,
			&reservation.Table.IsReserved,
			&reservation.Table.TableNumber,
			&reservation.Table.Restaurant.ID,
			&reservation.Table.Restaurant.Name,
			&reservation.Table.Restaurant.Address,
			&reservation.Table.Restaurant.Contact,
			&reservation.ReservationTime,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *ReservationRepo) Update(ctx context.Context, upReserv *domain.ReservationSql) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE reservations SET tableid = $1, reservationtime = $2, reservationdate = CURRENT_DATE WHERE id = $3"
	_, err = tx.Exec(ctx, query, upReserv.TableID, upReserv.ReservationTime, upReserv.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *ReservationRepo) GetAllByRestaurantId(ctx context.Context, restaurantId string) ([]*domain.ReservationStruct, error) {
	query := `SELECT reservations.id, reservations.userid, restables.id, restables.numberofseats,
       restables.isreserved, restables.tablenumber, restaurants.*, reservations.reservationtime
FROM reservations 
JOIN restables ON reservations.tableid = restables.id 
JOIN restaurants ON restables.restaurantid = restaurants.id 
WHERE restaurants.id = $1 AND reservations.reservationdate = CURRENT_DATE`

	rows, err := r.db.Query(ctx, query, restaurantId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	reservations := make([]*domain.ReservationStruct, 0)
	for rows.Next() {
		reservation := new(domain.ReservationStruct)
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.Table.ID,
			&reservation.Table.NumberOfSeats,
			&reservation.Table.IsReserved,
			&reservation.Table.TableNumber,
			&reservation.Table.Restaurant.ID,
			&reservation.Table.Restaurant.Name,
			&reservation.Table.Restaurant.Address,
			&reservation.Table.Restaurant.Contact,
			&reservation.ReservationTime,
		)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}
