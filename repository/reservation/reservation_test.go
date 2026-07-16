package reservation

import (
	"context"
	"dip/domain"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// setupPostgres starts a throwaway PostgreSQL, applies the migrations and
// returns a connected pool. It skips (rather than fails) when Docker is not
// available, so unit-only environments are unaffected.
func setupPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("reservista"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Skipf("cannot start postgres container (docker unavailable?): %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	t.Cleanup(pool.Close)

	files, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	if err != nil || len(files) == 0 {
		t.Fatalf("glob migrations: %v (found %d)", err, len(files))
	}
	sort.Strings(files)
	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			t.Fatalf("apply %s: %v", f, err)
		}
	}
	return pool
}

func seedTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) string {
	t.Helper()
	const restID = "11111111-1111-1111-1111-111111111111"
	const tableID = "22222222-2222-2222-2222-222222222222"
	if _, err := pool.Exec(ctx,
		`INSERT INTO restaurants (id, name, address, contact) VALUES ($1, 'r', 'a', 'c')`, restID); err != nil {
		t.Fatalf("seed restaurant: %v", err)
	}
	if _, err := pool.Exec(ctx,
		`INSERT INTO restables (id, numberofseats, isreserved, tablenumber, restaurantid)
		 VALUES ($1, 4, false, 1, $2)`, tableID, restID); err != nil {
		t.Fatalf("seed table: %v", err)
	}
	return tableID
}

func TestReservationRepo_OverlapConstraint(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test requires Docker")
	}
	ctx := context.Background()
	pool := setupPostgres(t)
	tableID := seedTable(t, ctx, pool)
	repo := NewReservationRepo(pool)

	base := time.Date(2026, 8, 1, 19, 0, 0, 0, time.UTC)
	booking := func(start time.Time) *domain.ReservationSql {
		return &domain.ReservationSql{
			UserID: "u", TableID: tableID,
			StartAt: start, EndsAt: start.Add(2 * time.Hour),
			PartySize: 2, Confirmed: false,
		}
	}

	if err := repo.Create(ctx, booking(base)); err != nil {
		t.Fatalf("first booking 19:00-21:00 should succeed: %v", err)
	}

	err := repo.Create(ctx, booking(base.Add(time.Hour)))
	if !errors.Is(err, domain.ErrTableOccupied) {
		t.Fatalf("overlapping booking 20:00-22:00: got %v, want ErrTableOccupied", err)
	}

	if err := repo.Create(ctx, booking(base.Add(2*time.Hour))); err != nil {
		t.Fatalf("adjacent booking 21:00-23:00 should succeed (half-open range): %v", err)
	}
}
