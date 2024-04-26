package table

import (
	"context"
	"dip/models"
	repo "dip/repository"
	"fmt"

	"github.com/gofrs/uuid"
)

type TableService struct {
	repo repo.Tables
}

func NewTableService(repo repo.Tables) *TableService {
	return &TableService{repo: repo}
}

func (s *TableService) GetById(ctx context.Context, id string) (*models.TableStruct, error) {
	newTableId, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(ctx, newTableId)
}

func (s *TableService) GetAll(ctx context.Context) ([]*models.TableStruct, error) {
	return s.repo.GetAll(ctx)
}

func (s *TableService) Create(ctx context.Context, res *models.TableInputSql) error {
	newRestaurantID, err := uuid.FromString(res.RestaurantID)
	if err != nil {
		return err
	}
	newTable := models.TableSql{
		NumberOfSeats: res.NumberOfSeats,
		TableNumber:   res.TableNumber,
		IsReserved:    res.IsReserved,
		RestaurantID:  newRestaurantID,
	}
	return s.repo.Create(ctx, &newTable)
}

func (s *TableService) UpdateById(ctx context.Context, upTable *models.UpdateTableInputSql) error {
	newTableId, err := uuid.FromString(upTable.TableID)
	if err != nil {
		return err
	}
	// (ctx, query, upTable.NumberOfSeats, upTable.IsReserved, upTable.TableNumber, upTable.TableID)
	newTable := models.TableSql{
		ID:            newTableId,
		NumberOfSeats: upTable.NumberOfSeats,
		TableNumber:   upTable.TableNumber,
		IsReserved:    upTable.IsReserved,
	}
	return s.repo.UpdateById(ctx, &newTable)
}

func (s *TableService) MarkOccupied(ctx context.Context, tableId string) error {
	newTableId, err := uuid.FromString(tableId)
	if err != nil {
		return err
	}

	table, err := s.repo.GetById(ctx, newTableId)
	if err != nil {
		return fmt.Errorf("markOccupied %v, %v", err, table)
	}

	if table.IsReserved {
		return fmt.Errorf("bad request is reserved occupied")
	}

	return s.repo.SetStatusById(ctx, &models.StatusTableInputSql{TableID: table.ID, IsReserved: true})
}

func (s *TableService) MarkVacant(ctx context.Context, tableId string) error {
	newTableId, err := uuid.FromString(tableId)
	if err != nil {
		return err
	}

	table, err := s.repo.GetById(ctx, newTableId)
	if err != nil {
		return fmt.Errorf("markVacant %v, %v", err, table)
	}

	if !table.IsReserved {
		return fmt.Errorf("bad request is reserved not occupied")
	}

	return s.repo.SetStatusById(ctx, &models.StatusTableInputSql{TableID: table.ID, IsReserved: false})
}

func (s *TableService) Delete(ctx context.Context, tableId string) error {
	newTableId, err := uuid.FromString(tableId)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, newTableId)
}

func (s *TableService) GetReserved(ctx context.Context, restid string) ([]*models.TableStruct, error) {
	newTableId, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}

	return s.repo.GetReserved(ctx, newTableId)
}

func (s *TableService) GetAvailable(ctx context.Context, restid string) ([]*models.TableStruct, error) {
	newTableId, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAvailable(ctx, newTableId)
}

func (s *TableService) GetAllByRestaurantId(ctx context.Context, restid string) ([]*models.TableStruct, error) {
	newTableId, err := uuid.FromString(restid)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllByRestaurantId(ctx, newTableId)
}
