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
	INSERT INTO reservations (userid, tableid, reservationtime, reservationdate, confirmed) 
VALUES ($1, $2, $3, CURRENT_DATE, $4)
RETURNING id;`
	args := []any{reservation.UserID, reservation.TableID, reservation.ReservationTime, reservation.Confirmed}

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
func (r *ReservationRepo) GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error) {
	query := `SELECT reservations.id, reservations.userid, restables.id, restables.numberofseats,
              restables.isreserved, restables.tablenumber, restaurants.*, reservations.reservationtime,
              reservations.reservationdate, reservations.confirmed
              FROM reservations 
              JOIN restables ON reservations.tableid = restables.id 
              JOIN restaurants ON restables.restaurantid = restaurants.id 
              WHERE reservations.userid = $1`

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
			&reservation.ReservationDate,
			&reservation.Confirmed,
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

func (r *ReservationRepo) GetById(ctx context.Context, resId uuid.UUID) (*domain.ReservationStruct, error) {
	query := `SELECT reservations.id, reservations.userid, restables.id, restables.numberofseats,
              restables.isreserved, restables.tablenumber, restaurants.*, reservations.reservationtime,
              reservations.reservationdate, reservations.confirmed
              FROM reservations 
              JOIN restables ON reservations.tableid = restables.id 
              JOIN restaurants ON restables.restaurantid = restaurants.id 
              WHERE reservations.id = $1`
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
		&reservation.ReservationDate,
		&reservation.Confirmed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *ReservationRepo) GetAllByRestaurantId(ctx context.Context, restaurantId string) ([]*domain.ReservationStruct, error) {
	query := `SELECT reservations.id, reservations.userid, restables.id, restables.numberofseats,
              restables.isreserved, restables.tablenumber, restaurants.*, reservations.reservationtime,
              reservations.reservationdate, reservations.confirmed
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
			&reservation.ReservationDate,
			&reservation.Confirmed,
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

	query := "UPDATE reservations SET tableid = $1, reservationtime = $2, reservationdate = CURRENT_DATE, confirmed = $3 WHERE id = $4"
	_, err = tx.Exec(ctx, query, upReserv.TableID, upReserv.ReservationTime, upReserv.Confirmed, upReserv.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *ReservationRepo) TableOccupied(ctx context.Context, tableID uuid.UUID, reservationTime string) error {
	query := `SELECT EXISTS (SELECT 1 FROM reservations WHERE tableid = $1 AND reservationtime = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, tableID, reservationTime).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrTableOccupied
	}
	return nil
}
