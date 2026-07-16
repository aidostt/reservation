package table

import (
	"context"
	"dip/domain"
	repo "dip/repository"
	"time"

	"github.com/gofrs/uuid"
)

type TableService struct {
	repo         repo.Tables
	turnDuration time.Duration
}

func NewTableService(repo repo.Tables, turnDuration time.Duration) *TableService {
	return &TableService{repo: repo, turnDuration: turnDuration}
}

func (s *TableService) GetById(ctx context.Context, id string) (*domain.TableStruct, error) {
	newTableId, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(ctx, newTableId)
}

func (s *TableService) GetAll(ctx context.Context) ([]*domain.TableStruct, error) {
	return s.repo.GetAll(ctx)
}

func (s *TableService) Create(ctx context.Context, res *domain.TableInputSql) error {
	newRestaurantID, err := uuid.FromString(res.RestaurantID)
	if err != nil {
		return err
	}
	newTable := domain.TableSql{
		NumberOfSeats: res.NumberOfSeats,
		TableNumber:   res.TableNumber,
		IsReserved:    res.IsReserved,
		RestaurantID:  newRestaurantID,
	}
	return s.repo.Create(ctx, &newTable)
}

func (s *TableService) UpdateById(ctx context.Context, upTable *domain.UpdateTableInputSql) error {
	newTableId, err := uuid.FromString(upTable.TableID)
	if err != nil {
		return err
	}
	newTable := domain.TableSql{
		ID:            newTableId,
		NumberOfSeats: upTable.NumberOfSeats,
		TableNumber:   upTable.TableNumber,
		IsReserved:    upTable.IsReserved,
	}
	return s.repo.UpdateById(ctx, &newTable)
}

func (s *TableService) Delete(ctx context.Context, tableId string) error {
	newTableId, err := uuid.FromString(tableId)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, newTableId)
}

func (s *TableService) GetReserved(ctx context.Context, restid string, startAt time.Time) ([]*domain.TableStruct, error) {
	newRestID, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}
	return s.repo.GetReserved(ctx, newRestID, startAt, startAt.Add(s.turnDuration))
}

func (s *TableService) GetAvailable(ctx context.Context, restid string, startAt time.Time) ([]*domain.TableStruct, error) {
	newRestID, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAvailable(ctx, newRestID, startAt, startAt.Add(s.turnDuration))
}

func (s *TableService) GetAllByRestaurantId(ctx context.Context, restid string) ([]*domain.TableStruct, error) {
	newTableId, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllByRestaurantId(ctx, newTableId)
}
