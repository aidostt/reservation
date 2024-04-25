package table

import (
	"context"
	"dip/models"
	repo "dip/repository"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TableService struct {
	repo repo.Tables
}

func NewTableService(repo repo.Tables) *TableService {
	return &TableService{repo: repo}
}

func (s *TableService) GetById(ctx context.Context, id pgtype.UUID) (*models.TableStruct, error) {
	return s.repo.GetById(ctx, id)
}

func (s *TableService) GetAll(ctx context.Context) ([]*models.TableStruct, error) {
	return s.repo.GetAll(ctx)
}

func (s *TableService) Create(ctx context.Context, res *models.TableSql) error {
	return s.repo.Create(ctx, res)
}

func (s *TableService) UpdateById(ctx context.Context, upTable *models.UpdateTableInputSql) error {
	return s.repo.UpdateById(ctx, upTable)
}

func (s *TableService) MarkOccupied(ctx context.Context, tableId pgtype.UUID) error {
	table, err := s.repo.GetById(ctx, tableId)
	if err != nil {
		return fmt.Errorf("markOccupied %v, %v", err, table)
	}

	if table.IsReserved {
		return fmt.Errorf("bad request is reserved occupied")
	}

	return s.repo.SetStatusById(ctx, &models.StatusTableInputSql{TableID: table.ID, IsReserved: true})
}

func (s *TableService) MarkVacant(ctx context.Context, tableId pgtype.UUID) error {
	table, err := s.repo.GetById(ctx, tableId)
	if err != nil {
		return fmt.Errorf("markVacant %v, %v", err, table)
	}

	if !table.IsReserved {
		return fmt.Errorf("bad request is reserved not occupied")
	}

	return s.repo.SetStatusById(ctx, &models.StatusTableInputSql{TableID: table.ID, IsReserved: false})
}

func (s *TableService) Delete(ctx context.Context, tableId pgtype.UUID) error {
	return s.repo.Delete(ctx, tableId)
}

func (s *TableService) GetReserved(ctx context.Context, restid pgtype.UUID) ([]*models.TableStruct, error) {
	return s.repo.GetReserved(ctx, restid)
}

func (s *TableService) GetAvailable(ctx context.Context, restid pgtype.UUID) ([]*models.TableStruct, error) {
	return s.repo.GetAvailable(ctx, restid)
}

func (s *TableService) GetAllByRestaurantId(ctx context.Context, restId pgtype.UUID) ([]*models.TableStruct, error) {
	return s.repo.GetAllByRestaurantId(ctx, restId)
}
