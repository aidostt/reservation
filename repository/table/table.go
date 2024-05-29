package table

import (
	"context"
	"dip/domain"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TableRepo struct {
	db *pgxpool.Pool
}

func NewRestaurantRepo(db *pgxpool.Pool) *TableRepo {
	return &TableRepo{
		db: db,
	}
}

func (r *TableRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.TableStruct, error) {
	query := `Select restables.id, restables.numberofseats, restables.isreserved,restables.tablenumber, restaurants.* 
from restables 
join restaurants on restables.restaurantId = restaurants.id 
where restables.id = $1`
	var table domain.TableStruct
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID,
		&table.NumberOfSeats,
		&table.IsReserved,
		&table.TableNumber,
		&table.Restaurant.ID,
		&table.Restaurant.Name,
		&table.Restaurant.Address,
		&table.Restaurant.Contact,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	return &table, nil
}

func (r *TableRepo) GetAll(ctx context.Context) ([]*domain.TableStruct, error) {
	query := `Select restables.id, restables.numberofseats, restables.isreserved,restables.tablenumber, restaurants.* 
from restables 
join restaurants on restables.restaurantId = restaurants.id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	tables := make([]*domain.TableStruct, 0)
	for rows.Next() {
		table := new(domain.TableStruct)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.Restaurant.ID,
			&table.Restaurant.Name,
			&table.Restaurant.Address,
			&table.Restaurant.Contact,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *TableRepo) Create(ctx context.Context, table *domain.TableSql) error {
	query := `
	INSERT INTO restables (numberofseats, isreserved, tablenumber, restaurantid) VALUES ($1, $2, $3, $4)
	RETURNING id;
	`
	args := []any{table.NumberOfSeats, table.IsReserved, table.TableNumber, table.RestaurantID}

	err := r.db.QueryRow(ctx, query, args...).Scan(&table.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return err
		}
		return err
	}
	return nil
}

func (r *TableRepo) SetStatusById(ctx context.Context, upTable *domain.StatusTableInputSql) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE restables SET isreserved = $1 WHERE id = $2"
	_, err = tx.Exec(ctx, query, upTable.IsReserved, upTable.TableID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *TableRepo) UpdateById(ctx context.Context, upTable *domain.TableSql) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE restables SET numberofseats = $1, isreserved = $2, tablenumber = $3  WHERE id = $4"
	_, err = tx.Exec(ctx, query, upTable.NumberOfSeats, upTable.IsReserved, upTable.TableNumber, upTable.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *TableRepo) GetAllByRestaurantId(ctx context.Context, restId uuid.UUID) ([]*domain.TableStruct, error) {
	query := `Select restables.id, restables.numberofseats, restables.isreserved,restables.tablenumber, restaurants.* 
from restables 
join restaurants on restables.restaurantId = restaurants.id 
where restables.restaurantid = $1`

	rows, err := r.db.Query(ctx, query, restId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	tables := make([]*domain.TableStruct, 0)
	for rows.Next() {
		table := new(domain.TableStruct)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.Restaurant.ID,
			&table.Restaurant.Name,
			&table.Restaurant.Address,
			&table.Restaurant.Contact,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *TableRepo) Delete(ctx context.Context, tableId uuid.UUID) error {
	query := `DELETE FROM restables WHERE id = $1`

	_, err := r.db.Exec(ctx, query, tableId)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepo) GetAvailable(ctx context.Context, restid uuid.UUID) ([]*domain.TableStruct, error) {
	query := `Select restables.id, restables.numberofseats, restables.isreserved,restables.tablenumber, restaurants.* 
from restables 
join restaurants on restables.restaurantId = restaurants.id 
where restables.restaurantid = $1 and restables.isreserved = false`

	rows, err := r.db.Query(ctx, query, restid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	tables := make([]*domain.TableStruct, 0)
	for rows.Next() {
		table := new(domain.TableStruct)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.Restaurant.ID,
			&table.Restaurant.Name,
			&table.Restaurant.Address,
			&table.Restaurant.Contact,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *TableRepo) GetReserved(ctx context.Context, restid uuid.UUID) ([]*domain.TableStruct, error) {
	query := `Select restables.id, restables.numberofseats, restables.isreserved,restables.tablenumber, restaurants.* 
from restables 
join restaurants on restables.restaurantId = restaurants.id 
where restables.restaurantid = $1 and restables.isreserved = true`

	rows, err := r.db.Query(ctx, query, restid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}
	tables := make([]*domain.TableStruct, 0)
	for rows.Next() {
		table := new(domain.TableStruct)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.Restaurant.ID,
			&table.Restaurant.Name,
			&table.Restaurant.Address,
			&table.Restaurant.Contact,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
