package table

import (
	"context"
	"dip/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *TableRepo) GetById(ctx context.Context, id pgtype.UUID) (*models.TableSql, error) {
	query := `SELECT * FROM restables WHERE id = $1`
	var table models.TableSql
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID,
		&table.NumberOfSeats,
		&table.IsReserved,
		&table.TableNumber,
		&table.RestaurantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found in db")
		}
		return nil, err
	}

	return &table, nil
}

func (r *TableRepo) GetAll(ctx context.Context) ([]*models.TableSql, error) {
	query := "SELECT * FROM restables"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	tables := make([]*models.TableSql, 0)
	for rows.Next() {
		table := new(models.TableSql)
		err := rows.Scan(&table.ID, &table.NumberOfSeats, &table.IsReserved, &table.TableNumber, &table.RestaurantID)
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

func (r *TableRepo) Create(ctx context.Context, table *models.TableSql) error {
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

func (r *TableRepo) SetStatusById(ctx context.Context, upTable *models.StatusTableInputSql) error {
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

func (r *TableRepo) UpdateById(ctx context.Context, upTable *models.UpdateTableInputSql) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE restables SET numberofseats = $1, isreserved = $2, tablenumber = $3  WHERE id = $4"
	_, err = tx.Exec(ctx, query, upTable.NumberOfSeats, upTable.IsReserved, upTable.TableNumber, upTable.TableID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *TableRepo) GetAllByRestaurantId(ctx context.Context, restId pgtype.UUID) ([]*models.TableSql, error) {
	query := `SELECT * FROM restables WHERE restaurantid = $1`

	rows, err := r.db.Query(ctx, query, restId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found in db")
		}
		return nil, err
	}
	tables := make([]*models.TableSql, 0)
	for rows.Next() {
		table := new(models.TableSql)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.RestaurantID,
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

func (r *TableRepo) Delete(ctx context.Context, tableId pgtype.UUID) error {
	query := `DELETE FROM restables WHERE id = $1`

	_, err := r.db.Exec(ctx, query, tableId)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepo) GetAvailable(ctx context.Context, restid pgtype.UUID) ([]*models.TableSql, error) {
	query := `SELECT * FROM restables WHERE isreserved = false AND restaurantid = $1`

	rows, err := r.db.Query(ctx, query, restid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found in db")
		}
		return nil, err
	}
	tables := make([]*models.TableSql, 0)
	for rows.Next() {
		table := new(models.TableSql)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.RestaurantID,
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

func (r *TableRepo) GetReserved(ctx context.Context, restid pgtype.UUID) ([]*models.TableSql, error) {
	query := `SELECT * FROM restables WHERE isreserved = true AND restaurantid = $1`

	rows, err := r.db.Query(ctx, query, restid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found in db")
		}
		return nil, err
	}
	tables := make([]*models.TableSql, 0)
	for rows.Next() {
		table := new(models.TableSql)
		err := rows.Scan(
			&table.ID,
			&table.NumberOfSeats,
			&table.IsReserved,
			&table.TableNumber,
			&table.RestaurantID,
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
