package table

import (
	"context"
	"dip/domain"
	"dip/repository/reservation"
	"dip/repository/testsupport"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setup(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test requires Docker")
	}
	pool, cleanup, err := testsupport.SetupPostgres(context.Background())
	if err != nil {
		t.Skipf("cannot start postgres (docker unavailable?): %v", err)
	}
	t.Cleanup(cleanup)
	return pool
}

func TestTableRepo_Availability(t *testing.T) {
	ctx := context.Background()
	pool := setup(t)

	const restID = "11111111-1111-1111-1111-111111111111"
	const tableA = "22222222-2222-2222-2222-222222222222"
	const tableB = "33333333-3333-3333-3333-333333333333"
	seed := []struct {
		q    string
		args []any
	}{
		{`INSERT INTO restaurants (id,name,address,contact) VALUES ($1,'R','A','C')`, []any{restID}},
		{`INSERT INTO restables (id,numberofseats,isreserved,tablenumber,restaurantid) VALUES ($1,4,false,1,$2)`, []any{tableA, restID}},
		{`INSERT INTO restables (id,numberofseats,isreserved,tablenumber,restaurantid) VALUES ($1,4,false,2,$2)`, []any{tableB, restID}},
	}
	for _, s := range seed {
		if _, err := pool.Exec(ctx, s.q, s.args...); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	// Book table A for 19:00-21:00 through the reservation repo, so the period
	// column and overlap constraint are exercised exactly as in production.
	base := time.Date(2026, 8, 1, 19, 0, 0, 0, time.UTC)
	if err := reservation.NewReservationRepo(pool).Create(ctx, &domain.ReservationSql{
		UserID: "u", TableID: tableA, StartAt: base, EndsAt: base.Add(2 * time.Hour), PartySize: 2,
	}); err != nil {
		t.Fatalf("book table A: %v", err)
	}

	repo := NewRestaurantRepo(pool)
	rest := uuid.FromStringOrNil(restID)

	numbers := func(tables []*domain.TableStruct) []int {
		got := make([]int, 0, len(tables))
		for _, tb := range tables {
			got = append(got, int(tb.TableNumber))
		}
		sort.Ints(got)
		return got
	}

	// During the booked slot: table A is taken, only B is available.
	avail, err := repo.GetAvailable(ctx, rest, base, base.Add(2*time.Hour))
	if err != nil {
		t.Fatalf("GetAvailable during booking: %v", err)
	}
	if got := numbers(avail); !reflect.DeepEqual(got, []int{2}) {
		t.Errorf("available during booking = %v, want [2]", got)
	}

	reserved, err := repo.GetReserved(ctx, rest, base, base.Add(2*time.Hour))
	if err != nil {
		t.Fatalf("GetReserved during booking: %v", err)
	}
	if got := numbers(reserved); !reflect.DeepEqual(got, []int{1}) {
		t.Errorf("reserved during booking = %v, want [1]", got)
	}

	// Adjacent slot (starts exactly when the booking ends): both free, because
	// the interval is half-open [start, end).
	after, err := repo.GetAvailable(ctx, rest, base.Add(2*time.Hour), base.Add(4*time.Hour))
	if err != nil {
		t.Fatalf("GetAvailable adjacent slot: %v", err)
	}
	if got := numbers(after); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("available in adjacent slot = %v, want [1 2]", got)
	}
}
