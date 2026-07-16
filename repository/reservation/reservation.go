package reservation

import (
	"context"
	"dip/domain"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgExclusionViolation is the SQLSTATE raised when an insert or update conflicts
// with an EXCLUDE constraint (here: an overlapping reservation for a table).
const pgExclusionViolation = "23P01"

type ReservationRepo struct {
	db *pgxpool.Pool
}

func NewReservationRepo(db *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{db: db}
}

// reservationSelect lists columns explicitly (rather than restaurants.*) so the
// scan order is stable when a joined table gains a column.
const reservationSelect = `
SELECT r.id, r.userid,
       t.id, t.numberofseats, t.isreserved, t.tablenumber,
       rest.id, rest.name, rest.address, rest.contact,
       r.start_at, r.party_size, r.confirmed
FROM reservations r
JOIN restables t ON r.tableid = t.id
JOIN restaurants rest ON t.restaurantid = rest.id`

func scanReservation(row pgx.Row, res *domain.ReservationStruct) error {
	return row.Scan(
		&res.ID, &res.UserID,
		&res.Table.ID, &res.Table.NumberOfSeats, &res.Table.IsReserved, &res.Table.TableNumber,
		&res.Table.Restaurant.ID, &res.Table.Restaurant.Name, &res.Table.Restaurant.Address, &res.Table.Restaurant.Contact,
		&res.StartAt, &res.PartySize, &res.Confirmed,
	)
}

func (r *ReservationRepo) Create(ctx context.Context, reservation *domain.ReservationSql) error {
	const query = `
		INSERT INTO reservations (userid, tableid, start_at, ends_at, party_size, confirmed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	err := r.db.QueryRow(ctx, query,
		reservation.UserID, reservation.TableID, reservation.StartAt, reservation.EndsAt,
		reservation.PartySize, reservation.Confirmed,
	).Scan(&reservation.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgExclusionViolation {
			return domain.ErrTableOccupied
		}
		return err
	}
	return nil
}

func (r *ReservationRepo) Delete(ctx context.Context, reservationId uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM reservations WHERE id = $1`, reservationId)
	return err
}

func (r *ReservationRepo) GetById(ctx context.Context, resId uuid.UUID) (*domain.ReservationStruct, error) {
	var res domain.ReservationStruct
	err := scanReservation(r.db.QueryRow(ctx, reservationSelect+` WHERE r.id = $1`, resId), &res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	return &res, nil
}

func (r *ReservationRepo) list(ctx context.Context, whereClause string, arg any) ([]*domain.ReservationStruct, error) {
	rows, err := r.db.Query(ctx, reservationSelect+whereClause, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := make([]*domain.ReservationStruct, 0)
	for rows.Next() {
		res := new(domain.ReservationStruct)
		if err := scanReservation(rows, res); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	return reservations, rows.Err()
}

func (r *ReservationRepo) GetAllByUserId(ctx context.Context, userId string) ([]*domain.ReservationStruct, error) {
	return r.list(ctx, ` WHERE r.userid = $1 ORDER BY r.start_at`, userId)
}

func (r *ReservationRepo) GetAllByRestaurantId(ctx context.Context, restaurantId string) ([]*domain.ReservationStruct, error) {
	return r.list(ctx, ` WHERE rest.id = $1 AND r.start_at::date = CURRENT_DATE ORDER BY r.start_at`, restaurantId)
}

func (r *ReservationRepo) Update(ctx context.Context, upReserv *domain.ReservationSql) error {
	const query = `
		UPDATE reservations
		SET tableid = $1, start_at = $2, ends_at = $3, party_size = $4, confirmed = $5
		WHERE id = $6`
	ct, err := r.db.Exec(ctx, query,
		upReserv.TableID, upReserv.StartAt, upReserv.EndsAt, upReserv.PartySize, upReserv.Confirmed, upReserv.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgExclusionViolation {
			return domain.ErrTableOccupied
		}
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFoundInDB
	}
	return nil
}
